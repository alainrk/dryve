package utils

import (
	"dryve/internal/config"
	"dryve/internal/datastruct"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type JWTClaims struct {
	UserID    uint
	UserEmail string
	*jwt.RegisteredClaims
}

// generateJWT generates a signed JWT using the given key and
// sets the claims for this server and the given user.
func GenerateJWT(config config.JWTConfig, user datastruct.User) (string, error) {

	claims := &JWTClaims{
		user.ID,
		user.Email,
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.TTLMins) * time.Minute)),
			ID:        uuid.New().String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    config.Issuer,
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// SigningMethodHS256 is a specific instance of SigningMethodHMAC
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.Key))
}

// extractJWT tries to retreive the token string from the
// Authorization request header: "Authorization: Bearer {JWT}"
func ExtractJWT(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	fmt.Println(bearer)
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

// verifyJWT verifies the given token and returns it.
func VerifyJWT(config config.JWTConfig, encodedJWT string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(encodedJWT, func(token *jwt.Token) (interface{}, error) {
		// Validate used algorithm (SigningMethodHMAC is the general type of above used HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Key), nil
	})

	if err != nil {
		return jwt.MapClaims{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// hashAndSaltPassword uses GenerateFromPassword to hash & salt the password.
//
// MinCost is an integer constant provided by bcrypt along with DefaultCost & MaxCost.
// The cost can be any value you want provided it isn't lower than MinCost.
//
// To init a password: https://go.dev/play/p/zZhbQbE48Cb
func HashAndSaltPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// verifyPassword verify the hashed-salted password in the
// database with the given one.
func VerifyPassword(dbPwd string, givenPwd string) bool {
	byteHash := []byte(dbPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(givenPwd))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// TODO improve this function
func IsValidPassword(pwd string) bool {
	return true
}
