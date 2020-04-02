# meso-license-service
Licence Verfication Service

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

## Verification service should return JSON. Example 
```JSON
{
    "firstName": "NAME",
    "lastName": "NAME",
    "licenseNumber": 111111,
    "licenseType": {
        "boardCode": 0,
        "licenseName": "Registered Nurse",
        "licenseCode": 224
    },
    "Status": " Current",
    "Expiration": "September 30, 2021",
    "Description": "",
    "SecondaryStatus": "",
    "Verify": true
}
