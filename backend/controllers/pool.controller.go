package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"live-pooling-backend/websocket"
	"live-pooling-backend/models"
	"live-pooling-backend/services"
)

func CreatePoll(c *gin.Context) {
	var poll models.Poll

	if err := c.ShouldBindJSON(&poll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if poll.Question == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Question is required",
		})
		return
	}

	if len(poll.Options) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least 2 options are required",
		})
		return
	}

	for _, option := range poll.Options {
		if option.ID == "" || option.Text == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Each option must have an ID and text",
			})
			return
		}
	}

	poll.CreatedAt = time.Now()

	collection := services.MongoClient.Database("live_polling").Collection("polls")

	result, err := collection.InsertOne(c, poll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create poll",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Poll created successfully",
		"id":      result.InsertedID,
	})
}

func GetPolls(c *gin.Context) {
	collection := services.MongoClient.Database("live_polling").Collection("polls")

	cursor, err := collection.Find(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch polls",
		})
		return
	}
	defer cursor.Close(c)

	polls := []models.Poll{}

	if err := cursor.All(c, &polls); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to decode polls",
		})
		return
	}

	c.JSON(http.StatusOK, polls)
}

func VotePoll(c *gin.Context) {
    pollID := c.Param("id")
    optionID := c.Param("optionId")

    objectID, err := primitive.ObjectIDFromHex(pollID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid poll ID",
        })
        return
    }

    collection := services.MongoClient.Database("live_polling").Collection("polls")

    filter := bson.M{
        "_id": objectID,
        "options.id": optionID,
    }

    update := bson.M{
        "$inc": bson.M{
            "options.$.votes": 1,
        },
    }

    result, err := collection.UpdateOne(c, filter, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to vote",
        })
        return
    }

    if result.MatchedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Poll or option not found",
        })
        return
    }
	ctx := c.Request.Context()

redisKey := "poll:" + pollID + ":option:" + optionID

err = services.RedisClient.Incr(ctx, redisKey).Err()
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "error": "Failed to update Redis vote count",
    })
    return
}
if websocket.CurrentHub != nil {
	websocket.CurrentHub.Broadcast(gin.H{
		"type":     "vote_update",
		"poll_id":  pollID,
		"option_id": optionID,
	})
}
    c.JSON(http.StatusOK, gin.H{
        "message": "Vote recorded successfully",
    })
}