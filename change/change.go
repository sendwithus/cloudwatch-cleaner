package setup

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func ChangeLogGroupsRetentionPolicy(cwl *cloudwatchlogs.CloudWatchLogs, group string) error {

	retention := int64(30)

	_, err := cwl.PutRetentionPolicy(&cloudwatchlogs.PutRetentionPolicyInput{
		LogGroupName:    &group,
		RetentionInDays: &retention,
	})
	if err != nil {
		log.Error("Could not setup the retention policy on the log group: %v", err)
	}

	return err
}
