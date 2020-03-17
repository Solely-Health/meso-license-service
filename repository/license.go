package repository

import (
	"github.com/beevik/guid"
)

type LicenseID string
type LicenseStatus int

//We may not use LicenseStatus because our verfy updates more than just status
const (
	Current LicenseStatus = iota + 1
	Active
	Suspended
)

type LicenseRepository interface {
	Store(lic *License) error
	Find(id LicenseID) (*License, error)
	FindAll() ([]*License, error)
}

type LicenseType struct {
	BoardCode   int    `json:"boardCode"`
	Name        string `json:"licenseName"`
	LicenseCode int    `json:"licenseCode"`
}

//TODO add licenseStatus type
type License struct {
	FirstName       string      `json:"firstName"`
	LastName        string      `json:"lastName"`
	Number          int         `json:"licenseNumber"`
	LicenseDesc     LicenseType `json:"licenseType"`
	Status          string
	Expiration      string
	Description     string
	SecondaryStatus string
	Verify          bool
	ID              LicenseID
}

func GenerateLicenseID() LicenseID {
	return LicenseID(guid.NewString())
}
