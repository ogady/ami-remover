package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {

	specified_date := flag.String("date", "", "YYYY-MM-DD. Deletes AMIs older than the specified date")
	isDryRun := flag.Bool("dry_run", false, "")

	flag.Parse()
	args := flag.Args()

	log.Printf("dry run mode: %v", *isDryRun)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	client := ec2.NewFromConfig(cfg)

	filter := types.Filter{
		Name:   aws.String("name"),
		Values: args,
	}
	input := ec2.DescribeImagesInput{
		Filters: []types.Filter{filter},
	}

	result, err := client.DescribeImages(context.TODO(), &input)
	if err != nil {
		fmt.Println("Got an error retrieving information about your Amazon EC2 instances:")
		fmt.Println(err)
		return
	}

	sd, err := time.Parse("20060102", *specified_date)
	if err != nil {
		log.Println("err: failed to parse specified date. Please specify in this format 'YYYY-MM-DD' ")
		fmt.Println(err)
		return
	}

	for _, r := range result.Images {
		rtime, err := time.Parse(time.RFC3339, *r.CreationDate)
		if err != nil {
			log.Println("err: failed to parse ami create date")
			fmt.Println(err)
			return
		}

		if rtime.Before(sd) {
			fmt.Println("AMI Name: " + *r.Name)
			fmt.Println("CreateDate:" + *r.CreationDate)
			delInput := ec2.DeregisterImageInput{
				ImageId: r.ImageId,
				DryRun:  isDryRun,
			}

			_, err = client.DeregisterImage(context.TODO(), &delInput)
			if err != nil {
				log.Println("err: failed to deregister ami")
				fmt.Println(err)
				return
			}
		}
	}
}
