package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})
	r.POST("/register", register)
	r.Run("localhost:12312")
}

func register(c *gin.Context) {
	kiosk_id := c.Query("kiosk_id")
	name := c.PostForm("name")
	age := c.PostForm("age")
	gender := c.PostForm("gender")
	// email := c.PostForm("email")
	fmt.Printf("kiosk_id: %s; name: %s; age: %s; gender: %s", kiosk_id, name, age, gender)

	
    c.JSON(200, gin.H{"message":"OK"})

}
