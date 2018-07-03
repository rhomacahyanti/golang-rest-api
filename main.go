package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Post struct {
	Id             int    `json:"id, omitempty"`
	Date           string `json:"date, omitempty"`
	Modified       string `json:"modified, omitempty"`
	Author         int    `json:"author, omitempty"`
	Title          string `json:"title, omitempty"`
	Slug           string `json:"slug, omitempty"`
	Excerpt        string `json:"excerpt, omitempty"`
	Permalink      string `json:"link, omitempty"`
	LinkApp        string `json:"link_apps, omitempty"`
	PromoCode      string `json:"coupon, omitempty"`
	StartDate      string `json:"start_date, omitempty"`
	EndDate        string `json:"end_date, omitempty"`
	DateText       string `json:"date_text, omitempty"`
	MinTransaction string `json:"min_transaction, omitempty"`
	PromoLink      string `json:"promo_link, omitempty"`
	ThumbnailImage string `json:"thumbnail_image, omitempty"`
	Categories     string `json:"categories, omitempty"`
}

type Category struct {
}

var posts []Post

func main() {

	//Connect to database
	db, err := sql.Open("mysql", "root:root@/promo")

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully connect to the database!")
	}

	db.Ping()

	defer db.Close()

	//Grab data from database
	posts = grabData(db)

	for _, post := range posts {
		fmt.Println("Post ID: ", post.Id, " Title: ", post.Title)
	}

	//Create API Router
	router := mux.NewRouter()
	router.HandleFunc("/posts", getAllPost).Methods("GET")
	router.HandleFunc("/post/{id}", getPost).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", router))
}

func grabData(db *sql.DB) []Post {
	rows, err := db.Query("SELECT wpp.ID, wpp.post_date, wpp.post_modified, wpp.post_author, wpp.post_title, wpp.post_name, wpp.post_excerpt, REPLACE( REPLACE( REPLACE( REPLACE( wpo.option_value, '%year%', DATE_FORMAT(wpp.post_date,'%Y') ), '%monthnum%', DATE_FORMAT(wpp.post_date, '%m') ), '%day%', DATE_FORMAT(wpp.post_date, '%d') ), '%postname%', wpp.post_name ) AS permalink, (SELECT wppm.meta_value FROM wp_postmeta wppm WHERE wppm.meta_key = 'kodepromo' AND wpp.ID = wppm.post_id LIMIT 1) AS promo_code, (SELECT wppm.meta_value FROM wp_postmeta wppm WHERE wppm.meta_key = 'start_date' AND wpp.ID = wppm.post_id LIMIT 1) AS start_date, (SELECT wppm.meta_value FROM wp_postmeta wppm WHERE wppm.meta_key = 'end_date' AND wpp.ID = wppm.post_id LIMIT 1) AS end_date, (SELECT wppm.meta_value FROM wp_postmeta wppm WHERE wppm.meta_key = 'min_transaction' AND wpp.ID = wppm.post_id LIMIT 1) AS min_transaction, (SELECT wppm.meta_value FROM wp_postmeta wppm WHERE wppm.meta_key = 'promo_link' AND wpp.ID = wppm.post_id LIMIT 1) AS promo_link, (SELECT wppm.meta_value FROM wp_postmeta wppm WHERE wppm.meta_key = 'thumbnail_image' AND wpp.ID = wppm.post_id LIMIT 1) AS thumbnail_image, (SELECT MAX(t.name) FROM `wp_terms` t JOIN `wp_term_taxonomy` tt ON(t.`term_id` = tt.`term_id`) JOIN `wp_term_relationships` ttr ON(ttr.`term_taxonomy_id` = tt.`term_taxonomy_id`) WHERE tt.`taxonomy` = 'category' AND ttr.`object_id` = wpp.ID) AS categories FROM wp_posts wpp JOIN wp_options wpo ON wpo.option_name = 'permalink_structure' WHERE wpp.post_type = 'post' AND wpp.post_status = 'publish'")
	checkError(err)

	posts := make([]Post, 0, 20)
	var post Post

	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Date, &post.Modified, &post.Author, &post.Title, &post.Slug, &post.Excerpt, &post.Permalink, &post.PromoCode, &post.StartDate, &post.EndDate, &post.MinTransaction, &post.PromoLink, &post.ThumbnailImage, &post.Categories)
		checkError(err)

		post.Permalink = "http://promo.test" + post.Permalink
		post.LinkApp = post.Permalink + "?flag_app=1"

		const layout = "2 Jan 2006"
		start, _ := time.Parse(layout, post.StartDate)
		end, _ := time.Parse(layout, post.EndDate)

		post.DateText = start.String() + " - " + end.String()
		checkError(err)

		posts = append(posts, post)
	}

	return posts
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

//Get All Post
func getAllPost(w http.ResponseWriter, r *http.Request) {
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
