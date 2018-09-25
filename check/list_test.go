package check

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockEC2 struct {
	ec2iface.EC2API
	mock.Mock
}

func (m *mockEC2) DescribeRegions(input *ec2.DescribeRegionsInput) (*ec2.DescribeRegionsOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*ec2.DescribeRegionsOutput), args.Error(1)
}

func TestListAllAwsRegions(t *testing.T) {
	elasticComputer2 := &mockEC2{}
	elasticComputer2.On("DescribeRegions", mock.Anything).Return(&ec2.DescribeRegionsOutput{
		Regions: []*ec2.Region{
			&ec2.Region{RegionName: aws.String("us-west-2")},
			&ec2.Region{RegionName: aws.String("us-east-2")},
		},
	}, nil)

	value, err := ListAllAwsRegions(elasticComputer2)

	assert.Nil(t, err)
	assert.Len(t, value, 2, "Expect two regions")
	assert.Equal(t, "us-west-2", value[0], "Expected us-west-2")
	assert.Equal(t, "us-east-2", value[1], "Expected us-east-2")
	elasticComputer2.AssertExpectations(t)
}

func (m *mockCWL) DescribeLogGroupsPages(input *cloudwatchlogs.DescribeLogGroupsInput, fn func(*cloudwatchlogs.DescribeLogGroupsOutput, bool) bool) error {
	args := m.Called(input)
	fn(args.Get(0).(*cloudwatchlogs.DescribeLogGroupsOutput), true)
	return args.Error(1)
}

func TestListLogGroups(t *testing.T) {
	cwl := &mockCWL{}

	cwl.On("DescribeLogGroupsPages", mock.Anything).Return(&cloudwatchlogs.DescribeLogGroupsOutput{
		LogGroups: []*cloudwatchlogs.LogGroup{
			&cloudwatchlogs.LogGroup{LogGroupName: aws.String("test-group")},
			&cloudwatchlogs.LogGroup{LogGroupName: aws.String("test-group-2")},
		},
	}, nil)

	value, err := ListLogGroups(cwl)

	assert.Nil(t, err, "Expected no error")
	assert.Len(t, value, 2, "Expect two buckets")
	assert.Equal(t, "test-group", value[0], "Expected first log group")
	assert.Equal(t, "test-group-2", value[1], "Expected second log group")
	cwl.AssertExpectations(t)
}

func (m *mockCWL) DescribeLogStreamsPages(input *cloudwatchlogs.DescribeLogStreamsInput, fn func(*cloudwatchlogs.DescribeLogStreamsOutput, bool) bool) error {
	args := m.Called(input)
	fn(args.Get(0).(*cloudwatchlogs.DescribeLogStreamsOutput), true)
	return args.Error(1)
}

func TestListLogStreams(t *testing.T) {
	cwl := &mockCWL{}

	cwl.On("DescribeLogStreamsPages", mock.Anything).Return(&cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: []*cloudwatchlogs.LogStream{
			&cloudwatchlogs.LogStream{LogStreamName: aws.String("test-log-stream")},
			&cloudwatchlogs.LogStream{LogStreamName: aws.String("test-log-stream-2")},
		},
	}, nil)

	value, err := ListLogStreams(cwl, "test-log-stream")

	assert.Nil(t, err, "Execpected no error")
	assert.Equal(t, "test-log-stream", value[0], "Expected first log stream")
	assert.Equal(t, "test-log-stream-2", value[1], "Expected second log stream")
	cwl.AssertExpectations(t)
}
