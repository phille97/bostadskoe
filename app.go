package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/mail"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/phille97/bostadskoe/notification"
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
	LastSeen int64
	Data     provider.ResidenceData
	Raw      interface{}
}

func handleTaskRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Header.Get("X-Appengine-Cron") != "true" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db, err := firestore.NewClient(ctx, "bostadskoe")
	if err != nil {
		log.Errorf(ctx, "could not connect to firestore: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	providers, err := allProviders(ctx)
	if err != nil {
		log.Errorf(ctx, "could not get providers: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allErrs := []error{}
	newProperties := []provider.ResidenceData{}
	var wg sync.WaitGroup

	for pname, p := range *providers {
		wg.Add(1)
		go func(pname string, p provider.Provider) {
			defer wg.Done()

			log.Debugf(ctx, "syncing %s", pname)

			collection := db.Collection(pname)
			now := time.Now().UnixNano()

			out := make(chan provider.Residence)
			errs := make(chan error)
			go p.CurrentResidences(out, errs)

			wg.Add(1)
			go func() {
				defer wg.Done()
				for err := range errs {
					log.Errorf(ctx, "recived error from p.CurrentResidences: %v", err)
					allErrs = append(allErrs, err)
				}
			}()

			for residence := range out {
				log.Debugf(ctx, "[%s] updating residence with ID %s", pname, residence.ID())

				toBeStored := StoredResidence{
					ID:       residence.ID(),
					LastSeen: now,
					Data:     residence.Data(),
					Raw:      residence,
				}

				item := collection.Doc(residence.ID())
				_, err = item.Get(ctx)
				if status.Code(err) == codes.NotFound {
					newProperties = append(newProperties, toBeStored.Data)
				}
				_, err = item.Set(ctx, toBeStored)
				if err != nil {
					log.Errorf(ctx, "could not save residence: %v", err)
					allErrs = append(allErrs, err)
				}
			}

			iter := collection.Where("LastSeen", "<", now).Documents(ctx)
			defer iter.Stop()
			for {
				doc, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					log.Errorf(ctx, "could not fetch obsolete residence for removal: %v", err)
					allErrs = append(allErrs, err)
					break
				}
				log.Debugf(ctx, "[%s] removing obsolete residence with ID %s", pname, doc.Ref.ID)
				_, err = doc.Ref.Delete(ctx)
				if err != nil {
					log.Errorf(ctx, "could not remove obsolete residence: %v", err)
					allErrs = append(allErrs, err)
				}
			}
		}(pname, p)
	}

	wg.Wait()

	if len(newProperties) > 0 {
		templateData := notification.NewPropertiesEmailTemplateData{
			RecipientFirstname: "Philip",
			Residences:         newProperties,
		}

		emailBodyBuf := new(bytes.Buffer)
		err = notification.NewPropertiesEmailTemplate.Execute(emailBodyBuf, templateData)
		if err != nil {
			log.Errorf(ctx, "error generating email from template: %v", err)
			allErrs = append(allErrs, err)
		}

		log.Debugf(ctx, "sending email with %d new properties", len(newProperties))

		msg := &mail.Message{
			Sender:   "Bostadskoe <no-reply@bostadskoe.appspotmail.com>",
			To:       []string{"Phi <phi+bostadskoe@qgr.se>"},
			Subject:  fmt.Sprintf("Found %d new rental apartments", len(newProperties)),
			HTMLBody: emailBodyBuf.String(),
		}
		if err := mail.Send(ctx, msg); err != nil {
			log.Errorf(ctx, "email failed to send: %v", err)
			allErrs = append(allErrs, err)
		}
	}

	response := struct {
		Errors []error
	}{
		Errors: allErrs,
	}

	log.Debugf(ctx, "finished sync, total errors: %d", len(response.Errors))

	if len(response.Errors) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}
