package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"promo-rest-api/connection"
	"promo-rest-api/queries"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Post struct {
	Id             int      `json:"id, omitempty"`
	Date           string   `json:"date, omitempty"`
	Modified       string   `json:"modified, omitempty"`
	Author         int      `json:"author, omitempty"`
	Title          string   `json:"title, omitempty"`
	Slug           string   `json:"slug, omitempty"`
	Excerpt        string   `json:"excerpt, omitempty"`
	Permalink      string   `json:"link, omitempty"`
	LinkApp        string   `json:"link_apps, omitempty"`
	PromoCode      string   `json:"coupon, omitempty"`
	StartDate      string   `json:"start_date, omitempty"`
	EndDate        string   `json:"end_date, omitempty"`
	DateText       string   `json:"date_text, omitempty"`
	MinTransaction string   `json:"min_transaction, omitempty"`
	PromoLink      string   `json:"promo_link, omitempty"`
	ThumbnailImage string   `json:"thumbnail_image, omitempty"`
	Categories     []string `json:"categories, omitempty"`
}

var queryPosts []queries.QueryPost

var posts []Post

var categories []string

func main() {

	//Connect to database
	db := connection.Connect()

	defer db.Close()

	//Grab data from database
	queryPosts = queries.QueryPosts(db)

	//Create API Router
	router := mux.NewRouter()
	router.HandleFunc("/posts", getAllPost).Methods("GET")
	router.HandleFunc("/post/{id}", getPost).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", router))
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

//Get All Post
func getAllPost(w http.ResponseWriter, r *http.Request) {
	//Connect to database
	db := connection.Connect()

	defer db.Close()

	for _, post := range queryPosts {
		categories = queries.QueryCategories(db, post)
		posts = append(posts, Post{Id: post.Id, Date: post.Date, Modified: post.Modified, Author: post.Author, Title: post.Title, Slug: post.Slug, Excerpt: post.Excerpt, Permalink: post.Permalink, PromoCode: post.PromoCode, StartDate: post.StartDate, EndDate: post.EndDate, MinTransaction: post.MinTransaction, PromoLink: post.PromoLink, ThumbnailImage: post.ThumbnailImage, Categories: categories})
	}
	json.NewEncoder(w).Encode(posts)
}

//Get Single Post
func getPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, post := range posts {
		paramsID, _ := strconv.Atoi(params["id"])
		if post.Id == paramsID {
			json.NewEncoder(w).Encode(post)
			return
		}
	}
	json.NewEncoder(w).Encode(&Post{})
}
