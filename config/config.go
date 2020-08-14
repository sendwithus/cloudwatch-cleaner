package config

import (
	"strings"

	"github.com/sendwithus/cloudwatch-cleaner/utils"
)

var (
	// BlockList is a list of logrgoup that will not get cleanup
	BlockList = strings.Split(utils.GetEnvString("BLOCKLIST", ""), ",")
)
