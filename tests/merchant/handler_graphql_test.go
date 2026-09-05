package merchant_test

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

type MerchantGraphqlTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler
}

func (s *MerchantGraphqlTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupMerchantService()
	s.handler = testhelper.NewTestHandler(
		testhelper.SetupServiceConnections(s.Conns),
		s.RedisClient(),
		s.Log,
	)
}

func (s *MerchantGraphqlTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *MerchantGraphqlTestSuite) gql(query string) map[string]interface{} {
	body, _ := json.Marshal(map[string]string{"query": query})
	req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

func (s *MerchantGraphqlTestSuite) Test1_FindAllMerchant() {
	query := `{ findAllMerchant(input: { search: "", page: 1, pageSize: 10 }) { status message } }`
	result := s.gql(query)
	s.Require().NotEmpty(result)
}

func (s *MerchantGraphqlTestSuite) Test2_SmokeQuery() {
	query := `{ __typename }`
	result := s.gql(query)
	s.Require().NotEmpty(result)
}

func TestMerchantGraphql(t *testing.T) {
	suite.Run(t, new(MerchantGraphqlTestSuite))
}
