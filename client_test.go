package redisgreen

import (
	"os"
	"strings"
	"testing"
)

var (
	ValidTokenClient   = Client{os.Getenv("TOKEN")}
	InvalidTokenClient = Client{"invalid"}
)

func init() {
	if len(ValidTokenClient.Token) == 0 {
		panic("TOKEN environment variable missing")
	}
	ApiUrl = "http://localhost:4100"
}

func TestListServers(t *testing.T) {
	// Valid token
	servers, err := ValidTokenClient.ListServers()
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) == 0 {
		t.Error("No servers returned")
	}

	// Invalid token
	_, err = InvalidTokenClient.ListServers()
	if err == nil {
		t.Error("Expected error, got none")
	}
	if !strings.Contains(err.Error(), "X-API-Token header was missing or incorrect") {
		t.Errorf("Unknown error returned: %#v", err)
	}
}

func TestCreateServer(t *testing.T) {
	// Valid server with slaves
	server, err := ValidTokenClient.CreateServer("test1", "dev", "us-east-1", 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(server.Slaves) != 2 {
		t.Errorf("Expected 2 slaves, got %d", len(server.Slaves))
	}

	// Invalid plan
	server, err = ValidTokenClient.CreateServer("test2", "oh no!", "us-east-1", 0)
	if err == nil {
		t.Error("Expected error, got none")
	}
	if !strings.Contains(err.Error(), `"oh no!" is not a valid plan`) {
		t.Errorf("Unknown error returned: %#v", err)
	}

	// Invalid token
	_, err = InvalidTokenClient.CreateServer("oops", "dev", "us-east-1", 0)
	if err == nil {
		t.Error("Expected error, got none")
	}
	if !strings.Contains(err.Error(), "X-API-Token header was missing or incorrect") {
		t.Errorf("Unknown error returned: %#v", err)
	}
}

func TestGetServer(t *testing.T) {
	createdServer, err := ValidTokenClient.CreateServer("test3", "dev", "us-east-1", 2)
	if err != nil {
		t.Fatal(err)
	}

	// Valid request
	server, err := ValidTokenClient.GetServer(createdServer.Id)
	if err != nil {
		t.Error(err)
	}
	if len(server.Slaves) != 2 {
		t.Error("Incomplete server record returned")
	}

	// 404
	_, err = ValidTokenClient.GetServer("oops")
	if err == nil {
		t.Error("Expected error, got none")
	}
	if !strings.Contains(err.Error(), "/servers/oops couldn't be found") {
		t.Errorf("Unknown error returned: %#v", err)
	}
}

func TestDeleteServer(t *testing.T) {
	createdServer, err := ValidTokenClient.CreateServer("test4", "dev", "us-east-1", 0)
	if err != nil {
		t.Fatal(err)
	}

	// Valid request
	err = ValidTokenClient.DeleteServer(createdServer.Id)
	if err != nil {
		t.Error(err)
	}

	// 404
	err = ValidTokenClient.DeleteServer("oops")
	if err == nil {
		t.Error("Expected error, got none")
	}
	if !strings.Contains(err.Error(), "/servers/oops couldn't be found") {
		t.Errorf("Unknown error returned: %#v", err)
	}
}
