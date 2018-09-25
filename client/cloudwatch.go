package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type CloudWatchClientIntr interface {
	CwClient(string) *cloudwatchlogs.CloudWatchLogs
}
type CloudWatchClient struct{}

func (cwc *CloudWatchClient) CwClient(region string) *cloudwatchlogs.CloudWatchLogs {

	sess := session.Must(session.NewSession())
	return cloudwatchlogs.New(sess, aws.NewConfig().WithRegion(region))
}

type CloudWatchClientMock struct{}

func (cwc *CloudWatchClientMock) CwClientMock(region string) *cloudwatchlogs.CloudWatchLogs {
	return cloudwatchlogs.New(MockSession)
}
