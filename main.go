package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

// let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database
var Articles []Article

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	// fmt.Fprintf(w, "Key: "+key)
	// Loop over all of our Articles
	// if the article.Id equals the key we pass in
	// return the article encoded as JSON
	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)
	// update our global Articles array to include
	// our new Article
	Articles = append(Articles, article)

	json.NewEncoder(w).Encode(article)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	// once again, we will need to parse the path parameters
	vars := mux.Vars(r)
	// we will need to extract the `id` of the article we
	// wish to delete
	id := vars["id"]

	// we then need to loop through all our articles
	for index, article := range Articles {
		// if our id path parameter matches one of our
		// articles
		if article.Id == id {
			// updates our Articles array to remove the
			// article
			json.NewEncoder(w).Encode(article)
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	// add our articles route and map it to our
	// returnAllArticles function like so
	myRouter.HandleFunc("/articles", returnAllArticles)
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	Articles = []Article{
		{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}
	handleRequests()
}