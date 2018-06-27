package main

import (
	"net/http"
	"html/template"
	"log"
)


func main() {

	http.HandleFunc("/", homeHandler(auth))
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))


	http.ListenAndServe(":8000", nil)
	

	
}

func auth(r *http.Request) error {
	return nil
}

func homeHandler(authCallback func(*http.Request) error) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := authCallback(r)
		if err != nil {
			return
		}

		t := template.New("home")
	
		t, err = t.ParseFiles("layouts/home.html")
		if err != nil {
			log.Fatal(err)
		}

		t.Execute(w, nil)

	})

}
