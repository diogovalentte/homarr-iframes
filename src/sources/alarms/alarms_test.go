package alarms

import "testing"

func TestNew(t *testing.T) {
	alarms, err := New()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if alarms == nil {
		t.Fatal("Expected alarms to be not nil")
	}
}
