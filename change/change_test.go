package change

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
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

type mockCWL struct {
	cloudwatchlogsiface.CloudWatchLogsAPI
	mock.Mock
}

func (m *mockCWL) DescribeLogGroups(input *cloudwatchlogs.DescribeLogGroupsInput) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*cloudwatchlogs.DescribeLogGroupsOutput), args.Error(1)
}

func (m *mockCWL) DescribeLogGroupsPages(input *cloudwatchlogs.DescribeLogGroupsInput, fn func(*cloudwatchlogs.DescribeLogGroupsOutput, bool) bool) error {
	args := m.Called(input)
	fn(args.Get(0).(*cloudwatchlogs.DescribeLogGroupsOutput), true)
	return args.Error(1)
}

func (m *mockCWL) PutRetentionPolicy(input *cloudwatchlogs.PutRetentionPolicyInput) (*cloudwatchlogs.PutRetentionPolicyOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*cloudwatchlogs.PutRetentionPolicyOutput), args.Error(1)
}

func (m *mockCWL) DescribeLogStreamsPages(input *cloudwatchlogs.DescribeLogStreamsInput, fn func(*cloudwatchlogs.DescribeLogStreamsOutput, bool) bool) error {
	args := m.Called(input)
	fn(args.Get(0).(*cloudwatchlogs.DescribeLogStreamsOutput), true)
	return args.Error(1)
}

func TestListRegions(t *testing.T) {
	ec2c := &mockEC2{}
	cwlc := &mockCWL{}
	c := &awsClient{
		ec2: ec2c,
		cwl: cwlc,
	}

	ec2c.On("DescribeRegions", mock.Anything).Return(&ec2.DescribeRegionsOutput{
		Regions: []*ec2.Region{
			&ec2.Region{RegionName: aws.String("us-west-2")},
			&ec2.Region{RegionName: aws.String("us-east-2")},
		},
	}, nil)

	value, err := c.ListRegions()

	assert.Nil(t, err)
	assert.Len(t, value, 2, "Expect two regions")
	assert.Equal(t, "us-west-2", value[0], "Expected us-west-2")
	assert.Equal(t, "us-east-2", value[1], "Expected us-east-2")

	ec2c.AssertExpectations(t)
}

func TestListGroups(t *testing.T) {
	ec2c := &mockEC2{}
	cwlc := &mockCWL{}
	c := &awsClient{
		ec2: ec2c,
		cwl: cwlc,
	}

	cwlc.On("DescribeLogGroupsPages", mock.Anything).Return(&cloudwatchlogs.DescribeLogGroupsOutput{
		LogGroups: []*cloudwatchlogs.LogGroup{
			&cloudwatchlogs.LogGroup{LogGroupName: aws.String("test-group")},
			&cloudwatchlogs.LogGroup{LogGroupName: aws.String("test-group-2")},
		},
	}, nil)

	value, err := c.ListGroups()

	assert.Nil(t, err, "Expected no error")
	assert.Len(t, value, 2, "Expect two buckets")
	assert.Equal(t, "test-group", value[0], "Expected first log group")
	assert.Equal(t, "test-group-2", value[1], "Expected second log group")

	cwlc.AssertExpectations(t)
}

func TestListStreams(t *testing.T) {
	ec2c := &mockEC2{}
	cwlc := &mockCWL{}
	c := &awsClient{
		ec2: ec2c,
		cwl: cwlc,
	}

	cwlc.On("DescribeLogStreamsPages", mock.Anything).Return(&cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: []*cloudwatchlogs.LogStream{
			&cloudwatchlogs.LogStream{LogStreamName: aws.String("test-log-stream")},
			&cloudwatchlogs.LogStream{LogStreamName: aws.String("test-log-stream-2")},
		},
	}, nil)

	value, err := c.ListStreams("test-log-stream")

	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, "test-log-stream", value[0], "Expected first log stream")
	assert.Equal(t, "test-log-stream-2", value[1], "Expected second log stream")

	cwlc.AssertExpectations(t)
}

func TestGetRetentionPolicy(t *testing.T) {
	// Arrange
	ec2c := &mockEC2{}
	cwlc := &mockCWL{}
	c := &awsClient{
		ec2: ec2c,
		cwl: cwlc,
	}
	cwlc.On("DescribeLogGroups", mock.Anything).Return(&cloudwatchlogs.DescribeLogGroupsOutput{
		LogGroups: []*cloudwatchlogs.LogGroup{
			{
				LogGroupName:    aws.String("test-group"),
				RetentionInDays: aws.Int64(20),
			},
		},
	}, nil)

	// Act
	value, err := c.GetRetentionPolicy("test-group")

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, int64(20), value)

	cwlc.AssertExpectations(t)
}

func TestSetRetentionPolicy(t *testing.T) {
	ec2c := &mockEC2{}
	cwlc := &mockCWL{}
	c := &awsClient{
		ec2: ec2c,
		cwl: cwlc,
	}

	cwlc.On("PutRetentionPolicy", mock.Anything).Return(&cloudwatchlogs.PutRetentionPolicyOutput{}, nil)

	err := c.SetRetentionPolicy(int64(30), "test-group")

	assert.Nil(t, err)

	cwlc.AssertExpectations(t)
}
