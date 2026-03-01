package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/database"
	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/internal/ai"
	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/models"
)

var movieCollection *mongo.Collection
var validate = validator.New()

func InitMovieController(aiService *ai.Service) *MovieController{
	movieCollection = database.OpenCollection("movies")
	return &MovieController{
		AIService: aiService,
	}
}

type MovieController struct {
	AIService *ai.Service
}

func GetMovies() gin.HandlerFunc{
	return func(c *gin.Context){
		// testing purpose
		// c.JSON(200, gin.H{"message":"List of Movies!"})

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var movies []models.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data."})
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies."})
			return 
		}

		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return 
		}

		var movie models.Movie

		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return 
		}

		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return 
		}

		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return 
		}

		count, _ := movieCollection.CountDocuments(ctx, bson.M{
			"imdb_id": movie.ImdbID,
		})

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Movie already exists"})
			return 
		}

		movie.ID = primitive.NewObjectID()

		result, err := movieCollection.InsertOne(ctx, movie)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return 
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Successfully added a movie",
			"id": result.InsertedID,
		})
	}

}

func (mc *MovieController) AdminReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
		defer cancel()

		imdb_id := c.Param("imdb_id")
		if imdb_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID."})
			return
		}

		var req struct {
			AdminReview string `json:"admin_review" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "admin_review is required."})
			return
		}

		var movie models.Movie
		if err := movieCollection.FindOne(ctx, bson.M{"imdb_id": imdb_id}).Decode(&movie); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		fmt.Println("1 - After movie fetch")
		rank, rawAIResponse, aiErr := mc.AIService.GenerateAdminRanking(ctx, req.AdminReview)
		fmt.Println("2 - After AI call")
		fmt.Println("After GenerateAdminRanking")
		fmt.Println("Rank:", rank)
		fmt.Println("Raw:", rawAIResponse)
		fmt.Println("Err:", aiErr)
		fmt.Println("3 - Before log insert")
		log := models.AIReviewLog{
			MovieID:       imdb_id,
			AdminReview:   req.AdminReview,
			AIRawResponse: rawAIResponse,
			ParsedRank:    rank,
			CreatedAt:     time.Now(),
		}
		if aiErr != nil {
			log.Status = "failed"
			log.ErrorMessage = aiErr.Error()
		}else {
			log.Status = "success"
		}
		fmt.Println("4 - After log insert")
		// Insert AI log
		aiLogsColl := database.OpenCollection("ai_review_logs")
		result, err := aiLogsColl.InsertOne(ctx, log)
		if err != nil {
			fmt.Println("Insert error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert AI log"})
			return
		}
		fmt.Println("Inserted ID:", result.InsertedID)

		if aiErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "AI review failed, please try again", "details": aiErr.Error()})
			return
		}
		fmt.Println("5 - Before ranking lookup")
		// fmt.Printf("AI Parsed Rank: %v\n", rank)
		// fmt.Printf("AI Parsed Rank Type: %T\n", rank)
		// Lookup ranking info
		rankingColl := database.OpenCollection("rankings")
		var rankDoc models.Ranking
		if err := rankingColl.FindOne(ctx, bson.M{"ranking_value": int32(rank)}).Decode(&rankDoc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ranking not found", "details": err.Error()})
			return
		}

		// Update movie
		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankDoc.RankingValue,
					"ranking_name":  rankDoc.RankingName,
				},
			},
		}

		if _, err := movieCollection.UpdateOne(
			ctx, 
			bson.M{"imdb_id": imdb_id}, 
			update,
			); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
			return
		}

		// Success response
		c.JSON(http.StatusOK, gin.H{
			"message": "Admin review added successfully",
			"rank":    rankDoc.RankingName,
		})
	}
}