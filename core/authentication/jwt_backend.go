package authentication

import (
	"api.jwt.auth/services/models"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"path/filepath"
)

type JWTAuthenticationBackend struct {
	privateKey []byte
	PublicKey  []byte
}

func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	authBack := new(JWTAuthenticationBackend)
	privateKeyPath, _ := filepath.Abs("./core/authentication/keys/private_key")
	publicKeyPath, _ := filepath.Abs("./core/authentication/keys/public_key.pub")
	authBack.privateKey, _ = ioutil.ReadFile(privateKeyPath)
	authBack.PublicKey, _ = ioutil.ReadFile(publicKeyPath)

	return authBack
}

func (backend *JWTAuthenticationBackend) GenerateToken() string {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	tokenString, _ := token.SignedString(backend.privateKey)
	return tokenString
}

func (backend *JWTAuthenticationBackend) Authenticate(user *models.User) bool {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testing"), 10)

	testUser := models.User{
		Username: "haku",
		Password: string(hashedPassword),
	}

	return user.Username == testUser.Username && bcrypt.CompareHashAndPassword([]byte(testUser.Password), []byte(user.Password)) == nil
}

func (backend *JWTAuthenticationBackend) Logout(token string) error {
	return nil
}
