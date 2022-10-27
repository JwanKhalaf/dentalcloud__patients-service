package patients

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type PatientStore struct {
	client    *dynamodb.Client
	tableName string
}

type PatientRepository interface {
	GetPatient(ctx context.Context, patientID string) (Patient, error)
}

func NewPatientStore() *PatientStore {
	dynamodbTableName, ok := os.LookupEnv("DYNAMODB_TABLENAME")
	if !ok {
		log.Fatal("the DYNAMODB_TABLENAME variable was not set!")
	}

	log.Printf("The DYNAMODB_TABLENAME variable is set to: %v", dynamodbTableName)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load sdk config: %v", err)
	}

	return &PatientStore{
		client:    dynamodb.NewFromConfig(cfg),
		tableName: dynamodbTableName,
	}
}

func (p *PatientStore) GetPatient(ctx context.Context, patientID string) (Patient, error) {
	patient := Patient{PatientID: patientID}
	response, err := p.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: patient.GetKey(), TableName: aws.String(p.tableName),
	})
	if err != nil {
		log.Printf("could not get find patient with id %q, here is why: %v\n", patientID, err)
	} else {
		if len(response.Item) == 0 {
			return Patient{}, fmt.Errorf("could not find patient with id %q in the database", patientID)
		}

		err = attributevalue.UnmarshalMap(response.Item, &patient)
		if err != nil {
			log.Printf("could not unmarshal response, here is why: %v\n", err)
		}
	}

	return patient, err
}
