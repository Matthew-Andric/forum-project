package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {
	db = startDB()
	defer db.Close()

	r := gin.Default()
	r.LoadHTMLGlob("./templates/*")
	r.Static("/static", "./static/")

	r.GET("/", getIndexHandler())
	r.GET("/board/:id", getBoardHandler())
	r.GET("/users/:id", getUserHandler())
	r.GET("/thread/:id", getThreadHandler())
	r.GET("/login", getLoginHandler())
	r.GET("/logout", getLogoutHandler())
	r.GET("/register", getRegistrationHandler())
	r.GET("/admin", getAdminPanelHandler())
	r.GET("/admin/boards", getAdminBoardsHandler())
	//r.GET("/admin/users", getAdminUsersPanel())

	r.POST("/login", postLoginHandler())
	r.POST("/register", postRegistrationHandler())
	r.POST("/thread/:id", postReplyHandler())
	r.POST("/delete/post/:id", deleteReplyHandler())
	r.POST("/edit/post/:id", editReplyHandler())
	r.POST("/board/:id", postThreadHandler())
	r.POST("/edit/thread/:id", editThreadHandler())
	r.POST("/edit/category/name", editCategoryNameHandler())
	r.POST("/edit/category/priority", editCategoryPriorityHandler())
	r.POST("/edit/subcategory/name", editSubCategoryNameHandler())
	r.POST("/edit/subcategory/category", editSubCategoryLocationHandler())
	r.POST("/edit/subcategory/permission", editSubCategoryPermissionHandler())
	r.POST("/edit/subcategory/priority", editSubCategoryPriorityHandler())
	r.POST("/delete/category", deleteCategoryHandler())
	r.POST("/add/category", addCategoryHandler())
	r.POST("/add/subcategory", addSubCategoryHandler())

	r.Run()
}

func saltPassword(pw string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}

	return string(hash)
}
