package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"github.com/phille97/bostadskoe/provider"
	"github.com/phille97/bostadskoe/provider/bostadstockholm"
	"github.com/phille97/bostadskoe/provider/senate"
)

func allProviders(ctx context.Context) (*map[string]provider.Provider, error) {
	var err error
	providers := map[string]provider.Provider{}

	httpClient := urlfetch.Client(ctx)

	providers["bostadstockholm"], err = bostadstockholm.New("https://bostad.stockholm.se", httpClient)
	if err != nil {
		return nil, err
	}

	providers["senate"], err = senate.New("http://senate.se", httpClient)
	if err != nil {
		return nil, err
	}

	return &providers, nil
}

func main() {
	http.HandleFunc("/tasks/refresh", handleTaskRefresh)
	http.HandleFunc("/", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=2592000")
	http.Error(w, "github.com/phille97/bostadskoe", http.StatusOK)
}

func handleTaskRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	stderr := log.New(os.Stderr, "", 0)

	if r.Header.Get("X-Appengine-Cron") != "true" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db, err := firestore.NewClient(ctx, "bostadskoe")
	if err != nil {
		stderr.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	providers, err := allProviders(ctx)
	if err != nil {
		stderr.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for name, provider := range *providers {
		collection := db.Collection(name)

		residenceSlice, err := provider.CurrentResidences()
		if err != nil {
			stderr.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, residence := range *residenceSlice {
			item := collection.Doc(residence.ID())
			_, err = item.Set(ctx, residence)
			if err != nil {
				stderr.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	http.Error(w, "done", http.StatusOK)
}
