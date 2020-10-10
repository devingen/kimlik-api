package integrationtests

import (
	"context"
	"github.com/devingen/api-core/database"
	"github.com/devingen/api-core/dvnruntime"
	coremodel "github.com/devingen/api-core/model"
	"github.com/devingen/api-core/util"
	"github.com/devingen/kimlik-api/controller"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/service"
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
	db, err := database.NewDatabaseWithURI("mongodb://localhost")
	if err != nil {
		log.Fatalf("Database connection failed %s", err.Error())
	}

	testSuite := &LoginTestSuite{
		controller: controller.NewServiceController(
			service.NewDatabaseService(db),
			json_web_token_service.NewTokenService(),
		),
		base: "dvn-kimlik-api-integration-test",
	}

	InsertTestData(db, testSuite.base)

	suite.Run(t, testSuite)
}

func (suite *LoginTestSuite) TestLoginWithNonExistingEmail() {
	response, _, err := suite.controller.LoginWithEmail(context.Background(),
		dvnruntime.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{\"email\": \"user-thats-not-in-system@devingen.io\", \"password\": \"selam\" }",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusNotFound, err.(*coremodel.DVNError).StatusCode)

	util.SaveResultFile("login-non-existent-email", err)
}

func (suite *LoginTestSuite) TestLoginWithWrongPassword() {
	response, _, err := suite.controller.LoginWithEmail(context.Background(),
		dvnruntime.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{\"email\": \"user1@devingen.io\", \"password\": \"this-is-not-my-password\" }",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusUnauthorized, err.(*coremodel.DVNError).StatusCode)

	util.SaveResultFile("login-wrong-email", err)
}

func (suite *LoginTestSuite) TestLoginSuccessful() {
	response, _, err := suite.controller.LoginWithEmail(context.Background(),
		dvnruntime.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{\"email\": \"user1@devingen.io\", \"password\": \"selam\" }",
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
		dvnruntime.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{\"email\": \"USER1@DEVINGEN.IO\", \"password\": \"selam\" }",
		},
	)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), response)

	loginResponse := response.(*dto.LoginWithEmailResponse)
	assert.Equal(suite.T(), loginResponse.UserID, "507f191e810c19729de860ea")
	assert.NotEmpty(suite.T(), loginResponse.JWT)

	util.SaveResultFile("login-successful-case-insensitive", response)
}

func (suite *LoginTestSuite) TestLoginSuccessfulWhiteSpaces() {
	response, _, err := suite.controller.LoginWithEmail(context.Background(),
		dvnruntime.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{\"email\": \"  user1@DEVINGEN.IO \", \"password\": \"selam\" }",
		},
	)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), response)

	loginResponse := response.(*dto.LoginWithEmailResponse)
	assert.Equal(suite.T(), loginResponse.UserID, "507f191e810c19729de860ea")
	assert.NotEmpty(suite.T(), loginResponse.JWT)

	util.SaveResultFile("login-successful-white-spaces", response)
}
