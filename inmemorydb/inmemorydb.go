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
	lr.mtx.Lock()
	defer lr.mtx.Unlock()
	if val, ok := lr.licenses[id]; ok {
		return val, nil
	}
	//TODO some sort of error handling
	//return nil, repository.LicenseFindError
	return &repository.License{}, nil
}
func (lr *licenseRepository) Update(id repository.LicenseID, status repository.LicenseStatus) (*repository.License, error) {
	lr.mtx.Lock()
	defer lr.mtx.Lock()
	c := make([]*repository.License, 0, len(lr.licenses))
	for _, val := range lr.licenses {
		c = append(c, val)
	}
	return &repository.License{}, nil
}

func NewLicenseRepository() repository.LicenseRepository {
	return &licenseRepository{
		licenses: make(map[repository.LicenseID]*repository.License),
	}
}
