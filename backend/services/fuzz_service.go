package services

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"wsinspect/backend/models"
	"wsinspect/backend/schemas"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type FuzzService struct {
	db *gorm.DB
}

func NewFuzzService(db *gorm.DB) *FuzzService {
	return &FuzzService{db: db}
}

func (s *FuzzService) CreateFuzzTest(req *schemas.CreateFuzzTestRequest) (*models.FuzzTest, error) {
	strategy := req.Strategy
	if strategy == "" {
		strategy = models.FuzzStrategyRandom
	}

	fuzzTest := &models.FuzzTest{
		SessionID: req.SessionID,
		Name:      req.Name,
		Strategy:  strategy,
		Template:  req.Template,
		Status:    models.FuzzStatusPending,
	}

	if err := s.db.Create(fuzzTest).Error; err != nil {
		return nil, err
	}

	return fuzzTest, nil
}

func (s *FuzzService) RunFuzzTest(id uint, targetURL string) error {
	var fuzzTest models.FuzzTest
	if err := s.db.First(&fuzzTest, id).Error; err != nil {
		return err
	}

	// Update status to running
	now := time.Now()
	fuzzTest.Status = models.FuzzStatusRunning
	fuzzTest.StartedAt = &now
	s.db.Save(&fuzzTest)

	// Generate fuzz cases based on strategy
	fuzzCases := s.generateFuzzCases(fuzzTest.Template, fuzzTest.Strategy, 50)

	// Connect to target
	serverConn, _, err := websocket.DefaultDialer.Dial(targetURL, nil)
	if err != nil {
		fuzzTest.Status = models.FuzzStatusFailed
		s.db.Save(&fuzzTest)
		return err
	}
	defer serverConn.Close()

	// Run fuzz tests
	for _, testCase := range fuzzCases {
		messageType := websocket.TextMessage
		if err := serverConn.WriteMessage(messageType, []byte(testCase)); err != nil {
			fuzzTest.FailCount++
		} else {
			fuzzTest.SuccessCount++
		}
		fuzzTest.TestCount++

		// Small delay between tests
		time.Sleep(10 * time.Millisecond)
	}

	// Update status to completed
	completedAt := time.Now()
	fuzzTest.Status = models.FuzzStatusCompleted
	fuzzTest.CompletedAt = &completedAt
	s.db.Save(&fuzzTest)

	return nil
}

func (s *FuzzService) generateFuzzCases(template string, strategy models.FuzzStrategy, count int) []string {
	cases := make([]string, 0, count)

	switch strategy {
	case models.FuzzStrategyRandom:
		cases = s.generateRandomFuzzCases(count)
	case models.FuzzStrategyMutation:
		cases = s.generateMutationFuzzCases(template, count)
	case models.FuzzStrategyBoundary:
		cases = s.generateBoundaryFuzzCases(count)
	case models.FuzzStrategyInvalid:
		cases = s.generateInvalidFuzzCases(count)
	default:
		cases = s.generateRandomFuzzCases(count)
	}

	return cases
}

func (s *FuzzService) generateRandomFuzzCases(count int) []string {
	cases := make([]string, 0, count)

	// Common fuzz patterns
	patterns := []string{
		`{"type": null}`,
		`{"type": ""}`,
		`{"type": 999999999999}`,
		`{}`,
		`{"type": "test", "data": ""}`,
		`{"type": "test", "amount": -1}`,
		`{"type": "test", "amount": 999999999}`,
		`{"invalid": true}`,
		`[]`,
		`"string"`,
		`12345`,
		`true`,
		`null`,
		`{"type": "` + randomString(1000) + `"}`,
	}

	for i := 0; i < count; i++ {
		cases = append(cases, patterns[i%len(patterns)])
	}

	return cases
}

func (s *FuzzService) generateMutationFuzzCases(template string, count int) []string {
	cases := make([]string, 0, count)

	// Try to parse as JSON and mutate
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(template), &data); err != nil {
		return s.generateRandomFuzzCases(count)
	}

	for i := 0; i < count; i++ {
		mutated := s.mutateJSON(data)
		jsonBytes, _ := json.Marshal(mutated)
		cases = append(cases, string(jsonBytes))
	}

	return cases
}

func (s *FuzzService) generateBoundaryFuzzCases(count int) []string {
	cases := []string{
		`{"value": 0}`,
		`{"value": -1}`,
		`{"value": 1}`,
		`{"value": 127}`,
		`{"value": 128}`,
		`{"value": 255}`,
		`{"value": 256}`,
		`{"value": 65535}`,
		`{"value": 65536}`,
		`{"value": 2147483647}`,
		`{"value": -2147483648}`,
		`{"value": ""}`,
		`{"value": "a"}`,
		`{"value": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"}`,
	}

	result := make([]string, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, cases[i%len(cases)])
	}

	return result
}

func (s *FuzzService) generateInvalidFuzzCases(count int) []string {
	cases := []string{
		`{`,
		`}`,
		`[`,
		`]`,
		`{{}`,
		`{{"key": "value"}}`,
		`{"key": }`,
		`{"key": : }`,
		`{"key": "value",}`,
		`{"key": "value" "key2": "value2"}`,
		`not json at all`,
		`{"nested": {"deep": {"value": }}}`,
		`\x00\x01\x02`,
		`{"data": "` + string([]byte{0, 1, 2, 255}) + `"}`,
	}

	result := make([]string, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, cases[i%len(cases)])
	}

	return result
}

func (s *FuzzService) mutateJSON(data map[string]interface{}) map[string]interface{} {
	mutated := make(map[string]interface{})
	for k, v := range data {
		mutated[k] = v

		// Randomly mutate some values
		if rand.Float32() > 0.5 {
			switch val := v.(type) {
			case string:
				mutated[k] = randomString(len(val))
			case float64:
				mutated[k] = val * rand.Float64() * 10
			case bool:
				mutated[k] = !val
			case nil:
				mutated[k] = "mutated"
			}
		}
	}

	// Add random keys
	if rand.Float32() > 0.5 {
		mutated[randomString(10)] = randomString(20)
	}

	return mutated
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *FuzzService) GetFuzzTest(id uint) (*models.FuzzTest, error) {
	var fuzzTest models.FuzzTest
	if err := s.db.First(&fuzzTest, id).Error; err != nil {
		return nil, err
	}
	return &fuzzTest, nil
}

func (s *FuzzService) ListFuzzTests(sessionID uint) ([]models.FuzzTest, error) {
	var tests []models.FuzzTest
	query := s.db.Order("created_at DESC")
	if sessionID > 0 {
		query = query.Where("session_id = ?", sessionID)
	}
	if err := query.Find(&tests).Error; err != nil {
		return nil, err
	}
	return tests, nil
}
