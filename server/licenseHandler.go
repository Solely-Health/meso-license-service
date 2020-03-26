package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/meso-org/meso-license-service/licenses"
	repo "github.com/meso-org/meso-license-service/repository"
)

var licenseSVC licenses.Service

type licenseHandler struct {
	s licenses.Service
}

func (h *licenseHandler) router() chi.Router {
	r := chi.NewRouter()

	r.Route("/license", func(chi.Router) {
		//r.Get("/", h.licenseStatus)
		r.Post("/", h.licenseVerify)
		/*
			if we were to add more sub routing:
			r.Route("/pattern", func(chi.Router) {
				r.Verb("/pattern", handlerFunc)
			})
		*/
	})
	r.Get("/ping", h.testPing)
	return r
}

func (h *licenseHandler) testPing(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var response = struct {
		Domain string `json:"domain"`
		Ping   string `json:"ping"`
	}{
		Domain: "license",
		Ping:   "pong",
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		encodeError(ctx, err, w)
		return
	}
}

func (h *licenseHandler) licenseVerify(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var newLicense repo.License
	var verifiedLicense repo.License
	var response struct {
		License repo.License `json:"license"`
	}
	var request struct {
		License repo.License `json:"license"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Printf("unable to decode json: %v", err)
	}
	fmt.Println("Here's the request string ", request.License)

	if err != nil {
		fmt.Fprintf(w, "Error reading body")
	}
	if err := json.Unmarshal(body, &newLicense); err != nil {
		log.Println(err)
	}
	if licenseSVC == nil {
		log.Println("err:licenseSVC nil")
	}
	verifiedLicense, err = licenseSVC.VerifyLicense(newLicense)
	if err != nil {
		log.Print(err)
	}
	licenseSVC.UpdateLicense(verifiedLicense)

	response.License = verifiedLicense
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
