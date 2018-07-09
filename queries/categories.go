package queries

import (
	"database/sql"
	"fmt"
)

func QueryCategories(db *sql.DB, post QueryPost) []string {

	//fmt.Println("Post ID: ", post.Id, " Title: ", post.Title)
	// SELECT MAXt.name FROM `wp_terms` t JOIN `wp_term_taxonomy` tt ON(t.`term_id` = tt.`term_id`) JOIN `wp_term_relationships` ttr ON(ttr.`term_taxonomy_id` = tt.`term_taxonomy_id`) WHERE tt.`taxonomy` = 'category' AND ttr.`object_id` = wpp.ID
	rows, err := db.Query("SELECT t.slug from wp_terms t, wp_term_taxonomy tt, wp_term_relationships tr WHERE t.term_id = tt.term_id AND tt.taxonomy = 'category' AND tt.term_taxonomy_id = tr.term_taxonomy_id and tr.object_id=?", post.Id)
	checkError(err)

	categories := make([]string, 0, 5)
	var cat string

	for rows.Next() {
		err = rows.Scan(&cat)
		checkError(err)

		fmt.Println("Post ID: ", post.Id, " Category: ", cat)
		categories = append(categories, cat)
	}
	fmt.Println("-------")

	return categories
}
