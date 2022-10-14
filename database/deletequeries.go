package database

import "fmt"

func DeleteReply(user *User, postid string) bool {
	sqlStatement := `DELETE FROM posts WHERE postid=$1 AND userid=$2;`

	_, err := db.Exec(sqlStatement, postid, user.Userid)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func DeleteCategory(catId string) bool {
	sqlStatement := `DELETE FROM categories WHERE categoryid=$1;`

	_, err := db.Exec(sqlStatement, catId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
