package test_test

import (
	"net/http"
	"testing"

	"encoding/json"

	"github.com/pilosa/pilosa/test"
)

func TestNewCluster(t *testing.T) {
	cluster := test.MustNewServerCluster(t, 3)
	response, err := http.Get("http://" + cluster.Servers[0].Server.Addr().String() + "/status")
	if err != nil {
		t.Fatalf("getting schema: %v", err)
	}
	dec := json.NewDecoder(response.Body)
	a := StatusResp{}
	err = dec.Decode(&a)
	if err != nil {
		t.Fatalf("decoding status response: %v", err)
	}

	bytes, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		t.Fatalf("encoding: %v", err)
	}

	if len(a.Status.Nodes) != 3 {
		t.Fatalf("wrong number of nodes in status: %s", bytes)
	}

	for i, node := range a.Status.Nodes {
		if node.State != "UP" {
			t.Fatalf("node %d should be up but is %s", i, node.State)
		}
	}
}

type StatusResp struct {
	Status struct {
		Nodes []struct {
			Host   string
			Schema string
			State  string
		}
	} `json:"status"`
}
