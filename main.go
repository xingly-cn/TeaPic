package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Homeindex
func HomeIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// AdminIndex
func AdminIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}

// LoginParms
type LoginParms struct {
	Username string
	Password string
}

// Login
func Login(c *gin.Context) {
	var loginParms LoginParms
	loginParms.Username = c.Query("username")
	loginParms.Password = c.Query("password")

	//Todo This to judge authority
	c.JSON(http.StatusOK, gin.H{
		"msg": loginParms,
	})
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
		// Pages
		HomePage.GET("index.go", HomeIndex)
		HomePage.GET("copyright.go", copyright)
	}

	// AdminPage Group
	AdminPage := r.Group("/admin/")
	{
		// Page
		AdminPage.GET("index.go", AdminIndex)
		// Interface
		AdminPage.GET("login.go", Login)
	}

	r.Run()
}
