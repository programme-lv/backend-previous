package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/lmittmann/tint"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/alexedwards/scs/v2"
	"github.com/programme-lv/backend/internal/environment"
	"github.com/programme-lv/backend/internal/graphql"
	"github.com/programme-lv/backend/internal/services/submissions"
	"github.com/programme-lv/director/msg"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/gomodule/redigo/redis"
)

const defaultPort = "3001"

func main() {
	log := slog.New(tint.NewHandler(os.Stdout, nil))

	conf := environment.ReadEnvConfig(log)
	conf.Print()

	log.Info("connecting to database...")
	sqlxDb := sqlx.MustConnect("postgres", conf.SqlxConnString)
	defer sqlxDb.Close()
	log.Info("successfully connected to database")

	log.Info("connecting to RabbitMQ...")
	rmqConn, err := amqp.Dial(conf.AMQPConnString)
	if err != nil {
		panic(err)
	}
	defer rmqConn.Close()
	log.Info("successfully connected to RabbitMQ")

	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", conf.RedisConnString)
		},
	}

	sessions := scs.New()
	sessions.Lifetime = 24 * time.Hour
	sessions.Store = redisstore.New(redisPool)

	log.Info("connecting to DigitalOcean Spaces...")
	urls, err := submissions.NewS3TestURLs(conf.DOSpacesKey, conf.DOSpacesSecret, "fra1", conf.S3Endpoint, conf.S3Bucket)
	if err != nil {
		panic(err)
	}
	if err := testTestURLs(urls); err != nil {
		log.Error("could not download test file", err)
		panic(err)
	}
	log.Info("successfully connected to DO Spaces")

	log.Info("connecting to \"director\" gRPC service...")
	conn, err := grpc.Dial(conf.DirectorEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("could not connect to director", err)
		panic(err)
		// log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	log.Info("successfully connected to \"director\" gRPC service")

	c := msg.NewDirectorClient(conn)
	md := metadata.New(map[string]string{"authorization": "...asdf"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	cc, err := c.EvaluateSubmission(ctx, &msg.EvaluationRequest{})
	if err != nil {
		log.Error("could not greet", err)
		panic(err)
	}

	for {
		res, err := cc.Recv()
		if err != nil {
			log.Error("could not greet", err)
			panic(err)
		}

		log.Info(fmt.Sprintf("%+v", res))
	}

	resolver := &graphql.Resolver{
		PostgresDB:     sqlxDb,
		SessionManager: sessions,
		Logger:         log,
		SubmissionRMQ:  rmqConn,
		TestURLs:       urls,
	}

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world :)"))
	})
	http.Handle("/query", sessions.LoadAndSave(srv))

	log.Info("Listening on port", defaultPort)
	log.Error(http.ListenAndServe(":"+defaultPort, nil).Error())
}

func testTestURLs(urls *submissions.S3TestURLs) error {
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
