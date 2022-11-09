package patients

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/jsii-runtime-go"
	"go.uber.org/zap"
)

type PatientStore struct {
	client    *dynamodb.Client
	tableName string
}

type PatientRepository interface {
	GetPatient(logger *zap.Logger, ctx context.Context, patientID string) (Patient, error)
	SearchPatients(logger *zap.Logger, ctx context.Context, searchTerm string) ([]PatientSearchResponseItem, error)
}

func NewPatientStore(logger *zap.Logger) *PatientStore {
	dynamodbTableName, ok := os.LookupEnv("DYNAMODB_TABLENAME")
	if !ok {
		logger.Fatal("the DYNAMODB_TABLENAME variable was not set!")
	}

	logger.Info("The DYNAMODB_TABLENAME variable is set", zap.String("DYNAMODB_TABLENAME", dynamodbTableName))

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.Fatal("unable to load sdk config", zap.Error(err))
	}

	return &PatientStore{
		client:    dynamodb.NewFromConfig(cfg),
		tableName: dynamodbTableName,
	}
}

func (p *PatientStore) GetPatient(logger *zap.Logger, ctx context.Context, patientID string) (Patient, error) {
	logger.Info("getting patient")
	patient := Patient{PatientID: patientID}
	response, err := p.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: patient.GetKey(), TableName: aws.String(p.tableName),
	})
	if err != nil {
		logger.Error("could not get find matching patient", zap.Error(err))
	} else {
		if len(response.Item) == 0 {
			return Patient{}, fmt.Errorf("could not find patient with id %q in the database", patientID)
		}

		err = attributevalue.UnmarshalMap(response.Item, &patient)
		if err != nil {
			logger.Error("could not unmarshal response", zap.Error(err))
		}
	}

	return patient, err
}

func (p *PatientStore) SearchPatients(logger *zap.Logger, ctx context.Context, searchTerm string) ([]PatientSearchResponseItem, error) {
	dentalPracticeID := "c9ec3cfe-9f2c-4d68-aec6-9c6a43bf9aec"
	lowerCaseSearchTerm := strings.ToLower(searchTerm)
	logger.Info("searching patients", zap.String("dentalPracticeID", dentalPracticeID))
	var patients []PatientSearchResponseItem
	response, err := p.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(p.tableName),
		IndexName:              jsii.String("name-index"),
		KeyConditionExpression: jsii.String("#_pk = :dpid and begins_with(#st, :st)"),
		ExpressionAttributeNames: map[string]string{
			"#_pk": "_pk",
			"#st":  "st",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":dpid": &types.AttributeValueMemberS{Value: fmt.Sprintf("dp#%v", dentalPracticeID)},
			":st":   &types.AttributeValueMemberS{Value: lowerCaseSearchTerm},
		},
	})
	if err != nil {
		logger.Error("could not find matching patients", zap.Error(err))
	} else {
		if len(response.Items) == 0 {
			return patients, nil
		}

		err = attributevalue.UnmarshalListOfMaps(response.Items, &patients)
		if err != nil {
			logger.Error("could not unmarshal response", zap.Error(err))
		}
	}

	return patients, err
}
