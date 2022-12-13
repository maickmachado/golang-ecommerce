package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maickmachado/golang-ecommerce/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	productCollection *mongo.Collection
	userCollection    *mongo.Collection
}

func NewApplication(productCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		productCollection: productCollection,
		userCollection:    userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.productCollection, app.userCollection, userQueryID, productID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, "sucess to add to cart")
	}
}

func RemoveItemFromCart(app *Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.RemoveProductFromCart(ctx, app.productCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, "Sucess removing product from cart")
	}
}

func GetItemFromCart() gin.HandlerFunc {}

func UpdateItemInCart() gin.HandlerFunc {}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")
		if userQueryID == "" {
			log.Panic("user id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("userID id is empty"))
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(http.StatusOK, "Sucess buy from cart")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.InstantBuyFromCart(ctx, app.productCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(http.StatusOK, "sucess to buy from cart")
	}
}
