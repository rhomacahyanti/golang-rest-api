package queries

import (
	"database/sql"
	"fmt"
	"time"
)

type QueryPost struct {
	Id             int
	Date           string
	Modified       string
	Author         int
	Title          string
	Slug           string
	Excerpt        string
	Permalink      string
	LinkApp        string
	PromoCode      string
	StartDate      string
	EndDate        string
	DateText       string
	MinTransaction string
	PromoLink      string
	ThumbnailImage string
}

var posts []QueryPost

func QueryPosts(db *sql.DB) []QueryPost {
	rows, err := db.Query("SELECT wpp.ID, wpp.post_date, wpp.post_modified, wpp.post_author, wpp.post_title, wpp.post_name, wpp.post_excerpt, REPLACE( REPLACE( REPLACE( REPLACE( wpo.option_value, '%year%', DATE_FORMAT(wpp.post_date,'%Y') ), '%monthnum%', DATE_FORMAT(wpp.post_date, '%m') ), '%day%', DATE_FORMAT(wpp.post_date, '%d') ), '%postname%', wpp.post_name ) AS permalink, (SELECT wppm.meta_value FROM wp_postmeta wppm WHERE wppm.meta_key = 'kodepromo' AND wpp.ID = wppm.post_id LIMIT 1) AS promo_code, (SELECT wppm.meta_value FROM wp_postmeta wppm WHERE wppm.meta_key = 'start_date' AND wpp.ID = wppm.post_id LIMIT 1) AS start_date, (SELECT wppm.meta_value FROM wp_postmeta wppm WHERE wppm.meta_key = 'end_date' AND wpp.ID = wppm.post_id LIMIT 1) AS end_date, (SELECT wppm.meta_value FROM wp_postmeta wppm WHERE wppm.meta_key = 'min_transaction' AND wpp.ID = wppm.post_id LIMIT 1) AS min_transaction, (SELECT wppm.meta_value FROM wp_postmeta wppm WHERE wppm.meta_key = 'promo_link' AND wpp.ID = wppm.post_id LIMIT 1) AS promo_link, (SELECT wppm.meta_value FROM wp_postmeta wppm WHERE wppm.meta_key = 'thumbnail_image' AND wpp.ID = wppm.post_id LIMIT 1) AS thumbnail_image FROM wp_posts wpp JOIN wp_options wpo ON wpo.option_name = 'permalink_structure' WHERE wpp.post_type = 'post' AND wpp.post_status = 'publish'")
	checkError(err)

	// posts := make([]Post, 0, 20)
	var post QueryPost

	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Date, &post.Modified, &post.Author, &post.Title, &post.Slug, &post.Excerpt, &post.Permalink, &post.PromoCode, &post.StartDate, &post.EndDate, &post.MinTransaction, &post.PromoLink, &post.ThumbnailImage)
		checkError(err)

		post.Permalink = "http://promo.test" + post.Permalink
		post.LinkApp = post.Permalink + "?flag_app=1"

		const layout = "2 Jan 2006"
		start, _ := time.Parse(layout, post.StartDate)
		end, _ := time.Parse(layout, post.EndDate)

		post.DateText = start.String() + " - " + end.String()
		checkError(err)

		fmt.Println(post.Title)

		posts = append(posts, post)
	}

	return posts
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
