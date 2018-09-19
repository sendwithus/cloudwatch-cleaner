package check

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func CheckLogGroupsRetentionPolicy(cwl *cloudwatchlogs.CloudWatchLogs, groupName string) (int64, error) {

	limit := int64(50)

	r, err := cwl.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: &groupName,
		Limit:              &limit,
	})
	if err != nil {
		log.WithError(err).Error("Could not describe the log group")
		return int64(-1), err
	}

	for _, group := range r.LogGroups {
		if *group.LogGroupName == groupName {
			if group.RetentionInDays == nil {
				return int64(0), nil
			}
			return *group.RetentionInDays, nil
		}
	}
	return int64(-1), errors.New("OH GOD WHY!??")
}
