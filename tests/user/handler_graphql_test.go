package user_test

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

type UserGraphqlTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler
}

func (s *UserGraphqlTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.handler = testhelper.NewTestHandler(
		testhelper.SetupServiceConnections(s.Conns),
		s.RedisClient(),
		s.Log,
	)
}

func (s *UserGraphqlTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *UserGraphqlTestSuite) gql(query string) map[string]interface{} {
	body, _ := json.Marshal(map[string]string{"query": query})
	req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

func (s *UserGraphqlTestSuite) Test1_FindAllUsers() {
	query := `{ findAllUsers(input: { search: "", page: 1, pageSize: 10 }) { status message } }`
	result := s.gql(query)
	s.Require().NotEmpty(result)
}

func (s *UserGraphqlTestSuite) Test2_CreateUser() {
	query := `mutation { createUser(input: { firstname: "GraphQL", lastname: "User", email: "gql.user@test.com", password: "password123", confirmPassword: "password123" }) { status message } }`
	result := s.gql(query)
	s.Require().NotEmpty(result)
}

func (s *UserGraphqlTestSuite) Test3_SmokeQuery() {
	query := `{ __typename }`
	result := s.gql(query)
	s.Require().NotEmpty(result)
}

func TestUserGraphql(t *testing.T) {
	suite.Run(t, new(UserGraphqlTestSuite))
}
