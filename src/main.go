package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

// ExternalService represents any external service (S3, Database, API, etc.)
type ExternalService interface {
	Connect(ctx context.Context) error
	Ping(ctx context.Context) error
	GetData(ctx context.Context, key string) (string, error)
	PutData(ctx context.Context, key string, value string) error
	ListKeys(ctx context.Context) ([]string, error)
}

// MockService simulates an external service
type MockService struct {
	name         string
	responseTime time.Duration
	failureRate  float32
	data         map[string]string
}

// NewMockService creates a new mock service
func NewMockService(name string, responseTime time.Duration, failureRate float32) *MockService {
	return &MockService{
		name:         name,
		responseTime: responseTime,
		failureRate:  failureRate,
		data:         make(map[string]string),
	}
}

// Connect simulates connecting to the service
func (m *MockService) Connect(ctx context.Context) error {
	time.Sleep(m.responseTime)
	if m.shouldFail() {
		return fmt.Errorf("failed to connect to %s", m.name)
	}
	fmt.Printf("✓ Connected to %s\n", m.name)
	return nil
}

// Ping simulates a health check
func (m *MockService) Ping(ctx context.Context) error {
	time.Sleep(m.responseTime / 2)
	if m.shouldFail() {
		return fmt.Errorf("%s is not responding", m.name)
	}
	return nil
}

// GetData retrieves data from the mock service
func (m *MockService) GetData(ctx context.Context, key string) (string, error) {
	time.Sleep(m.responseTime)
	if m.shouldFail() {
		return "", fmt.Errorf("failed to get data from %s", m.name)
	}
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key %s not found", key)
}

// PutData stores data in the mock service
func (m *MockService) PutData(ctx context.Context, key string, value string) error {
	time.Sleep(m.responseTime)
	if m.shouldFail() {
		return fmt.Errorf("failed to put data to %s", m.name)
	}
	m.data[key] = value
	return nil
}

// ListKeys returns all keys in the mock service
func (m *MockService) ListKeys(ctx context.Context) ([]string, error) {
	time.Sleep(m.responseTime)
	if m.shouldFail() {
		return nil, fmt.Errorf("failed to list keys from %s", m.name)
	}
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func (m *MockService) shouldFail() bool {
	return rand.Float32() < m.failureRate
}

// ServiceConfig holds configuration for a service
type ServiceConfig struct {
	Name         string
	Type         string
	ResponseTime time.Duration
	FailureRate  float32
}

// LoadServiceConfig loads service configuration from environment
func LoadServiceConfig() []ServiceConfig {
	// Simulate different services with different characteristics
	return []ServiceConfig{
		{
			Name:         "S3-like Storage",
			Type:         os.Getenv("STORAGE_TYPE"),
			ResponseTime: 100 * time.Millisecond,
			FailureRate:  0.01, // 1% failure rate
		},
		{
			Name:         "Database",
			Type:         os.Getenv("DB_TYPE"),
			ResponseTime: 50 * time.Millisecond,
			FailureRate:  0.02, // 2% failure rate
		},
		{
			Name:         "External API",
			Type:         os.Getenv("API_TYPE"),
			ResponseTime: 200 * time.Millisecond,
			FailureRate:  0.05, // 5% failure rate
		},
	}
}

func main() {
	// Initialize random number generator
	// Note: As of Go 1.20, rand.Seed is deprecated and not needed
	// The random number generator is automatically seeded
	ctx := context.Background()

	fmt.Println("=== Integration Testing Demo ===")
	fmt.Println("This simulates integration with external services")
	fmt.Println()

	configs := LoadServiceConfig()
	services := make([]ExternalService, 0, len(configs))

	// Initialize services
	for _, cfg := range configs {
		if cfg.Type == "" {
			cfg.Type = "mock"
		}
		fmt.Printf("Initializing %s service (%s)...\n", cfg.Name, cfg.Type)
		svc := NewMockService(cfg.Name, cfg.ResponseTime, cfg.FailureRate)
		services = append(services, svc)
	}

	fmt.Println("\n--- Running Integration Tests ---")

	// Test each service
	for i, svc := range services {
		cfg := configs[i]
		fmt.Printf("\nTesting %s:\n", cfg.Name)

		// Test connection
		if err := svc.Connect(ctx); err != nil {
			log.Printf("  ✗ Connection failed: %v", err)
			continue
		}

		// Test ping
		if err := svc.Ping(ctx); err != nil {
			log.Printf("  ✗ Ping failed: %v", err)
			continue
		}
		fmt.Printf("  ✓ Ping successful\n")

		// Test data operations
		testKey := fmt.Sprintf("test-key-%d", time.Now().Unix())
		testValue := fmt.Sprintf("test-value-%s", cfg.Name)

		if err := svc.PutData(ctx, testKey, testValue); err != nil {
			log.Printf("  ✗ Put data failed: %v", err)
			continue
		}
		fmt.Printf("  ✓ Data stored successfully\n")

		retrieved, err := svc.GetData(ctx, testKey)
		if err != nil {
			log.Printf("  ✗ Get data failed: %v", err)
			continue
		}
		if retrieved != testValue {
			log.Printf("  ✗ Data mismatch: expected %s, got %s", testValue, retrieved)
			continue
		}
		fmt.Printf("  ✓ Data retrieved successfully\n")

		keys, err := svc.ListKeys(ctx)
		if err != nil {
			log.Printf("  ✗ List keys failed: %v", err)
			continue
		}
		fmt.Printf("  ✓ Listed %d keys\n", len(keys))
	}

	fmt.Println("\n=== Integration Tests Complete ===")
}