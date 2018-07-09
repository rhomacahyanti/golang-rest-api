package queries

import (
	"database/sql"
	"fmt"
)

func QueryTags(db *sql.DB, post QueryPost) []string {

	//fmt.Println("Post ID: ", post.Id, " Title: ", post.Title)
	// SELECT MAXt.name FROM `wp_terms` t JOIN `wp_term_taxonomy` tt ON(t.`term_id` = tt.`term_id`) JOIN `wp_term_relationships` ttr ON(ttr.`term_taxonomy_id` = tt.`term_taxonomy_id`) WHERE tt.`taxonomy` = 'category' AND ttr.`object_id` = wpp.ID
	rows, err := db.Query("SELECT t.slug from wp_terms t, wp_term_taxonomy tt, wp_term_relationships tr WHERE t.term_id = tt.term_id AND tt.taxonomy = 'post_tag' AND tt.term_taxonomy_id = tr.term_taxonomy_id and tr.object_id=?", post.Id)
	checkError(err)

	tags := make([]string, 0, 5)
	var tag string

	for rows.Next() {
		err = rows.Scan(&tag)
		checkError(err)

		fmt.Println("Post ID: ", post.Id, " Tag: ", tag)
		tags = append(tags, tag)
	}
	fmt.Println("-------")

	return tags
}
