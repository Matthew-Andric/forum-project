package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func GetCategories(user *User) []Category {
	var categories []Category
	sqlStatement := `SELECT categories.categoryid, categories.name, subcategoryid, subcategories.name, permissionlevel 
					FROM subcategories JOIN categories ON subcategories.categoryid = categories.categoryid 
					WHERE permissionlevel <=$1;`
	p := 0
	if user != nil {
		p = user.PermissionLevel
	}

	rows, err := db.Query(sqlStatement, p)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	m := make(map[string]Category)
	var o []string

	for rows.Next() {
		var cat Category
		var subcat Subcategory

		err = rows.Scan(&cat.CategoryId, &cat.CategoryName, &subcat.SubCategoryId, &subcat.SubCategoryName, &subcat.PermissionLevel)
		if err != nil {
			panic(err)
		}

		if value, ok := m[cat.CategoryName]; ok {
			value.Subcategories = append(value.Subcategories, Subcategory{subcat.SubCategoryId, subcat.SubCategoryName, subcat.PermissionLevel})
			m[cat.CategoryName] = value
		} else {
			cat.Subcategories = append(cat.Subcategories, Subcategory{subcat.SubCategoryId, subcat.SubCategoryName, subcat.PermissionLevel})
			o = append(o, cat.CategoryName)
			m[cat.CategoryName] = cat
		}
	}

	for _, v := range o {
		categories = append(categories, m[v])
	}

	return categories
}

func GetSubCategoryDetails(id string) Subcategory {
	sqlStatement := `SELECT subcategoryid, name, permissionlevel FROM subcategories WHERE subcategoryid=$1;`
	var subcat Subcategory

	row := db.QueryRow(sqlStatement, id)

	switch err := row.Scan(&subcat.SubCategoryId, &subcat.SubCategoryName, &subcat.PermissionLevel); err {
	case sql.ErrNoRows:
		fmt.Println("No rows returned")
	case nil:

	default:
		panic(err)
	}

	return subcat
}

func GetThreadsSummary(id string) []ThreadSummary {
	var threads []ThreadSummary
	sqlStatement := `SELECT threadid, threadname, users.userid, username, createdate, lastupdated FROM threads JOIN users ON threads.userid = users.userid WHERE subcategoryid=$1;`

	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var threadid, userid int
		var threadname, username string
		var createdate, lastupdated time.Time

		err = rows.Scan(&threadid, &threadname, &userid, &username, &createdate, &lastupdated)
		if err != nil {
			panic(err)
		}

		threads = append(threads, ThreadSummary{threadid, threadname, userid, username, createdate, lastupdated})
	}

	return threads
}

func GetUserDetails(id string) User {
	var user User
	id_int, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("User ID invalid")
		return User{}
	}

	sqlStatement := `SELECT username, creationdate, permissionlevel FROM users WHERE userid=$1;`

	row := db.QueryRow(sqlStatement, id)
	var username string
	var creationdate time.Time
	var permissionlevel int
	switch err := row.Scan(&username, &creationdate, &permissionlevel); err {
	case sql.ErrNoRows:
		fmt.Println("User not found")
	case nil:
		fmt.Printf("User %q found", username)
		user = User{id_int, username, creationdate, permissionlevel}
	default:
		panic(err)
	}

	return user
}

func GetUserPosts(id string) []ProfilePost {
	var posts []ProfilePost
	sqlStatement := `SELECT postid, posts.threadid, threads.threadname, posts.userid, users.username, posttext, postdate, posts.editdate 
					FROM posts JOIN users ON users.userid = posts.userid JOIN threads ON posts.threadid = threads.threadid
					WHERE posts.userid=$1;`

	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var postId, threadId, userId int
		var threadName, username, postText string
		var postDate time.Time
		var editDate sql.NullTime

		err = rows.Scan(&postId, &threadId, &threadName, &userId, &username, &postText, &postDate, &editDate)
		if err != nil {
			panic(err)
		}

		if editDate.Valid {
			posts = append(posts, ProfilePost{postId, threadId, threadName, userId, username, postText, postDate, &editDate.Time})
		} else {
			posts = append(posts, ProfilePost{postId, threadId, threadName, userId, username, postText, postDate, nil})
		}
	}

	return posts
}

func GetThreadDetails(id string) Thread {
	var thread Thread
	sqlStatement := `SELECT threadid, threadname, threads.userid, username, threadtext, createdate, editdate FROM threads JOIN users ON users.userid = threads.userid WHERE threadid=$1`

	row := db.QueryRow(sqlStatement, id)
	switch err := row.Scan(&thread.Threadid, &thread.ThreadName, &thread.UserId, &thread.Username, &thread.ThreadText, &thread.CreateDate, &thread.EditDate); err {
	case sql.ErrNoRows:
		fmt.Println("Thread not found")
	case nil:
		fmt.Println("Thread found")
	default:
		panic(err)
	}

	return thread
}

func GetThreadPosts(id string) []ThreadPost {
	var posts []ThreadPost
	sqlStatement := `SELECT postid, posts.userid, username, posttext, postdate, editdate FROM posts JOIN users ON posts.userid = users.userid WHERE threadid=$1 ORDER BY postdate;`

	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var post ThreadPost

		err = rows.Scan(&post.PostId, &post.UserId, &post.Username, &post.PostText, &post.PostDate, &post.EditDate)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)
	}

	return posts
}

func GetCategoriesAdmin() []FullCategory {
	var categories []FullCategory
	sqlStatement := `SELECT categories.categoryid, categories.name, categories.priority, subcategoryid, subcategories.name, permissionlevel, subcategories.priority FROM categories LEFT JOIN subcategories ON
					categories.categoryid = subcategories.categoryid;`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	m := make(map[string]FullCategory)
	var o []string

	for rows.Next() {
		var cat FullCategory
		var subCatId, permission, subCatPriority sql.NullInt32
		var subCatName sql.NullString

		err = rows.Scan(&cat.CategoryId, &cat.CategoryName, &cat.Priority, &subCatId, &subCatName, &permission, &subCatPriority)
		if err != nil {
			panic(err)
		}

		if value, ok := m[cat.CategoryName]; ok {
			value.Subcategories = append(value.Subcategories, FullSubCategory{int(subCatId.Int32), subCatName.String, int(permission.Int32), int(subCatPriority.Int32)})
			m[cat.CategoryName] = value
		} else {
			if subCatId.Valid && subCatName.Valid && permission.Valid {
				cat.Subcategories = append(cat.Subcategories, FullSubCategory{int(subCatId.Int32), subCatName.String, int(permission.Int32), int(subCatPriority.Int32)})
			}

			o = append(o, cat.CategoryName)
			m[cat.CategoryName] = cat
		}
	}

	for _, v := range o {
		categories = append(categories, m[v])
	}

	return categories
}

func GetThreadIdFromPostId(postId string) string {
	sqlStatement := `SELECT threadid FROM posts WHERE postid=$1;`
	threadId := ""

	err := db.QueryRow(sqlStatement, postId).Scan(&threadId)
	if err != nil {
		fmt.Println(err)
	}

	return threadId
}
