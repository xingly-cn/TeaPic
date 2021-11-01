package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Homeindex
func HomeIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func AdminIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}

// copyright
func copyright(c *gin.Context) {
	c.HTML(http.StatusOK, "copyright.html", nil)
}

func main() {
	/**
	*	Hello TeaPic
	*	A Simple Pic Bed By XiaoNianXin
	*	Created By 2021-10-31 16:53
	**/

	r := gin.Default()

	// Resource Loading
	r.LoadHTMLGlob("page/*")

	// 404 Handler
	r.NoRoute(func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"msg": "TeaPic 404."}) })

	// HomePage Group
	HomePage := r.Group("/user/")
	{
		HomePage.GET("index.go", HomeIndex)
		HomePage.GET("copyright.go", copyright)
	}

	// AdminPage Group
	AdminPage := r.Group("/admin/")
	{
		AdminPage.GET("index.go", AdminIndex)
	}

	r.Run()
}
