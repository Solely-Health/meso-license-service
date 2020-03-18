package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	inmem "github.com/meso-org/meso-license-service/inmemorydb"
	"github.com/meso-org/meso-license-service/licenses"
	repo "github.com/meso-org/meso-license-service/repository"
)

var licenseSVC licenses.Service

func main() {
	var (
		inmemorydb = true
	)
	var (
		licenseRepository repo.LicenseRepository
	)
	if inmemorydb {
		licenseRepository = inmem.NewLicenseRepository()
	} else {
		//other db
	}
	licenseSVC = licenses.NewService(licenseRepository)

	log.Println("Started License service")
	router := mux.NewRouter()
	router.HandleFunc("/license", licenseRequest)
	//router.HandleFunc("/ping", ping)

	//for local testing
	log.Fatal(http.ListenAndServe(":6060", router))
}

func licenseRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var newLicense repo.License
	var verifiedLicense repo.License
	if err != nil {
		fmt.Fprintf(w, "Error reading body")
	}
	if err := json.Unmarshal(body, &newLicense); err != nil {
		log.Println(err)
	}
	if licenseSVC != nil {
		log.Println("err:licenseSVC nil")
	}
	verifiedLicense, err = licenseSVC.VerifyLicense(newLicense)
	if err != nil {
		log.Print(err)
	}
	licenseSVC.UpdateLicense(verifiedLicense)
	//return struct back as json
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(newLicense)
}
