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
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
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
	cp := codepipeline.New(cfg)

	cfreq := cf.CreateInvalidationRequest(&cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(distributionID),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(time.Now().String()),
			Paths: &cloudfront.Paths{
				Quantity: aws.Int64(1),
				Items:    []string{"/*"},
			},
		},
	})
	_, err = cfreq.Send(context.TODO())
	if err != nil {
		log.Println(err)
		cpreq := cp.PutJobFailureResultRequest(&codepipeline.PutJobFailureResultInput{
			FailureDetails: &codepipeline.FailureDetails{
				Message: aws.String(err.Error()),
			},
			JobId: aws.String(jobID),
		})
		_, err = cpreq.Send(context.TODO())
	} else {
		log.Println(distributionID, "invalidated")
		cpreq := cp.PutJobSuccessResultRequest(&codepipeline.PutJobSuccessResultInput{
			JobId: aws.String(jobID),
		})
		_, err = cpreq.Send(context.TODO())
	}
	if err != nil {
		log.Println(err)
	}
	return err
}

func main() {
	lambda.Start(handler)
}
