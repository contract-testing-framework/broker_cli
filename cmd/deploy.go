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
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the Signet broker to a new ECS Fargate cluster",
	Long:  `Deploy the Signet broker to a new ECS Fargate cluster`,
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

		fmt.Println(colorGreen + "Deploying" + colorReset + " - deploying the Signet broker to a new ECS Fargate cluster, this will take a few minutes...")

		if err := waitForDeploymentDone(cfClient); err != nil {
			return err
		}

		fmt.Println("\n" + colorGreen + "Deployed Successfully" + colorReset)

		printURLofELB(cfClient, cfg)

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

func printURLofELB(cfClient *cloudformation.Client, cfg aws.Config) error {
	name := stackName
	dsrInput := &cloudformation.DescribeStackResourcesInput{StackName: &name}

	dsrOutput, err := cfClient.DescribeStackResources(context.TODO(), dsrInput)
	if err != nil {
		fmt.Println("Cannot display the URL of the ELB in front of the Signet broker cluster - check AWS console for the ELB's URL")
		return nil
	}

	var elbArn string
	for _, resource := range dsrOutput.StackResources {
		if *resource.LogicalResourceId == "LoadBalancer" {
			elbArn = *resource.PhysicalResourceId
		}
	}

	elbClient := elasticloadbalancingv2.NewFromConfig(cfg)
	dlbInput := elasticloadbalancingv2.DescribeLoadBalancersInput{LoadBalancerArns: []string{elbArn}}
	dlbOutput, err := elbClient.DescribeLoadBalancers(context.TODO(), &dlbInput)
	if err != nil {
		fmt.Println("Cannot display the URL of the ELB in front of the Signet broker cluster - check AWS console for the ELB's URL")
		return nil
	}

	dnsName := *dlbOutput.LoadBalancers[0].DNSName
	fmt.Println("Signet broker is exposed through an Elastic Load Balancer at " + colorBlue + "http://" + dnsName + colorReset)
	fmt.Println("\nAdd a TLS certificate to the ELB to enable HTTPS")

	return nil
}
	
func init() {
	RootCmd.AddCommand(deployCmd)
}