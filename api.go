package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

type Session struct {
	ID       string
	Username string
}

type User struct {
	UserID   string
	Username string
	Password string
}

type Budget struct {
	Balances []Balance `json:"balances"`
	Items    []Item    `json:"items"`
}

type Balance struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Name    string             `json:"name"`
	Balance int                `json:"balance"`
}

type Item struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id"`
	Name   string             `json:"name"`
	Amount int                `json:"amount"`
	Date   int                `json:"date"`
}

func (s *server) budget(w http.ResponseWriter, req *http.Request) {
	var budget Budget
	var items []Item
	cursor, err := s.mongoClient.Database("bank").Collection("items").Find(context.Background(), bson.D{}, options.Find())
	if err != nil {
		s.serveError(w, err, "failed to find budget cursor in budget")
		return
	}

	if err = cursor.All(context.Background(), &items); err != nil {
		s.serveError(w, err, "failed to find budget in budget")
		return
	}

	budget.Items = items
	budget.Balances = []Balance{
		{
			ID:      primitive.NewObjectID(),
			Name:    "WF",
			Balance: 1900,
		},
	}

	budgetJSON, err := json.Marshal(&budget)
	if err != nil {
		s.serveError(w, err, "failed to marshal budget in budget")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(budgetJSON)
}

func (s *server) addItem(w http.ResponseWriter, req *http.Request) {
	// Decode request into an item and save it
	var i Item
	err := json.NewDecoder(req.Body).Decode(&i)
	if err != nil {
		s.serveError(w, err, "failed to unmarhsal item in addItem")
		return
	}

	i.ID = primitive.NewObjectID()

	iBson, err := bson.Marshal(&i)
	if err != nil {
		s.serveError(w, err, "failed to marhsal bson item in addItem")
		return
	}

	_, err = s.mongoClient.Database("bank").Collection("items").InsertOne(context.Background(), iBson)
	if err != nil {
		s.serveError(w, err, "failed to save item in addItem")
		return
	}
}

func (s *server) login(w http.ResponseWriter, req *http.Request) {
	conn := s.redisPool.Get()

	// TODO: Validate user against db

	// Ideally this would be the username from the db, not from the request
	session := Session{
		//Username: username,
		ID: fmt.Sprintf("%v", uuid.New()),
	}

	// Save the session
	sessionBytes, err := json.Marshal(&session)
	if err != nil {
		s.serveError(w, err, "failed to marshal session in login")
		return
	}

	_, err = conn.Do("SET", session.ID, sessionBytes)
	if err != nil {
		s.serveError(w, err, "failed to save session in login")
		return
	}

	http.SetCookie(w, &http.Cookie{Name: cookieString, Value: session.ID, Path: "/"})
}

func (s *server) createAccount(w http.ResponseWriter, req *http.Request) {
	// Get user from body
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		s.serveError(w, err, "failed to unmarshal user in createAccount")
		return
	}

	user.ID = uuid.New()

	// Encrypt password
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	// LEFTOFF

}
