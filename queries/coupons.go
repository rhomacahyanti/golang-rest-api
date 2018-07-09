package queries

import "database/sql"

func QueryMultipleCoupon(db *sql.DB, post QueryPost) []string {
	rows, err := db.Query("SELECT meta_value FROM wp_postmeta WHERE meta_key = 'promo_codes' AND post_id = ?", post.Id)
	checkError(err)

	coupons := make([]string, 0, 5)
	var coupon string

	for rows.Next() {
		err = rows.Scan(&coupon)
		checkError(err)

		coupons = append(coupons, coupon)
	}

	return coupons
}
