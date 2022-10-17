package main

import (
	"flag"

	"forum/database"
	handler "forum/handler"

	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	devMode := flag.Bool("dev", false, "Enables developer mode, this will disable security settings for easy use during testing")
	database.StartDB()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)
	conf := secure.DefaultConfig()
	if !*devMode {
		r.Use(secure.New(conf))
	}
	r.LoadHTMLGlob("./templates/*")
	r.Static("/static", "./static/")
	r.MaxMultipartMemory = 5 << 20

	r.GET("/", handler.GetIndexHandler())
	r.GET("/board/:id", handler.GetBoardHandler())
	r.GET("/users/:id", handler.GetUserHandler())
	r.GET("/thread/:id", handler.GetThreadHandler())
	r.GET("/login", handler.GetLoginHandler())
	r.GET("/logout", handler.GetLogoutHandler())
	r.GET("/register", handler.GetRegistrationHandler())
	r.GET("/admin", handler.GetAdminPanelHandler())
	r.GET("/admin/boards", handler.GetAdminBoardsHandler())
	r.GET("/usersettings", handler.GetUserSettings())

	r.POST("/login", handler.PostLoginHandler())
	r.POST("/register", handler.PostRegistrationHandler())
	r.POST("/thread/:id", handler.PostReplyHandler())
	r.POST("/delete/post/:id", handler.DeleteReplyHandler())
	r.POST("/edit/post/:id", handler.EditReplyHandler())
	r.POST("/board/:id", handler.PostThreadHandler())
	r.POST("/edit/thread/:id", handler.EditThreadHandler())
	r.POST("/edit/category/name", handler.EditCategoryNameHandler())
	r.POST("/edit/category/priority", handler.EditCategoryPriorityHandler())
	r.POST("/edit/subcategory/name", handler.EditSubCategoryNameHandler())
	r.POST("/edit/subcategory/category", handler.EditSubCategoryLocationHandler())
	r.POST("/edit/subcategory/permission", handler.EditSubCategoryPermissionHandler())
	r.POST("/edit/subcategory/priority", handler.EditSubCategoryPriorityHandler())
	r.POST("/delete/category", handler.DeleteCategoryHandler())
	r.POST("/add/category", handler.AddCategoryHandler())
	r.POST("/add/subcategory", handler.AddSubCategoryHandler())
	r.POST("/update/profilepicture", handler.UploadImageHandler())
	r.POST("/update/password", handler.UpdatePasswordHandler())

	r.Run()
	database.CloseDB()
}
