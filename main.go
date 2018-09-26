package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/techdroplabs/cloudwatch-cleaner/change"
)

func main() {
	client := change.New()
	err := run(client)
	log.WithError(err).Error("Run returned:")
}

func run(client change.Client) error {
	client.SetRegion("us-west-2") // Need to call this first to init clients.

	regions, err := client.ListRegions()
	if err != nil {
		return err
	}

	for _, region := range regions {
		client.SetRegion(region)
		groups, _ := client.ListGroups()

		for _, group := range groups {
			retention, err := client.GetRetentionPolicy(group)
			log.WithError(err).Error("GetRetentionPolicy returned:")
			if retention != 30 {
				err := client.SetRetentionPolicy(group)
				log.WithError(err).Error("SetRetentionPolicy returned:")
			}
		}
	}
	return nil
}
