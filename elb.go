package goqueryaws

import (
	// "fmt"

	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

/*
GetELBsByTag queries AWS ELB by tags and returns a slice of JSON
encoded service.elb.LoadBalancerDescription.
*/
func GetELBsByTag(tagKey string, tagValue string) ([][]byte, error) {

	// see https://github.com/aws/aws-sdk-go/blob/master/aws/session/session.go#L30-L33
	// for session struct
	svc := elb.New(session.Must(session.NewSession()))

	var elbs [][]byte

	err := svc.DescribeLoadBalancersPages(nil,
		func(page *elb.DescribeLoadBalancersOutput, lastPage bool) bool {
			for _, elb_desc := range page.LoadBalancerDescriptions {

				// Describes the tags associated with an elb.
				param := elb.DescribeTagsInput{
					LoadBalancerNames: []*string{
						aws.String(*elb_desc.LoadBalancerName),
					},
				}
				result, err := svc.DescribeTags(&param)

				if err != nil {
					return false
				}

				for _, tag_desc := range result.TagDescriptions {
					for _, tag := range tag_desc.Tags {
						if tagKey == *tag.Key && tagValue == *tag.Value {
							item, err := json.Marshal(elb_desc)
							if err != nil {
								// Todo: count the errors then return it as an error.
							}
							elbs = append(elbs, item)
							break
						}
					}
				}
			}
			if lastPage {
				return false
			} else {
				return true
			}
		})
	return elbs, err
}
