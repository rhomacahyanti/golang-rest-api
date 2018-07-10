package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"promo-rest-api/connection"
	"promo-rest-api/queries"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Post struct {
	Id                int            `json:"id, omitempty"`
	Date              string         `json:"date, omitempty"`
	Modified          string         `json:"modified, omitempty"`
	Author            int            `json:"author, omitempty"`
	Title             string         `json:"title, omitempty"`
	Slug              string         `json:"slug, omitempty"`
	Excerpt           string         `json:"excerpt, omitempty"`
	Permalink         string         `json:"link, omitempty"`
	LinkApps          string         `json:"link_apps, omitempty"`
	PromoCode         string         `json:"coupon, omitempty"`
	MultiplePromoCode MultipleCoupon `json:"multiple_coupon, omitempty"`
	StartDate         string         `json:"start_date, omitempty"`
	EndDate           string         `json:"end_date, omitempty"`
	DateText          string         `json:"date_text, omitempty"`
	MinTransaction    string         `json:"min_transaction, omitempty"`
	AppLink           string         `json:"app_link, omitempty"`
	PromoLink         string         `json:"promo_link, omitempty"`
	Images            PostImages     `json:"images, omitempty"`
	Categories        []string       `json:"categories, omitempty"`
	Tags              []string       `json:"tags, omitempty"`
}

type PostImages struct {
	ThumbnailImage string `json:"thumbnail, omitempty"`
	FeaturedImage  string `json:"banner, omitempty"`
}

type MultipleCoupon struct {
	TotalCoupon int      `json:"total_coupon, omitempty"`
	Data        []string `json:"data, omitempty"`
}

type Header struct {
	TotalData   int     `json:"total_data, omitempty"`
	ProcessTime float64 `json:"process_time, omitempty"`
}

type RestResponse struct {
	ApiHeader Header `json:"header, omitempty"`
	Data      []Post `json:"data, omitempty"`
}

var queryPosts []queries.QueryPost

var posts []Post

var categories []string

var tags []string

var images PostImages

var header Header

var restResponse RestResponse

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
	start := time.Now()

	//Connect to database
	db := connection.Connect()

	defer db.Close()

	totalData := 0

	for _, post := range queryPosts {
		//Get post image
		thumbnail := queries.QueryThumbnailImage(db, post)
		banner := queries.QueryFeatureImage(db, post)

		images.ThumbnailImage = thumbnail
		images.FeaturedImage = banner

		//Get post categories
		categories = queries.QueryCategories(db, post)

		//Get post tags
		tags = queries.QueryTags(db, post)

		//Get multiple coupon
		coupons := queries.QueryMultipleCoupon(db, post)

		//Total Coupon
		var totalCoupons int
		for i := 0; i < len(coupons); i++ {
			totalCoupons++
		}

		posts = append(posts, Post{Id: post.Id, Date: post.Date, Modified: post.Modified, Author: post.Author, Title: post.Title, Slug: post.Slug, Excerpt: post.Excerpt, Permalink: post.Permalink, LinkApps: post.LinkApps, PromoCode: post.PromoCode, MultiplePromoCode: MultipleCoupon{TotalCoupon: totalCoupons, Data: coupons}, StartDate: post.StartDate, EndDate: post.EndDate, DateText: post.DateText, MinTransaction: post.MinTransaction, AppLink: post.AppLink, PromoLink: post.PromoLink, Images: images, Categories: categories, Tags: tags})
		totalData++
	}

	//Time elapsed calculation
	time.Sleep(time.Second * 2)
	elapsed := time.Since(start)

	//Api Header
	header.ProcessTime = float64(math.Floor(elapsed.Seconds()*100) / 100)
	header.TotalData = totalData

	//Rest Response
	restResponse.ApiHeader = header
	restResponse.Data = posts

	json.NewEncoder(w).Encode(restResponse)
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
