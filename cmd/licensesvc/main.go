package main

import (
	"fmt"
	"net/http"

	inmem "github.com/meso-org/meso-license-service/inmemorydb"
	"github.com/meso-org/meso-license-service/licenses"
	repo "github.com/meso-org/meso-license-service/repository"
	server "github.com/meso-org/meso/server"
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
	srv := server.New(licensesSVC)
	fmt.Println("about to start License server")
	http.ListenAndServe(":6060", srv)
}
