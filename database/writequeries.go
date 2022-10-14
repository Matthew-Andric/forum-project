package database

import "fmt"

func EditCategoryName(id, name string) bool {
	sqlStatement := `UPDATE categories SET name=$1 WHERE categoryid=$2;`

	_, err := db.Exec(sqlStatement, name, id)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Category name updated:", id, name)

	return true
}

func EditSubCategoryLocation(subCatId, catId string) bool {
	sqlStatement := `UPDATE subcategories SET categoryid=$1 WHERE subcategoryid=$2;`

	_, err := db.Exec(sqlStatement, catId, subCatId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Sub Category location updated:", subCatId, catId)

	return true
}

func EditSubCategoryPermission(subCatId, permission string) bool {
	sqlStatement := `UPDATE subcategories SET permissionlevel=$1 WHERE subcategoryid=$2;`

	_, err := db.Exec(sqlStatement, permission, subCatId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Sub Category permission updated ID:", subCatId, "permission:", permission)

	return true
}

func PostReply(user *User, post string, threadid string) bool {
	sqlStatement := `INSERT INTO posts(threadid, userid, posttext, postdate) VALUES ($1, $2, $3, NOW());`

	_, err := db.Exec(sqlStatement, threadid, user.Userid, post)
	if err != nil {
		fmt.Println("Error posting reply:", err)
		return false
	}

	sqlStatement = `UPDATE threads SET lastupdated = postdate FROM posts WHERE postid = (SELECT MAX(postid) FROM posts) AND posts.threadid=$1;`

	_, err = db.Exec(sqlStatement, threadid)
	if err != nil {
		fmt.Println(err)
	}

	return true
}

func EditReply(user *User, postid string, post string) bool {
	sqlStatement := `UPDATE posts SET posttext=$1, editdate=NOW() WHERE userid=$2 AND postid=$3;`

	_, err := db.Exec(sqlStatement, post, user.Userid, postid)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func PostThread(user *User, subcatid string, title string, post string) (int, bool) {
	sqlStatement := `INSERT INTO threads(subcategoryid, userid, threadname, threadtext, createdate, lastupdated)
					VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING threadid;`

	newid := 0
	err := db.QueryRow(sqlStatement, subcatid, user.Userid, title, post).Scan(&newid)
	if err != nil {
		fmt.Println(err)
		return 0, false
	}

	fmt.Println("new thread:", newid)

	return int(newid), true
}

func EditThread(user *User, threadid string, name string, post string) bool {
	sqlStatement := `UPDATE threads SET threadname = $1, threadtext = $2 WHERE userid=$3 AND threadid=$4;`

	_, err := db.Exec(sqlStatement, name, post, user.Userid, threadid)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Thread edited id:", threadid)

	return true
}

func AddCategory(name, priority string) bool {
	sqlStatement := `INSERT INTO categories(name, priority) VALUES ($1, $2);`

	_, err := db.Exec(sqlStatement, name, priority)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func AddSubCategory(name, catId, permission, priority string) bool {
	sqlStatement := `INSERT INTO subcategories(name, categoryid, permissionlevel, priority) VALUES ($1, $2, $3, $4);`

	_, err := db.Exec(sqlStatement, name, catId, permission, priority)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func EditCategoryPriority(catId, priority string) bool {
	sqlStatement := `UPDATE categories SET priority=$1 WHERE categoryid=$2;`

	_, err := db.Exec(sqlStatement, priority, catId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func EditSubCategoryName(subCatId, name string) bool {
	sqlStatement := `UPDATE subcategories SET name=$1 WHERE subcategoryid=$2;`

	_, err := db.Exec(sqlStatement, name, subCatId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func EditSubCategoryPriority(subCatId, priority string) bool {
	sqlStatement := `UPDATE subcategories SET priority=$1 WHERE subcategoryid=$2;`

	_, err := db.Exec(sqlStatement, priority, subCatId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
