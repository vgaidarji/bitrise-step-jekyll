package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

// Config ...
type Config struct {
	WorkDir string `env:"work_dir,dir"`
}

func initializeConfig() (Config, error) {
	var config Config
	if err := stepconf.Parse(&config); err != nil {
		fmt.Printf("Error parsing step config: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(config)

	workDir, err := pathutil.AbsPath(config.WorkDir)
	if err != nil {
		fmt.Printf("Error normalizing workdir path: %s", err)
	}

	exists, err := pathutil.IsDirExists(workDir)
	if err != nil {
		fmt.Printf("Error validating workdir `%s`: %s", workDir, err)
	}
	if !exists {
		fmt.Printf("Specified path `%s` does not exist", workDir)
	}
	return config, nil
}

func main() {
	config, _ := initializeConfig()

	fmt.Println()
	log.Infof("Building Jekyll project")

	cmdBundleInstall := exec.Command("bundle", "install")
	cmdBundleInstall.Dir = config.WorkDir
	cmdBundleInstallResult, err := cmdBundleInstall.Output()
	if err != nil {
		fmt.Printf("Failed to install dependencies using bundle, error: %#v | output: %s", err, cmdBundleInstallResult)
		os.Exit(1)
	}
	fmt.Printf(string(cmdBundleInstallResult))

	cmdJekyllBuild := exec.Command("bundle", "exec", "jekyll", "build")
	cmdJekyllBuild.Dir = config.WorkDir
	cmdJekyllBuildResult, err := cmdJekyllBuild.Output()
	if err != nil {
		fmt.Printf("Failed to build Jekyll project, error: %#v | output: %s", err, cmdJekyllBuildResult)
		os.Exit(1)
	}
	fmt.Printf(string(cmdJekyllBuildResult))

	siteDir := config.WorkDir + "/_site"
	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "JEKYLL_GENERATED_SITE_FOLDER", "--value", siteDir).CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
		os.Exit(1)
	}

	log.Donef("  Done")
}
