package handler

import (
	"fmt"
	"forum/database"
	"forum/util"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ADMIN     = 2
	MODERATOR = 1
	USER      = 0
)

func PostLoginHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		fmt.Println("Username:", username, "Password:", password)

		id, ok := database.ValidateLogIn(username, password)
		if !ok {
			location := url.URL{Path: "/login", RawQuery: "result=Invalid credentials"}
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		encoded, err := util.S.Encode("session", map[string]interface{}{"User": username, "Date": time.Now().String()})
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

		database.SaveSession(encoded, id)
		http.SetCookie(c.Writer, cookie)
		location := url.URL{Path: "/", RawQuery: "result=Signed in successfully"}
		c.Redirect(http.StatusFound, location.RequestURI())
	}

	return gin.HandlerFunc(fn)
}

func PostRegistrationHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		cpassword := c.PostForm("confirmpassword")

		var location url.URL

		if password != cpassword {
			location = url.URL{Path: "/register", RawQuery: "result=Passwords do not match"}
		} else if database.RegisterAccount(username, util.SaltPassword(password)) {
			location = url.URL{Path: "/login", RawQuery: "result=Registration successful, please sign in"}
		} else {
			location = url.URL{Path: "/register", RawQuery: "result=Error signing in"}
		}

		c.Redirect(http.StatusFound, location.RequestURI())
	}

	return gin.HandlerFunc(fn)
}

func PostReplyHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
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

		if !database.CanAccessThread(user, id) {
			location := url.URL{Path: "/"}
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		result := ""
		if !database.PostReply(user, reply, id) {
			result = "result=Error posting reply"
		}

		location := url.URL{Path: "/thread/" + id, RawQuery: result}
		c.Redirect(http.StatusFound, location.RequestURI())
	}

	return gin.HandlerFunc(fn)
}

func EditReplyHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		location := url.URL{Path: "/"}

		user := database.ValidateSession(c)
		if user == nil {
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		id := c.Param("id")
		threadId := database.GetThreadIdFromPostId(id)
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

		if !database.EditReply(user, id, post) {
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		fmt.Println("Successfully edited post")
		c.Redirect(http.StatusFound, location.RequestURI())
	}

	return gin.HandlerFunc(fn)
}

func PostThreadHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
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

		newid, success := database.PostThread(user, id, title, post)
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

func AddCategoryHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			name, priority := c.PostForm("addcatname"), c.PostForm("addcatpriority")

			if !database.AddCategory(name, priority) {
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

func AddSubCategoryHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			location := url.URL{Path: "/"}
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		} else {
			name, catId := c.PostForm("addsubcatname"), c.PostForm("addsubcatparent")
			permission, priority := c.PostForm("addsubcatpermission"), c.PostForm("addsubcatpriority")

			if !database.AddSubCategory(name, catId, permission, priority) {
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

func UploadImageHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
		if user == nil {
			location := url.URL{Path: "/"}
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		file, err := c.FormFile("pfp")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "No file uploaded"})
			return
		}

		filePath := "static/media/profile/" + uuid.New().String() + filepath.Ext(file.Filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error saving file"})
			return
		}

		fileType, err := util.ValidateFileType(filePath)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error validating file type"})
			os.Remove(filePath)
			return
		}

		fmt.Println("file type:", fileType)
		if fileType != "image/jpeg" && fileType != "image/png" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid file type, must be an image"})
			os.Remove(filePath)
			return
		}

		if !database.UpdateProfilePicture(user, filePath) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "database error"})
			os.Remove(filePath)
			return
		}

		location := url.URL{Path: "/"}
		c.Redirect(http.StatusFound, location.RequestURI())
	}

	return gin.HandlerFunc(fn)
}
