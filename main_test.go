package main

import (
	"testing"

	"github.com/techdroplabs/cloudwatch-cleaner/check"

	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	check.Client
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

func TestRun(t *testing.T) {
	c := &mockClient{}

	c.On("ListRegions").Return([]string{"us-west-2"}, nil)
	c.On("SetRegion").Return()
	c.On("ListGroups").Return([]string{"group-name"}, nil)

	run(c)

	c.AssertExpectations(t)
}
