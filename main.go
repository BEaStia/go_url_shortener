package main;

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/gin-gonic/gin"
	"crypto/md5"
)

type Record struct {
	gorm.Model
	Url string
	Shortened string
}

type Config struct {
	db *gorm.DB
	dialect string
	dburl string
}

var config Config

func (conf Config) establishConnection() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	} else {
		conf.db = db
	}
}

func crypt(url string) string {
	crypter := md5.New()
	crypter.Write([]byte(url))
	cryptedUrl := string(crypter.Sum(nil))
	config.db.Create(&Record{Url: url, Shortened: cryptedUrl})
	return cryptedUrl
}

func decrypt(encryptedUrl string) string {
	var record Record
	config.db.First(&record, "Shortened = ?", encryptedUrl)
	return record.Url
}

func main() {
	config.establishConnection()
	config.db.AutoMigrate(&Record{})
	defer config.db.Close()

	r := gin.Default()
	r.GET("/c/:param", func(c *gin.Context) {
		url := c.Param("param")
		c.JSON(200, gin.H{
			"crypted": crypt(url),
		})
	})

	r.GET("/d/:param", func(c *gin.Context) {
		encryptedUrl := c.Param("param")
		c.JSON(200, gin.H{
			"decrypted": decrypt(encryptedUrl),
		})
	})
	r.Run()

}

