package patients

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Patient struct {
	PatientID                         string `dynamodbav:"pid" json:"patient_id"`
	Title                             string `dynamodbav:"t" json:"title"`
	FirstName                         string `dynamodbav:"fn" json:"first_name"`
	MiddleName                        string `dynamodbav:"mn" json:"middle_name"`
	LastName                          string `dynamodbav:"ln" json:"last_name"`
	NationalInsuranceNumber           string `dynamodbav:"ni" json:"national_insurance_number"`
	Email                             string `dynamodbav:"e" json:"email"`
	Gender                            string `dynamodbav:"g" json:"gender"`
	DateOfBirth                       string `dynamodbav:"dob" json:"date_of_birth"`
	AddressLine1                      string `dynamodbav:"al1" json:"address_line_1"`
	AddressLine2                      string `dynamodbav:"al2" json:"address_line_2"`
	City                              string `dynamodbav:"c" json:"city"`
	County                            string `dynamodbav:"cty" json:"county"`
	PostCode                          string `dynamodbav:"pc" json:"post_code"`
	Country                           string `dynamodbav:"ctry" json:"country"`
	MobilePhone                       string `dynamodbav:"mp" json:"mobile_phone"`
	HomePhone                         string `dynamodbav:"hp" json:"home_phone"`
	WorkPhone                         string `dynamodbav:"wp" json:"work_phone"`
	EmergencyContactFullName          string `dynamodbav:"ecfn" json:"emergency_contact_full_name"`
	EmergencyContactPhone             string `dynamodbav:"ecp" json:"emergency_contact_phone"`
	EmergencyContactRelationToPatient string `dynamodbav:"ecrtp" json:"emergency_contact_relation_to_patient"`
	Ethnicity                         string `dynamodbav:"eth" json:"ethnicity"`
	Occupation                        string `dynamodbav:"o" json:"occupation"`
	AcquisitionSource                 string `dynamodbav:"as" json:"acquisition_source"`
	AssignedDentist                   string `dynamodbav:"ad" json:"assigned_dentist"`
	AssignedHygienist                 string `dynamodbav:"ah" json:"assigned_hygienist"`
	Active                            bool   `dynamodbav:"a" json:"active"`
	CreatedAt                         string `dynamodbav:"ca" json:"created_at"`
	ModifiedAt                        string `dynamodbav:"ma" json:"modified_at"`
}

type PatientSearchResponseItem struct {
	PatientID   string `dynamodbav:"pid" json:"patient_id"`
	FirstName   string `dynamodbav:"fn" json:"first_name"`
	MiddleName  string `dynamodbav:"mn" json:"middle_name"`
	LastName    string `dynamodbav:"ln" json:"last_name"`
	DateOfBirth string `dynamodbav:"dob" json:"date_of_birth"`
	Email       string `dynamodbav:"e" json:"email"`
	MobilePhone string `dynamodbav:"mp" json:"mobile_phone"`
	PostCode    string `dynamodbav:"pc" json:"post_code"`
}

// returns the composite primary key of the patient in a format that can be
// sent to dynamo.
func (p Patient) GetKey() map[string]types.AttributeValue {
	patientID, err := attributevalue.Marshal(fmt.Sprintf("p#%v", p.PatientID))
	if err != nil {
		panic(err)
	}

	dentalPracticeID, err := attributevalue.Marshal(fmt.Sprintf("dp#%v", "c9ec3cfe-9f2c-4d68-aec6-9c6a43bf9aec"))
	if err != nil {
		panic(err)
	}

	return map[string]types.AttributeValue{"_pk": dentalPracticeID, "_sk": patientID}
}
