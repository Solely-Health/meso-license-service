# meso-license-service
licence verfication service

## JSON POST request Example 
```JSON
{
	"firstName":"name",
	"lastName":"name",
	"licenseNumber":0,
	"licenseType":{
		"boardCode":0,
		"licenseName":"Registered Nurse",
		"licenseCode":224
	}
}

```
- First/Last name are not case sensitive. It gets converted to uppercase for verification
- licenseName is case sensitive