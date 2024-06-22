package main

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/google/uuid"
	"github.com/programme-lv/backend/config"
	"github.com/programme-lv/backend/internal/common/database/dospaces"
	"github.com/programme-lv/backend/internal/eval"
	submission2 "github.com/programme-lv/backend/internal/eval"
	mygraphql "github.com/programme-lv/backend/internal/graphql"
	"github.com/programme-lv/backend/internal/lang"
	"github.com/programme-lv/backend/internal/task"
	"github.com/programme-lv/backend/internal/user"
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
	"github.com/programme-lv/director/msg"

	"github.com/gomodule/redigo/redis"
)

var logger *slog.Logger

func main() {
	//defer profile.Start().Stop()

	logger = newColorfulLogger()
	slog.SetDefault(logger)

	conf := loadConfigFromEnvFile()

	pgDB := connectToPostgres(conf.Postgres.Host, conf.Postgres.User, conf.Postgres.Password,
		conf.Postgres.DBName, conf.Postgres.SSLMode, conf.Postgres.Port)

	sessions := initRedisAuthSessionStore(conf.Redis.Host, conf.Redis.Port)

	spaces := intS3Store(conf.S3.Endpoint, conf.S3.Key, conf.S3.Secret, conf.S3.Bucket)

	director := connectToTestDirector(conf.Director.Endpoint, conf.Director.AuthKey)

	userSrv := user.NewService(pgDB)
	taskSrv := task.NewService(userSrv, pgDB)
	submSrv := eval.NewService(pgDB, taskSrv)
	languages := lang.NewService(pgDB)

	gqlResolver := &mygraphql.Resolver{
		Languages:      languages,
		UserSrv:        userSrv,
		TaskSrv:        taskSrv,
		SubmSrv:        submSrv,
		SessionManager: sessions,
		Logger:         logger,
		TestURLs:       spaces,
		DirectorConn: &mygraphql.AuthDirectorConn{
			GRPCClient: director,
			Password:   conf.Director.AuthKey,
		},
	}

	srv := handler.NewDefaultServer(mygraphql.NewExecutableSchema(mygraphql.Config{Resolvers: gqlResolver}))
	srv.AddTransport(&transport.Websocket{})
	srv.AroundOperations(timeElapsedMiddleware)
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

func timeElapsedMiddleware(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	id := uuid.New()
	reqLogger := logger.With("request_id", id)
	rawQuery := graphql.GetOperationContext(ctx).RawQuery
	//rawQuery = strings.Replace(rawQuery, "\n", "", -1)
	//rawQuery = strings.Replace(rawQuery, "\t", "", -1)
	//rawQuery = strings.Replace(rawQuery, " ", "", -1)
	rawQuery = shortenStr(rawQuery)
	reqLogger.Info("received request", "query", rawQuery, "variables", fmt.Sprintf("%+v", graphql.GetOperationContext(ctx).Variables))
	start := time.Now()
	nxt := next(ctx)
	return func(ctxInner context.Context) *graphql.Response {
		res := nxt(ctxInner)
		elapsed := time.Since(start)
		resJson, err := res.Data.MarshalJSON()
		if err != nil {
			panic(err)
		}
		reqLogger.Info("request completed", "elapsed", elapsed, "result", string(resJson), "errors", fmt.Sprintf("%+v", res.Errors.Unwrap()))
		return res
	}
}

func shortenStr(str string) string {
	res := ""

	for i := 0; i < len(str); i++ {
		if str[i] != ' ' {
			res += string(str[i])
		} else {
			if i == 0 || str[i] != str[i-1] {
				res += string(str[i])
			}
		}
	}

	return res
}

func testTestURLs(urls submission2.TestDownloadURLProvider) error {
	url, err := urls.GetTestDownloadURL("test")
	if err != nil {
		return err
	}

	// try to download the file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("could not close response body", "error", err)
		}
	}(resp.Body)

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
	ctx2, f := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer f()
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
	connStr := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	logger.Info("connecting to PostgreSQL")
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
			return redis.Dial("tcp", fmt.Sprintf("%s:%v", redisHost, redisPort))
		},
	}

	conn, err := redisPool.Dial()
	if err != nil {
		logger.Error("could not connect to Redis", "error", err)
		panic(err)
	}
	defer func(conn redis.Conn) {
		err = conn.Close()
		if err != nil {
			logger.Error("could not close connection to Redis", "error", err)
		}
	}(conn)

	// Set a test value
	_, err = conn.Do("SET", "test_key", "test_value")
	if err != nil {
		logger.Error("could not set value in Redis", "error", err)
	}

	// Get the test value
	value, err := redis.String(conn.Do("GET", "test_key"))
	if err != nil {
		logger.Error("could not get value from Redis", "error", err)
	}

	if value != "test_value" {
		logger.Error("unexpected value from Redis", "value", value)
	}

	sessions := scs.New()
	sessions.Lifetime = 24 * time.Hour
	sessions.Store = redisstore.New(redisPool)

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
