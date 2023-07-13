package cmd

import (
	// "bufio"
	"errors"
	"fmt"
	"context"
	"io/ioutil"
	"log"
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
	Long:  `Deploy the Signet broker to you AWS account on ECS with Fargate
	
	flags:

	-s -—silent                (bool) silence docker's status updates as it provisions AWS infrastructure
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// silent = viper.GetBool("deploy.silent")

		cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
			return errors.New("unable to load SDK config: " + err.Error())
    }
		
		// need to get template from node module...
		templateFile, err := ioutil.ReadFile("cftemplate4.yaml")
		if err != nil {
			return errors.New("unable to load CloudFormation template: " + err.Error())
		}
		template := string(templateFile)
		
		csInput := &cloudformation.CreateStackInput{
			StackName: aws.String(stackName),
			TemplateBody: aws.String(template),
			Capabilities: []types.Capability{"CAPABILITY_IAM"},
		}
		
		cfClient := cloudformation.NewFromConfig(cfg)
		_, err = cfClient.CreateStack(context.TODO(), csInput)
		if err != nil {
			return errors.New("unable to create CloudFormation stack: " + err.Error())
		}

		fmt.Println(colorGreen + "Deploying" + colorReset + " - deploying Signet broker to your AWS cloud using ECS with Fargate, this may take a few minutes...")

		waiter := cloudformation.NewStackCreateCompleteWaiter(cfClient)
		waitDuration, err := time.ParseDuration("15m") // 15 minutes
		if err != nil {
			return err
		}

		name := stackName
		dsInput := &cloudformation.DescribeStacksInput{StackName: &name}
		if err := waiter.Wait(context.TODO(), dsInput, waitDuration); err != nil { // changed this from WaitForOutput before committing - untested
			return errors.New("error while waiting for CloudFormation stack to be created - it is likely that Signet CLI simply timed out. Check your AWS console for the status of the deployment: " + err.Error())
		}

		fmt.Println(colorGreen + "Deployed Successfully")

		return nil
	},
}
	
func init() {
	RootCmd.AddCommand(deploy_sdk)
}