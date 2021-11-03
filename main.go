package main

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
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
	Username string `form:"username"`
	Password string `form:"password"`
}

// Login
func Login(c *gin.Context) {
	var loginParms LoginParms
	c.ShouldBind(&loginParms)

	//Todo This to judge authority
	c.JSON(http.StatusOK, gin.H{
		"msg": loginParms,
	})
}

// UploadConfig
func UploadConfig() (*oss.Bucket, error) {
	client, _ := oss.New("https://cdn.xingly.cn/", "LTAI4GK67v43NPGCeundD6wq", "LZdknJ9ZCB3MFM7CJopF7NK4LIds2Dg", oss.UseCname(true), oss.EnableCRC(true))
	bucket, _ := client.Bucket("xingly")
	return bucket, nil
}

// Upload
func Upload(c *gin.Context) {
	t, _ := UploadConfig()
	data, _ := c.FormFile("file")
	uuid := xid.New().String()
	UrlPath := "http://cdn.xingly.cn/" + uuid + data.Filename[len(data.Filename)-4:]
	dataHander, _ := data.Open()
	defer dataHander.Close()
	fileByte, _ := ioutil.ReadAll(dataHander)
	t.PutObject(uuid+data.Filename[len(data.Filename)-4:], bytes.NewReader(fileByte))
	c.JSON(http.StatusOK, gin.H{
		"name": data.Filename,
		"url":  UrlPath,
	})
}

// Copyright
func Copyright(c *gin.Context) {
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
		HomePage.GET("copyright.go", Copyright)
		// Interface
		HomePage.POST("upload.go", Upload)
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
