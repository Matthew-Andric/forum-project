package handler

import (
	"fmt"
	"forum/database"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func EditCategoryPriorityHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			catId, priority := c.PostForm("prioritycategoryid"), c.PostForm("categorypriority")

			if !database.EditCategoryPriority(catId, priority) {
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

func EditSubCategoryNameHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			subCatId, name := c.PostForm("editsubcatnameid"), c.PostForm("editsubcatname")

			if !database.EditSubCategoryName(subCatId, name) {
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

func EditSubCategoryPriorityHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			subCatId, priority := c.PostForm("subcatpriorityid"), c.PostForm("subcatpriority")

			if !database.EditSubCategoryPriority(subCatId, priority) {
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

func EditThreadHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
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

		if !database.EditThread(user, id, title, post) {
			fmt.Println("Error editing post")
			return
		}

		c.Redirect(http.StatusFound, location.RequestURI())
	}

	return gin.HandlerFunc(fn)
}

func EditCategoryNameHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			id, name := c.PostForm("categoryid"), c.PostForm("categoryname")

			if !database.EditCategoryName(id, name) {
				fmt.Println("error updating category name")
				return
			}

			location := url.URL{Path: "/admin/boards"}
			c.Redirect(http.StatusFound, location.RequestURI())
		}
	}

	return gin.HandlerFunc(fn)
}

func EditSubCategoryLocationHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			subCatId, catId := c.PostForm("subcategory"), c.PostForm("category")

			if !database.EditSubCategoryLocation(subCatId, catId) {
				fmt.Println("Error updating subcategory location")
				return
			}

			location := url.URL{Path: "/admin/boards"}
			c.Redirect(http.StatusFound, location.RequestURI())
		}
	}

	return gin.HandlerFunc(fn)
}

func EditSubCategoryPermissionHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user := database.ValidateSession(c)
		if user == nil || user.PermissionLevel != ADMIN {
			c.HTML(http.StatusBadRequest, "404", nil)
		} else {
			subCatId, permission := c.PostForm("subcategory"), c.PostForm("permission")

			if !database.EditSubCategoryPermission(subCatId, permission) {
				fmt.Println("Error updating subcategory permissions")
				return
			}

			location := url.URL{Path: "/admin/boards"}
			c.Redirect(http.StatusFound, location.RequestURI())
		}
	}

	return gin.HandlerFunc(fn)
}
