package authentication

import (
	"api.jwt.auth/core/redis"
	"api.jwt.auth/services/models"
	"code.google.com/p/go-uuid/uuid"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"path/filepath"
	"time"
)

type JWTAuthenticationBackend struct {
	privateKey []byte
	PublicKey  []byte
}

const (
	tokenDuration = 72
	expireOffset  = 3600
)

func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	authBack := new(JWTAuthenticationBackend)
	privateKeyPath, _ := filepath.Abs("./core/authentication/keys/private_key")
	publicKeyPath, _ := filepath.Abs("./core/authentication/keys/public_key.pub")
	authBack.privateKey, _ = ioutil.ReadFile(privateKeyPath)
	authBack.PublicKey, _ = ioutil.ReadFile(publicKeyPath)

	return authBack
}

func (backend *JWTAuthenticationBackend) GenerateToken(user *models.User) string {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenDuration)).Unix()
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["sub"] = user.UUID
	tokenString, _ := token.SignedString(backend.privateKey)
	return tokenString
}

func (backend *JWTAuthenticationBackend) Authenticate(user *models.User) bool {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testing"), 10)

	testUser := models.User{
		UUID:     uuid.New(),
		Username: "haku",
		Password: string(hashedPassword),
	}

	return user.Username == testUser.Username && bcrypt.CompareHashAndPassword([]byte(testUser.Password), []byte(user.Password)) == nil
}

func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	redisConn := redis.Connect()
	return redisConn.SetValue(tokenString, tokenString, backend.getTokenRemainingValidity(token.Claims["exp"]))
}

func (backend *JWTAuthenticationBackend) IsInBlacklist(token string) bool {
	redisConn := redis.Connect()
	redisToken, _ := redisConn.GetValue(token)

	if redisToken == nil {
		return false
	}

	return true
}
