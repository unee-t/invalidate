package main

import (
	"context"
	"encoding/json"

	"github.com/apex/log"
	jsonhandler "github.com/apex/log/handlers/json"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, evt json.RawMessage) error {
	// cfg, err := external.LoadDefaultAWSConfig()
	// if err != nil {
	// 	log.WithError(err).Error("failed to load AWS config")
	// 	return err
	// }
	log.WithField("raw", string(evt)).Info("incoming")
	return nil
}

func main() {
	log.SetHandler(jsonhandler.Default)
	lambda.Start(handler)
}
