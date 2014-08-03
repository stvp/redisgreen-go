package redisgreen

import (
	"strings"
	"testing"
)

func TestListMonitors(t *testing.T) {
	monitors, err := ValidTokenClient.ListMonitors()
	if err != nil {
		t.Fatal(err)
	}
	if len(monitors) == 0 {
		t.Error("No monitors returned")
	}
}

func TestCreateMonitor(t *testing.T) {
	// Valid monitor
	monitor, err := ValidTokenClient.CreateMonitor("monitor1", "redis://foobar.redisgreen.net:12345")
	if err != nil {
		t.Fatal(err)
	}
	if monitor.Name != "monitor1" {
		t.Error("Failed to load JSON response correctly")
	}

	// Invalid monitor
	monitor, err = ValidTokenClient.CreateMonitor("monitor2", "oops")
	if err == nil {
		t.Fatal("Expected error, got none")
	}
	if !strings.Contains(err.Error(), `URL should be formatted as`) {
		t.Errorf("Unknown error returned: %#v", err)
	}
}

func TestGetMonitor(t *testing.T) {
	createdMonitor, err := ValidTokenClient.CreateMonitor("monitor2", "redis://foobar.redisgreen.net:12345")
	if err != nil {
		t.Fatal(err)
	}

	// Valid request
	monitor, err := ValidTokenClient.GetMonitor(createdMonitor.Id)
	if err != nil {
		t.Error(err)
	}
	if monitor.Name != createdMonitor.Name {
		t.Errorf("Returned montior %#v instead of %#v", monitor.Name, createdMonitor.Name)
	}

	// 404
	_, err = ValidTokenClient.GetMonitor("oops")
	if err == nil {
		t.Error("Expected error, got none")
	}
	if !strings.Contains(err.Error(), "/monitors/oops couldn't be found") {
		t.Errorf("Unknown error returned: %#v", err)
	}
}

func TestDeleteMonitor(t *testing.T) {
	createdMonitor, err := ValidTokenClient.CreateMonitor("monitor3", "redis://foobar.redisgreen.net:12345")
	if err != nil {
		t.Fatal(err)
	}

	// Valid request
	err = ValidTokenClient.DeleteMonitor(createdMonitor.Id)
	if err != nil {
		t.Error(err)
	}

	// 404
	err = ValidTokenClient.DeleteMonitor("oops")
	if err == nil {
		t.Error("Expected error, got none")
	}
	if !strings.Contains(err.Error(), "/monitors/oops couldn't be found") {
		t.Errorf("Unknown error returned: %#v", err)
	}
}
