// package controller

// import (
// 	"errors"
// 	"fmt"
// 	"go/token"
// 	"regexp"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/VillagerV/blogbackend/database"
// 	"github.com/VillagerV/blogbackend/models"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/utils"
// )

// func validateEmail(email string) bool {
// 	Re := regexp.MustCompile(`[a-z0-9._%+\-]+@[a-z0-9._%+\-]+\.[a-z0-9._%+\-]`)
// 	return Re.MatchString(email)
// }

// func Register(c *fiber.Ctx) error {
// 	var data map[string]interface{}
// 	if err := c.BodyParser(&data); err != nil {
// 		fmt.Println("Unable to parse body")
// 	}

// 	// Check if password is less than 6 characters
// 	if len(data["password"].(string)) <= 6 {
// 		c.Status(400)
// 		return c.JSON(fiber.Map{
// 			"message": "Password must be greater than 6 characters",
// 		})
// 	}

// 	// Check if email is valid
// 	if !validateEmail(strings.TrimSpace(data["email"].(string))) {
// 		c.Status(400)
// 		return c.JSON(fiber.Map{
// 			"message": "Invalid Email Address",
// 		})
// 	}

// 	// First, try to find a user with the given email
// 	var user models.User
// 	database.DB.Where(models.User{Email: strings.TrimSpace(data["email"].(string))}).Assign(models.User{
// 		FirstName: data["first_name"].(string),
// 		LastName:  data["last_name"].(string),
// 		Phone:     data["phone"].(string),
// 		Email:     strings.TrimSpace(data["email"].(string)),
// 	}).FirstOrCreate(&user)

// 	if database.DB.Error != nil {
// 		c.Status(500)
// 		return c.JSON(fiber.Map{
// 			"message": "Unable to create account",
// 		})
// 	}

// 	//Original code from the tutorial
// 	// If the user doesn't exist, create a new one
// 	// if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 	// 	user = models.User{
// 	// 		FirstName: data["first name"].(string),
// 	// 		LastName:  data["last name"].(string),
// 	// 		Phone:     data["phone"].(string),
// 	// 		Email:     strings.TrimSpace(data["email"].(string)),
// 	// 	}
// 	// 	user.SetPassword(data["password"].(string))
// 	// 	result = database.DB.Create(&user)
// 	// 	if result.Error != nil {
// 	// 		c.Status(500)
// 	// 		return c.JSON(fiber.Map{
// 	// 			"message": "Error creating user",
// 	// 		})
// 	// 	}
// 	// }

// 	c.Status(200)
// 	return c.JSON(fiber.Map{
// 		"message": "Account created successfully",
// 	})

// 	func Login(c *fiber.Ctx)error  {
// 		var data map[string]string

// 		if err := c.BodyParser(&data); err != nil {
// 			fmt.Println("Unable to parse body")
// 		}

// 		var user models.User
// 		database.DB.Where("email=?", data["email"]).First(&user)
// 		if user.Id == 0{
// 			c.Status(404)
// 			return c.JSON(fiber.Map{
// 				"message":"Email Address does not exist, kindly create an account",
// 			})
// 		}
// 		if err:=user.ComparePassword(data["password"]); err != nil{
// 			c.Status(400)
// 			return c.JSON(fiber.Map{
// 				"message":"incorrect password",
// 			})
// 		}
// 		token,err:=util.GenerateJwt(strconv.Itoa(int(user.Id)),)
// 		if err !=nil{
// 			c.Status(fiber.StatusInternalServerError)
// 			return nil
// 		}

//			cookie := fiber.Cookie{
//				Name:"jwt",
//				Value:token,
//				Expires: time.Now().Add(time.Hour*24),
//				HTTPOnly: true,
//			}
//			c.Cookie(&cookie)
//			return c.JSON(fiber.Map{
//				"message":"You have successfully logged in",
//				"user":user,
//			})
//	}
package controller

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"

	"github.com/VillagerV/blogbackend/database"
	"github.com/VillagerV/blogbackend/models"
	"github.com/VillagerV/blogbackend/util"
)

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`[a-z0-9._%+\-]+@[a-z0-9._%+\-]+\.[a-z0-9._%+\-]`)
	return Re.MatchString(email)
}

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
	}
	//Check if password is less than 6 characters
	if len(data["password"].(string)) <= 6 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Password must be greater than 6 character",
		})
	}

	if !validateEmail(strings.TrimSpace(data["email"].(string))) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid Email Address",
		})

	}
	//Check if email already exist in database
	database.DB.Where("email=?", strings.TrimSpace(data["email"].(string))).First(&userData)
	if userData.Id != 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Email already exist",
		})

	}
	user := models.User{
		FirstName: data["first_name"].(string),
		LastName:  data["last_name"].(string),
		Phone:     data["phone"].(string),
		Email:     strings.TrimSpace(data["email"].(string)),
	}
	user.SetPassword(data["password"].(string))
	err := database.DB.Create(&user)
	if err != nil {
		log.Println(err)
	}
	c.Status(200)
	return c.JSON(fiber.Map{
		"user":    user,
		"message": "Account created successfullys",
	})

}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
	}
	var user models.User
	database.DB.Where("email=?", data["email"]).First(&user)
	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "Email Address doesn't exit, kindly create an account",
		})
	}
	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}
	token, err := util.GenerateJwt(strconv.Itoa(int(user.Id)))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "you have successfully login",
		"user":    user,
	})

}

type Claims struct {
	jwt.StandardClaims
}
