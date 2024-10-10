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
	xksProxyUriPath     string
	xksProxyUriEndpoint string
	accessKeyId         string
	secretAccessKey     string
	customKeyStoreName  string
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		SetupXks(context.TODO())
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.PersistentFlags().StringVarP(&customKeyStoreName, "xks-name", "n", "aws-xks", "name for the external key store")
	setupCmd.PersistentFlags().StringVarP(&xksProxyUriPath, "uri", "p", "", "HTTP PATH for the XKS Proxy")
	setupCmd.PersistentFlags().StringVarP(&xksProxyUriEndpoint, "endpoint", "e", "", "HTTP Endpoint (domain) for the XKS Proxy")
	setupCmd.PersistentFlags().StringVarP(&accessKeyId, "access-key", "a", "", "Access Key Id")
	setupCmd.PersistentFlags().StringVarP(&secretAccessKey, "secret-access-key", "s", "", "secret access key")
}

func SetupXks(ctx context.Context) {
	kmsClient := kms.NewFromConfig(libswitch.GetDefaultAwsClientConfig())
	logger.Infof("creating AWS XKS ...")
	createCustomKeyStoreInput := &kms.CreateCustomKeyStoreInput{
		CustomKeyStoreName:   aws.String(customKeyStoreName),
		CustomKeyStoreType:   types.CustomKeyStoreTypeExternalKeyStore,
		XksProxyConnectivity: types.XksProxyConnectivityTypePublicEndpoint,
		XksProxyUriPath:      aws.String(xksProxyUriPath),
		XksProxyUriEndpoint:  aws.String(xksProxyUriEndpoint),
		XksProxyAuthenticationCredential: &types.XksProxyAuthenticationCredentialType{
			AccessKeyId:        aws.String(accessKeyId),
			RawSecretAccessKey: aws.String(secretAccessKey),
		},
	}
	createCustomerKeyStoreOutput, err := kmsClient.CreateCustomKeyStore(ctx, createCustomKeyStoreInput)
	if err != nil {
		logger.Fatalf("error creating XKS. error = %v", err)
	}
	keyStoreId := aws.ToString(createCustomerKeyStoreOutput.CustomKeyStoreId)
	logger.Infof("created custom key store with Id: %s", keyStoreId)

	logger.Infof("attempting to connect to XKS : %s", keyStoreId)

	_, err = kmsClient.ConnectCustomKeyStore(ctx, &kms.ConnectCustomKeyStoreInput{CustomKeyStoreId: &keyStoreId})
	if err != nil {
		logger.Fatalf("error connecting to XKS. error = %v", err)
	}
	logger.Infof("done. you can create KMS Keys with this CKS now")
}
