package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/xid"
)

//---------------------------------------视图映射---------------------------------------------

// Homeindex
func HomeIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// AdminIndex
func AdminIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}

// UserAdmin
func UserAdmin(c *gin.Context) {
	c.HTML(http.StatusOK, "user.html", nil)
}

// UserBoard
func UserBoard(c *gin.Context) {
	c.HTML(http.StatusOK, "userBoard.html", nil)
}

// AdminBoard
func AdminBoard(c *gin.Context) {
	c.HTML(http.StatusOK, "adminBoard.html", nil)
}

//---------------------------------------配置类---------------------------------------------

// UploadConfig
func UploadConfig() (*oss.Bucket, error) {
	client, _ := oss.New("https://cdn.xingly.cn/", "LTAI4GK67v43NPGCeundD6wq", "ZdknJ9ZCB3MFM7CJopF7NK4LIds2Dg", oss.UseCname(true), oss.EnableCRC(true))
	bucket, _ := client.Bucket("xingly")
	return bucket, nil
}

// MysqlConfig
func MysqlConfig() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:XNXxnx520@(8.142.199.134:3306)/teapic?charset=utf8mb4")
	return db, err
}

//---------------------------------------通用接口---------------------------------------------

// listPic
func listPic(c *gin.Context) {
	// 判断是否有权限
	token := c.Query("token")
	user, _, _ := JwtParseUser(token)
	Auth := ""
	if user == nil {
		c.JSON(http.StatusOK, nil)
		return
	} else {
		if user.Status == "0" {
			Auth = user.Username
		} else if user.Status == "2" {
			c.JSON(http.StatusOK, nil)
			return
		}
	}
	// 获取图片信息并封装
	t, _ := UploadConfig()
	files, _ := t.ListObjects(oss.Prefix(Auth), oss.MaxKeys(1000))
	count := len(files.Objects)
	picData := []map[string]string{}
	for i := 0; i < count; i++ {
		temp := map[string]string{
			"uid":  strconv.Itoa(i + 1),
			"name": files.Objects[i].Key,
			"size": strconv.FormatInt(files.Objects[i].Size, 10),
			"time": files.Objects[i].LastModified.Format("2006-01-02 15:04:05"),
		}
		picData = append(picData, temp)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"count": count,
		"data":  picData,
		"user":  user,
	})
}

// delPic
func delPic(c *gin.Context) {
	file := c.Query("name")
	t, _ := UploadConfig()
	err := t.DeleteObject(file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err,
		})
		return 
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除[" + file + "]成功",
	})
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

// Jwt-rz
func Jwt_rz(c *gin.Context) {
	token := c.Query("token")
	user, _, _ := JwtParseUser(token)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"Data": user,
	})
}

// LoginParms
type LoginParms struct {
	Username string
	Password string
}

// Login
func ALogin(c *gin.Context) {
	var loginParms LoginParms
	loginParms.Username = c.Query("username")

	user := User{
		Uuid:     xid.New().String(),
		Username: loginParms.Username,
		Password: "Hello Bitch",
		Phone:    "",  // 查数据库获得
		Status:   "1", // 查数据库获得
		Picday:   100,
	}

	token, _ := jwtGenerateToken(&user, time.Hour*24*365)

	//Todo This to judge authority
	c.JSON(http.StatusOK, gin.H{
		"t":     loginParms,
		"code":  200,
		"token": token,
	})
}

//---------------------------------------版权信息---------------------------------------------

// Copyright
func Copyright(c *gin.Context) {
	c.HTML(http.StatusOK, "copyright.html", nil)
}

//---------------------------------------JWT 身份验证---------------------------------------------

// JWT Secret
var secret = "4680bbfa2dc33ab4f5d657658156a075"

// User Entity
type User struct {
	Uuid     string
	Username string
	Password string
	Phone    string
	Status   string
	Picday   int
}

// My Payload
type userClaims struct {
	jwt.StandardClaims
	*User
}

// Generate Jwt
func jwtGenerateToken(m *User, d time.Duration) (string, error) {
	m.Password = ""
	expireTime := time.Now().Add(d)
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expireTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        m.Uuid,
		Issuer:    "WWW.XINGLY.CN",
	}
	uClaims := userClaims{
		StandardClaims: stdClaims,
		User:           m,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, uClaims).SignedString([]byte(secret))
	return token, err
}

// Parse Jwt
func JwtParseUser(token string) (*User, int64, error) {
	if token == "" {
		return nil, 0, errors.New("Token is not null")
	}
	claims := userClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, 0, err
	}

	return claims.User, claims.ExpiresAt, err
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
		HomePage.GET("admin.go", UserAdmin)
		HomePage.GET("userBoard.go", UserBoard)
		HomePage.GET("adminBoard.go", AdminBoard)
		// Interface
		HomePage.POST("upload.go", Upload)
	}

	// AdminPage Group
	AdminPage := r.Group("/admin/")
	{
		// Page
		AdminPage.GET("index.go", AdminIndex)
		// Interface
		AdminPage.GET("login.go", ALogin)
		// Interface
		AdminPage.GET("jwt-rz", Jwt_rz)
		// Interface
		AdminPage.GET("listPic", listPic)
		// Interface
		AdminPage.GET("delPic", delPic)
	}
	r.Run(":8080")
}
