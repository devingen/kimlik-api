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

type RegistrationTestSuite struct {
	suite.Suite
	controller controller.IServiceController
	base       string
}

func TestRegistration(t *testing.T) {
	db, err := database.NewDatabaseWithURI("mongodb://localhost")
	if err != nil {
		log.Fatalf("Database connection failed %s", err.Error())
	}
	testSuite := &RegistrationTestSuite{
		controller: controller.NewServiceController(
			service.NewDatabaseService(db),
			json_web_token_service.NewTokenService(),
		),
		base: "dvn-kimlik-api-integration-test",
	}

	InsertTestData(db, testSuite.base)

	suite.Run(t, testSuite)
}

func (suite *RegistrationTestSuite) TestRegisterWithExistingEmail() {
	response, _, err := suite.controller.RegisterWithEmail(context.Background(),
		dvnruntime.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{ \"firstName\": \"User\", \"lastName\": \"New\", \"email\": \"user1@devingen.io\", \"password\": \"123456\"}",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusConflict, err.(*coremodel.DVNError).StatusCode)

	util.SaveResultFile("register-conflict", err)
}

func (suite *RegistrationTestSuite) TestRegisterWithExistingEmailCaseSensitive() {
	response, _, err := suite.controller.RegisterWithEmail(context.Background(),
		dvnruntime.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{ \"firstName\": \"User\", \"lastName\": \"New\", \"email\": \"USER1@DEVINGEN.IO\", \"password\": \"123456\"}",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusConflict, err.(*coremodel.DVNError).StatusCode)

	util.SaveResultFile("register-conflict-case-sensitive", err)
}

func (suite *RegistrationTestSuite) TestRegisterWithInvalidEmail() {
	response, _, err := suite.controller.RegisterWithEmail(context.Background(),
		dvnruntime.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{ \"firstName\": \"User\", \"lastName\": \"New\", \"email\": \"user1 @devingen.io\", \"password\": \"selam\"}",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, err.(*coremodel.DVNError).StatusCode)

	util.SaveResultFile("register-conflict-invalid-email", err)
}

func (suite *RegistrationTestSuite) TestRegisterSuccessful() {
	response, _, err := suite.controller.RegisterWithEmail(context.Background(),
		dvnruntime.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{ \"firstName\": \"User\", \"lastName\": \"Second\", \"email\": \"user2@devingen.io\", \"password\": \"123456\"}",
		},
	)
	assert.Nil(suite.T(), err)
	registerResponse := response.(*dto.RegisterWithEmailResponse)

	assert.NotNil(suite.T(), registerResponse)
	assert.NotEmpty(suite.T(), registerResponse.UserID)
	assert.NotEmpty(suite.T(), registerResponse.JWT)

	util.SaveResultFile("register-successful", err)
}
