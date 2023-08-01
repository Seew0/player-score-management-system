package server

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seew0/player-score-management-system/db"
	"github.com/seew0/player-score-management-system/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Storage *mongo.Collection
var err error

func init() {
	Storage, err = db.InitDB()
	if err != nil {
		log.Fatalf("failed to connect to db:  %v ", err)
	}
}
func createPlayer(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")

	var newPlayer model.Player
	if err := c.ShouldBindJSON(&newPlayer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid req"})
		return
	}

	if newPlayer.Name == "" || len(newPlayer.Name) > 15 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name"})
		return
	}

	if newPlayer.Country == "" || len(newPlayer.Country) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country code"})
		return
	}

	newPlayer.ID = primitive.NewObjectID()

	_, err := Storage.InsertOne(context.Background(), newPlayer)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create player"})
		return
	}

	c.JSON(http.StatusCreated, newPlayer)
}

func updatePlayer(c *gin.Context) {
	id := c.Param("id")
	log.Println(id)
	upId, erre := primitive.ObjectIDFromHex(id)
	if erre != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}
	var updatedPlayer model.Player
	if err := c.ShouldBindJSON(&updatedPlayer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if updatedPlayer.Name == "" || len(updatedPlayer.Name) > 15 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name"})
		return
	}

	filter := bson.M{"_id": upId}
	update := bson.M{"$set": bson.M{"name": updatedPlayer.Name, "score": updatedPlayer.Score}}

	_, err := Storage.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update player"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Player updated successfully", "updated": updatedPlayer})
}

func deletePlayer(c *gin.Context) {
	id := c.Param("id")

	delId, erre := primitive.ObjectIDFromHex(id)
	if erre != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	filter := bson.M{"_id": delId}
	_, err := Storage.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete player"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Player deleted successfully"})
}

func getAllPlayers(c *gin.Context) {
	options := options.Find().SetSort(bson.M{"score": -1})

	cur, err := Storage.Find(context.Background(), bson.D{}, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
		return
	}
	defer cur.Close(context.Background())

	var players []model.Player
	for cur.Next(context.Background()) {
		var player model.Player
		if err := cur.Decode(&player); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode players"})
			return
		}
		players = append(players, player)
	}

	c.JSON(http.StatusOK, players)
}

func getPlayerByRank(c *gin.Context) {
	rankParam := c.Param("val")

	rank, err := strconv.ParseInt(rankParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rank"})
		return
	}

	var players []model.Player
	options := options.Find().SetSort(bson.M{"score": -1}).SetLimit(1).SetSkip(rank)
	cur, err := Storage.Find(context.Background(), bson.D{}, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player"})
		return
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var player model.Player
		if err := cur.Decode(&player); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode player"})
			return
		}
		players = append(players, player)
	}

	if len(players) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	c.JSON(http.StatusOK, players)
}

func getRandomPlayer(c *gin.Context) {
	options := options.Aggregate()
	cur, err := Storage.Aggregate(context.Background(), bson.A{bson.M{"$sample": bson.M{"size": 1}}}, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player"})
		return
	}
	defer cur.Close(context.Background())

	var players []model.Player
	for cur.Next(context.Background()) {
		var player model.Player
		if err := cur.Decode(&player); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode player"})
			return
		}
		players = append(players, player)
	}

	c.JSON(http.StatusOK, players)
}
