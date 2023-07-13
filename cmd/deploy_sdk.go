package cmd

import (
	"errors"
	"fmt"
	"context"
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

const stackName = "signetbroker"

var deploy_sdk = &cobra.Command{
	Use:   "deploy_sdk",
	Short: "Deploy the Signet broker to you AWS account on ECS with Fargate",
	Long:  `Deploy the Signet broker to you AWS account on ECS with Fargate`,
	RunE: func(cmd *cobra.Command, args []string) error {
		template, err := getCloudFormationTemplate()
		if err != nil {
			return err
		}
		
		csInput := &cloudformation.CreateStackInput{
			StackName: aws.String(stackName),
			TemplateBody: aws.String(template),
			Capabilities: []types.Capability{"CAPABILITY_IAM"},
		}

		cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
			return errors.New("unable to load SDK config - have you configured your aws cli with `aws configure`? " + err.Error())
    }
		
		cfClient := cloudformation.NewFromConfig(cfg)
		_, err = cfClient.CreateStack(context.TODO(), csInput)
		if err != nil {
			return errors.New("unable to create CloudFormation stack: " + err.Error())
		}

		fmt.Println(colorGreen + "Deploying" + colorReset + " - deploying Signet broker to your AWS cloud using ECS with Fargate, this will take a few minutes...")

		if err := waitForDeploymentDone(cfClient); err != nil {
			return err
		}

		fmt.Println(colorGreen + "Deployed Successfully")

		return nil
	},
}

func getCloudFormationTemplate() (string, error) {
	signetRoot, err := getNpmPkgRoot()
	if err != nil {
		return "", errors.New("unable to find signet-cli global npm package: " + err.Error())
	}
	templatePath := signetRoot + "/cftemplate.yaml"
	
	templateFile, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", errors.New("unable to load CloudFormation template: " + err.Error())
	}
	
	return string(templateFile), nil
}

func waitForDeploymentDone(cfClient *cloudformation.Client) error {
	waiter := cloudformation.NewStackCreateCompleteWaiter(cfClient)
	waitDuration, err := time.ParseDuration("15m")
	if err != nil {
		return err
	}

	name := stackName
	dsInput := &cloudformation.DescribeStacksInput{StackName: &name}
	if err := waiter.Wait(context.TODO(), dsInput, waitDuration); err != nil { // changed this from WaitForOutput before committing - untested
		return errors.New("error while waiting for CloudFormation stack to be created - it is likely that Signet CLI simply timed out. Check your AWS console for the status of the deployment: " + err.Error())
	}

	return nil
}
	
func init() {
	RootCmd.AddCommand(deploy_sdk)
}