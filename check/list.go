package check

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

func ListAllAwsRegions(elasticComputer2 ec2iface.EC2API) ([]string, error) {

	regions := []string{}

	region, err := elasticComputer2.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		log.Error("Could not get aws regions: %v", err)
	}
	for _, r := range region.Regions {
		regions = append(regions, *r.RegionName)
	}

	return regions, err
}

func ListLogGroups(cwl cloudwatchlogsiface.CloudWatchLogsAPI) ([]string, error) {

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

func ListLogStreams(cwl cloudwatchlogsiface.CloudWatchLogsAPI, groupName string) ([]string, error) {

	LogStreams := []string{}
	Descending := true

	err := cwl.DescribeLogStreamsPages(
		&cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName: &groupName,
			Descending:   &Descending,
		},
		func(res *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) bool {
			for _, r := range res.LogStreams {
				LogStreams = append(LogStreams, *r.LogStreamName)
			}
			return true
		},
	)
	return LogStreams, err
}
