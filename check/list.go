package check

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/techdroplabs/cloudwatch-cleaner/client"
)

func ListAllAwsRegions() ([]string, error) {

	ec2Client := client.Ec2Client()
	regions := []string{}

	region, err := ec2Client.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		log.Error("Could not get aws regions: %v", err)
	}
	for _, r := range region.Regions {
		regions = append(regions, *r.RegionName)
	}

	return regions, err
}

func ListLogGroups(cwl *cloudwatchlogs.CloudWatchLogs) ([]string, error) {

	logGroups := []string{}

	err := cwl.DescribeLogGroupsPages(
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
