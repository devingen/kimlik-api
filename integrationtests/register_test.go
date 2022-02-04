package integrationtests

import (
	"context"
	"fmt"
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

type RegistrationTestSuite struct {
	suite.Suite
	controller controller.IServiceController
	base       string
}

func TestRegistration(t *testing.T) {
	db, err := database.New("mongodb://localhost")
	if err != nil {
		log.Fatalf("Database connection failed %s", err.Error())
	}
	testSuite := &RegistrationTestSuite{
		controller: service_controller.New(
			mongods.New("dvn-kimlik-api-integration-test", db),
			json_web_token_service.New("sample-jwt-sign-key"),
		),
		base: "dvn-kimlik-api-integration-test",
	}

	InsertTestData(db, testSuite.base)

	suite.Run(t, testSuite)
}

func (suite *RegistrationTestSuite) TestRegisterWithExistingEmail() {
	response, _, err := suite.controller.RegisterWithEmail(context.Background(),
		core.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{ \"firstName\": \"User\", \"lastName\": \"New\", \"email\": \"user1@devingen.io\", \"password\": \"123456\"}",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusConflict, err.(*core.DVNError).StatusCode)

	util.SaveResultFile("register-conflict", err)
}

func (suite *RegistrationTestSuite) TestRegisterWithExistingEmailCaseSensitive() {
	response, _, err := suite.controller.RegisterWithEmail(context.Background(),
		core.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{ \"firstName\": \"User\", \"lastName\": \"New\", \"email\": \"USER1@DEVINGEN.IO\", \"password\": \"123456\"}",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusConflict, err.(*core.DVNError).StatusCode)

	util.SaveResultFile("register-conflict-case-sensitive", err)
}

func (suite *RegistrationTestSuite) TestRegisterWithInvalidEmail() {
	response, _, err := suite.controller.RegisterWithEmail(context.Background(),
		core.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{ \"firstName\": \"User\", \"lastName\": \"New\", \"email\": \"user1 @devingen.io\", \"password\": \"selam\"}",
		},
	)

	assert.Nil(suite.T(), response)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, err.(*core.DVNError).StatusCode)

	util.SaveResultFile("register-conflict-invalid-email", err)
}

func (suite *RegistrationTestSuite) TestRegisterSuccessful() {
	response, _, err := suite.controller.RegisterWithEmail(context.Background(),
		core.Request{
			PathParameters: map[string]string{
				"base": suite.base,
			},
			Body: "{ \"firstName\": \"User\", \"lastName\": \"Second\", \"email\": \"user2@devingen.io\", \"password\": \"123456\"}",
		},
	)
	assert.Nil(suite.T(), err)
	fmt.Println(err)
	fmt.Println(response)
	registerResponse := response.(*dto.RegisterWithEmailResponse)

	assert.NotNil(suite.T(), registerResponse)
	assert.NotEmpty(suite.T(), registerResponse.UserID)
	assert.NotEmpty(suite.T(), registerResponse.JWT)

	util.SaveResultFile("register-successful", err)
}
