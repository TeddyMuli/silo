package spaces

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client

func InitS3() {
	// Get the DigitalOcean Spaces URL from the environment variables
	spacesURL := os.Getenv("D_O_SPACES_URL")
	if spacesURL == "" {
		log.Println("Error, missing DigitalOcean Spaces URL in environment variables")
		return
	}

	// Get the AWS credentials from environment variables
	accessKey := os.Getenv("D_O_ACCESS_KEY_ID")
	secretKey := os.Getenv("D_O_SECRET_ACCESS_KEY")
	if accessKey == "" || secretKey == "" {
		log.Println("Error, missing Spaces credentials in environment variables")
		return
	}

	// Use static credentials provider
	staticCredentials := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""))

	// Load the configuration with static credentials and region
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(staticCredentials),
	)
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// Initialize the S3 client with credentials and the DigitalOcean Spaces endpoint
	s3Client := s3.New(s3.Options{
		Credentials:       cfg.Credentials,
		Region:            cfg.Region,
		EndpointResolver:  s3.EndpointResolverFromURL("https://" + spacesURL),
	})

	// Set the S3 client globally
	S3Client = s3Client
}

func DeleteFile(fileName string) error {
	// Create the input for the delete request
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("D_O_SPACES_URL")),
		Key:    aws.String(fileName),
	}

	_, err := S3Client.DeleteObject(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to delete file %s: %v", fileName, err)
		return err
	}

	return nil
}
