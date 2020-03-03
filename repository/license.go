package repository

import (
	"github.com/beevik/guid"
)

type LicenseID string
type LicenseStatus int

const (
	Current LicenseStatus = iota + 1
	Active
	Suspended
)

type LicenseRepository interface {
	Store(lic *License)
	Find(id LicenseID) (*License, error)
	Update(id LicenseID, status LicenseStatus) (*License, error)
}

type LicenseType struct {
	BoardCode   int    `json:"boardCode"`
	Name        string `json:"licenseName"`
	LicenseCode int    `json:"licenseCode"`
}

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
