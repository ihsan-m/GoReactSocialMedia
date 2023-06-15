package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/your_project/config"
	"github.com/your_project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// isValidEmail checks if a given string is a valid email address
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !isValidEmail(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if len(user.Password) < 8 ||
		!strings.ContainsAny(user.Password, "abcdefghijklmnopqrstuvwxyz") ||
		!strings.ContainsAny(user.Password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") ||
		!strings.ContainsAny(user.Password, "0123456789") {
		http.Error(w, "Password must be at least 8 characters, and contain at least one lowercase letter, one uppercase letter, and one number", http.StatusBadRequest)
		return
	}

	if len(user.Username) > 15 {
		http.Error(w, "Username must be at most 15 characters", http.StatusBadRequest)
		return
	}

	if len(user.FullName) > 50 {
		http.Error(w, "Full name must be at most 50 characters", http.StatusBadRequest)
		return
	}

	user.Password, err = HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	client := config.Connect()
	userCollection := client.Database("your_database_name").Collection("users")
	_, err = userCollection.InsertOne(r.Context(), user)
	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var login struct {
		EmailOrUsername string `json:"email_or_username"`
		Password        string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client := config.Connect()
	userCollection := client.Database("your_database_name").Collection("users")

	filter := bson.M{"$or": []bson.M{{"email": login.EmailOrUsername}, {"username": login.EmailOrUsername}}}
	var user models.User
	err = userCollection.FindOne(r.Context(), filter).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Invalid email/username or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Error finding user", http.StatusInternalServerError)
		}
		return
	}

	if !CheckPasswordHash(login.Password, user.Password) {
		http.Error(w, "Invalid email/username or password", http.StatusUnauthorized)
		return
	}

	// Generate a session token, store it in a cookie, and save it to your session storage
	// You might want to use a package like "github.com/gorilla/sessions" to manage sessions

	w.WriteHeader(http.StatusOK)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Invalidate the session token, remove it from the session storage, and clear the session cookie
	// You might want to use a package like "github.com/gorilla/sessions" to manage sessions

	w.WriteHeader(http.StatusOK)
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with a hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
