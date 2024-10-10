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
	"os"

	"github.com/spf13/cobra"
)

var (
	debug                          = true
	region                         = ""
	kmsKeyArn                      = ""
	I_KNOW_EXACTLY_WHAT_I_AM_DOING = false
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "keyswitch",
	Short: "quickly switch the underlying AWS KMS encryption key for AWS resources",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "v", true, "enable debug logs")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "AWS Region to run script on")
	rootCmd.PersistentFlags().StringVarP(&kmsKeyArn, "key", "k", "", "KMS Key ID to use for encryption")
	rootCmd.PersistentFlags().BoolVar(&I_KNOW_EXACTLY_WHAT_I_AM_DOING, "i-know-what-i-am-doing-and-no-one-other-than-me-is-responsible-for-any-data-loss", false, "this flag skips all user confirmation asks. NEVER USE THIS")
}
