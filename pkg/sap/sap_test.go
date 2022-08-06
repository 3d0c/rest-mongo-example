package sap

import (
	"fmt"
	"testing"
)

func TestParseTimeStamp(t *testing.T) {
	obtained, err := ParseTimeStamp("20150430155806")
	if err != nil {
		t.Fatalf("Unexpected error - %s\n", err)
	}

	fmt.Printf("obtained - %v\n", obtained)
}
