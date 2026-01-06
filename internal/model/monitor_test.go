package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimeWithTimeZone_MarshalJSON(t *testing.T) {
	tm := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tz := TimeWithTimeZone(tm)

	data, err := json.Marshal(tz)
	if err != nil {
		t.Errorf("MarshalJSON() error = %v, want nil", err)
	}

	expected := `"2024-01-01T12:00:00Z"`
	if string(data) != expected {
		t.Errorf("MarshalJSON() = %s, want %s", string(data), expected)
	}
}

func TestTimeWithTimeZone_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "valid time",
			data:    []byte(`"2024-01-01T12:00:00Z"`),
			wantErr: false,
		},
		{
			name:    "invalid format",
			data:    []byte(`"invalid"`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tz TimeWithTimeZone
			err := json.Unmarshal(tt.data, &tz)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeWithTimeZone_Value(t *testing.T) {
	tm := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tz := TimeWithTimeZone(tm)

	value, err := tz.Value()
	if err != nil {
		t.Errorf("Value() error = %v, want nil", err)
	}

	if value.(time.Time) != tm {
		t.Errorf("Value() = %v, want %v", value, tm)
	}
}

func TestTimeWithTimeZone_Scan(t *testing.T) {
	tm := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{
			name:    "time.Time",
			value:   tm,
			wantErr: false,
		},
		{
			name:    "[]byte",
			value:   []byte("2024-01-01 12:00:00"),
			wantErr: false,
		},
		{
			name:    "string",
			value:   "2024-01-01 12:00:00",
			wantErr: false,
		},
		{
			name:    "nil",
			value:   nil,
			wantErr: false,
		},
		{
			name:    "invalid type",
			value:   123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tz TimeWithTimeZone
			err := tz.Scan(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeWithTimeZone_Scan_InvalidTimeFormat(t *testing.T) {
	var tz TimeWithTimeZone
	err := tz.Scan([]byte("invalid time"))
	if err == nil {
		t.Error("Scan() should return error for invalid time format")
	}
}

func TestTimeWithTimeZone_ToTime(t *testing.T) {
	tm := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tz := TimeWithTimeZone(tm)

	result := tz.ToTime()
	if result != tm {
		t.Errorf("ToTime() = %v, want %v", result, tm)
	}
}

func TestSystemMetrics_TableName(t *testing.T) {
	metrics := SystemMetrics{}
	if metrics.TableName() != "system_metrics" {
		t.Errorf("TableName() = %s, want system_metrics", metrics.TableName())
	}
}

func TestSystemMetrics_Fields(t *testing.T) {
	metrics := SystemMetrics{
		ID:             "metrics-1",
		CPUUsage:       75.5,
		CPUCores:       4,
		GoroutineCount: 100,
		MemoryAlloc:    1024,
		RequestCount:   1000,
		ErrorCount:     10,
		AverageLatency: 200,
	}

	if metrics.ID != "metrics-1" {
		t.Errorf("ID = %s, want metrics-1", metrics.ID)
	}
	if metrics.CPUUsage != 75.5 {
		t.Errorf("CPUUsage = %f, want 75.5", metrics.CPUUsage)
	}
	if metrics.CPUCores != 4 {
		t.Errorf("CPUCores = %d, want 4", metrics.CPUCores)
	}
	if metrics.GoroutineCount != 100 {
		t.Errorf("GoroutineCount = %d, want 100", metrics.GoroutineCount)
	}
	if metrics.RequestCount != 1000 {
		t.Errorf("RequestCount = %d, want 1000", metrics.RequestCount)
	}
	if metrics.ErrorCount != 10 {
		t.Errorf("ErrorCount = %d, want 10", metrics.ErrorCount)
	}
	if metrics.AverageLatency != 200 {
		t.Errorf("AverageLatency = %d, want 200", metrics.AverageLatency)
	}
}

func TestLogEntry_TableName(t *testing.T) {
	entry := LogEntry{}
	if entry.TableName() != "log_entries" {
		t.Errorf("TableName() = %s, want log_entries", entry.TableName())
	}
}

func TestLogEntry_Fields(t *testing.T) {
	entry := LogEntry{
		ID:         "log-1",
		Level:      "info",
		Service:    "test-service",
		Message:    "Test message",
		Method:     "GET",
		Path:       "/api/test",
		StatusCode: 200,
		Latency:    100,
		IPAddress:  "127.0.0.1",
	}

	if entry.ID != "log-1" {
		t.Errorf("ID = %s, want log-1", entry.ID)
	}
	if entry.Level != "info" {
		t.Errorf("Level = %s, want info", entry.Level)
	}
	if entry.Service != "test-service" {
		t.Errorf("Service = %s, want test-service", entry.Service)
	}
	if entry.Message != "Test message" {
		t.Errorf("Message = %s, want Test message", entry.Message)
	}
	if entry.Method != "GET" {
		t.Errorf("Method = %s, want GET", entry.Method)
	}
	if entry.Path != "/api/test" {
		t.Errorf("Path = %s, want /api/test", entry.Path)
	}
	if entry.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", entry.StatusCode)
	}
	if entry.Latency != 100 {
		t.Errorf("Latency = %d, want 100", entry.Latency)
	}
	if entry.IPAddress != "127.0.0.1" {
		t.Errorf("IPAddress = %s, want 127.0.0.1", entry.IPAddress)
	}
}
