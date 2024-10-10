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

	"github.com/HarshVaragiya/keyswitch/libswitch"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/spf13/cobra"
)

var (
	description      string
	xksKeyId         string
	customKeyStoreId string
)

// kmsCmd represents the kms command
var kmsCmd = &cobra.Command{
	Use:   "kms",
	Short: "Setup KMS Symmetric keys for Ransomware Simultation",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		SetupKmsKeyWithCks(context.TODO())
	},
}

func init() {
	rootCmd.AddCommand(kmsCmd)
	kmsCmd.PersistentFlags().StringVarP(&description, "description", "d", "AWS XKS Key", "description for the KMS Key")
	kmsCmd.PersistentFlags().StringVar(&xksKeyId, "xks-key-id", "thekey", "AWS XKS Key Id (present in dockerfile and settings.toml)")
	kmsCmd.PersistentFlags().StringVar(&customKeyStoreId, "cks", "aws-xks", "the custom key store id for the generated XKS")
}

func SetupKmsKeyWithCks(ctx context.Context) {
	kmsClient := kms.NewFromConfig(libswitch.GetDefaultAwsClientConfig())
	logger.Infof("attempting to create a symmetric key now")
	createKeyInput := &kms.CreateKeyInput{
		CustomKeyStoreId: aws.String(customKeyStoreId),
		MultiRegion:      aws.Bool(false),
		Description:      aws.String(description),
		KeyUsage:         types.KeyUsageTypeEncryptDecrypt,
		KeySpec:          types.KeySpecSymmetricDefault,
		Origin:           types.OriginTypeExternalKeyStore,
		XksKeyId:         aws.String(xksKeyId),
	}
	createKeyOutput, err := kmsClient.CreateKey(ctx, createKeyInput)
	if err != nil {
		logger.Fatalf("error creating KMS key with XKS. error = %v", err)
	}
	kmsKeyArn := aws.ToString(createKeyOutput.KeyMetadata.Arn)
	logger.Infof("created KMS key Backed by XKS: %s", kmsKeyArn)
	logger.Infof("done. you can now use this key to hold data for ransom!")
}
