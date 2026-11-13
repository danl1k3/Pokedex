package api

import (
	"fmt"
	"testing"
	"time"
)

var client = NewClient(5 * time.Second)

func TestDebugApiCall(t *testing.T) {
	res, err := client.GetAreas()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	

}
