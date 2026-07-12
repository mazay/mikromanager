package internal

import (
	"testing"

	"github.com/go-routeros/routeros/v3/proto"
	"github.com/mazay/mikromanager/db"
)

func TestCpuResourceParse(t *testing.T) {
	tests := []struct {
		name        string
		outputs     []*proto.Sentence
		wantErr     bool
		expectedLen int
	}{
		{
			name:        "Empty outputs",
			outputs:     []*proto.Sentence{},
			wantErr:     false,
			expectedLen: 0,
		},
		{
			name: "Single output with valid data",
			outputs: []*proto.Sentence{
				{
					Map: map[string]interface{}{
						".id":   "1",
						"cpu":   "0",
						"load":  "5",
						"irq":   "2",
						"disk":  "1",
					},
				},
			},
			wantErr:     false,
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CpuResource{}
			result, err := c.Parse(tt.outputs)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			if len(result) != tt.expectedLen {
				t.Errorf("Expected %d results, got %d", tt.expectedLen, len(result))
			}
		})
	}
}

func TestSentencesToCpuResources(t *testing.T) {
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
			name: "Valid CPU resources data",
			outputs: []*proto.Sentence{
				{
					Map: map[string]interface{}{
						".id":  "1",
						"cpu":  "0",
						"load": "5",
						"irq":  "2",
						"disk": "1",
					},
				},
				{
					Map: map[string]interface{}{
						".id":  "2",
						"cpu":  "1",
						"load": "7",
						"irq":  "3",
						"disk": "2",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid CPU resource data (missing core)",
			outputs: []*proto.Sentence{
				{
					Map: map[string]interface{}{
						".id":  "1",
						"cpu":  "", // missing core identifier
						"load": "5",
						"irq":  "2",
						"disk": "1",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SentencesToCpuResources(tt.outputs)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("SentencesToCpuResources() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			// For non-empty input and no error, ensure we have a valid result
			if !tt.wantErr && result == nil {
				t.Error("Expected non-nil results for valid inputs")
			}
		})
	}
}

func TestGetCpuResources(t *testing.T) {
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
			name:     "Basic cpu resources retrieval",
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