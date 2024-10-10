/*
Copyright © 2024 Harsh Varagiya

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"sync"

	"github.com/HarshVaragiya/keyswitch/libebs"
	"github.com/HarshVaragiya/keyswitch/liblogger"
	"github.com/spf13/cobra"
)

var (
	logger = liblogger.Log
)

// ebsCmd represents the ebs command
var ebsCmd = &cobra.Command{
	Use:   "ebs",
	Short: "re-encrypt ebs volumes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		sanityCheckVariables()
		PerformEbsReencryption()
	},
}

func init() {
	rootCmd.AddCommand(ebsCmd)
}

func PerformEbsReencryption() {
	ctx := context.Background()
	cks := libebs.NewComputeKeySwitcher(region)

	// just list all volumes - no harm done.
	logger.Infof("fetching list of all EBS volumes ...")
	existingVolumes, err := cks.GetAllEbsVolumes(ctx)
	if err != nil {
		logger.Fatalf("error fetching EBS volumes. error = %v", err)
	}
	volumeIds := libebs.GetEbsVolumeIds(existingVolumes)

	// now create snapshots for the existing volumes - this incurs additional cost $$$
	getUserConfirmation("create ALL volume snapshots")
	logger.Infof("creating snapshots for the volumes ...")
	snapshotIds, err := cks.CreateSnapshotsRequestForVolume(ctx, volumeIds)
	if err != nil {
		logger.Fatalf("error fetching EBS volumes. error = %v", err)
	}

	encryptWg := &sync.WaitGroup{}
	availableSnapshotIds := make(chan string, 100)
	encryptedVolumeIds := make(chan string, 100)
	// now we wait for the snapshots to become available
	logger.Infof("waiting for snapshots to become available ...")
	encryptWg.Add(1)
	go cks.WaitForAvailableSnapshots(ctx, snapshotIds, availableSnapshotIds, encryptWg)

	// creating a new EBS volume using the snapshot
	getUserConfirmation("encrypt snapshots with given KMS key")
	logger.Infof("encrypting snapshots and creating new volumes ...")
	encryptWg.Add(1)
	go cks.CreateEncryptedVolumeFromSnapshot(ctx, availableSnapshotIds, kmsKeyArn, encryptedVolumeIds, encryptWg)

	logger.Infof("waiting for encryption process to finish ...")
	encryptWg.Wait()
	logger.Infof("encryption process has finished.")

	deleteWg := &sync.WaitGroup{}
	// now we can go ahead and delete all the snapshots in the region
	getUserConfirmation("delete ALL snapshots")
	logger.Infof("deleting snapshots ...")
	deleteWg.Add(1)
	if err := cks.DeleteAllSnapshots(ctx, deleteWg); err != nil {
		logger.Fatalf("error deleting snapshots. error = %v", err)
	}

	logger.Infof("waiting for delete processes to finish ...")
	deleteWg.Wait()
	logger.Infof("done")
}
