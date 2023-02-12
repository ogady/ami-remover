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
	ctx := context.Background()
	specified_date := flag.String("date", "", "YYYY-MM-DD. Deletes AMIs older than the specified date")
	isDryRun := flag.Bool("dry_run", false, "")

	flag.Parse()
	args := flag.Args()

	log.Printf("dry run mode: %v", *isDryRun)

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
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

	result, err := client.DescribeImages(ctx, &input)
	if err != nil {
		log.Fatalf("Got an error retrieving information about your Amazon EC2 instances. err: %v", err)
	}

	sd, err := time.Parse("20060102", *specified_date)
	if err != nil {
		log.Fatalf("failed to parse specified date. Please specify in this format 'YYYY-MM-DD'. err: %v", err)
	}

	for _, r := range result.Images {
		rtime, err := time.Parse(time.RFC3339, *r.CreationDate)
		if err != nil {
			log.Fatalf("failed to parse ami create date. err: %v", err)
		}

		// 引数で指定した日付より前であれば削除する
		if rtime.Before(sd) {
			fmt.Println("AMI Name: " + *r.Name)
			fmt.Println("CreateDate:" + *r.CreationDate)
			deregImageInput := ec2.DeregisterImageInput{
				ImageId: r.ImageId,
				DryRun:  isDryRun,
			}

			_, err = client.DeregisterImage(ctx, &deregImageInput)
			if err != nil {
				// 消せない場合もlogだけ出力して処理を続行する
				log.Printf("failed to deregister ami. AMI Name = %s err: %v", *r.Name, err)
			}

			// AMIの削除後にEBS Snapshotを削除する
			for _, blockDeviceMapping := range r.BlockDeviceMappings {
				delEBSSnapshotInput := ec2.DeleteSnapshotInput{
					SnapshotId: blockDeviceMapping.Ebs.SnapshotId,
					DryRun:     isDryRun,
				}
				_, err = client.DeleteSnapshot(ctx, &delEBSSnapshotInput)
				if err != nil {
					// 消せない場合もlogだけ出力して処理を続行する
					log.Printf("failed to deregister ami. EBS SnapshotID = %s err: %v", *blockDeviceMapping.Ebs.SnapshotId, err)
				}
			}
		}
	}
}
