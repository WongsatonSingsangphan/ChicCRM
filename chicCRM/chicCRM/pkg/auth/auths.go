package auth

import (
	"chicCRM/modules/users/login/models"
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateTokenI(email string) (string, error) {
	var jwtSecret = []byte("thenilalive")
	expirationTime := time.Now().Add(24 * time.Hour)

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = email
	claims["requires_action"] = "change_password" // fixed the typo and added this field to the token
	claims["exp"] = expirationTime.Unix()         // Token expiration time

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func IsTokenMatched(db *sql.DB, token string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM organize_member_credential WHERE orgmbcr_blacklist_token = $1", token).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func IsTokenBlacklisted(db *sql.DB, token string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM organize_member_credential WHERE orgmbcr_blacklist_token = $1", token).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func CreateJWT(claims models.JwtResponse) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("thenilalive"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func CreateJWTTeamlead(claims models.JwtResponseTeamleadSecuredog) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("thenilalive"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
