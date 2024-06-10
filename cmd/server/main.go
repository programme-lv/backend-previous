package main

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/programme-lv/backend/config"
	"github.com/programme-lv/backend/internal/database/dospaces"
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

func main() {
	conf := loadConfigFromEnvFile()
	logger := newColorfulLogger()
	pgDB := mustConnectToPostgres()
	logger.Info("connecting to database...")
	sqlxDb := sqlx.MustConnect("postgres", conf.SqlxConnString)
	defer func(sqlxDb *sqlx.DB) {
		err := sqlxDb.Close()
		if err != nil {
			logger.Error("could not close database connection", "error", err)
		}
	}(sqlxDb)
	logger.Info("successfully connected to database")

	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", conf.RedisConnString)
		},
	}

	sessions := scs.New()
	sessions.Lifetime = 24 * time.Hour
	sessions.Store = redisstore.New(redisPool)

	logger.Info("connecting to DigitalOcean Spaces...")
	spacesConnParams := dospaces.DOSpacesConnParams{
		AccessKey: conf.DOSpacesKey,
		SecretKey: conf.DOSpacesSecret,
		Region:    "fra1",
		Endpoint:  conf.S3Endpoint,
		Bucket:    conf.S3Bucket,
	}
	spacesConn, err := dospaces.NewDOSpacesConn(spacesConnParams)
	if err != nil {
		panic(err)
	}

	if err := testTestURLs(spacesConn); err != nil {
		logger.Error("could not download test file", "error", err)
		panic(err)
	}
	logger.Info("successfully connected to DO Spaces")

	logger.Info("connecting to \"director\" gRPC service...")
	gConn, err := grpc.Dial(conf.DirectorEndpoint, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	if err != nil {
		logger.Error("could not connect to director", "error", err)
		panic(err)
		// logger.Fatalf("did not connect: %v", err)
	}
	defer gConn.Close()
	logger.Info("successfully connected to \"director\" gRPC service")

	logger.Info("testing connection to tester")
	err = testConnToDirector(gConn, conf.DirectorAuthKey)
	if err != nil {
		logger.Error("could not test connection to tester", "error", err)
		panic(err)
	}
	logger.Info("successfully tested connection to tester")

	userRepo := repository.NewUserRepoPostgreSQLImpl(sqlxDb)
	userSrv := service.NewUserService(userRepo, slog.Default())

	resolver := &graphql.Resolver{
		UserSrv:        userSrv,
		AuthState:      nil,
		PostgresDB:     sqlxDb,
		SessionManager: sessions,
		Logger:         logger,
		// SubmissionRMQ:  rmqConn,
		TestURLs: spacesConn,
		DirectorConn: &graphql.AuthDirectorConn{
			GRPCClient: msg.NewDirectorClient(gConn),
			Password:   conf.DirectorAuthKey,
		},
	}

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))
	srv.AddTransport(&transport.Websocket{}) // <---- This is the important part!

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world :)"))
	})
	http.Handle("/query", sessions.LoadAndSave(srv))

	logger.Info("Listening on", "port", defaultPort)
	logger.Error(http.ListenAndServe(":"+defaultPort, nil).Error())
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

func loadConfigFromEnvFile() *config.EnvConfig {
	conf := config.ReadEnvConfig(slog.Default())
	conf.Print()
	return conf
}

func connectToPostgres(host, port, user, password, dbname string) *sqlx.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, "disable")
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}
