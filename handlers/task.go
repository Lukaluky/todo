package handlers

import (
	"context"
    "net/http"
    "time"
    "todo/db"
    "todo/models"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// === CREATE ===
func CreateTask(c *gin.Context) {
	var t models.Task

	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uint)
	t.ID = primitive.NewObjectID()
	t.UserID = userID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := db.TasksColl.InsertOne(ctx, t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, t)
}

// === LIST ===
func GetTasks(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := db.TasksColl.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)

	var tasks []models.Task
	if err := cur.All(ctx, &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// === GET BY ID ===
func GetTaskByID(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task models.Task
	err = db.TasksColl.FindOne(ctx, bson.M{"_id": oid, "user_id": userID}).Decode(&task)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// === UPDATE ===
func UpdateTask(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var body struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{}
	if body.Title != nil {
		update["title"] = *body.Title
	}
	if body.Completed != nil {
		update["completed"] = *body.Completed
	}
	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Обновляем и возвращаем обновлённый документ
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Task
	err = db.TasksColl.FindOneAndUpdate(ctx,
		bson.M{"_id": oid, "user_id": userID},
		bson.M{"$set": update},
		opts,
	).Decode(&updated)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// === DELETE ===
func DeleteTask(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := db.TasksColl.DeleteOne(ctx, bson.M{"_id": oid, "user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
