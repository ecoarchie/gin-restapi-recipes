// Recipes API
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
// Contact: Artur Yusupov <eco.eco@mail.com> https://artur.com
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// swagger:parameters recipes newRecipe
type Recipe struct {
	//swagger:ignore
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
}

var recipes []Recipe
var err error
var ctx context.Context
var client *mongo.Client
var collection *mongo.Collection

func init() {
	// recipes = make([]Recipe, 0)
	// file, _ := os.ReadFile("recipes.json")
	// _ = json.Unmarshal([]byte(file), &recipes)

	ctx = context.Background()
	client, err = mongo.Connect(
		ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")),
	)
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	// var listOfRecipes []interface{}
	// for _, recipe := range recipes {
	// 	listOfRecipes = append(listOfRecipes, recipe)
	// }
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	// insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("Inserted recipes: ", len(insertManyResult.InsertedIDs))
}

// swagger:operation POST /recipes recipes createRecipe
// Create new recipe
//
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	  description: Recipe created
//	'400':
//	  description: Invalid input
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err = collection.InsertOne(ctx, recipe)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while inserting a new recipe",
		})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//		description: Successful operation
func ListRecipesHanlder(c *gin.Context) {
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println("here")
		return
	}
	defer cursor.Close(ctx)

	recipes := make([]Recipe, 0)
	for cursor.Next(ctx) {
		var recipe Recipe
		cursor.Decode(&recipe)
		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, recipes)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update an existing recipe
//
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	  description: Successful operation
//	'400':
//	  description: Invalid input
//	'404':
//	  description: Invalid recipe ID
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	objectID, _ := primitive.ObjectIDFromHex(id)
	_, err = collection.UpdateOne(ctx, bson.M{
		"_id": objectID,
	}, bson.D{
		{"$set", bson.D{
			{"name", recipe.Name},
			{"instructions", recipe.Instructions},
			{"ingredients", recipe.Ingredients},
			{"tags", recipe.Tags},
		}},
	})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// Delete recipe by ID
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	  description: recipe has been deleted
//	'404':
//	  description: recipe not found
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, _ := primitive.ObjectIDFromHex(id)
	_, err := collection.DeleteOne(ctx, bson.M{
		"_id": objectID,
	})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted",
	})
}

// swagger:operation GET /recipes/search recipes searchRecipes
// Search recipes by tag
// ---
// produces:
// - application/json
// parameters:
//   - name: tag
//     in: query
//     description: tag to search for in recipes
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	  description: recipes found
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)

	cursor, err := collection.Find(ctx, bson.M{
		"tags": tag,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println("here")
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var recipe Recipe
		cursor.Decode(&recipe)
		listOfRecipes = append(listOfRecipes, recipe)
	}

	c.JSON(http.StatusOK, listOfRecipes)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHanlder)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("recipes/search", SearchRecipesHandler)
	router.Run()
}
