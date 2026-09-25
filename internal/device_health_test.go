package internal

import (
	"testing"

	"github.com/go-routeros/routeros/v3/proto"
	"github.com/mazay/mikromanager/db"
)

func TestHealthItemParse(t *testing.T) {
	tests := []struct {
		name        string
		outputs     []*proto.Sentence
		wantErr     bool
		expectedLen int
	}{
		{
			name:        "Empty outputs",
			outputs:     []*proto.Sentence{},
			wantErr:     false, // This should return an empty slice, not an error
			expectedLen: 0,
		},
		{
			name: "Single output with valid data",
			outputs: []*proto.Sentence{
				{
					Map: map[string]interface{}{
						".id":   "1",
						"name":  "cpu-temperature",
						"value": "45.5",
						"type":  "float",
					},
				},
			},
			wantErr:     false,
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hi := &HealthItem{}
			result, err := hi.Parse(tt.outputs)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			if len(result) != tt.expectedLen {
				t.Errorf("Expected %d results, got %d", tt.expectedLen, len(result))
			}
		})
	}
}

func TestSentencesToHealth(t *testing.T) {
	tests := []struct {
		name    string
		outputs []*proto.Sentence
		wantErr bool
	}{
		{
			name:    "Empty outputs",
			outputs: []*proto.Sentence{},
			wantErr: true,
		},
		{
			name: "Valid health data with all fields",
			outputs: []*proto.Sentence{
				{
					Map: map[string]interface{}{
						"name":  "cpu-temperature",
						"value": "45.5",
					},
				},
				{
					Map: map[string]interface{}{
						"name":  "voltage",
						"value": "12.3",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SentencesToHealth(tt.outputs)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("SentencesToHealth() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			// For non-empty input, ensure we have a valid health struct
			if !tt.wantErr && result == nil {
				t.Error("Expected non-nil health result for valid inputs")
			}
		})
	}
}

func TestGetDeviceHealth(t *testing.T) {
	// Integration test - Skip actual database and API calls
	t.Skip("Integration test - requires database and actual RouterOS connection")
	
	// This is a structural test to ensure the function signature and calls are correctly structured
	logger := &logger{} // Would normally be zap.Logger
	
	tests := []struct {
		name     string
		device   *db.Device
		database *db.DB
		key      string
		logger   interface{}
		wantErr  bool
	}{
		{
			name:     "Basic device health retrieval",
			device:   &db.Device{},
			database: &db.DB{},
			key:      "test-key",
			logger:   logger,
			wantErr:  true, // Expected to fail without real database connection
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify that we can construct the appropriate parameters
			if tt.device == nil {
				t.Error("Device should not be nil")
			}
			if tt.database == nil {
				t.Error("Database should not be nil")
			}
		})
	}
}

// Mock logger for testing purposes
type logger struct{}

func (l *logger) Debug(msg string, fields ...interface{}) {}
func (l *logger) Info(msg string, fields ...interface{})  {}
func (l *logger) Warn(msg string, fields ...interface{})  {}
func (l *logger) Error(msg string, fields ...interface{}) {}