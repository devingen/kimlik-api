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

type RegistrationTestSuite struct {
	suite.Suite
	controller controller.ServiceController
	base       string
}

func TestRegistration(t *testing.T) {
	db, err := database.NewDatabaseWithURI("mongodb://localhost")
	if err != nil {
		log.Fatalf("Database connection failed %s", err.Error())
	}
	testSuite := &RegistrationTestSuite{
		controller: controller.NewServiceController(service.NewDatabaseService(db)),
		base:       "dvn-kimlik-api-integration-test",
	}

	InsertTestData(db, testSuite.base)

	suite.Run(t, testSuite)
}

func (suite *RegistrationTestSuite) TestRegisterWithExistingEmail() {
	response, err := suite.controller.RegisterWithEmail(
		suite.base,
		"",
		"",
		"",
		&dto.RegisterWithEmailRequest{
			Email:     "user1@devingen.io",
			FirstName: "User",
			LastName:  "New",
			Password:  "selam",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusConflict, err.(*coremodel.DVNError).StatusCode)

	util.SaveResultFile("register-conflict", err)
}

func (suite *RegistrationTestSuite) TestRegisterWithExistingEmailCaseSensitive() {
	response, err := suite.controller.RegisterWithEmail(
		suite.base,
		"",
		"",
		"",
		&dto.RegisterWithEmailRequest{
			Email:     "USER1@DEVINGEN.IO",
			FirstName: "User",
			LastName:  "New",
			Password:  "selam",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusConflict, err.(*coremodel.DVNError).StatusCode)

	util.SaveResultFile("register-conflict-case-sensitive", err)
}

func (suite *RegistrationTestSuite) TestRegisterWithExistingEmailWhiteSpaces() {
	response, err := suite.controller.RegisterWithEmail(
		suite.base,
		"",
		"",
		"",
		&dto.RegisterWithEmailRequest{
			Email:     " user1@DEVINGEN.IO  ",
			FirstName: "User",
			LastName:  "New",
			Password:  "selam",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusConflict, err.(*coremodel.DVNError).StatusCode)

	util.SaveResultFile("register-conflict-white-spaces", err)
}

func (suite *RegistrationTestSuite) TestRegisterWithInvalidEmail() {
	response, err := suite.controller.RegisterWithEmail(
		suite.base,
		"",
		"",
		"",
		&dto.RegisterWithEmailRequest{
			Email:     "user1 @devingen.io",
			FirstName: "User",
			LastName:  "New",
			Password:  "selam",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, err.(*coremodel.DVNError).StatusCode)
	assert.Equal(suite.T(), "invalid-email", err.(*coremodel.DVNError).Message)

	util.SaveResultFile("register-conflict-invalid-email", err)
}

func (suite *RegistrationTestSuite) TestRegisterSuccessful() {
	response, err := suite.controller.RegisterWithEmail(
		suite.base,
		"",
		"",
		"",
		&dto.RegisterWithEmailRequest{
			Email:     "user2@devingen.io",
			FirstName: "User",
			LastName:  "Second",
			Password:  "selam",
		},
	)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.NotEmpty(suite.T(), response.UserID)
	assert.NotEmpty(suite.T(), response.JWT)

	util.SaveResultFile("register-successful", err)
}
