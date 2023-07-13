package cmd

import (
	// "bufio"
	// "errors"
	"fmt"
	// "os/exec"
	// "regexp"
	"context"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	// "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

// var silent bool
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

		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }
		
		templateFile, err := ioutil.ReadFile("cftemplate4.yaml")
		if err != nil {
			log.Fatalf("unable to load template file, %v", err)
		}
		template := string(templateFile)
		
		csInput := &cloudformation.CreateStackInput{
			StackName: aws.String(stackName),
			TemplateBody: aws.String(template),
			Capabilities: []types.Capability{"CAPABILITY_IAM"},
		}
		
		cfClient := cloudformation.NewFromConfig(cfg)
		output, err := cfClient.CreateStack(context.TODO(), csInput)
		if err != nil {
			log.Fatalf("failed to create stack, %v", err)
		}

		fmt.Println("Output from cf.CreateStack()")
		fmt.Println(output.StackId)
		fmt.Println(output.ResultMetadata)

		return nil
	},
}
	
func init() {
	RootCmd.AddCommand(deploy_sdk)

	// deploy_sdk.Flags().BoolVarP(&silent, "silent", "s", false, "silence docker's status updates as it provisions AWS infrastructure")

	// viper.BindPFlag("deploy.silent", deploy_sdk.Flags().Lookup("silent"))
}