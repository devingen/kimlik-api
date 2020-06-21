package integrationtests

import (
	"github.com/devingen/api-core/database"
	coremodel "github.com/devingen/api-core/model"
	"github.com/devingen/api-core/util"
	"github.com/devingen/kimlik-api/controller"
	"github.com/devingen/kimlik-api/dto"
	"github.com/devingen/kimlik-api/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"testing"
)

type LoginTestSuite struct {
	suite.Suite
	controller controller.ServiceController
	base       string
}

func TestLogin(t *testing.T) {
	db, err := database.NewDatabaseWithURI("mongodb://localhost")
	if err != nil {
		log.Fatalf("Database connection failed %s", err.Error())
	}
	testSuite := &LoginTestSuite{
		controller: controller.NewServiceController(service.NewDatabaseService(db)),
		base:       "dvn-kimlik-api-integration-test",
	}

	InsertTestData(db, testSuite.base)

	suite.Run(t, testSuite)
}

func (suite *LoginTestSuite) TestLoginWithNonExistingEmail() {
	response, err := suite.controller.LoginWithEmail(
		suite.base,
		"",
		"",
		"",
		&dto.LoginWithEmailRequest{
			Email:    "user-thats-not-in-system@devingen.io",
			Password: "selam",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusNotFound, err.(*coremodel.DVNError).StatusCode)

	util.SaveResultFile("login-non-existent-email", err)
}

func (suite *LoginTestSuite) TestLoginWithWrongPassword() {
	response, err := suite.controller.LoginWithEmail(
		suite.base,
		"",
		"",
		"",
		&dto.LoginWithEmailRequest{
			Email:    "user1@devingen.io",
			Password: "this-is-not-my-password",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusUnauthorized, err.(*coremodel.DVNError).StatusCode)

	util.SaveResultFile("login-wrong-email", err)
}

func (suite *LoginTestSuite) TestLoginSuccessful() {
	response, err := suite.controller.LoginWithEmail(
		suite.base,
		"",
		"",
		"",
		&dto.LoginWithEmailRequest{
			Email:    "user1@devingen.io",
			Password: "selam",
		},
	)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), response.UserID, "507f191e810c19729de860ea")
	assert.NotEmpty(suite.T(), response.JWT)

	util.SaveResultFile("login-successful", response)
}

func (suite *LoginTestSuite) TestLoginSuccessfulCaseInsensitive() {
	response, err := suite.controller.LoginWithEmail(
		suite.base,
		"",
		"",
		"",
		&dto.LoginWithEmailRequest{
			Email:    "USER1@DEVINGEN.IO",
			Password: "selam",
		},
	)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), response.UserID, "507f191e810c19729de860ea")
	assert.NotEmpty(suite.T(), response.JWT)

	util.SaveResultFile("login-successful-case-insensitive", response)
}

func (suite *LoginTestSuite) TestLoginSuccessfulWhiteSpaces() {
	response, err := suite.controller.LoginWithEmail(
		suite.base,
		"",
		"",
		"",
		&dto.LoginWithEmailRequest{
			Email:    "  user1@DEVINGEN.IO ",
			Password: "selam",
		},
	)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), response.UserID, "507f191e810c19729de860ea")
	assert.NotEmpty(suite.T(), response.JWT)

	util.SaveResultFile("login-successful-white-spaces", response)
}
