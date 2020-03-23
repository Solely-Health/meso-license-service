package licenses

import (
	"testing"

	"github.com/go-test/deep"

	"github.com/meso-org/meso-license-service/inmemorydb"
	"github.com/meso-org/meso-license-service/repository"
)

func TestLicenseVerify(t *testing.T) {
	var (
		firstname = "RUBY"
		lastname  = "ABRANTES"
		licnumber = 633681
		lictype   = repository.LicenseType{
			BoardCode:   0,
			Name:        "Registered Nurse",
			LicenseCode: 224,
		}
		id = repository.GenerateLicenseID()
	)

	inputObject := repository.License{
		FirstName:       firstname,
		LastName:        lastname,
		Number:          licnumber,
		LicenseDesc:     lictype,
		Status:          "",
		Expiration:      "",
		Description:     "",
		SecondaryStatus: "",
		Verify:          false,
		ID:              id,
	}

	expectedLicenseObject := repository.License{
		FirstName:       firstname,
		LastName:        lastname,
		Number:          licnumber,
		LicenseDesc:     lictype,
		Status:          "Current",
		Expiration:      "September 30, 2021",
		Description:     "",
		SecondaryStatus: "",
		Verify:          true,
		ID:              id,
	}

	mockLicenseRepo := inmemorydb.NewLicenseRepository()
	mockService := NewService(mockLicenseRepo)

	lic, err := mockService.VerifyLicense(inputObject)
	if err != nil {
		t.Fatalf("Failed verification test: %v", err)
	}
	if diff := deep.Equal(expectedLicenseObject, lic); diff != nil {
		t.Fatalf("Return statement failed: %v", diff)
	}
}
