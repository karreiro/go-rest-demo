package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// User represents a user
type User struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:123@/users_development?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	db.AutoMigrate(&User{})
}

func main() {
	router := gin.Default()
	router.POST("/users", createUser)
	router.GET("/users", listUsers)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)
	router.Run(":3000")
}

// API: curl -d 'first_name=Guilherme&last_name=Carreiro&username=karreiro' -X POST http://localhost:3000/users
func createUser(c *gin.Context) {

	user := User{
		FirstName: c.PostForm("first_name"),
		LastName:  c.PostForm("last_name"),
		Username:  c.PostForm("username"),
	}

	db.Save(&user)

	c.JSON(http.StatusCreated, gin.H{
		"message": "The User was successfully created!",
		"userId":  user.ID,
	})
}

// API: curl -X GET http://localhost:3000/users
func listUsers(c *gin.Context) {

	var users []User

	db.Find(&users)

	c.JSON(http.StatusOK, users)
}

// API: curl -d 'first_name=Guilherme2&last_name=Carreiro2&username=karreiro2' -X PUT http://localhost:3000/users/2
func updateUser(c *gin.Context) {

	var user User
	var status int
	var message interface{}

	db.First(&user, c.Param("id"))

	if user.ID != 0 {
		db.Model(&user).Update("first_name", c.PostForm("first_name"))
		db.Model(&user).Update("last_name", c.PostForm("last_name"))
		db.Model(&user).Update("username", c.PostForm("username"))

		status = http.StatusOK
		message = "The User was successfully updated!"
	} else {
		status = http.StatusNotFound
		message = "The User could not be found."
	}

	c.JSON(status, gin.H{"message": message})
}

// API: curl -X DELETE http://localhost:3000/users/1
func deleteUser(c *gin.Context) {

	var user User
	var status int
	var message interface{}

	db.First(&user, c.Param("id"))

	if user.ID != 0 {
		db.Delete(&user)
		status = http.StatusOK
		message = "The User was successfully removed!"
	} else {
		status = http.StatusNotFound
		message = "The User could not be found."
	}

	c.JSON(status, gin.H{"message": message})
}
