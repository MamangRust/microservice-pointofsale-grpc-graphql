package order_item_test

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

type OrderItemGraphqlTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler
}

func (s *OrderItemGraphqlTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupOrderItemService()
	s.handler = testhelper.NewTestHandler(
		testhelper.SetupServiceConnections(s.Conns),
		s.RedisClient(),
		s.Log,
	)
}

func (s *OrderItemGraphqlTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *OrderItemGraphqlTestSuite) gql(query string) map[string]interface{} {
	body, _ := json.Marshal(map[string]string{"query": query})
	req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

func (s *OrderItemGraphqlTestSuite) Test1_FindAllOrderItem() {
	query := `{ findAllOrderItem(input: { search: "", page: 1, pageSize: 10 }) { status message } }`
	result := s.gql(query)
	s.Require().NotEmpty(result)
}

func (s *OrderItemGraphqlTestSuite) Test2_SmokeQuery() {
	query := `{ __typename }`
	result := s.gql(query)
	s.Require().NotEmpty(result)
}

func TestOrderItemGraphql(t *testing.T) {
	suite.Run(t, new(OrderItemGraphqlTestSuite))
}
