package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/sendwithus/cloudwatch-cleaner/change"
	"github.com/sendwithus/cloudwatch-cleaner/config"
)

func main() {
	client := change.New()
	err := run(client)
	if err != nil {
		log.WithError(err).Error("Run returned:")
	}
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
			for _, whiteList := range config.WhiteList {
				if whiteList == group {
					log.Info("Group in white list, skipping.")
				} else {
					retention, err := client.GetRetentionPolicy(group)
					if err != nil {
						log.WithError(err).Error("GetRetentionPolicy returned:")
					}
					retentionDays, err := client.ConvertToInt64(os.Getenv("RETENTION_DAYS"))
					if err != nil {
						log.WithError(err).Error("ConvertToInt64 returned:")
					}
					if retention != retentionDays {
						err := client.SetRetentionPolicy(retentionDays, group)
						if err != nil {
							log.WithError(err).Error("SetRetentionPolicy returned:")
						}
					}
				}
			}
		}
	}
	return nil
}
