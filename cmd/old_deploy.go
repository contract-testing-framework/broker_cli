package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ecsContext string
var silent bool
var destroy bool

var oldDeployCmd = &cobra.Command{
	Use:   "oldDeploy",
	Short: "Deploy the Signet broker to you AWS account on ECS with Fargate",
	Long:  `Deploy the Signet broker to you AWS account on ECS with Fargate
	
	flags:

	-c --ecs-context           the name of the local docker ecs context with AWS credentials
	
	-s -—silent                (bool) silence docker's status updates as it provisions AWS infrastructure

	-d --destroy               (bool) causes the Signet broker to be torn down from AWS instead of deployed
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ecsContext = viper.GetString("deploy.ecs-context")
		silent = viper.GetBool("deploy.silent")

		if err := checkEcsContextExists(ecsContext); err != nil {
			return err
		}

		currentDockerContext, err := cacheCurrentDockerContext()
		if err != nil {
			return err
		}
	
		if exec.Command("docker", "context", "use", ecsContext).Run(); err != nil {
			return err
		}

		fmt.Println("Info: deploying Signet broker to your AWS cloud using ECS with Fargate, this may take a few minutes...")
		if err = deploySignet(silent, destroy); err != nil {
			return err
		}

		fmt.Println(colorGreen + "Signet deployed" + colorReset + " - the Signet broker has been deployed to your AWS cloud")

		// it is okay if this command is not successful
		_ = exec.Command("docker", "context", "use", currentDockerContext).Run()

		// if err = exec.Command("docker", "context", "use", currentDockerContext).Run(); err != nil {
		// 	fmt.Println("Errored on line 56")
		// 	return err
		// }

		return nil
	},
}

func checkEcsContextExists(context string) error {
	checkContextCmd := exec.Command("docker", "context", "list")
	stdoutStderr, err := checkContextCmd.CombinedOutput()
	checkContextOutput := string(stdoutStderr)
	
	if err != nil && len(checkContextOutput) == 0 {
		return errors.New("failed to check for existing docker context. Is docker running? - " + err.Error())
	}
	
	found, err := regexp.MatchString(context, checkContextOutput)
	if err != nil {
		return err
	}
	
	if !found {
		return errors.New("docker context not found for --ecs-context. Please run 'docker context create ecs <context-name>' and follow the prompts to configure your AWS credentials")
	}
	
	return nil
}

func cacheCurrentDockerContext() (string, error) {
	getContextCmd := exec.Command("docker", "context", "show")
	stdoutStderr, err := getContextCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	if len(stdoutStderr) == 0 {
		return "", errors.New("could not get current docker context")
	}

	return string(stdoutStderr[:len(stdoutStderr) - 1]), nil
}

func deploySignet(silent, destory bool) error {
	signetRoot, err := getNpmPkgRoot()
	if err != nil {
		return err
	}

	upDown := "up"
	if destory {
		upDown = "down"
	}
	
	signetUpCmd := exec.Command("docker", "compose", "--project-name=signetbroker", "--file=" + signetRoot + "/docker-compose.yml", upDown)
	stderr, err := signetUpCmd.StderrPipe()
	if err != nil {
		return err
	}

	if !silent {
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		}()
	}
	
	err = signetUpCmd.Run()
	if err != nil {
		return err
	}
	
	return nil
}
	
func init() {
	RootCmd.AddCommand(oldDeployCmd)

	oldDeployCmd.Flags().StringVarP(&ecsContext, "ecs-context", "c", "", "the name of the local docker ecs context with AWS credentials")
	oldDeployCmd.Flags().BoolVarP(&silent, "silent", "s", false, "silence docker's status updates as it provisions AWS infrastructure")
	oldDeployCmd.Flags().BoolVarP(&destroy, "destroy", "d", false, "causes the Signet broker to be torn down from AWS instead of deployed")

	viper.BindPFlag("deploy.ecs-context", oldDeployCmd.Flags().Lookup("ecs-context"))
	viper.BindPFlag("deploy.silent", oldDeployCmd.Flags().Lookup("silent"))
}