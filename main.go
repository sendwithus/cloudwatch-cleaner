package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/techdroplabs/cloudwatch-cleaner/check"
)

func main() {
	client := check.New()
	run(client)
}

func run(client check.Client) error {
	client.SetRegion("us-west-2") // Need to call this first to init clients.

	regions, err := client.ListRegions()
	if err != nil {
		return err
	}

	for _, region := range regions {
		client.SetRegion(region)
		groupNames, _ := client.ListGroups()

		for _, groupName := range groupNames {
			retention, err := client.GetRetentionPolicy(groupName)
			if retention != 30 {
				client.SetRetentionPolicy(groupName)
			}
			log.WithError(err).Error("CheckLogGroupsRetentionPolicy returned:")

			// 	logStreamNames, err := check.ListLogStreams(c, groupName)
			// 	log.WithError(err).Error("ListLogStreams returned:")
			// 	fmt.Println(logStreamNames)
		}
	}
	return nil
}
