package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
	// utils "github.com/contract-testing-framework/broker_cli/utils"
)

const contextName = "signetecscontext"

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the Signet broker to you AWS account on ECS with Fargate",
	Long:  `Deploy the Signet broker to you AWS account on ECS with Fargate`,
	RunE: func(cmd *cobra.Command, args []string) error {

		checkContextCmd := exec.Command("docker", "context", "list")
		stdoutStderr, err := checkContextCmd.CombinedOutput()
		checkContextOutput := string(stdoutStderr)

		if err != nil && len(checkContextOutput) == 0 {
			return errors.New("failed to check for existing docker context. Is docker running? - " + err.Error())
		}

		found, err := regexp.MatchString(contextName, checkContextOutput)
		if err != nil {
			return err
		}

		if !found {
			fmt.Println("Docker context not found. Please run 'docker context create ecs signetecscontext' and setup your AWS credentials")
			return nil
		}

		ctx := context.Background()

		changeCtxCmd := exec.Command("docker", "context", "use", contextName)
		err = changeCtxCmd.Start()
		if err != nil {
			return err
		}
		err = changeCtxCmd.Wait()
		if err != nil {
			return err
		}

		err = RunCmd(ctx, "docker compose -p signet-broker -f signet-broker/docker-compose.yml up")
		if err != nil {
			fmt.Println("THIS ONE")
			return err
		}

		if err != nil {
			log.Fatal(err)
		}

		log.Println("done")

		// var wg sync.WaitGroup
		// wg.Add(1)

		// ecsUpCmd := exec.Command("docker", "compose", "up", "--context", contextName)
		// stdout, _ := ecsUpCmd.StdoutPipe()

		// scanner := bufio.NewScanner(stdout)
		// go func() {
		// 	for scanner.Scan() {
		// 		log.Printf("out: %s", scanner.Text())
		// 	}
		// 	wg.Done()
		// }()
		// // https://gobyexample.com/waitgroups
		// err = ecsUpCmd.Start()
		// if err != nil {
		// 	return errors.New("docker failed to deploy Signet broker to ECS" + err.Error())
		// }

		// wg.Wait()

		// // scanner := bufio.NewScanner(stdout)
		// // scanner.Split(bufio.ScanWords)
		// // for scanner.Scan() {
		// //     m := scanner.Text()
		// //     fmt.Println(m)
		// // }

		// ecsUpCmd.Wait()

		fmt.Println("HEY, SIGNET IS DONE DEPLOYING")

		return nil
	},
}

// https://gist.github.com/hivefans/ffeaf3964924c943dd7ed83b406bbdea

func init() {
	RootCmd.AddCommand(deployCmd)
}

/*
? check if there is a myecscontext docker context
if not, run docker context create
	after the user adds creds, return control to signet cli
run docker compose up

*/

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"io"
// 	"os"
// 	"os/exec"
// 	"strings"
// )

// func execPrint() {
// 	cmdName := "ping 127.0.0.1"
// 	cmdArgs := strings.Fields(cmdName)

// 	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
// 	stdout, _ := cmd.StdoutPipe()
// 	cmd.Start()
// 	oneByte := make([]byte, 100)
// 	num := 1
// 	for {
// 		_, err := stdout.Read(oneByte)
// 		if err != nil {
// 			fmt.Printf(err.Error())
// 			break
// 		}
// 		r := bufio.NewReader(stdout)
// 		line, _, _ := r.ReadLine()
// 		fmt.Println(string(line))
// 		num = num + 1
// 		if num > 3 {
// 			os.Exit(0)
// 		}
// 	}

// 	cmd.Wait()
// }

func RunCmd(ctx context.Context, cmdstr string) error {
	args := strings.Fields(cmdstr)
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			log.Printf("out: %s", scanner.Text())
		}
		wg.Done()
	}()

	if err = cmd.Start(); err != nil {
		return err
	}

	wg.Wait()

	return cmd.Wait()
}

// func main() {
//   cmd := exec.Command("ping", "127.0.0.1")
//   stdout, err := cmd.StdoutPipe()
//   if err != nil {
//     log.Fatal(err)
//   }
//   cmd.Start()

//   buf := bufio.NewReader(stdout) // Notice that this is not in a loop
//   num := 1
//   for {
//     line, _, _ := buf.ReadLine()
//     if num > 3 {
//       os.Exit(0)
//     }
//     num += 1
//     fmt.Println(string(line))
//   }
// }
