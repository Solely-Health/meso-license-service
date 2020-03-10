package licenses

import (
	"github.com/meso-org/meso-license-service/repository"
)

type Service interface {
	StoreLicense(lic repository.License) (repository.License, error)
	UpdateLicense(lic repository.License) error
	VerfifyLicense(lic repository.License) error
}

type service struct {
	licenses repository.LicenseRepository
}

func (s *service) StoreLicense(lic repository.License) (repository.License, error) {
	s.licenses.Store(&lic)
	return lic, nil
}

func (s *service) UpdateLicense(lic repository.License) error {
	// ObjectL := new(repository.License)
	// ObjectL, err:= s.licenses.Find(lic.ID)
	// if err != nil {

	// }
	// ObjectL = &lic
	// return nil
	return nil
}

func (s *service) VerfifyLicense()

func NewService(licenseRepository repository.LicenseRepository) Service {
	return &service{
		licenses: licenseRepository,
	}
}
