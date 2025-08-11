package main

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("go-app"),
		newrelic.ConfigLicense("xxxxxxxxxxxxxxxxxxxxxxx"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		panic(err)
	}
	// MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	db := client.Database("testdb")
	itemsCol := db.Collection("items")

	// Insert sample data if collection is empty
	ctxCheck, cancelCheck := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCheck()
	count, err := itemsCol.CountDocuments(ctxCheck, bson.M{})
	if err == nil && count == 0 {
		_, _ = itemsCol.InsertOne(ctxCheck, bson.M{"name": "Sample Item"})
	}

	router := gin.Default()
	// Add the nrgin middleware before other middlewares or routes:
	router.Use(nrgin.Middleware(app))

	// Example endpoint: /api/data (list items from MongoDB, with New Relic segment)
	router.GET("/api/data", func(c *gin.Context) {
		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			StartTime:  newrelic.StartSegmentNow(txn),
			Product:    newrelic.DatastoreMongoDB,
			Collection: "items",
			Operation:  "find",
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cur, err := itemsCol.Find(ctx, bson.M{})
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(ctx)
		var items []bson.M
		if err := cur.All(ctx, &items); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	// Example endpoint: /api/add (insert item to MongoDB, with New Relic segment)
	router.POST("/api/add", func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			StartTime:  newrelic.StartSegmentNow(txn),
			Product:    newrelic.DatastoreMongoDB,
			Collection: "items",
			Operation:  "insert",
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := itemsCol.InsertOne(ctx, bson.M{"name": req.Name})
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"insertedID": res.InsertedID, "name": req.Name})
	})

	router.Run(":8080")

}
