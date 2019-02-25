package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bitrise-io/go-utils/log"
)

func main() {
	fmt.Println()
	log.Infof("Building Jekyll project")

	cmdBundleInstallResult, err := exec.Command("bundle", "install").CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to install dependencies using bundle, error: %#v | output: %s", err, cmdBundleInstallResult)
		os.Exit(1)
	}
	fmt.Printf(string(cmdBundleInstallResult))

	cmdJekyllBuildResult, err := exec.Command("bundle", "exec", "jekyll", "build").CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to build Jekyll project, error: %#v | output: %s", err, cmdJekyllBuildResult)
		os.Exit(1)
	}
	fmt.Printf(string(cmdJekyllBuildResult))

	//
	// --- Step Outputs: Export Environment Variables for other Steps:
	// You can export Environment Variables for other Steps with
	//  envman, which is automatically installed by `bitrise setup`.
	// A very simple example:
	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "JEKYLL_GENERATED_SITE_FOLDER", "--value", "the value you want to share").CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
		os.Exit(1)
	}
	// You can find more usage examples on envman's GitHub page
	//  at: https://github.com/bitrise-io/envman

	log.Donef("  Done")
}
