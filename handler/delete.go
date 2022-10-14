package handler

import (
	"fmt"
	"forum/database"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func DeleteReplyHandler() gin.HandlerFunc {
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
		}

		location.Path = "/thread/" + threadId
		if !database.DeleteReply(user, id) {
			c.Redirect(http.StatusFound, location.RequestURI())
			return
		}

		fmt.Println("Successfully deleted post")
		c.Redirect(http.StatusFound, location.RequestURI())
	}

	return gin.HandlerFunc(fn)
}

func DeleteCategoryHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			catId := c.PostForm("deletecategoryid")

			if !database.DeleteCategory(catId) {
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
