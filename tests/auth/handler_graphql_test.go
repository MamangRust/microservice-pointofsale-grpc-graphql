package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MamangRust/monolith-graphql-pointofsale-apigateway/testhelper"
	"github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type AuthGraphqlTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler
}

func (s *AuthGraphqlTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupAuthService()
	s.handler = testhelper.NewTestHandler(
		testhelper.SetupServiceConnections(s.Conns),
		s.RedisClient(),
		s.Log,
	)
}

func (s *AuthGraphqlTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *AuthGraphqlTestSuite) gql(query string) map[string]interface{} {
	body, _ := json.Marshal(map[string]string{"query": query})
	req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

func (s *AuthGraphqlTestSuite) Test1_RegisterUser() {
	query := `mutation { registerUser(input: { firstname: "Auth", lastname: "Test", email: "auth.graphql@test.com", password: "password123", confirmPassword: "password123" }) { status message } }`
	result := s.gql(query)
	s.Require().NotEmpty(result)
}

func (s *AuthGraphqlTestSuite) Test2_SmokeQuery() {
	query := `{ __typename }`
	result := s.gql(query)
	s.Require().NotEmpty(result)
}

func TestAuthGraphql(t *testing.T) {
	suite.Run(t, new(AuthGraphqlTestSuite))
}
