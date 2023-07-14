package cmd

import (
	"errors"
	"fmt"
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

var undeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Tear down the Signet broker deployment on AWS ECS",
	Long:  `Tear down the Signet broker deployment on AWS ECS`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dsInput := &cloudformation.DeleteStackInput{StackName: aws.String(stackName)}

		cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
			return errors.New("unable to load SDK config - have you configured your aws cli with `aws configure`? " + err.Error())
    }
		
		cfClient := cloudformation.NewFromConfig(cfg)
		_, err = cfClient.DeleteStack(context.TODO(), dsInput)
		if err != nil {
			return errors.New("unable to delete CloudFormation stack: " + err.Error())
		}

		fmt.Println(colorGreen + "Undeploying" + colorReset + " - tearing down the Signet broker ECS Cluster, this will take a few minutes...")

		if err := waitForUndeploymentDone(cfClient); err != nil {
			return err
		}

		fmt.Println(colorGreen + "Undeployed Successfully")

		return nil
	},
}

func waitForUndeploymentDone(cfClient *cloudformation.Client) error {
	waiter := cloudformation.NewStackDeleteCompleteWaiter(cfClient)
	waitDuration, err := time.ParseDuration("15m")
	if err != nil {
		return err
	}

	name := stackName
	dsInput := &cloudformation.DescribeStacksInput{StackName: &name}
	if err := waiter.Wait(context.TODO(), dsInput, waitDuration); err != nil { // changed this from WaitForOutput before committing - untested
		return errors.New("error while waiting for CloudFormation stack to be deleted - it is likely that Signet CLI simply timed out. Check your AWS console for the status of the teardown operations: " + err.Error())
	}

	return nil
}
	
func init() {
	RootCmd.AddCommand(undeployCmd)
}