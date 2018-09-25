package check

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type Client interface {
	SetRegion(region string)
	ListRegions() ([]string, error)
	ListGroups() ([]string, error)
	ListStreams(group string) ([]string, error)
	GetRetentionPolicy(group string) (int64, error)
	SetRetentionPolicy(group string) error
}

func New() Client {
	return &awsClient{}
}

type awsClient struct {
	ec2 ec2iface.EC2API
	cwl cloudwatchlogsiface.CloudWatchLogsAPI
}

func (c *awsClient) SetRegion(region string) {
	sess := session.Must(session.NewSession())
	c.ec2 = ec2.New(sess, aws.NewConfig().WithRegion(region))
	c.cwl = cloudwatchlogs.New(sess, aws.NewConfig().WithRegion(region))
}

func (c *awsClient) ListRegions() ([]string, error) {
	regions := []string{}
	region, err := c.ec2.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		log.Error("Could not get aws regions: %v", err)
	}
	for _, r := range region.Regions {
		regions = append(regions, *r.RegionName)
	}
	return regions, err
}

func (c *awsClient) ListGroups() ([]string, error) {
	logGroups := []string{}
	err := c.cwl.DescribeLogGroupsPages(
		&cloudwatchlogs.DescribeLogGroupsInput{},
		func(res *cloudwatchlogs.DescribeLogGroupsOutput, lastPage bool) bool {
			for _, r := range res.LogGroups {
				logGroups = append(logGroups, *r.LogGroupName)
			}
			return true
		},
	)
	return logGroups, err
}

func (c *awsClient) ListStreams(groupName string) ([]string, error) {
	logStreams := []string{}
	descending := true

	err := c.cwl.DescribeLogStreamsPages(
		&cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName: &groupName,
			Descending:   &descending,
		},
		func(res *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) bool {
			for _, r := range res.LogStreams {
				logStreams = append(logStreams, *r.LogStreamName)
			}
			return true
		},
	)
	return logStreams, err
}

func (c *awsClient) GetRetentionPolicy(groupName string) (int64, error) {
	r, err := c.cwl.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: &groupName,
	})
	if err != nil {
		log.WithError(err).Error("Could not describe the log group")
		return int64(-1), err
	}

	for _, group := range r.LogGroups {
		if *group.LogGroupName == groupName {
			if group.RetentionInDays == nil {
				return int64(0), nil
			}
			return *group.RetentionInDays, nil
		}
	}
	return int64(-1), errors.New("OH GOD WHY!??")
}

func (c *awsClient) SetRetentionPolicy(group string) error {
	retention := int64(30)

	_, err := c.cwl.PutRetentionPolicy(&cloudwatchlogs.PutRetentionPolicyInput{
		LogGroupName:    &group,
		RetentionInDays: &retention,
	})
	if err != nil {
		log.Error("Could not setup the retention policy on the log group: %v", err)
	}
	return err
}
