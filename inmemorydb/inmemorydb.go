package inmemorydb

import (
	"sync"

	"github.com/meso-org/meso-license-service/repository"
)

type licenseRepository struct {
	mtx      sync.Mutex
	licenses map[repository.LicenseID]*repository.License
}

func (lr *licenseRepository) Store(lic *repository.License) error {
	lr.mtx.Lock()
	defer lr.mtx.Unlock()
	lr.licenses[lic.ID] = lic
	return nil
}
func (lr *licenseRepository) Find(id repository.LicenseID) (*repository.License, error) {
	return &repository.License{}, nil
}
func (lr *licenseRepository) Update(id repository.LicenseID, status repository.LicenseStatus) (*repository.License, error) {
	return &repository.License{}, nil
}

func NewLicenseRepository() repository.LicenseRepository {
	return &licenseRepository{
		licenses: make(map[repository.LicenseID]*repository.License),
	}
}
