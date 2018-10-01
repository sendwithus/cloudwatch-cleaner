package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/techdroplabs/cloudwatch-cleaner/change"
)

type mockClient struct {
	change.Client
	mock.Mock
}

func (m *mockClient) SetRegion(region string) {
	m.Called()
}

func (m *mockClient) ListRegions() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockClient) ListGroups() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockClient) ListStreams(group string) ([]string, error) {
	args := m.Called(group)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockClient) GetRetentionPolicy(group string) (int64, error) {
	args := m.Called(group)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockClient) ConvertToInt64(days string) (int64, error) {
	args := m.Called(days)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockClient) SetRetentionPolicy(days int64, group string) error {
	args := m.Called(group)
	return args.Error(0)
}

// Test Run with a retention value of 10 days to trigger the function SetRetentionPolicy
func TestRun(t *testing.T) {
	c := &mockClient{}

	c.On("SetRegion").Return()
	c.On("ListRegions").Return([]string{"us-west-2"}, nil)
	c.On("ListGroups").Return([]string{"group-name"}, nil)
	c.On("GetRetentionPolicy", mock.Anything).Return(int64(10), nil)
	c.On("ConvertToInt64", mock.Anything).Return(int64(30), nil)
	c.On("SetRetentionPolicy", mock.Anything).Return(nil)

	err := run(c)

	assert.Nil(t, err, "Expected no error")
	c.AssertExpectations(t)
}

// Test Run with a retention value of 30 days without using the function SetRetentionPolicy
func TestRun30(t *testing.T) {

	c := &mockClient{}

	c.On("SetRegion").Return()
	c.On("ListRegions").Return([]string{"us-west-2"}, nil)
	c.On("ListGroups").Return([]string{"group-name"}, nil)
	c.On("GetRetentionPolicy", mock.Anything).Return(int64(30), nil)
	c.On("ConvertToInt64", mock.Anything).Return(int64(30), nil)

	err := run(c)

	assert.Nil(t, err, "Expected no error")
	c.AssertExpectations(t)
}
