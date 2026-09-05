package category_test

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

type CategoryGraphqlTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler
}

func (s *CategoryGraphqlTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupCategoryService()
	s.handler = testhelper.NewTestHandler(
		testhelper.SetupServiceConnections(s.Conns),
		s.RedisClient(),
		s.Log,
	)
}

func (s *CategoryGraphqlTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *CategoryGraphqlTestSuite) gql(query string) map[string]interface{} {
	body, _ := json.Marshal(map[string]string{"query": query})
	req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

func (s *CategoryGraphqlTestSuite) Test1_FindAllCategory() {
	query := `{ findAllCategory(input: { search: "", page: 1, pageSize: 10 }) { status message } }`
	result := s.gql(query)
	s.Require().NotEmpty(result)
}

func (s *CategoryGraphqlTestSuite) Test2_CreateCategory() {
	query := `mutation { createCategory(input: { name: "GQL Category", description: "Test" }) { status message } }`
	result := s.gql(query)
	s.Require().NotEmpty(result)
}

func (s *CategoryGraphqlTestSuite) Test3_SmokeQuery() {
	query := `{ __typename }`
	result := s.gql(query)
	s.Require().NotEmpty(result)
}

func TestCategoryGraphql(t *testing.T) {
	suite.Run(t, new(CategoryGraphqlTestSuite))
}
