package check

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	cwliface "github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockCWL struct {
	cwliface.CloudWatchLogsAPI
	mock.Mock
}

func (m *mockCWL) DescribeLogGroups(input *cloudwatchlogs.DescribeLogGroupsInput) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*cloudwatchlogs.DescribeLogGroupsOutput), args.Error(1)
}

func TestCheckLogGroupsRetentionPolicy(t *testing.T) {
	// Arrange
	cwl := &mockCWL{}
	cwl.On("DescribeLogGroups", mock.Anything).Return(&cloudwatchlogs.DescribeLogGroupsOutput{
		LogGroups: []*cloudwatchlogs.LogGroup{
			{
				LogGroupName:    aws.String("test-group"),
				RetentionInDays: aws.Int64(20),
			},
		},
	}, nil)

	// Act
	value, err := CheckLogGroupsRetentionPolicy(cwl, "test-group")

	// Assert
	require.Nil(t, err)
	require.Equal(t, int64(20), value)
	cwl.AssertExpectations(t)
}
