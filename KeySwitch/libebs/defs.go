package libebs

import (
	"time"

	"github.com/HarshVaragiya/keyswitch/liblogger"
	"github.com/aws/aws-sdk-go-v2/service/ebs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

var (
	logger                      = liblogger.Log
	PAGINATION_SLEEP_TIME       = time.Second * 1
	POLLING_TIME                = time.Second * 30
	DELETE_API_RATE_LIMIT_DELAY = time.Second * 2
)

type ComputeKeySwitcher struct {
	Region    string
	ebsClient *ebs.Client
	ec2Client *ec2.Client
}
