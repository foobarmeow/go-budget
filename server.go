package main

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type server struct {
	devMode     bool
	port        string
	redisURL    string
	mongoURL    string
	mongoClient *mongo.Client
	redisPool   *redis.Pool
}

func (s *server) auth() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var session Session
			var sessionCookie string

			conn := s.redisPool.Get()

			// Read cookie
			cookies := req.Cookies()

			for i := range cookies {
				if cookies[i].Name == cookieString {
					sessionCookie = cookies[i].Value
				}
			}

			if sessionCookie == "" {
				log.Error("Session cookie was empty")
				w.WriteHeader(401)
				return
			}

			// Read session
			sessionString, err := redis.String(conn.Do("GET", sessionCookie))
			if err != nil {
				log.Error(err)
				w.WriteHeader(500)
				return
			}

			err = json.Unmarshal([]byte(sessionString), &session)
			if err != nil {
				log.Error(err)
				w.WriteHeader(500)
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

func (s *server) cors() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:9090")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			next.ServeHTTP(w, req)
		})
	}
}

func (s *server) serveError(w http.ResponseWriter, err error, message string) {
	log.WithFields(log.Fields{message: message}).Error(err)
	w.WriteHeader(http.StatusInternalServerError)
}
