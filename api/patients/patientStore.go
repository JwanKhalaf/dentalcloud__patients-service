package patients

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/jsii-runtime-go"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PatientStore struct {
	client    *dynamodb.Client
	tableName string
}

type PatientRepository interface {
	CreatePatient(logger *zap.Logger, ctx context.Context, patient CreatePatientRequest) (CreatePatientResponse, error)
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

func (p *PatientStore) CreatePatient(logger *zap.Logger, ctx context.Context, patient CreatePatientRequest) (CreatePatientResponse, error) {
	// generate the unique patient id
	patient.PatientID = uuid.New().String()

	item, err := attributevalue.MarshalMap(patient)
	if err != nil {
		logger.Error("could not marshal the create patient request for dynamodb", zap.Error(err))
	}

	partitionKey := patient.GetKey()["_pk"]
	sortKey := patient.GetKey()["_sk"]
	item["_pk"] = partitionKey
	item["_sk"] = sortKey
	item["et"] = &types.AttributeValueMemberS{Value: "patient"}
	item["ca"] = &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)}
	item["a"] = &types.AttributeValueMemberBOOL{Value: true}

	_, err = p.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(p.tableName), Item: item,
	})
	if err != nil {
		logger.Error("could not add new patient to dynamodb table", zap.Error(err))
	}

	sharedAttributesMap := map[string]types.AttributeValue{
		"_pk": &types.AttributeValueMemberS{Value: "dp#c9ec3cfe-9f2c-4d68-aec6-9c6a43bf9aec"},
		"pid": &types.AttributeValueMemberS{Value: patient.PatientID},
		"fn":  &types.AttributeValueMemberS{Value: patient.FirstName},
		"mn":  &types.AttributeValueMemberS{Value: patient.MiddleName},
		"ln":  &types.AttributeValueMemberS{Value: patient.LastName},
		"dob": &types.AttributeValueMemberS{Value: patient.DateOfBirth},
		"e":   &types.AttributeValueMemberS{Value: patient.Email},
		"mp":  &types.AttributeValueMemberS{Value: patient.MobilePhone},
		"pc":  &types.AttributeValueMemberS{Value: patient.PostCode},
		"et":  &types.AttributeValueMemberS{Value: "search-item"},
	}

	// for first name
	sharedAttributesMap["_sk"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("p#%v#fn", patient.PatientID)}
	sharedAttributesMap["st"] = &types.AttributeValueMemberS{Value: strings.ToLower(patient.FirstName)}

	_, err = p.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(p.tableName), Item: sharedAttributesMap,
	})
	if err != nil {
		logger.Error("could not add first name for the name gsi", zap.Error(err))
	}

	// clear the first name attribute items
	delete(sharedAttributesMap, "_sk")
	delete(sharedAttributesMap, "st")

	// for first name
	sharedAttributesMap["_sk"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("p#%v#ln", patient.PatientID)}
	sharedAttributesMap["st"] = &types.AttributeValueMemberS{Value: strings.ToLower(patient.LastName)}

	_, err = p.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(p.tableName), Item: sharedAttributesMap,
	})
	if err != nil {
		logger.Error("could not add last name for the name gsi", zap.Error(err))
	}

	return CreatePatientResponse{PatientID: patient.PatientID}, err
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
		err = attributevalue.UnmarshalListOfMaps(response.Items, &patients)
		if err != nil {
			logger.Error("could not unmarshal response", zap.Error(err))
		}
	}

	return patients, err
}
