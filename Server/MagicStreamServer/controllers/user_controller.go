package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/database"
	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/models"
	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

func InitUserController(){
	userCollection = database.OpenCollection("users")
}

func HashPassword(password string)(string, error){
	hashPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password), 
		bcrypt.DefaultCost,
	)
	
	if err != nil {
		return "", err
	}

	return string(hashPassword), nil
}

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}


		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return 
		}

		//Check if user already exists
		count, err := userCollection.CountDocuments(ctx, bson.M{
			"email": user.Email,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return 
		}

		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		}

		user.ID = primitive.NewObjectID()
		// user.UserID = user.ID.Hex()
		user.Password = hashedPassword
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		_, err = userCollection.InsertOne(ctx, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failedto create user"})
			return 
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User registered successfully",
			"user_id": user.ID,
		})
		
	}
}

func LoginUser() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var userLogin models.UserLogin
		if err := c.ShouldBindJSON(&userLogin);err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input."})
			return 
		}

		if err := validate.Struct(userLogin);err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error.", "details": err.Error()})
			return 
		}

		var foundUser models.User
		if err := userCollection.FindOne(ctx, bson.M{"email": userLogin.Email}).Decode(&foundUser);err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password."})
			return 
		}

		err := bcrypt.CompareHashAndPassword(
			[]byte(foundUser.Password), 
			[]byte(userLogin.Password),
		)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password."})
			return 
		}

		//Generate Tokens
		accessToken, refreshToken, err := utils.GenerateAllTokens(
			foundUser.ID.Hex(),
			foundUser.Role,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens."})
			return 
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": accessToken,
			"refresh_token": refreshToken,
			"user":gin.H{
				"id": foundUser.ID.Hex(),
				"email": foundUser.Email,
				"first_name": foundUser.FirstName,
				"last_name": foundUser.LastName,
				"role": foundUser.Role,
			},
		})
		
	}
} 