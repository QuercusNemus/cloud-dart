package main

// import (
// 	"testing"

// 	"github.com/aws/aws-cdk-go/awscdk/v2"
// 	assertions "github.com/aws/aws-cdk-go/awscdk/v2/assertions"
// 	"github.com/aws/jsii-runtime-go"
// )

// example tests. To run these tests, uncomment this file along with the
// example resource in cloud-dart_test.go
// func TestCloudDartStack(t *testing.T) {
// 	// GIVEN
// 	app := awscdk.NewApp(nil)

// 	// WHEN
// 	stack := NewCloudDartStack(app, "MyStack", nil)

// 	// THEN
// 	template := assertions.Template_FromStack(stack)

// 	template.HasResourceProperties(jsii.String("AWS::SQS::Queue"), map[string]interface{}{
// 		"VisibilityTimeout": 300,
// 	})
// }
