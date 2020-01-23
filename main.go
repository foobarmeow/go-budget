package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"net/http"
	"time"
)

const cookieString string = "budget-session"

type Session struct {
	ID       string
	Username string
}

type Budget struct {
	Items []Item `json:"items"`
}

type Item struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id"`
	Name   string             `json:"name"`
	Amount int                `json:"amount"`
	Date   int                `json:"date"`
}

type server struct {
	port        string
	redisURL    string
	mongoURL    string
	mongoClient *mongo.Client
	redisPool   *redis.Pool
}

func main() {
	var err error
	s := server{
		redisURL: "localhost:6379",
		mongoURL: "mongodb://localhost:27017",
	}

	// Read port
	flag.StringVar(&s.port, "port", "8080", "port to listen on")

	s.mongoClient, err = mongo.NewClient(options.Client().ApplyURI(s.mongoURL))
	if err != nil {
		log.Fatal(err)
	}

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

func (s *server) budget(w http.ResponseWriter, req *http.Request) {
	var budget Budget
	var items []Item
	cursor, err := s.mongoClient.Database("bank").Collection("items").Find(context.Background(), bson.D{}, options.Find())
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}

	if err = cursor.All(context.Background(), &items); err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}

	budget.Items = items

	budgetJSON, err := json.Marshal(&budget)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}

	w.Write([]byte(budgetJSON))
}

func (s *server) addItem(w http.ResponseWriter, req *http.Request) {
	// Decode request into an item and save it
	var i Item
	err := json.NewDecoder(req.Body).Decode(&i)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}

	iBson, err := bson.Marshal(&i)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}

	_, err = s.mongoClient.Database("bank").Collection("items").InsertOne(context.Background(), iBson)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}
}

func (s *server) login(w http.ResponseWriter, req *http.Request) {
	conn := s.redisPool.Get()

	// Parse username and password from request
	username := req.Form.Get("username")
	if username == "" {
		return
	}

	password := req.Form.Get("password")
	if password == "" {
		return
	}

	// TODO: Validate user against db

	// Ideally this would be the username from the db, not from the request
	session := Session{
		Username: username,
		ID:       fmt.Sprintf("%v", uuid.New()),
	}

	// Save the session
	sessionBytes, err := json.Marshal(&session)
	if err != nil {
		log.Error(err)

		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}

	_, err = conn.Do("SET", session.ID, sessionBytes)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	http.SetCookie(w, &http.Cookie{Name: cookieString, Value: session.ID, Path: "/"})
}

func homeHandler(authCallback func(*http.Request) (*Session, int, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, code, err := authCallback(r)
		if err != nil {
			log.Error(err)
			w.WriteHeader(code)
			w.Write([]byte(err.Error()))
			return
		}

		log.Println("Got request to home from", session.Username)

		t := template.New("home")

		t, err = t.ParseFiles("layouts/home.html")
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(code)

		t.Execute(w, nil)

	})
}
