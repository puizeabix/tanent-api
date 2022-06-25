package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"zeabix.com/tanent-api/account"
)

func main() {
	var (
		httpAddr    = flag.String("http.addr", ":8080", "HTTP listen address")
		mongoUrl    = getEnv("MONGO_CONNNECTION_URL", "mongodb://localhost:27017") //flag.String("mongo.url", "mongodb://localhost:27017", "Connection URL for mongodb")
		mongoDbname = getEnv("MONGO_DATABASE_NAME", "tanentv1")                    //flag.String("mongo.dbname", "blogs", "Mongo Database name")
		mongoCol    = getEnv("MONGO_COLLECTION_NAME", "accounts")                  //flag.String("mongo.colname", "blogs", "Mongo Collection name")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	ctx := context.TODO()
	client, err := makeMongoClient(ctx, mongoUrl)

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		logger.Log("Unable to connect to DB, shutdown")
		panic("Unable to connect to DB")
	}

	col := client.Database(mongoDbname).Collection(mongoCol)

	if err != nil {
		panic(err)
	}

	fieldKeys := []string{"method"}

	var as account.Service
	{
		as = account.NewAccountService(col)
		as = account.NewLoggingMiddleware(logger, as)
		as = account.NewInstrumentingMiddleware(
			kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "api",
				Subsystem: "tanent",
				Name:      "request_count",
				Help:      "Number of requests recieved",
			}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "api",
				Subsystem: "tanent",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
			as,
		)
	}

	var h http.Handler
	{
		h = account.MakeHTTPServerHandler(as, logger)
	}

	mux := http.NewServeMux()

	mux.Handle("/tanent/v1/", h)
	http.Handle("/", accessControl(mux))

	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()

	logger.Log("exit", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func makeMongoClient(ctx context.Context, url string) (*mongo.Client, error) {
	return mongo.Connect(ctx, options.Client().ApplyURI(url))
}

func getEnv(env string, def string) string {
	e := os.Getenv(env)
	if e == "" {
		return def
	}

	return e
}
