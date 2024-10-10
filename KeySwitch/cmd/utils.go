package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func sanityCheckVariables() {
	if region == "" {
		logger.Fatalf("AWS Region string cannot be empty")
	}
	if kmsKeyArn == "" {
		logger.Fatalf("AWS KMS Key ARN cannot be empty")
	}
	setDebugLevel()
}

func setDebugLevel() {
	if debug {
		logger.Level = logrus.DebugLevel
	}
}

func getUserConfirmation(action string) {
	if I_KNOW_EXACTLY_WHAT_I_AM_DOING {
		logger.Warnf("SKIPPING USER CONFIRMATION for %s", action)
		return
	}
	fmt.Printf("Please confirm if you wish to proceed with '%s' (YES/NO) ? \n > ", action)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if input == "YES\n" {
		logger.Warnf("User confirmation recieved. I sure hope you know what you are doing.")
	} else {
		logger.Fatal("User confirmation declined.")
	}
}
