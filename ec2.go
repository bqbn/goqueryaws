package goqueryaws

import (
	"fmt"
	"strings"

	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

/*
GetInstancesByTag queries AWS EC2 by tags and returns a slice of JSON
encoded service.ec2.Instance.
*/
func GetInstancesByTag(tagKey string, tagValue string) ([][]byte, error) {

	// see https://github.com/aws/aws-sdk-go/blob/master/aws/session/session.go#L30-L33
	// for session struct
	sess := session.Must(session.NewSession())

	ec2Svc := ec2.New(sess)

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String(strings.Join([]string{"tag", tagKey}, ":")),
				Values: []*string{
					aws.String(tagValue),
				},
			},
		},
	}

	result, err := ec2Svc.DescribeInstances(params)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances with %s: %v",
			params, err)
	}

	// instances: slice of JSON encoded service.ec2.Instance
	var instances [][]byte
	for _, rsv := range result.Reservations {
		for _, inst := range rsv.Instances {
			instance, err := json.Marshal(inst)
			if err != nil {
				// Todo: count the errors then return it as an error.
			}
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

/*
GetAllTags retrieves all tags for the specified resource type.
*/
func GetAllTags(resource_type string) ([]*ec2.TagDescription, error) {
	sess := session.Must(session.NewSession())

	ec2Svc := ec2.New(sess)

	params := &ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("resource-type"),
				Values: []*string{
					aws.String(resource_type),
				},
			},
		},
	}

	var tags []*ec2.TagDescription

	err := ec2Svc.DescribeTagsPages(params,
		func(page *ec2.DescribeTagsOutput, lastPage bool) bool {
			for _, tag := range page.Tags {
				tags = append(tags, tag)
			}
			if lastPage {
				return false
			} else {
				return true
			}
		})

	if err != nil {
		return nil, fmt.Errorf("failed to describe tags pages with %s: %v",
			params, err)
	}

	return tags, nil

}
