package change

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
)

func ChangeLogGroupsRetentionPolicy(cwl cloudwatchlogsiface.CloudWatchLogsAPI, group string) error {

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
