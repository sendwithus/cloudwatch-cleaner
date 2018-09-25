package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/techdroplabs/cloudwatch-cleaner/change"
	"github.com/techdroplabs/cloudwatch-cleaner/check"
	"github.com/techdroplabs/cloudwatch-cleaner/client"
)

func main() {
	cwc := client.CloudWatchClient{}
	ec2c := client.ElasticCompute2Client{}
	run(&cwc, &ec2c)
}

func run(cwc client.CloudWatchClientIntr, ec2c client.Ec2ClientIntr) {
	ec2Client := ec2c.Ec2Client()
	regions, err := check.ListAllAwsRegions(ec2Client)
	if err != nil {
		log.WithError(err).Error("could not list the aws regions")
	}

	for _, region := range regions {
		cwcClient := cwc.CwClient(region)
		groupNames, _ := check.ListLogGroups(cwcClient)

		for _, groupName := range groupNames {
			retention, err := check.CheckLogGroupsRetentionPolicy(cwcClient, groupName)
			if retention != 30 {
				change.ChangeLogGroupsRetentionPolicy(cwcClient, groupName)
			}
			log.WithError(err).Error("CheckLogGroupsRetentionPolicy returned:")

			// 	logStreamNames, err := check.ListLogStreams(c, groupName)
			// 	log.WithError(err).Error("ListLogStreams returned:")
			// 	fmt.Println(logStreamNames)
		}
	}
}
