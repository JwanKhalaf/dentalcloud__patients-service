package patients

type Patient struct {
	PatientID string `dynamodbav:"patient_id"`
	FirstName string `dynamodbav:"first_name"`
	LastName  string `dynamodbav:"last_name"`
}
