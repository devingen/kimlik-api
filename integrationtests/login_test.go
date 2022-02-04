package integrationtests

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/database"
	"github.com/devingen/api-core/util"
	"github.com/devingen/kimlik-api/controller"
	service_controller "github.com/devingen/kimlik-api/controller/service-controller"
	mongods "github.com/devingen/kimlik-api/data-service/mongo-data-service"
	"github.com/devingen/kimlik-api/dto"
	json_web_token_service "github.com/devingen/kimlik-api/token-service/json-web-token-service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"testing"
)

type LoginTestSuite struct {
	suite.Suite
	controller controller.IServiceController
	base       string
}

func TestLogin(t *testing.T) {
	db, err := database.New("mongodb://localhost")
	if err != nil {
		log.Fatalf("Database connection failed %s", err.Error())
	}

	testSuite := &LoginTestSuite{
		controller: service_controller.New(
			mongods.New("dvn-kimlik-api-integration-test", db),
			json_web_token_service.New("sample-jwt-sign-key"),
		),
		base: "dvn-kimlik-api-integration-test",
	}

	InsertTestData(db, testSuite.base)

	suite.Run(t, testSuite)
}

func (suite *LoginTestSuite) TestLoginWithNonExistingEmail() {
	response, _, err := suite.controller.LoginWithEmail(context.Background(),
		core.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{\"email\": \"user-thats-not-in-system@devingen.io\", \"password\": \"123456\" }",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusNotFound, err.(*core.DVNError).StatusCode)

	util.SaveResultFile("login-non-existent-email", err)
}

func (suite *LoginTestSuite) TestLoginWithWrongPassword() {
	response, _, err := suite.controller.LoginWithEmail(context.Background(),
		core.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{\"email\": \"user1@devingen.io\", \"password\": \"this-is-not-my-password\" }",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusUnauthorized, err.(*core.DVNError).StatusCode)

	util.SaveResultFile("login-wrong-email", err)
}

func (suite *LoginTestSuite) TestLoginSuccessful() {
	response, _, err := suite.controller.LoginWithEmail(context.Background(),
		core.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{\"email\": \"user1@devingen.io\", \"password\": \"123456\" }",
		},
	)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), response)

	loginResponse := response.(*dto.LoginWithEmailResponse)
	assert.Equal(suite.T(), loginResponse.UserID, "507f191e810c19729de860ea")
	assert.NotEmpty(suite.T(), loginResponse.JWT)

	util.SaveResultFile("login-successful", response)
}

func (suite *LoginTestSuite) TestLoginSuccessfulCaseInsensitive() {
	response, _, err := suite.controller.LoginWithEmail(context.Background(),
		core.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{\"email\": \"USER1@DEVINGEN.IO\", \"password\": \"123456\" }",
		},
	)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), response)

	loginResponse := response.(*dto.LoginWithEmailResponse)
	assert.Equal(suite.T(), loginResponse.UserID, "507f191e810c19729de860ea")
	assert.NotEmpty(suite.T(), loginResponse.JWT)

	util.SaveResultFile("login-successful-case-insensitive", response)
}
