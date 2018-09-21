package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/techdroplabs/cloudwatch-cleaner/change"
	"github.com/techdroplabs/cloudwatch-cleaner/check"
	"github.com/techdroplabs/cloudwatch-cleaner/client"
)

func main() {

	regions, err := check.ListAllAwsRegions()
	if err != nil {
		log.WithError(err).Error("could not list the aws regions")
	}

	for _, region := range regions {
		c := client.CwClient(region)
		groupNames, _ := check.ListLogGroups(c)

		for _, groupName := range groupNames {
			retention, err := check.CheckLogGroupsRetentionPolicy(c, groupName)
			if retention != 30 {
				setup.ChangeLogGroupsRetentionPolicy(c, groupName)
			}
			log.WithError(err).Error("CheckLogGroupsRetentionPolicy returned:")

			// logStreamNames, err := check.ListLogStreams(c, groupName)
			// log.WithError(err).Error("ListLogStreams returned:")
			// fmt.Println(logStreamNames)
		}
	}
}
