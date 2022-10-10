package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"
)

var (
	config YAMLFile
	s      securecookie.SecureCookie
)

type YAMLFile struct {
	Database Database `yaml:"database"`
	Secret   Secret   `yaml:"secret"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type Secret struct {
	Hash  string `yaml:"hash"`
	Block string `yaml:"block"`
}

type Category struct {
	CategoryId    int
	CategoryName  string
	Subcategories []Subcategory
}

type Subcategory struct {
	SubCategoryId   int
	SubCategoryName string
	PermissionLevel int
}

type FullCategory struct {
	CategoryId    int
	CategoryName  string
	Subcategories []FullSubCategory
	Priority      int
}

type FullSubCategory struct {
	SubCategoryId   int
	SubCategoryName string
	PermissionLevel int
	Priority        int
}

type ThreadSummary struct {
	Threadid    int
	ThreadName  string
	Userid      int
	Username    string
	CreateDate  time.Time
	LastUpdated time.Time
}

type User struct {
	Userid          int
	Username        string
	CreationDate    time.Time
	PermissionLevel int
}

type ProfilePost struct {
	Postid       int
	Threadid     int
	ThreadName   string
	Userid       int
	Username     string
	PostText     string
	PostDate     time.Time
	PostEditDate *time.Time
}

type ThreadPost struct {
	PostId   int
	UserId   int
	Username string
	PostText string
	PostDate time.Time
	EditDate *time.Time
}

type Thread struct {
	Threadid   int
	ThreadName string
	UserId     int
	Username   string
	ThreadText string
	CreateDate time.Time
	EditDate   *time.Time
}

func startDB() *sql.DB {
	err := parseConfig()
	if err != nil {
		panic(err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("DB Opened")

	s = *securecookie.New([]byte(config.Secret.Hash), []byte(config.Secret.Block))

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("DB pinged")

	return db
}

func getCategories(user *User) []Category {
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

func getSubCategoryDetails(id string) Subcategory {
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

func getThreadsSummary(id string) []ThreadSummary {
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

func getUserDetails(id string) User {
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

func getUserPosts(id string) []ProfilePost {
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

func getThreadDetails(id string) Thread {
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

func getThreadPosts(id string) []ThreadPost {
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

func validateLogIn(username string, password string) (int, bool) {
	sqlStatement := `SELECT userid, password FROM users WHERE username = $1;`
	var pw string
	var id int

	row := db.QueryRow(sqlStatement, username)
	switch err := row.Scan(&id, &pw); err {
	case sql.ErrNoRows:
		fmt.Println("Invalid credentials")
	case nil:
		if bcrypt.CompareHashAndPassword([]byte(pw), []byte(password)) == nil {
			return id, true
		} else {
			fmt.Println("Invalid credentials")
		}
	default:
		panic(err)
	}

	return 0, false
}

func saveSession(s string, id int) bool {
	sqlStatement := `INSERT INTO sessions(session, userid, startdate) VALUES ($1, $2, NOW());`

	_, err := db.Exec(sqlStatement, s, id)
	if err != nil {
		panic(err)
	}

	return true
}

func validateSession(c *gin.Context) *User {
	var user User

	cookie, err := c.Cookie("session")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	sqlStatement := `SELECT users.userid, username, creationdate, permissionlevel FROM users JOIN sessions ON users.userid = sessions.userid WHERE session=$1;`

	row := db.QueryRow(sqlStatement, cookie)
	switch err := row.Scan(&user.Userid, &user.Username, &user.CreationDate, &user.PermissionLevel); err {
	case sql.ErrNoRows:
		fmt.Println("session not found")
		return nil
	case nil:
		fmt.Println("session found")
		return &user
	default:
		panic(err)
	}
}

func deleteSession(s string) bool {
	sqlStatement := `DELETE FROM sessions WHERE session=$1;`

	_, err := db.Exec(sqlStatement, s)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func registerAccount(username string, password string) bool {
	sqlStatement := `INSERT INTO users(username, password, creationdate, permissionlevel) VALUES ($1, $2, NOW(), 0);`

	_, err := db.Exec(sqlStatement, username, password)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func postReply(user *User, post string, threadid string) bool {
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

func deleteReply(user *User, postid string) bool {
	sqlStatement := `DELETE FROM posts WHERE postid=$1 AND userid=$2;`

	_, err := db.Exec(sqlStatement, postid, user.Userid)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func editReply(user *User, postid string, post string) bool {
	sqlStatement := `UPDATE posts SET posttext=$1, editdate=NOW() WHERE userid=$2 AND postid=$3;`

	_, err := db.Exec(sqlStatement, post, user.Userid, postid)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func postThread(user *User, subcatid string, title string, post string) (int, bool) {
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

func editThread(user *User, threadid string, name string, post string) bool {
	sqlStatement := `UPDATE threads SET threadname = $1, threadtext = $2 WHERE userid=$3 AND threadid=$4;`

	_, err := db.Exec(sqlStatement, name, post, user.Userid, threadid)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Thread edited id:", threadid)

	return true
}

func canAccessThread(user *User, threadid string) bool {
	sqlStatement := `SELECT permissionlevel FROM subcategories JOIN threads ON subcategories.subcategoryid = threads.subcategoryid
					WHERE threadid=$1;`

	var p int
	err := db.QueryRow(sqlStatement, threadid).Scan(&p)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if user.PermissionLevel < p {
		return false
	}

	return true
}

func editCategoryName(id, name string) bool {
	sqlStatement := `UPDATE categories SET name=$1 WHERE categoryid=$2;`

	_, err := db.Exec(sqlStatement, name, id)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Category name updated:", id, name)

	return true
}

func editSubCategoryLocation(subCatId, catId string) bool {
	sqlStatement := `UPDATE subcategories SET categoryid=$1 WHERE subcategoryid=$2;`

	_, err := db.Exec(sqlStatement, catId, subCatId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Sub Category location updated:", subCatId, catId)

	return true
}

func editSubCategoryPermission(subCatId, permission string) bool {
	sqlStatement := `UPDATE subcategories SET permissionlevel=$1 WHERE subcategoryid=$2;`

	_, err := db.Exec(sqlStatement, permission, subCatId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Sub Category permission updated ID:", subCatId, "permission:", permission)

	return true
}

func deleteCategory(catId string) bool {
	sqlStatement := `DELETE FROM categories WHERE categoryid=$1;`

	_, err := db.Exec(sqlStatement, catId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func getCategoriesAdmin() []FullCategory {
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

func addCategory(name, priority string) bool {
	sqlStatement := `INSERT INTO categories(name, priority) VALUES ($1, $2);`

	_, err := db.Exec(sqlStatement, name, priority)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func addSubCategory(name, catId, permission, priority string) bool {
	sqlStatement := `INSERT INTO subcategories(name, categoryid, permissionlevel, priority) VALUES ($1, $2, $3, $4);`

	_, err := db.Exec(sqlStatement, name, catId, permission, priority)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func editCategoryPriority(catId, priority string) bool {
	sqlStatement := `UPDATE categories SET priority=$1 WHERE categoryid=$2;`

	_, err := db.Exec(sqlStatement, priority, catId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func editSubCategoryName(subCatId, name string) bool {
	sqlStatement := `UPDATE subcategories SET name=$1 WHERE subcategoryid=$2;`

	_, err := db.Exec(sqlStatement, name, subCatId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func editSubCategoryPriority(subCatId, priority string) bool {
	sqlStatement := `UPDATE subcategories SET priority=$1 WHERE subcategoryid=$2;`

	_, err := db.Exec(sqlStatement, priority, subCatId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func getThreadIdFromPostId(postId string) string {
	sqlStatement := `SELECT threadid FROM posts WHERE postid=$1;`
	threadId := ""

	err := db.QueryRow(sqlStatement, postId).Scan(&threadId)
	if err != nil {
		fmt.Println(err)
	}

	return threadId
}

func parseConfig() error {
	cfg, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(cfg, &config)
	return err
}
