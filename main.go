package main

import (
	"context"
	"flag"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

const cookieString string = "budget-session"
const devAssetServer string = "http://localhost:9090/assets"

func main() {
	var err error
	s := server{
		redisURL: "localhost:6379",
		mongoURL: "mongodb://localhost:27017",
		devMode:  true,
	}

	// Read port
	flag.StringVar(&s.port, "port", "9042", "port to listen on")

	s.mongoClient, err = mongo.NewClient(options.Client().ApplyURI(s.mongoURL))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = s.mongoClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Create redis pool
	s.redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", s.redisURL)
			if err != nil {
				return nil, err
			}

			if _, err := c.Do("SELECT", "1"); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}

	router := mux.NewRouter()
	router.Path("/login").HandlerFunc(s.login)

	api := router.PathPrefix("/api").Subrouter()
	api.Methods("GET").Path("/budget").HandlerFunc(s.budget)
	api.Methods("POST").Path("/item").HandlerFunc(s.addItem)

	api.Use(s.cors())
	api.Use(s.auth())

	h := &http.Server{
		Addr:    ":" + s.port,
		Handler: router,
	}

	log.Info("Listening on :", s.port)

	if err := h.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
