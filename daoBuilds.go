package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DaoBuilds serializes to/from AWS DynamoDB
type DaoBuilds struct {
	svc *dynamodb.DynamoDB
}

type BuildRecordKey struct {
	Build string
}

type BuildRecord struct {
	BuildRecordKey
	Name        string
	Namespace   string
	Number      int
	LastUpdated string
}

type BuildCounterRecord struct {
	BuildRecordKey
	CurrentCount int
}

var svc = dynamodb.New(session.New(&aws.Config{Region: awsRegion}))

// NewDaoBuilds are used to perform CRUD operations on Builds
func NewDaoBuilds() (*DaoBuilds, error) {
	if svc != nil {
		return &DaoBuilds{svc: svc}, nil
	}
	return nil, fmt.Errorf("unable to construct DaoBuilds object with a <nil> DynamoDB attribute")
}

func formatBuildKey(namespace string, build string, number int) string {
	return namespace + "/" + build + "/" + strconv.Itoa(number)
}

// use DynamoDB's "atomic counter" "feature" to get the next build number
// for each namespace/buildName pair
func (dao *DaoBuilds) incrementBuildCounter(build string) (int, error) {
	o := &BuildRecordKey{
		Build: build,
	}

	item, err := dynamodbattribute.MarshalMap(o)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	params := &dynamodb.UpdateItemInput{
		TableName:        aws.String("jeeves.dev.buildCounters"),
		Key:              item,
		UpdateExpression: aws.String("ADD CurrentCount :incrementValue"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":incrementValue": {
				N: aws.String("1"),
			},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}

	resp, err := dao.svc.UpdateItem(params)
	if err != nil {
		log.Fatalf("updating %v encountered: %v", *params, err)
		return 0, err
	}

	out := &BuildCounterRecord{}
	err = dynamodbattribute.UnmarshalMap(resp.Attributes, out)

	return out.CurrentCount, err
}

func (dao *DaoBuilds) Persist(obj *Build) error {

	number, err := dao.incrementBuildCounter(obj.Namespace + "/" + obj.Build)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("Build Counter for %s/%s: %d\n", obj.Namespace, obj.Build, number)

	o := &BuildRecord{
		BuildRecordKey: BuildRecordKey{
			Build: formatBuildKey(obj.Namespace, obj.Build, obj.Number),
		},
		Name:        obj.Build,
		Namespace:   obj.Namespace,
		Number:      obj.Number,
		LastUpdated: time.Now().Format(time.RFC3339Nano),
	}

	item, err := dynamodbattribute.MarshalMap(o)
	if err != nil {
		log.Fatal(err)
		return err
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String("jeeves.dev.builds"),
		Item:      item,
	}

	_, err = dao.svc.PutItem(params)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (dao *DaoBuilds) Fetch(namespace string, buildName string, number int) (*Build, error) {

	o := &BuildRecordKey{
		Build: formatBuildKey(namespace, buildName, number),
	}

	item, err := dynamodbattribute.MarshalMap(o)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	params := &dynamodb.GetItemInput{
		TableName: aws.String("jeeves.dev.builds"),
		Key:       item,
	}

	resp, err := dao.svc.GetItem(params)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	out := &BuildRecord{}
	err = dynamodbattribute.UnmarshalMap(resp.Item, out)

	updated, err := time.Parse(time.RFC3339Nano, out.LastUpdated)
	if err != nil {
		updated = time.Time{}
	}

	result := &Build{
		Build:       out.Name,
		Namespace:   out.Namespace,
		Number:      out.Number,
		LastUpdated: updated,
	}
	return result, nil
}
