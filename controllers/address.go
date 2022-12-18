package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maickmachado/golang-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAdress() gin.HandlerFunc {}

func EditAdress() gin.HandlerFunc {}

func EditWorkAdress() gin.HandlerFunc {}

func DeleteAdress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")
		if userQueryID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index"})
			c.Abort()
			return
		}

		addresses := make([]models.AdressDetails, 0)
		userIDHex, err := primitive.ObjectIDFromHex(userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal error")
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "id", Value: userIDHex}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, "wrong command")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "success deleted")
	}
}
