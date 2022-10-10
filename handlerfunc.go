package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ADMIN     = 2
	MODERATOR = 1
	USER      = 0
)

func getIndexHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var params map[string]interface{}
		user := validateSession(c)
		categoryArray := getCategories(user)
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

func getBoardHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		subid := c.Param("id")
		subcat := getSubCategoryDetails(subid)
		threads := getThreadsSummary(subid)
		c.HTML(http.StatusOK, "board", map[string]interface{}{"subcat": subcat, "threads": threads, "user": user})
	}

	return gin.HandlerFunc(fn)
}

func getUserHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		userid := c.Param("id")
		User := getUserDetails(userid)
		UserPosts := getUserPosts(userid)
		c.HTML(http.StatusOK, "userprofile", map[string]interface{}{"user": User, "posts": UserPosts})
	}

	return gin.HandlerFunc(fn)
}

func getThreadHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		threadid := c.Param("id")
		Thread := getThreadDetails(threadid)
		Posts := getThreadPosts(threadid)
		c.HTML(http.StatusOK, "thread", map[string]interface{}{"thread": Thread, "posts": Posts, "user": user})
	}

	return gin.HandlerFunc(fn)
}

func getLoginHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if user := validateSession(c); user != nil {
			location := url.URL{Path: "/"}
			c.Redirect(http.StatusFound, location.RequestURI())
		}
		c.HTML(http.StatusOK, "login", map[string]interface{}{"result": c.Query("result")})
	}

	return gin.HandlerFunc(fn)
}

func getLogoutHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		cookie, err := c.Cookie("session")
		var result string
		if err != nil {
			result = "Already signed out"
			location := url.URL{Path: "/", RawQuery: result}
			c.Redirect(http.StatusFound, location.RequestURI())
		}

		if deleteSession(cookie) {
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

func getRegistrationHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.HTML(http.StatusOK, "register", nil)
	}

	return gin.HandlerFunc(fn)
}

func postLoginHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		fmt.Println("Username:", username, "Password:", password)

		id, ok := validateLogIn(username, password)
		if !ok {
			location := url.URL{Path: "/login", RawQuery: "result=Invalid credentials"}
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		encoded, err := s.Encode("session", map[string]interface{}{"User": username, "Date": time.Now().String()})
		if err != nil {
			fmt.Println(err)
			location := url.URL{Path: "/", RawQuery: "result=Error signing in, please try again"}
			c.Redirect(http.StatusFound, location.RequestURI())
		}

		fmt.Println("Encoded string:", encoded)
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}

		saveSession(encoded, id)
		http.SetCookie(c.Writer, cookie)
		location := url.URL{Path: "/", RawQuery: "result=Signed in successfully"}
		c.Redirect(http.StatusFound, location.RequestURI())
	}

	return gin.HandlerFunc(fn)
}

func postRegistrationHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		cpassword := c.PostForm("confirmpassword")

		var location url.URL

		if password != cpassword {
			location = url.URL{Path: "/register", RawQuery: "result=Passwords do not match"}
		} else if registerAccount(username, saltPassword(password)) {
			location = url.URL{Path: "/login", RawQuery: "result=Registration successful, please sign in"}
		} else {
			location = url.URL{Path: "/register", RawQuery: "result=Error signing in"}
		}

		c.Redirect(http.StatusFound, location.RequestURI())
	}

	return gin.HandlerFunc(fn)
}

func postReplyHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		id := c.Param("id")

		if user == nil {
			location := url.URL{Path: "/thread/" + id, RawQuery: "result=Not signed in"}
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		reply := c.PostForm("postreply")
		if reply == "" {
			location := url.URL{Path: "/thread/" + id, RawQuery: "result=Post input cant be empty"}
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		if !canAccessThread(user, id) {
			location := url.URL{Path: "/"}
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		result := ""
		if !postReply(user, reply, id) {
			result = "result=Error posting reply"
		}

		location := url.URL{Path: "/thread/" + id, RawQuery: result}
		c.Redirect(http.StatusFound, location.RequestURI())
	}

	return gin.HandlerFunc(fn)
}

func deleteReplyHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		location := url.URL{Path: "/"}

		user := validateSession(c)
		if user == nil {
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		id := c.Param("id")
		threadId := getThreadIdFromPostId(id)
		if threadId == "" {
			c.Redirect(http.StatusFound, location.RequestURI())
		}

		location.Path = "/thread/" + threadId
		if !deleteReply(user, id) {
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		fmt.Println("Successfully deleted post")
		c.Redirect(http.StatusFound, location.RequestURI())
	}

	return gin.HandlerFunc(fn)
}

func editReplyHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		location := url.URL{Path: "/"}

		user := validateSession(c)
		if user == nil {
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		id := c.Param("id")
		threadId := getThreadIdFromPostId(id)
		if threadId == "" {
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		location.Path = "/thread/" + threadId
		post := c.PostForm("editreply")
		if post == "" {
			fmt.Println("edit must not be empty")
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		if !editReply(user, id, post) {
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		fmt.Println("Successfully edited post")
		c.Redirect(http.StatusFound, location.RequestURI())
	}

	return gin.HandlerFunc(fn)
}

func postThreadHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		if user == nil {
			fmt.Println("Invalid request: User not signed in")
			return
		}

		id := c.Param("id")
		title := c.PostForm("threadtitle")
		post := c.PostForm("threadtext")
		if title == "" || post == "" {
			location := url.URL{Path: "/board/" + id}
			c.Redirect(http.StatusFound, location.RequestURI())
		}

		newid, success := postThread(user, id, title, post)
		if !success {
			fmt.Println("Error posting thread")
			location := url.URL{Path: "/board/" + id}
			c.Redirect(http.StatusFound, location.RequestURI())
		}

		//maybe change to string instead of int ID
		fmt.Println("Thread created with ID:", newid)
		location := url.URL{Path: "/thread/" + fmt.Sprint(newid)}
		c.Redirect(http.StatusFound, location.RequestURI())
	}

	return gin.HandlerFunc(fn)
}

func editThreadHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		if user == nil {
			location := url.URL{Path: "/"}
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		id := c.Param("id")
		title := c.PostForm("threadtitle")
		post := c.PostForm("threadpost")

		location := url.URL{Path: "/thread/" + id}
		if title == "" || post == "" {
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		if !editThread(user, id, title, post) {
			fmt.Println("Error editing post")
			return
		}

		c.Redirect(http.StatusFound, location.RequestURI())
	}

	return gin.HandlerFunc(fn)
}

func getAdminPanelHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusNotFound, "404", nil)
		} else {
			c.HTML(http.StatusFound, "admin", map[string]interface{}{"user": user})
		}
	}

	return gin.HandlerFunc(fn)
}

func getAdminBoardsHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusNotFound, "404", nil)
		} else {
			categories := getCategoriesAdmin()
			c.HTML(http.StatusFound, "adminboards", map[string]interface{}{"user": user, "categories": categories, "permission": []string{"User", "Moderator", "Admin"}})
		}
	}

	return gin.HandlerFunc(fn)
}

func editCategoryNameHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			id, name := c.PostForm("categoryid"), c.PostForm("categoryname")

			if !editCategoryName(id, name) {
				fmt.Println("error updating category name")
				return
			}

			location := url.URL{Path: "/admin/boards"}
			c.Redirect(http.StatusFound, location.RequestURI())
		}
	}

	return gin.HandlerFunc(fn)
}

func editSubCategoryLocationHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			subCatId, catId := c.PostForm("subcategory"), c.PostForm("category")

			if !editSubCategoryLocation(subCatId, catId) {
				fmt.Println("Error updating subcategory location")
				return
			}

			location := url.URL{Path: "/admin/boards"}
			c.Redirect(http.StatusFound, location.RequestURI())
		}
	}

	return gin.HandlerFunc(fn)
}

func editSubCategoryPermissionHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			subCatId, permission := c.PostForm("subcategory"), c.PostForm("permission")

			if !editSubCategoryPermission(subCatId, permission) {
				fmt.Println("Error updating subcategory permissions")
				return
			}

			location := url.URL{Path: "/admin/boards"}
			c.Redirect(http.StatusFound, location.RequestURI())
		}
	}

	return gin.HandlerFunc(fn)
}

func deleteCategoryHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			catId := c.PostForm("deletecategoryid")

			if !deleteCategory(catId) {
				fmt.Println("Error deleting category")
				return
			}

			fmt.Println("Successfully deleted category")

			location := url.URL{Path: "/admin/boards"}
			c.Redirect(http.StatusFound, location.RequestURI())
		}
	}

	return gin.HandlerFunc(fn)
}

func addCategoryHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			name, priority := c.PostForm("addcatname"), c.PostForm("addcatpriority")

			if !addCategory(name, priority) {
				fmt.Println("Error adding category")
				return
			}

			fmt.Println("Category successfully added")

			location := url.URL{Path: "/admin/boards"}
			c.Redirect(http.StatusFound, location.RequestURI())
		}
	}

	return gin.HandlerFunc(fn)
}

func addSubCategoryHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			name, catId := c.PostForm("addsubcatname"), c.PostForm("addsubcatparent")
			permission, priority := c.PostForm("addsubcatpermission"), c.PostForm("addsubcatpriority")

			if !addSubCategory(name, catId, permission, priority) {
				fmt.Println("Error adding category")
				return
			}

			fmt.Println("Subcategory successfully added")
			location := url.URL{Path: "/admin/boards"}
			c.Redirect(http.StatusFound, location.RequestURI())
		}
	}

	return gin.HandlerFunc(fn)
}

func editCategoryPriorityHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			catId, priority := c.PostForm("prioritycategoryid"), c.PostForm("categorypriority")

			if !editCategoryPriority(catId, priority) {
				fmt.Println("Error editing category priority")
				return
			}

			fmt.Println("Category priority changed")
			location := url.URL{Path: "/admin/boards"}
			c.Redirect(http.StatusFound, location.RequestURI())
		}
	}

	return gin.HandlerFunc(fn)
}

func editSubCategoryNameHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			subCatId, name := c.PostForm("editsubcatnameid"), c.PostForm("editsubcatname")

			if !editSubCategoryName(subCatId, name) {
				fmt.Println("Error editing subcategory name")
				return
			}

			fmt.Println("Successfully edited subcategory name")
			location := url.URL{Path: "/admin/boards"}
			c.Redirect(http.StatusFound, location.RequestURI())
		}
	}

	return gin.HandlerFunc(fn)
}

func editSubCategoryPriorityHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := validateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			subCatId, priority := c.PostForm("subcatpriorityid"), c.PostForm("subcatpriority")

			if !editSubCategoryPriority(subCatId, priority) {
				fmt.Println("Error editing subcategory priority")
				return
			}

			fmt.Println("Successfully edited subcategory priority")
			location := url.URL{Path: "/admin/boards"}
			c.Redirect(http.StatusFound, location.RequestURI())
		}
	}

	return gin.HandlerFunc(fn)
}
