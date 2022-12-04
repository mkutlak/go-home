package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {
	var (
		instanceId string
		err        error
	)

	ctx := context.Background()

	if instanceId, err = createEC2(ctx, "eu-central-1"); err != nil {
		fmt.Printf("createEC2 error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Instance ID: %s\n", instanceId)
}

func createEC2(ctx context.Context, region string) (string, error) {
	var (
		keyPairOutput *ec2.CreateKeyPairOutput
		err           error
	)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))

	if err != nil {
		return "", fmt.Errorf("unable to load SDK config: %v", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)

	keyPairs, err := ec2Client.DescribeKeyPairs(ctx, &ec2.DescribeKeyPairsInput{
		KeyNames: []string{"go-aws-mkutlak-demo"},
	})

	if err != nil && !strings.Contains(err.Error(), "InvalidKeyPair.NotFound") {
		return "", fmt.Errorf("unable to describe keypairs: %v", err)
	}

	if keyPairs == nil || len(keyPairs.KeyPairs) == 0 {
		keyPairOutput, err = ec2Client.CreateKeyPair(ctx, &ec2.CreateKeyPairInput{
			KeyName: aws.String("go-aws-mkutlak-demo"),
		})

		if err != nil {
			return "", fmt.Errorf("unable to create KeyPair: %v", err)
		}

		os.WriteFile("go-aws-mkutlak-demo.priv", []byte(*keyPairOutput.KeyMaterial), 0600)
	}

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
		Owners: []string{"099720109477"},
	})

	if err != nil {
		return "", fmt.Errorf("unable to create KeyPair: %v", err)
	}

	if len(imageOutput.Images) == 0 {
		return "", fmt.Errorf("imageOutput Images is 0 length")
	}

	// imageOutput.Images[0].ImageId
	instance, err := ec2Client.RunInstances(ctx, &ec2.RunInstancesInput{
		ImageId:      aws.String(*imageOutput.Images[0].ImageId),
		KeyName:      aws.String("go-aws-mkutlak-demo"),
		InstanceType: types.InstanceTypeT3Micro,
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
	})

	if err != nil {
		return "", fmt.Errorf("unable to create ec2 instance: %v", err)
	}

	if len(instance.Instances) == 0 {
		return "", fmt.Errorf("instance.Instances is 0 length")
	}

	return *instance.Instances[0].ImageId, nil
}
