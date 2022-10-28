package patients

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Patient struct {
	PatientID                         string `dynamodbav:"patient_id" json:"patientId"`
	Title                             string `dynamodbav:"title" json:"title"`
	FirstName                         string `dynamodbav:"first_name" json:"firstName"`
	MiddleName                        string `dynamodbav:"middle_name" json:"middleName"`
	LastName                          string `dynamodbav:"last_name" json:"lastName"`
	NationalInsuranceNumber           string `dynamodbav:"national_insurance_number" json:"nationalInsuranceNumber"`
	Email                             string `dynamodbav:"email" json:"email"`
	Gender                            string `dynamodbav:"gender" json:"gender"`
	DateOfBirth                       string `dynamodbav:"date_of_birth" json:"dateOfBirth"`
	AddressLine1                      string `dynamodbav:"address_line_1" json:"addressLine1"`
	AddressLine2                      string `dynamodbav:"address_line_2" json:"addressLine2"`
	City                              string `dynamodbav:"city" json:"city"`
	County                            string `dynamodbav:"county" json:"county"`
	PostCode                          string `dynamodbav:"post_code" json:"postCode"`
	Country                           string `dynamodbav:"country" json:"country"`
	MobilePhone                       string `dynamodbav:"mobile_phone" json:"mobilePhone"`
	HomePhone                         string `dynamodbav:"home_phone" json:"homePhone"`
	WorkPhone                         string `dynamodbav:"work_phone" json:"workPhone"`
	EmergencyContactFullName          string `dynamodbav:"emergency_contact_full_name" json:"emergencyContactFullName"`
	EmergencyContactPhone             string `dynamodbav:"emergency_contact_phone" json:"emergencyContactPhone"`
	EmergencyContactRelationToPatient string `dynamodbav:"emergency_contact_relation_to_patient" json:"emergencyContactRelationToPatient"`
	Ethnicity                         string `dynamodbav:"ethnicity" json:"ethnicity"`
	Occupation                        string `dynamodbav:"occupation" json:"occupation"`
	AcquisitionSource                 string `dynamodbav:"acquisition_source" json:"acquisitionSource"`
	AssignedDentist                   string `dynamodbav:"assigned_dentist" json:"assignedDentist"`
	AssignedHygienist                 string `dynamodbav:"assigned_hygienist" json:"assignedHygienist"`
	CreatedAt                         string `dynamodbav:"created_at" json:"createdAt"`
	ModifiedAt                        string `dynamodbav:"modified_at" json:"modifiedAt"`
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
