package queries

import (
	"database/sql"
	"fmt"
)

func QueryFeatureImage(db *sql.DB, post QueryPost) string {
	var featuredImage string

	sqlStatement := `SELECT guid FROM wp_posts WHERE post_type='attachment' and id=(SELECT meta_value FROM wp_postmeta WHERE meta_key='_thumbnail_id' and post_id = ? LIMIT 1)`

	row := db.QueryRow(sqlStatement, post.Id)

	err := row.Scan(&featuredImage)
	checkError(err)

	fmt.Println("Featured Image: ", featuredImage)

	return featuredImage
}

func QueryThumbnailImage(db *sql.DB, post QueryPost) string {
	var thumbnailImage string

	sqlStatement := `SELECT wppm.meta_value FROM wp_postmeta wppm WHERE wppm.meta_key = 'thumbnail_image' AND wppm.post_id = ?  LIMIT 1`

	row := db.QueryRow(sqlStatement, post.Id)

	err := row.Scan(&thumbnailImage)
	checkError(err)

	fmt.Println("Thumbnail Image: ", thumbnailImage)

	return thumbnailImage
}
