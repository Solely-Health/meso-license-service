package licenses

import (
	"github.com/meso-org/meso-license-service/repository"
)

type Service interface {
	StoreLicense(lic repository.License) (repository.License, error)
}

type service struct {
	licenses repository.LicenseRepository
}

func (s *service) StoreLicense(lic repository.License) (repository.License, error) {
	s.licenses.Store(&lic)
	return lic, nil
}

func NewService(licenseRepository repository.LicenseRepository) Service {
	return &service{
		licenses: licenseRepository,
	}
}
