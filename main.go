package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

// https://github.com/made2591/immutable.templates/blob/master/templates/static-website/lib/invalidation-lambda/index.js
func handler(ctx context.Context, evt events.CodePipelineEvent) error {
	jobID := evt.CodePipelineJob.ID
	distributionID := evt.CodePipelineJob.Data.ActionConfiguration.Configuration.UserParameters
	log.Println(jobID, distributionID)
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return err
	}
	cf := cloudfront.New(cfg)
	req := cf.CreateInvalidationRequest(&cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(distributionID),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(time.Now().String()),
			Paths: &cloudfront.Paths{
				Quantity: aws.Int64(1),
				Items:    []string{"/*"},
			},
		},
	})
	_, err = req.Send(context.TODO())
	if err != nil {
		return err
	}
	log.Println(distributionID, "invalidated")
	return nil
}

func main() {
	lambda.Start(handler)
}
