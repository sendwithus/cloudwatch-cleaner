package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/mock"
)

type mockEC2 struct {
	ec2iface.EC2API
	mock.Mock
}

type mockCWL struct {
	cloudwatchlogsiface.CloudWatchLogsAPI
	mock.Mock
}

func (m *mockEC2) DescribeRegions(input *ec2.DescribeRegionsInput) (*ec2.DescribeRegionsOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*ec2.DescribeRegionsOutput), args.Error(1)
}

func TestRun(t *testing.T) {
	cwl := &mockCWL{}
	elasticComputer2 := &mockEC2{}

	elasticComputer2.On("DescribeRegions", mock.Anything).Return(&ec2.DescribeRegionsOutput{
		Regions: []*ec2.Region{
			&ec2.Region{RegionName: aws.String("us-west-2")},
			&ec2.Region{RegionName: aws.String("us-east-2")},
		},
	}, nil)

	run(cwl, elasticComputer2)

	elasticComputer2.AssertExpectations(t)
}
