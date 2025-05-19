package controllers

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"api-sec-go/config"
	"api-sec-go/models"
	"api-sec-go/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, password, name, and optional plan and type
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "User registration data"
// @Success 201 {object} models.Message  "User registered successfully"
// @Failure 400 {object} models.Message  "Invalid input or email format"
// @Failure 409 {object} models.Message  "User with this email already exists"
// @Failure 500 {object} models.Message  "Internal server error"
// @Router /register [post]
func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.IsValidEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	var existingUser models.User
	err := config.DB.Collection("users").FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "A user with this email already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while checking for existing user"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = nil

	if user.Type == "" || user.Type != string(models.UserTypeAdmin) {
		user.Type = string(models.UserTypeUser)
	}

	if user.Plan == "" || user.Plan != string(models.UserPlanPremium) {
		user.Plan = string(models.UserPlanFree)
	}

	_, err = config.DB.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, models.Message{Message: "User registered successfully"})

}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "User login credentials"
// @Success 200 {object} models.Token "JWT token returned"
// @Failure 400 {object} models.Message  "Invalid input"
// @Failure 401 {object} models.Message "Invalid credentials or password"
// @Failure 500 {object} models.Message  "Internal server error"
// @Router /login [post]
func Login(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := config.DB.Collection("users").FindOne(context.TODO(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    user.Email,
		"name":     user.Name,
		"plan":     user.Plan,
		"type":     user.Type,
		"googleId": user.GoogleID,
		"exp":      time.Now().Add(2 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}

// Update godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile (name, plan, password)
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
//
//	@Param user body object true "User update data" schema(type=object, properties={
//	   email: {type: string, example: "user@example.com"},
//	   name: {type: string, example: "New Name"},
//	   plan: {type: string, example: "free"},
//	   password: {type: string, example: "newpassword123"}
//	})
//
// @Success 200 {object} models.Message  "User updated successfully"
// @Failure 400 {object} models.Message  "Invalid input or no valid fields"
// @Failure 401 {object} models.Message  "Invalid token or claims"
// @Failure 403 {object} models.Message  "Unauthorized to update this user"
// @Failure 404 {object} models.Message  "User not found"
// @Failure 500 {object} models.Message  "Internal server error"
// @Router /auth/update [put]
func Update(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	emailToken := strings.ToLower(strings.TrimSpace(claims["email"].(string)))

	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name"`
		Plan     string `json:"plan"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailBody := strings.ToLower(strings.TrimSpace(input.Email))

	var user models.User
	err = config.DB.Collection("users").FindOne(context.TODO(), bson.M{"email": emailBody}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if emailToken != emailBody {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this user"})
		return
	}

	updateFields := bson.M{}

	if input.Name != "" {
		updateFields["name"] = input.Name
	}

	if input.Plan != "" && (input.Plan == string(models.UserPlanFree) || input.Plan == string(models.UserPlanPremium)) {
		updateFields["plan"] = input.Plan
	}

	if input.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		updateFields["password"] = string(hashed)
	}

	updateFields["updatedAt"] = time.Now()

	if len(updateFields) == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	_, err = config.DB.Collection("users").UpdateOne(
		context.TODO(),
		bson.M{"email": emailBody},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusCreated, models.Message{Message: "User updated successfully"})

}
