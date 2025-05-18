package launcher

import (
	"database/sql"
	"gosbrw/database"
	"gosbrw/database/structs"
	password_service "gosbrw/services/launcher/password"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckRequiredHeaders(c *gin.Context) (bool, []string) {
	requiredHeaders := []string{
		"X-GameLauncherCertificate",
		"X-GameLauncherHash",
		"X-HiddenHWID",
		"X-HWID",
		"os-version",
		"X-UserAgent",
	}

	var missingHeaders []string
	for _, headerName := range requiredHeaders {
		if c.GetHeader(headerName) == "" {
			missingHeaders = append(missingHeaders, headerName)
		}
	}

	if len(missingHeaders) > 0 {
		return false, missingHeaders
	}

	return true, nil
}

func ParseHTTPRequestHeaders(req *http.Request) map[string]string {
	parseHeaders := make(map[string]string)
	if req != nil {
		for name, values := range req.Header {
			if len(values) > 0 {
				parseHeaders[name] = values[len(values)-1]
			}
		}
	}
	return parseHeaders
}

func UserRegistration(c *gin.Context, request structs.UserRegistrationRequest) (int, map[string]any) {
	allHeadersPresent, missingHeaders := CheckRequiredHeaders(c)
	if !allHeadersPresent {
		return http.StatusBadRequest, gin.H{"error": "Missing required headers", "missing_headers": missingHeaders}
	}
	validPassword := password_service.ValidatePassword(request.Password)
	if !validPassword {
		return http.StatusBadRequest, gin.H{"error": "Password does not meet requirements: must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, and one number."}
	}
	RawEmail := strings.ToLower(request.Email)
	PlusAliasRegex := regexp.MustCompile(`\+[^@]+@`)
	RawEmail = PlusAliasRegex.ReplaceAllString(RawEmail, "@")
	emailFormatRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailFormatRegex.MatchString(RawEmail) {
		return http.StatusBadRequest, gin.H{"error": "invalid email format"}
	}
	UserEmail := RawEmail
	_, err := database.GetUserByEmail(UserEmail)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			log.Printf("Error checking if email exists '%s': %v", UserEmail, err)
			return http.StatusInternalServerError, gin.H{"error": "Error checking email existence"}
		}
	} else {
		return http.StatusConflict, gin.H{"error": "Email already exists"}
	}

	RawPassword := request.Password
	HashedPassword, err := password_service.HashPassword(RawPassword)
	if err != nil {
		log.Printf("Error hashing password for email '%s': %v", UserEmail, err)
		return http.StatusInternalServerError, gin.H{"error": "Internal server error during password hashing"}
	}

	User := structs.User{
		Email:    UserEmail,
		Password: HashedPassword,
		DateCreated: time.Now().Format(time.RFC3339),
	}

	UserID, err := database.CreateNewUser(User)
	if err != nil {
		log.Printf("Error creating new user for email '%s': %v", UserEmail, err)
		return http.StatusInternalServerError, gin.H{"error": "Failed to create user account"}
	}

	log.Printf("Successfully registered user with email: %s, UserID: %d", UserEmail, UserID)
	return http.StatusCreated, gin.H{"message": "User registered successfully", "userID": UserID, "email": UserEmail}
}

func UserLogin(c *gin.Context, request structs.UserLoginRequest) (int, map[string]any) {
	allHeadersPresent, missingHeaders := CheckRequiredHeaders(c)
	if !allHeadersPresent {
		return http.StatusBadRequest, gin.H{"error": "Missing required headers", "missing_headers": missingHeaders}
	}
	validPassword := password_service.ValidatePassword(request.Password)
	if !validPassword {
		return http.StatusBadRequest, gin.H{"error": "Password does not meet requirements: must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, and one number."}
	}
	RawEmail := strings.ToLower(request.Email)
	PlusAliasRegex := regexp.MustCompile(`\+[^@]+@`)
	RawEmail = PlusAliasRegex.ReplaceAllString(RawEmail, "@")
	emailFormatRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailFormatRegex.MatchString(RawEmail) {
		return http.StatusBadRequest, gin.H{"error": "invalid email format"}
	}
	UserEmail := RawEmail

	User, err := database.GetUserByEmail(UserEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, gin.H{"error": "Email not found"}
		} else {
			log.Printf("Error checking if email exists '%s': %v", UserEmail, err)
			return http.StatusInternalServerError, gin.H{"error": "Error checking email existence"}
		}
	}

	RawPassword := request.Password
	HashedPassword, err := password_service.HashPassword(RawPassword)
	if err != nil {
		log.Printf("Error hashing password for email '%s': %v", UserEmail, err)
		return http.StatusInternalServerError, gin.H{"error": "Internal server error during password hashing"}
	}
	PasswordMatch := password_service.CheckPasswordHash(RawPassword, HashedPassword)
	if !PasswordMatch {
		return http.StatusUnauthorized, gin.H{"error": "Invalid password"}
	}

	expiresAt := time.Now().Add(20 * time.Minute).Format(time.RFC3339)
	token := database.GenerateToken(User.ID, expiresAt)
	if token == "" {
		return http.StatusInternalServerError, gin.H{"error": "Failed to generate security token"}
	}

	return http.StatusAccepted, gin.H{"token": token, "userId": User.ID}
}
