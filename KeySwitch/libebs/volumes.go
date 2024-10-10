package libebs

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/sirupsen/logrus"
)

func (comp *ComputeKeySwitcher) GetAllEbsVolumes(ctx context.Context) ([]types.Volume, error) {
	input := &ec2.DescribeVolumesInput{}
	nextToken := ""
	volumes := make([]types.Volume, 0)
	for {
		input = &ec2.DescribeVolumesInput{}
		if nextToken != "" {
			input.NextToken = aws.String(nextToken)
		}
		volumesOutput, err := comp.ec2Client.DescribeVolumes(ctx, input)
		if err != nil {
			err := fmt.Errorf("error listing ebs volumes. error = %v", err)
			logger.WithFields(logrus.Fields{"region": comp.Region}).Error(err)
			return volumes, err
		}
		for _, volume := range volumesOutput.Volumes {
			volumeId := aws.ToString(volume.VolumeId)
			logger.WithFields(logrus.Fields{"region": comp.Region, "volumeId": volumeId}).Debugf("discovered new volume")

		}
		volumes = append(volumes, volumesOutput.Volumes...)
		if volumesOutput.NextToken == nil {
			logger.WithFields(logrus.Fields{"region": comp.Region}).Infof("found all volumes. count = %d", len(volumes))
			break
		}
		nextToken = aws.ToString(volumesOutput.NextToken)
	}
	return volumes, nil
}

func (comp *ComputeKeySwitcher) CreateEncryptedVolumeFromSnapshot(ctx context.Context, snapshotIds chan string, kmsKeyArn string, volumeIds chan string, wg *sync.WaitGroup) error {
	defer wg.Done()
	for snapshotId := range snapshotIds {
		input := &ec2.CreateVolumeInput{
			SnapshotId:       aws.String(snapshotId),
			AvailabilityZone: aws.String(fmt.Sprintf("%sa", comp.Region)),
			Encrypted:        aws.Bool(true),
			KmsKeyId:         aws.String(kmsKeyArn),
		}
		out, err := comp.ec2Client.CreateVolume(ctx, input)
		if err != nil {
			err := fmt.Errorf("error creating encrypted snapshot. error = %v", err)
			logger.WithFields(logrus.Fields{"region": comp.Region, "snapshotId": snapshotId}).Error(err)
			continue
		}
		volumeId := aws.ToString(out.VolumeId)
		volumeIds <- volumeId
		logger.WithFields(logrus.Fields{"region": comp.Region, "snapshotId": snapshotId, "volumeId": volumeId}).Infof("requested creation of encrypted volume")
	}
	return nil
}

func (comp *ComputeKeySwitcher) DeleteVolumes(ctx context.Context, volumes []string, wg *sync.WaitGroup) error {
	defer wg.Done()
	for _, volumeId := range volumes {
		_, err := comp.ec2Client.DeleteVolume(
			ctx,
			&ec2.DeleteVolumeInput{VolumeId: aws.String(volumeId)})
		if err != nil {
			err := fmt.Errorf("error deleting volume. error = %v", err)
			logger.WithFields(logrus.Fields{"region": comp.Region, "volumeId": volumeId}).Error(err)
			continue
		}
		logger.WithFields(logrus.Fields{"region": comp.Region, "volumeId": volumeId}).Infof("deleted volume")
	}
	return nil
}

func (comp *ComputeKeySwitcher) DetachAllEbsVolumes(ctx context.Context) error {
	resp, err := comp.GetAllEbsVolumes(ctx)
	if err != nil {
		err := fmt.Errorf("error describing EBS volumes: %v", err)
		logger.WithFields(logrus.Fields{"region": comp.Region}).Error(err)
		return err
	}
	for _, volume := range resp {
		if volume.State == types.VolumeStateInUse {
			volumeId := aws.ToString(volume.VolumeId)
			_, err := comp.ec2Client.DetachVolume(context.Background(), &ec2.DetachVolumeInput{
				VolumeId: aws.String(volumeId),
			})
			if err != nil {
				err := fmt.Errorf("error detaching volume. error = %v", err)
				logger.WithFields(logrus.Fields{"region": comp.Region, "volumeId": volumeId}).Error(err)
				continue
			}
			logger.WithFields(logrus.Fields{"region": comp.Region, "volumeId": volumeId}).Info("detached ebs volume")
		}
	}
	return nil
}
