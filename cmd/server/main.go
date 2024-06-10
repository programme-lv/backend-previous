package main

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/programme-lv/backend/config"
	"github.com/programme-lv/backend/internal/components/user"
	"github.com/programme-lv/backend/internal/database/dospaces"
	"github.com/rs/zerolog/log"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/lmittmann/tint"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/alexedwards/scs/v2"
	"github.com/programme-lv/backend/internal/graphql"
	"github.com/programme-lv/backend/internal/services/submissions"
	"github.com/programme-lv/director/msg"

	"github.com/gomodule/redigo/redis"
)

const defaultPort = "3001"

var logger *slog.Logger

func main() {
	conf := loadConfigFromEnvFile()
	logger = newColorfulLogger()

	pgDB := connectToPostgres(conf.Postgres.Host, conf.Postgres.User, conf.Postgres.Password,
		conf.Postgres.DBName, conf.Postgres.SSLMode, conf.Postgres.Port)

	sessions := initRedisAuthSessionStore(conf.Redis.Host, conf.Redis.Port)

	spaces := intS3Store(conf.S3.Endpoint, conf.S3.Key, conf.S3.Secret, conf.S3.Bucket)

	director := connectToTestDirector(conf.Director.Endpoint, conf.Director.AuthKey)

	userSrv := user.NewService(pgDB)

	gqlResolver := &graphql.Resolver{
		UserSrv:        userSrv,
		PostgresDB:     pgDB,
		SessionManager: sessions,
		Logger:         logger,
		TestURLs:       spaces,
		DirectorConn: &graphql.AuthDirectorConn{
			GRPCClient: director,
			Password:   conf.Director.AuthKey,
		},
	}

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: gqlResolver}))
	srv.AddTransport(&transport.Websocket{})
	http.Handle("/query", sessions.LoadAndSave(srv))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello world :)"))
		if err != nil {
			logger.Error("could not write response", "error", err)
			return
		}
	})

	address := fmt.Sprintf(":%v", conf.Server.Port)
	logger.Info("listening on", "address", address)
	err := http.ListenAndServe(address, nil)
	logger.Error("http listener has returned an error", "error", err)
}

func testTestURLs(urls submissions.TestDownloadURLProvider) error {
	url, err := urls.GetTestDownloadURL("test")
	if err != nil {
		return err
	}

	// try to download the file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// compare the value to "test"
	buf := make([]byte, 4)
	_, err = resp.Body.Read(buf)
	if err != nil {
		return err
	}

	if string(buf) != "test" {
		return fmt.Errorf("expected 'test', got '%s'", string(buf))
	}

	return nil
}

func testConnToDirector(conn *grpc.ClientConn, directorAuthKey string) error {
	c := msg.NewDirectorClient(conn)
	md := metadata.New(map[string]string{"authorization": directorAuthKey})
	ctx2, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	ctx := metadata.NewOutgoingContext(ctx2, md)
	cc, err := c.EvaluateSubmission(ctx, &msg.EvaluationRequest{})
	if err != nil {
		return err
	}

	for {
		res, err := cc.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		slog.Debug(fmt.Sprintf("%+v", res))
	}

	return nil
}

func newColorfulLogger() *slog.Logger {
	opts := &tint.Options{
		Level: slog.LevelDebug,
	}
	return slog.New(tint.NewHandler(os.Stdout, opts))
}

func loadConfigFromEnvFile() *config.Config {
	conf, err := config.LoadConfig(".env")
	if err != nil {
		logger.Error("could not load config", "error", err)
		panic(err)
	}
	conf.Print()
	return conf
}

func connectToPostgres(host, user, password, dbname, sslmode string, port int) *sqlx.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	logger.Info("connecting to PostgreSQL", "connection_string", connStr)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		logger.Error("could not connect to PostgreSQL", "error", err)
		panic(err)
	}
	logger.Info("successfully connected to PostgreSQL")

	return db
}

func initRedisAuthSessionStore(redisHost string, redisPort int) *scs.SessionManager {
	logger.Info("connecting to Redis to store sessions...")
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", redisHost, redisPort))
		},
	}

	sessions := scs.New()
	sessions.Lifetime = 24 * time.Hour
	sessions.Store = redisstore.New(redisPool)

	dial, err := redisPool.Dial()
	if err != nil {
		return nil
	}
	err = dial.Close()
	if err != nil {
		logger.Error("could not close connection to Redis", "error", err)
		panic(err)
	}

	logger.Info("successfully connected to Redis")

	return sessions
}

func intS3Store(endpoint, key, secret, bucket string) *dospaces.DOSpacesS3ObjStorage {
	logger.Info("connecting to DigitalOcean Spaces...")
	spacesConnParams := dospaces.DOSpacesConnParams{
		AccessKey: key,
		SecretKey: secret,
		Region:    "fra1",
		Endpoint:  endpoint,
		Bucket:    bucket,
	}
	spacesConn, err := dospaces.NewDOSpacesConn(spacesConnParams)
	if err != nil {
		logger.Error("could not connect to DO Spaces", "error", err)
		panic(err)
	}

	if err = testTestURLs(spacesConn); err != nil {
		logger.Error("could not download test file", "error", err)
		log.Error().Msgf("could not download test file: %v", err)
		panic(err)
	}
	logger.Info("successfully connected to DO Spaces")

	return spacesConn
}

func connectToTestDirector(endpoint, authKey string) msg.DirectorClient {
	logger.Info("connecting to \"director\" gRPC service...")
	gConn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	if err != nil {
		logger.Error("could not connect to director", "error", err)
		panic(err)
		// logger.Fatalf("did not connect: %v", err)
	}
	defer func(gConn *grpc.ClientConn) {
		err := gConn.Close()
		if err != nil {
			logger.Error("could not close connection to director", "error", err)
			panic(err)
		}
	}(gConn)
	logger.Info("successfully connected to \"director\" gRPC service")

	logger.Info("testing connection to tester")
	err = testConnToDirector(gConn, authKey)
	if err != nil {
		logger.Error("could not test connection to tester", "error", err)
		panic(err)
	}
	logger.Info("successfully tested connection to tester")

	return msg.NewDirectorClient(gConn)
}
