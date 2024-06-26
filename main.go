package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"os"
)

var (
	instanceId string
	err        error
)

func main() {
	ctx := context.TODO()
	if instanceId, err = createEc2(ctx, "us-east-1"); err != nil {
		fmt.Printf("CreateEC2 error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Instance id: %s\n", instanceId)
}

func createEc2(ctx context.Context, region string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("Unable to load SDK config: %v", err)
	}
	ec2Client := ec2.NewFromConfig(cfg)
	_, err = ec2Client.CreateKeyPair(ctx, &ec2.CreateKeyPairInput{
		KeyName: aws.String("go-aws-sdk"),
	})

	imageOutput, err := ec2Client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("name"),
				Values: []string{"ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"},
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: []string{"hvm"},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("DescribeImages error: %s", err)
	}
	if len(imageOutput.Images) == 0 {
		return "", fmt.Errorf("imageOutput.Images is of 0 lenght error")
	}

	instance, err := ec2Client.RunInstances(ctx, &ec2.RunInstancesInput{
		ImageId:      imageOutput.Images[0].ImageId,
		KeyName:      aws.String("go-aws-sdk"),
		InstanceType: types.InstanceTypeT3Micro,
		MaxCount:     aws.Int32(1),
		MinCount:     aws.Int32(1),
	})
	if err != nil {
		return "", fmt.Errorf("RunInstances error: %s", err)
	}
	if len(instance.Instances) == 0 {
		return "", fmt.Errorf("instance.Instances is of 0 lenght error")
	}

	return *instance.Instances[0].InstanceId, nil
}
