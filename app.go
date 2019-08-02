package main

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/appengine"

	"github.com/phille97/bostadskoe/provider"
	"github.com/phille97/bostadskoe/provider/bostadstockholm"
)

var (
	providers map[string]provider.Provider
	db        *firestore.Client
)

func main() {
	var err error
	ctx := appengine.BackgroundContext()

	db, err = firestore.NewClient(ctx, "bostadskoe")
	if err != nil {
		panic(err)
	}

	providers["bostadstockholm"], err = bostadstockholm.New()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/tasks/refresh", handleTaskRefresh)
	http.HandleFunc("/", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "github.com/phille97/bostadskoe")
}

func handleTaskRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	for name, provider := range providers {
		collection := db.Collection(name)

		residenceSlice, err := provider.CurrentResidences()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, err.Error())
			return
		}

		for _, residence := range *residenceSlice {
			item := collection.Doc(residence.ID())
			_, err = item.Set(ctx, residence)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintln(w, err.Error())
				return
			}
		}
	}

	fmt.Fprintln(w, "Done")
}
