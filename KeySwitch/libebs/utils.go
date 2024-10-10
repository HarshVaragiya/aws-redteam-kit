package libebs

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func GetEbsVolumeIds(volumes []types.Volume) []string {
	volumeIds := make([]string, len(volumes))
	for i := 0; i < len(volumes); i++ {
		volumeIds[i] = aws.ToString(volumes[i].VolumeId)
	}
	return volumeIds
}
