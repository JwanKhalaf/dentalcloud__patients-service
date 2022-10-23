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
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
	response, err := p.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(p.tableName),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: patientID},
		},
	})
	if err != nil {
		return Patient{}, fmt.Errorf("could not get item from the dynamodb table: %w", err)
	}

	var patient Patient
	err = attributevalue.UnmarshalMap(response.Item, &patient)

	return patient, err
}
