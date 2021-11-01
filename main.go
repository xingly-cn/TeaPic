package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// index
func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
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
		// index
		HomePage.GET("index.go", index)

		// copyright
		HomePage.GET("copyright.go", copyright)
	}

	r.Run()
}
