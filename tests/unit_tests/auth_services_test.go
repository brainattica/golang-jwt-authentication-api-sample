package unit_tests

import (
	"api.jwt.auth/services"
	"api.jwt.auth/services/models"
	"api.jwt.auth/settings"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	. "gopkg.in/check.v1"
	"net/http"
	"os"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type AuthenticationServicesTestSuite struct{}

var _ = Suite(&AuthenticationServicesTestSuite{})
var t *testing.T

func (s *AuthenticationServicesTestSuite) SetUpSuite(c *C) {
	os.Setenv("GO_ENV", "tests")
	settings.Init()
}

func (suite *AuthenticationServicesTestSuite) TestLogin(c *C) {
	user := models.User{
		Username: "haku",
		Password: "testing",
	}
	response, token := services.Login(&user)
	assert.Equal(t, http.StatusOK, response)
	assert.NotEmpty(t, token)
}

func (suite *AuthenticationServicesTestSuite) TestLoginIncorrectPassword(c *C) {
	user := models.User{
		Username: "haku",
		Password: "Password",
	}
	response, token := services.Login(&user)
	assert.Equal(t, http.StatusUnauthorized, response)
	assert.Empty(t, token)
}

func (suite *AuthenticationServicesTestSuite) TestLoginIncorrectUsername(c *C) {
	user := models.User{
		Username: "Username",
		Password: "testing",
	}
	response, token := services.Login(&user)
	assert.Equal(t, http.StatusUnauthorized, response)
	assert.Empty(t, token)
}

func (suite *AuthenticationServicesTestSuite) TestLoginEmptyCredentials(c *C) {
	user := models.User{
		Username: "",
		Password: "",
	}
	response, token := services.Login(&user)
	assert.Equal(t, http.StatusUnauthorized, response)
	assert.Empty(t, token)
}

func (suite *AuthenticationServicesTestSuite) TestRefreshToken(c *C) {
	user := models.User{
		Username: "haku",
		Password: "testing",
	}
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenString, err := authBackend.GenerateToken(user.UUID)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	assert.Nil(t, err)

	newToken := services.RefreshToken(token)
	assert.NotEmpty(t, newToken)
}

func (suite *AuthenticationServicesTestSuite) TestRefreshTokenInvalidToken(c *C) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	newToken := services.RefreshToken(token)
	assert.Empty(t, newToken)
}

func (suite *AuthenticationServicesTestSuite) TestLogout(c *C) {
	user := models.User{
		Username: "haku",
		Password: "testing",
	}
	authBackend := auth.InitJWTAuthenticationBackend()
	tokenString, err := authentication.GenerateToken(user.UUID)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})

	err = services.Logout(tokenString, token)
	assert.Nil(t, err)
}
