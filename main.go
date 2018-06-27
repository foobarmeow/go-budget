package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"html/template"
	"log"
	"net/http"
)

const cookieString string = "budget-session"

type Session struct {
	Username string
}

func main() {

	var port string

	flag.StringVar(&port, "port", ":8080", "port to listen on")

	rClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   9,
	})

	http.HandleFunc("/", homeHandler(auth(rClient)))
	http.HandleFunc("/login", loginHandler(rClient))

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))

	log.Println("Listening on", port)

	http.ListenAndServe(port, nil)
}

func auth(client *redis.Client) func(r *http.Request) (*Session, int, error) {

	return func(r *http.Request) (*Session, int, error) {

		var session Session

		// Read cookie
		sessionCookie, err := r.Cookie(cookieString)
		if err != nil || sessionCookie.Value == "" {
			return nil, 401, fmt.Errorf("No session found")
		}

		// Read session
		sessionString, err := client.Get(sessionCookie.Value).Result()
		if err != nil {
			return nil, 500, err
		}

		err = json.Unmarshal([]byte(sessionString), &session)
		if err != nil {
			return nil, 500, err
		}

		return &session, 200, nil
	}
}

func loginHandler(client *redis.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s := Session{"caleb"}

		sessionBytes, err := json.Marshal(&s)

		_, err = client.Set("9801", sessionBytes, 0).Result()
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(500)
			return
		}

		http.Redirect(w, r, "/", 301)
	}
}

func homeHandler(authCallback func(*http.Request) (*Session, int, error)) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, code, err := authCallback(r)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(code)
			return
		}

		http.SetCookie(w, &http.Cookie{Name: "budget-session", Value: "9801"})

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
