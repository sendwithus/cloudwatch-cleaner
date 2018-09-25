package client

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Ec2ClientIntr interface {
	Ec2Client() *ec2.EC2
}

type ElasticCompute2Client struct{}

func (ec2c *ElasticCompute2Client) Ec2Client() *ec2.EC2 {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return ec2.New(sess)
}

type ElasticCompute2ClientMock struct{}

func (cwc *ElasticCompute2ClientMock) Ec2Client() *ec2.EC2 {
	return ec2.New(MockSession)
}
