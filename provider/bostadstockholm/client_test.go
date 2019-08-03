package bostadstockholm

import (
	"net/http"
	"testing"
)

func TestCurrentResidences(t *testing.T) {
	c, err := New("https://bostad.stockholm.se", http.DefaultClient)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	residences, err := c.CurrentResidences()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if residences == nil {
		t.Fatal("err: residences is null")
	}
}
