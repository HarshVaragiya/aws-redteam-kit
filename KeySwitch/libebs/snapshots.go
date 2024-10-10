package libebs

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/HarshVaragiya/keyswitch/libswitch"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ebs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/sirupsen/logrus"
)

func NewComputeKeySwitcher(region string) *ComputeKeySwitcher {
	computeConfig := libswitch.GetDefaultAwsClientConfig().Copy()
	computeConfig.Region = region
	ebsC := ebs.NewFromConfig(computeConfig)
	ec2C := ec2.NewFromConfig(computeConfig)
	logger.Infof("creating computeKeySwitcher instance for region: %s", region)
	return &ComputeKeySwitcher{
		Region:    region,
		ebsClient: ebsC,
		ec2Client: ec2C,
	}
}

func (comp *ComputeKeySwitcher) CreateSnapshotsRequestForVolume(ctx context.Context, volumeIds []string) ([]string, error) {
	snapshotIds := make([]string, 0)
	for _, volId := range volumeIds {
		input := &ec2.CreateSnapshotInput{
			VolumeId: aws.String(volId),
		}
		result, err := comp.ec2Client.CreateSnapshot(ctx, input)
		if err != nil {
			err := fmt.Errorf("couldn't create snapshot for volume %s: %v", volId, err)
			logger.WithFields(logrus.Fields{"region": comp.Region, "volumeId": volId}).Error(err)
			continue
		}
		snapshotId := aws.ToString(result.SnapshotId)
		logger.WithFields(logrus.Fields{"region": comp.Region, "volumeId": volId, "snapshotId": snapshotId}).Debugf("created snapshot request")
		snapshotIds = append(snapshotIds, snapshotId)
	}
	return snapshotIds, nil
}

func (comp *ComputeKeySwitcher) WaitForAvailableSnapshots(ctx context.Context, snapshotIds []string, availableSnapshotIds chan string, wg *sync.WaitGroup) error {
	defer wg.Done()
	waitingSnapshotsCount := len(snapshotIds)
	sentSnapshotIds := make([]string, 0)

	for {
		snapshots, err := comp.GetAllSnapshots(ctx)
		if err != nil {
			return err
		}

		for _, snapshot := range snapshots {
			snapshotId := aws.ToString(snapshot.SnapshotId)
			if slices.Contains(snapshotIds, snapshotId) {
				// this is a snapshot we are interested in ..
				if snapshot.State == types.SnapshotStateCompleted {
					if !slices.Contains(sentSnapshotIds, snapshotId) {
						// new snapshot is available
						availableSnapshotIds <- snapshotId
						sentSnapshotIds = append(sentSnapshotIds, snapshotId)
						logger.WithFields(logrus.Fields{"region": comp.Region}).Debugf("found new snapshot: %s", snapshotId)
					} else {
						logger.WithFields(logrus.Fields{"region": comp.Region, "snapshotId": snapshotId}).Debugf("snapshot progress: %s", aws.ToString(snapshot.Progress))
					}
				}
			}
		}

		if len(sentSnapshotIds) == waitingSnapshotsCount {
			logger.WithFields(logrus.Fields{"region": comp.Region}).Info("all snapshots are available")
			break
		}
		time.Sleep(POLLING_TIME)
	}
	close(availableSnapshotIds)
	return nil

}

func (comp *ComputeKeySwitcher) GetAllSnapshots(ctx context.Context) ([]types.Snapshot, error) {
	snapshots := make([]types.Snapshot, 0)
	nextToken := ""
	for {
		input := &ec2.DescribeSnapshotsInput{
			OwnerIds: []string{"self"},
		}
		if nextToken != "" {
			input.NextToken = &nextToken
		}
		out, err := comp.ec2Client.DescribeSnapshots(ctx, input)
		if err != nil {
			err := fmt.Errorf("error getting snapshot details. error = %v", err)
			logger.WithFields(logrus.Fields{"region": comp.Region}).Error(err)
			return snapshots, err
		}
		snapshots = append(snapshots, out.Snapshots...)
		if out.NextToken == nil {
			logger.WithFields(logrus.Fields{"region": comp.Region}).Debugf("listed all snapshots. count = %d", len(snapshots))
			break
		} else {
			nextToken = aws.ToString(out.NextToken)
		}
		time.Sleep(PAGINATION_SLEEP_TIME)
	}
	return snapshots, nil
}

func (comp *ComputeKeySwitcher) DeleteAllSnapshots(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()
	snapshots, err := comp.GetAllSnapshots(ctx)
	if err != nil {
		logger.WithFields(logrus.Fields{"region": comp.Region}).Errorf("error listing snapshots ...")
		return err
	}
	for _, snapshot := range snapshots {
		snapshotId := aws.ToString(snapshot.SnapshotId)
		_, err := comp.ec2Client.DeleteSnapshot(ctx, &ec2.DeleteSnapshotInput{SnapshotId: aws.String(snapshotId)})
		if err != nil {
			err := fmt.Errorf("error deleting snapshot. error = %v", err)
			logger.WithFields(logrus.Fields{"region": comp.Region, "snapshotId": snapshotId}).Error(err)
			continue
		}
		logger.WithFields(logrus.Fields{"region": comp.Region, "snapshotId": snapshotId}).Infof("deleted snapshot")
		time.Sleep(DELETE_API_RATE_LIMIT_DELAY)
	}
	return nil
}
