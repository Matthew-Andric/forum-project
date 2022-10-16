package handler

import (
	"forum/database"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func GetIndexHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var params map[string]interface{}
		user := database.ValidateSession(c)
		categoryArray := database.GetCategories(user)
		result := c.Query("result")
		params = map[string]interface{}{
			"user":       user,
			"categories": categoryArray,
			"result":     result,
		}
		c.HTML(http.StatusOK, "index", params)
	}

	return gin.HandlerFunc(fn)
}

func GetBoardHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
		subid := c.Param("id")
		subcat := database.GetSubCategoryDetails(subid)
		threads := database.GetThreadsSummary(subid)
		c.HTML(http.StatusOK, "board", map[string]interface{}{"subcat": subcat, "threads": threads, "user": user})
	}

	return gin.HandlerFunc(fn)
}

func GetUserHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		userid := c.Param("id")
		User := database.GetUserDetails(userid)
		UserPosts := database.GetUserPosts(userid)
		c.HTML(http.StatusOK, "userprofile", map[string]interface{}{"user": User, "posts": UserPosts})
	}

	return gin.HandlerFunc(fn)
}

func GetThreadHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
		threadid := c.Param("id")
		Thread := database.GetThreadDetails(threadid)
		Posts := database.GetThreadPosts(threadid)
		c.HTML(http.StatusOK, "thread", map[string]interface{}{"thread": Thread, "posts": Posts, "user": user})
	}

	return gin.HandlerFunc(fn)
}

func GetLoginHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if user := database.ValidateSession(c); user != nil {
			location := url.URL{Path: "/"}
			c.Redirect(http.StatusFound, location.RequestURI())
		}
		c.HTML(http.StatusOK, "login", map[string]interface{}{"result": c.Query("result")})
	}

	return gin.HandlerFunc(fn)
}

func GetLogoutHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		cookie, err := c.Cookie("session")
		var result string
		if err != nil {
			result = "Already signed out"
			location := url.URL{Path: "/", RawQuery: result}
			c.Redirect(http.StatusFound, location.RequestURI())
		}

		if database.DeleteSession(cookie) {
			result = "result=Successfully signed out"
		} else {
			result = "result=Error signing out"
		}

		ck := &http.Cookie{
			Name:   "session",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}

		http.SetCookie(c.Writer, ck)
		location := url.URL{Path: "/", RawQuery: result}
		c.Redirect(http.StatusFound, location.RequestURI())
	}

	return gin.HandlerFunc(fn)
}

func GetRegistrationHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.HTML(http.StatusOK, "register", nil)
	}

	return gin.HandlerFunc(fn)
}

func GetAdminPanelHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusNotFound, "404", nil)
		} else {
			c.HTML(http.StatusFound, "admin", map[string]interface{}{"user": user})
		}
	}

	return gin.HandlerFunc(fn)
}

func GetAdminBoardsHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusNotFound, "404", nil)
		} else {
			categories := database.GetCategoriesAdmin()
			c.HTML(http.StatusFound, "adminboards", map[string]interface{}{"user": user, "categories": categories, "permission": []string{"User", "Moderator", "Admin"}})
		}
	}

	return gin.HandlerFunc(fn)
}

func GetUserSettings() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
		if user == nil {
			location := url.URL{Path: "/"}
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		c.HTML(http.StatusFound, "usersettings", map[string]interface{}{"user": user, "result": c.Query("result")})
	}

	return gin.HandlerFunc(fn)
}
