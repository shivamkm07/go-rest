package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"os"
	"path/filepath"
	"regexp"
	"time"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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

var app_dir string = "/home/shivamkm07/codes/example-node-js-app"

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

func find_file(regex string, dir string) string{

	re := regexp.MustCompile(regex);
	var file string
	
	walk := func(fn string, fi os.FileInfo, err error) error {
		if re.MatchString(fn) == false {
			return nil
		}
		file = fn
		return nil
	}
	filepath.Walk(dir, walk)
	return file

}
func handleNodeServer(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	pid := vars["pid"]
	app := "kill"
	a0 := "-USR2"
	a1 := pid

    cmd := exec.Command(app,a0,a1)
    stdout, err := cmd.Output()

    if err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Println(string(stdout))
	file := find_file(`.*snapshot` , app_dir )
	fmt.Printf(file)
	time.Sleep(time.Second)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.ServeFile(w, r, file)

	e := os.Remove(file)
    if e != nil {
        log.Fatal(e)
    }
}

type Command struct {
	Id      int `json:"id"`
	Method   string `json:"method"`
}
func profile(w http.ResponseWriter, r *http.Request){
	var Dialer websocket.Dialer
	ws, _, _ := Dialer.Dial("ws://127.0.0.1:9229/a12ec86c-f4ee-4b4f-b28f-8322ec50da8d", nil)
	var profileEnableCommand = &Command{
		Id: 1,
		Method: "Profiler.enable"}
	res, _ := json.Marshal(profileEnableCommand)
	ws.WriteMessage(websocket.TextMessage, []byte(res))
	_, message, _ := ws.ReadMessage()
	fmt.Fprintf(w, string(message))
	var profileStartCommand = &Command{
		Id: 2,
		Method: "Profiler.start"}
	res2, _ := json.Marshal(profileStartCommand)
	ws.WriteMessage(websocket.TextMessage, []byte(res2))
	_, message2, _ := ws.ReadMessage()
	fmt.Fprintf(w, string(message2))
	var profileStopCommand = &Command{
		Id: 3,
		Method: "Profiler.stop"}
	res3, _ := json.Marshal(profileStopCommand)
	ws.WriteMessage(websocket.TextMessage, []byte(res3))
	_, message3, _ := ws.ReadMessage()
	fmt.Fprintf(w, string(message3))	
	ws.Close()
	// fmt.Fprintf(w, string(res))
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
	myRouter.HandleFunc("/server/node/{pid}",handleNodeServer)
	myRouter.HandleFunc("/profiles", profile)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	Articles = []Article{
		{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}
	handleRequests()
}
