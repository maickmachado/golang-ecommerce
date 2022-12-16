package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/maickmachado/golang-ecommerce/database"
	"github.com/maickmachado/golang-ecommerce/models"
	"github.com/maickmachado/golang-ecommerce/tokens"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Package validator implements value validations for structs and individual fields based on tags.
var Validate = validator.New()
var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func VerifyPassword(userPassword, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "invalid password"
		valid = false
	}
	return valid, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.UserModel
		//bindings = ligação
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.PhoneNumber})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already exists"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()

		token, refreshToken, _ := tokens.TokenGenerator(*user.Email, *user.FistName, *user.LastName, user.UserID)
		user.Token = &token
		user.RefreshToken = &refreshToken
		user.UserCart = make([]models.UserCart, 0)
		user.AddressDetails = make([]models.AdressDetails, 0)
		user.Order = make([]models.Order, 0)
		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "the user was not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, gin.H{"status": "success signup"})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.UserModel
		var foundUser models.UserModel
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err = UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}

		PasswordIsValid, msg, err := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		token, refreshToken, _ := tokens.TokenGenerator(*foundUser.Email, *foundUser.FistName, *foundUser.LastName, foundUser.UserID)
		defer cancel()

		tokens.UpdateAllTokens(token, refreshToken, foundUser.UserID)

		c.JSON(http.StatusFound, foundUser)
	}
}

func Logout() gin.HandlerFunc {}

func Register() gin.HandlerFunc {}

func ResetPassword() gin.HandlerFunc {}

func ChangePassword() gin.HandlerFunc {}

func ChangeEmail() gin.HandlerFunc {}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "something went wrong")
			return
		}
		err = cursor.All(ctx, &productList)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		err = cursor.Err()
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, "Invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(http.StatusOK, productList)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProduct []models.Product
		queryParam := c.Query("name")

		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "invalid search query"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchQuerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, "not found")
			return
		}

		err = searchQuerydb.All(ctx, &searchProduct)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, "invalid search query")
		}

		defer searchQuerydb.Close(ctx)
		err = searchQuerydb.Err()
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, "Invalid")
			return
		}

		defer cancel()
		c.IndentedJSON(http.StatusOK, searchProduct)
	}
}

func ProductViewerAdmin() gin.HandlerFunc {}
