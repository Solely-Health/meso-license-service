package licenses

import (
	"github.com/meso-org/meso-license-service/repository"
)

type Service interface {
	StoreLicense(lic repository.License)
}

type service struct {
	licenses repository.License
}

func (s *service) StoreLicense(lic repository.License) (repository.License, error) {

}

func NewService(licenseRepository repository.License) Service {
	return &service{
		licenses: licenseRepository,
	}
}
