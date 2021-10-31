package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// index
func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func main() {
	/**
	*	Hello TeaPic
	*	A Simple Pic Bed By XiaoNianXin
	*	Created By 2021-10-31 16:53
	**/

	r := gin.Default()

	// Resource Loading
	r.LoadHTMLFiles("./index.html")

	// 404 Handler
	r.NoRoute(func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"msg": "TeaPic 404."}) })

	// HomePage Group
	HomePage := r.Group("/")
	{
		HomePage.GET("", index)
		HomePage.GET("index", index)
		HomePage.GET("index.html", index)
		HomePage.GET("index.php", index)
		HomePage.GET("index.go", index)
	}

	

	r.Run()
}
