package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
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
	http.Error(w, "github.com/phille97/bostadskoe", http.StatusOK)
}

type StoredResidence struct {
	ID       string
	LastSeen time.Time
	Data     provider.Residence
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

	allErrs := sync.Map{}
	var wg sync.WaitGroup

	for pname, p := range *providers {
		wg.Add(1)
		go func(pname string, p provider.Provider) {
			defer wg.Done()

			out := make(chan provider.Residence)
                        errs := make(chan error)
                        defer close(errs)

                        go func() {
                            list := []error{}
                            for err := range errs {
                                list = append(list, err)
                            }
                            if len(list) > 0 {
                                allErrs.Store(pname, list)
                            }
                        }()

			go p.CurrentResidences(out, errs)

			collection := db.Collection(pname)

			now := time.Now()

			for residence := range out {
				item := collection.Doc(residence.ID())
				_, err = item.Set(ctx, StoredResidence{
					ID:       residence.ID(),
					LastSeen: now,
					Data:     residence,
				})
				if err != nil {
					errs <- err
				}
			}

			iter := collection.Where("LastSeen", "<", time.Now()).Documents(ctx)
			defer iter.Stop()
			for {
				doc, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					errs <- err
					break
				}
				_, err = doc.Ref.Delete(ctx)
				if err != nil {
					errs <- err
				}
			}
		}(pname, p)
	}

	wg.Wait()

	response := struct{
            Errors *sync.Map,
        }{
            Errors: &allErrs,
        }

	if len(response.Errors) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}
