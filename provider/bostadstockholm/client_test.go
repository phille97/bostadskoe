package bostadstockholm

import (
	"net/http"
	"testing"

	"github.com/phille97/bostadskoe/provider"
)

func TestCurrentResidences(t *testing.T) {
	c, err := New("https://bostad.stockholm.se", http.DefaultClient)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	residences := make(chan provider.Residence)
	errs := make(chan error)

	go c.CurrentResidences(residences, errs)

	for residence := range residences {
		t.Log(residence)
	}

	close(errs)

	for err := range errs {
		t.Fatal(err.Error())
	}
}
