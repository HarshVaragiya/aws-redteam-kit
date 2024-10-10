package libebs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/sirupsen/logrus"
)

func (comp *ComputeKeySwitcher) ListAllEc2Instances(ctx context.Context) ([]string, error) {
	instanceIds := make([]string, 0)
	nextToken := ""
	for {
		input := &ec2.DescribeInstancesInput{}
		if nextToken != "" {
			input.NextToken = aws.String(nextToken)
		}
		out, err := comp.ec2Client.DescribeInstances(ctx, input)
		if err != nil {
			err := fmt.Errorf("error fetching list of EC2 instances. error = %v", err)
			logger.WithFields(logrus.Fields{"region": comp.Region}).Error(err)
			return instanceIds, err
		}

		for _, reservation := range out.Reservations {
			for _, instance := range reservation.Instances {
				instanceIds = append(instanceIds, aws.ToString(instance.InstanceId))
			}
		}

		if out.NextToken == nil {
			logger.WithFields(logrus.Fields{"region": comp.Region}).Infof("located %d instances", len(instanceIds))
			break
		} else {
			nextToken = aws.ToString(out.NextToken)
		}

	}
	return instanceIds, nil
}

func (comp *ComputeKeySwitcher) StopAllEc2Instances(instanceIDs []string) error {
	_, err := comp.ec2Client.StopInstances(context.Background(), &ec2.StopInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		err := fmt.Errorf("error making the stop instances API call. error = %v", err)
		logger.WithFields(logrus.Fields{"region": comp.Region}).Error(err)
		return err
	}
	logger.WithFields(logrus.Fields{"region": comp.Region}).Infof("done shutting down ec2 instances ...")
	return nil
}
