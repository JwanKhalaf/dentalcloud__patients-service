package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/jsii-runtime-go"

	"github.com/aws/constructs-go/constructs/v10"
)

const stackName string = "dentalcloud--patients-service"

type PatientsServiceAppStackProps struct {
	awscdk.StackProps
}

func NewPatientsServiceAppStack(scope constructs.Construct, id string, props *PatientsServiceAppStackProps) awscdk.Stack {
	var sprops awscdk.StackProps

	if props != nil {
		sprops = props.StackProps
	}

	// create a new stack
	stack := awscdk.NewStack(scope, &id, &sprops)

	// create a dynamodb table
	table := awsdynamodb.NewTable(stack, jsii.String("dentalcloud"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("_pk"),
			Type: awsdynamodb.AttributeType_STRING},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("_sk"),
			Type: awsdynamodb.AttributeType_STRING},
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
	})

	// add a global secondary index based on name
	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName:        jsii.String("name-index"),
		PartitionKey:     &awsdynamodb.Attribute{Name: jsii.String("_pk"), Type: awsdynamodb.AttributeType_STRING},
		SortKey:          &awsdynamodb.Attribute{Name: jsii.String("st"), Type: awsdynamodb.AttributeType_STRING},
		NonKeyAttributes: jsii.Strings("pid", "fn", "mn", "ln", "e", "mp", "dob", "pc"),
		ProjectionType:   awsdynamodb.ProjectionType_INCLUDE,
	})

	// bundling options to make go fast
	bundlingOptions := &awscdklambdagoalpha.BundlingOptions{
		GoBuildFlags: &[]*string{jsii.String(`-ldflags "-s -w" -tags lambda.norpc`)},
	}

	// creating the aws lambda for creating a patient
	createPatientHandler := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("CreatePatientFunction"), &awscdklambdagoalpha.GoFunctionProps{
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../api/patients/create/lambda"),
		Environment:  &map[string]*string{"DYNAMODB_TABLENAME": table.TableName()},
		Bundling:     bundlingOptions,
		MemorySize:   jsii.Number(1024),
		Timeout:      awscdk.Duration_Millis(jsii.Number(15000)),
	})

	// grant dynamodb read write permissions to the create patient lambda
	table.GrantReadWriteData(createPatientHandler)

	// creating the aws lambda for getting a patient
	getPatientHandler := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("GetPatientFunction"), &awscdklambdagoalpha.GoFunctionProps{
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../api/patients/get/lambda"),
		Environment:  &map[string]*string{"DYNAMODB_TABLENAME": table.TableName()},
		Bundling:     bundlingOptions,
		MemorySize:   jsii.Number(1024),
		Timeout:      awscdk.Duration_Millis(jsii.Number(15000)),
	})

	// grant dynamodb read write permissions to the get patient lambda
	table.GrantReadWriteData(getPatientHandler)

	// creating the aws lambda for finding a patient
	searchPatientsHandler := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("SearchPatientsFunction"), &awscdklambdagoalpha.GoFunctionProps{
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../api/patients/search/lambda"),
		Environment:  &map[string]*string{"DYNAMODB_TABLENAME": table.TableName()},
		Bundling:     bundlingOptions,
		MemorySize:   jsii.Number(1024),
		Timeout:      awscdk.Duration_Millis(jsii.Number(15000)),
	})

	// grant dynamodb read write permissions to the search patients lambda
	table.GrantReadWriteData(searchPatientsHandler)

	// create a new http patientsApi gateway
	patientsApi := awscdkapigatewayv2alpha.NewHttpApi(stack, jsii.String("PatientsApi"), &awscdkapigatewayv2alpha.HttpApiProps{})

	// add route for creating a patient
	patientsApi.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:    jsii.String("/patients"),
		Methods: &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_POST},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("createPatientLambdaIntegration"), createPatientHandler, &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{
			PayloadFormatVersion: awscdkapigatewayv2alpha.PayloadFormatVersion_VERSION_2_0(),
		}),
	})

	// add route for getting a patient
	patientsApi.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:    jsii.String("/patients/{patient-id}"),
		Methods: &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("getPatientLambdaIntegration"), getPatientHandler, &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{
			PayloadFormatVersion: awscdkapigatewayv2alpha.PayloadFormatVersion_VERSION_2_0(),
		}),
	})

	// add route for searching patients
	patientsApi.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:    jsii.String("/patients"),
		Methods: &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("searchPatientsLambdaIntegration"), searchPatientsHandler, &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{
			PayloadFormatVersion: awscdkapigatewayv2alpha.PayloadFormatVersion_VERSION_2_0(),
		}),
	})

	// output the lambda url to the console
	awscdk.NewCfnOutput(stack, jsii.String("PatientsApiUrl"), &awscdk.CfnOutputProps{Value: patientsApi.Url()})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewPatientsServiceAppStack(app, stackName, &PatientsServiceAppStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
