package patients

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Patient struct {
	PatientID                         string `dynamodbav:"patient_id"`
	Title                             string `dynamodbav:"title"`
	FirstName                         string `dynamodbav:"first_name"`
	MiddleName                        string `dynamodbav:"middle_name"`
	LastName                          string `dynamodbav:"last_name"`
	NationalInsuranceNumber           string `dynamodbav:"national_insurance_number"`
	Email                             string `dynamodbav:"email"`
	Gender                            string `dynamodbav:"gender"`
	DateOfBirth                       string `dynamodbav:"date_of_birth"`
	AddressLine1                      string `dynamodbav:"address_line_1"`
	AddressLine2                      string `dynamodbav:"address_line_2"`
	City                              string `dynamodbav:"city"`
	County                            string `dynamodbav:"county"`
	PostCode                          string `dynamodbav:"post_code"`
	Country                           string `dynamodbav:"country"`
	MobilePhone                       string `dynamodbav:"mobile_phone"`
	HomePhone                         string `dynamodbav:"home_phone"`
	WorkPhone                         string `dynamodbav:"work_phone"`
	EmergencyContactFullName          string `dynamodbav:"emergency_contact_full_name"`
	EmergencyContactPhone             string `dynamodbav:"emergency_contact_phone"`
	EmergencyContactRelationToPatient string `dynamodbav:"emergency_contact_relation_to_patient"`
	Ethnicity                         string `dynamodbav:"enthnicity"`
	Occupation                        string `dynamodbav:"occupation"`
	AcquisitionSource                 string `dynamodbav:"acquisition_source"`
}

// returns the composite primary key of the patient in a format that can be
// sent to dynamo.
func (p Patient) GetKey() map[string]types.AttributeValue {
	patientID, err := attributevalue.Marshal(fmt.Sprintf("p#%v", p.PatientID))
	if err != nil {
		panic(err)
	}
	dentalPracticeID, err := attributevalue.Marshal(fmt.Sprintf("dp#%v", 1))
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"pk": dentalPracticeID, "sk": patientID}
}
