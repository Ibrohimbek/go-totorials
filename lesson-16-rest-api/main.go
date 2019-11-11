package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Article is a main type for articles
type Article struct {
	Id      string `json: "id"`
	Title   string `json: "title"`
	Desc    string `json: "desc"`
	Content string `json: "content"`
}

// Articles - let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database
type Articles []Article

var articles Articles

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("All articles endpoint requested!")

	err := json.NewEncoder(w).Encode(articles)
	if err != nil {
		panic(err)
	}
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Single article enpoint hit!")
	vars := mux.Vars(r)
	key := vars["id"]

	for _, article := range articles {
		if article.Id == key {
			err := json.NewEncoder(w).Encode(article)
			if err != nil {
				panic(err)
			}
		}
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article Article
	err := json.Unmarshal(reqBody, &article)
	if err != nil {
		panic(err)
	}

	articles = append(articles, article)
	err = json.NewEncoder(w).Encode(article)
	if err != nil {
		panic(err)
	}
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for index, article := range articles {
		if article.Id == key {
			articles = append(articles[:index], articles[index+1:]...)
		}
	}
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)

	var a Article
	err := json.Unmarshal(reqBody, &a)
	if err != nil {
		panic(err)
	}
	fmt.Println(a)

	for i, article := range articles {
		if article.Id == key {
			articles[i] = a
		}
	}

	err = json.NewEncoder(w).Encode(articles)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "Homepage Enpoint Hit")
	if err != nil {
		panic(err)
	}

}

func handleRequests() {
	myRoute := mux.NewRouter().StrictSlash(true)
	myRoute.HandleFunc("/", homePage)
	myRoute.HandleFunc("/articles", createNewArticle).Methods("POST")
	myRoute.HandleFunc("/articles", returnAllArticles)
	myRoute.HandleFunc("/articles/{id}", deleteArticle).Methods("DELETE")
	myRoute.HandleFunc("/articles/{id}", updateArticle).Methods("PUT")
	myRoute.HandleFunc("/articles/{id}", returnSingleArticle)

	log.Fatal(http.ListenAndServe(":8081", myRoute))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers!")

	articles = Articles{
		Article{Id: "1", Title: "Rest API", Desc: "Rest API in Golang", Content: "Test service to build REST API in Golang!"},
		Article{Id: "2", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		Article{Id: "3", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}
	handleRequests()
}
