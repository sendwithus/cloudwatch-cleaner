package change

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCWL struct {
	cloudwatchlogsiface.CloudWatchLogsAPI
	mock.Mock
}

func (m *mockCWL) PutRetentionPolicy(input *cloudwatchlogs.PutRetentionPolicyInput) (*cloudwatchlogs.PutRetentionPolicyOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*cloudwatchlogs.PutRetentionPolicyOutput), args.Error(1)
}

func TestChangeLogGroupsRetentionPolicy(t *testing.T) {
	cwl := &mockCWL{}

	cwl.On("PutRetentionPolicy", mock.Anything).Return(&cloudwatchlogs.PutRetentionPolicyOutput{}, nil)

	err := ChangeLogGroupsRetentionPolicy(cwl, "test-group")

	assert.Nil(t, err)
	cwl.AssertExpectations(t)
}
