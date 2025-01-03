package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

var (
	s3Client     *s3.Client
	dynamoClient *dynamodb.Client
)

type Record struct {
	ImageID        string `json:"imageId"  dynamodbav:"imageId"`
	OriginalName   string `json:"originalName"  dynamodbav:"originalName"`
	ThumbnailImage string `json:"thumbnailImage"  dynamodbav:"thumbnailImage"`
	UploadDate     string `json:"uploadDate"  dynamodbav:"uploadDate"`
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Panicf("unable to load SDK config, %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	dynamoClient = dynamodb.NewFromConfig(cfg)
}

func uploadRecordToDynamo(ctx context.Context, tableName string, item Record) error {
	av, err1 := attributevalue.MarshalMap(item)

	if err1 != nil {
		return err1
	}

	input := dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	if _, err := dynamoClient.PutItem(ctx, &input); err != nil {
		return err
	}

	return nil
}

func downloadFileToS3(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	input := s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	object, err := s3Client.GetObject(ctx, &input)

	if err != nil {
		return nil, err
	}

	return object.Body, nil
}

func uploadToS3(ctx context.Context, bucket, key string, body io.Reader) error {
	input := s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	}

	_, err := s3Client.PutObject(ctx, &input)

	if err != nil {
		return err
	}

	return nil
}

func resizeImageToThumbnail(input io.ReadCloser, maxWidth, maxHeight uint) (io.Reader, error) {
	defer input.Close()

	var buf bytes.Buffer

	_, err := io.Copy(&buf, input)
	if err != nil {
		return nil, fmt.Errorf("Error copying content: %w", err)
	}

	img, _, err := image.Decode(&buf)

	if err != nil {
		return nil, fmt.Errorf("Error decoding image: %w", err)
	}

	thumbnail := resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3)

	var outputBuf bytes.Buffer

	err = jpeg.Encode(&outputBuf, thumbnail, nil)

	if err != nil {
		return nil, fmt.Errorf("Error encoding image: %w", err)
	}

	return &outputBuf, nil
}

func Handler(ctx context.Context, event events.S3Event) error {
	tableName := os.Getenv("TABLE_NAME")

	for _, record := range event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

		file, err := downloadFileToS3(ctx, bucket, key)

		if err != nil {
			log.Fatalf("Error processing file: %v", err)
			return err
		}

		thumbnailFile, err := resizeImageToThumbnail(file, 128, 128)

		if err != nil {
			return err
		}

		thumbnailKey := fmt.Sprintf("thumbnails/thumbnail-%s.jpg", uuid.New().String())

		err = uploadToS3(ctx, bucket, thumbnailKey, thumbnailFile)

		if err != nil {
			log.Fatalf("Error processing file: %v", err)
			return err
		}

		record := Record{
			ImageID:        uuid.New().String(),
			OriginalName:   key,
			ThumbnailImage: thumbnailKey,
			UploadDate:     time.Now().UTC().Format(time.RFC3339),
		}

		err = uploadRecordToDynamo(ctx, tableName, record)

		if err != nil {
			log.Fatalf("Error processing file: %v", err)
			return err
		}

		log.Printf("Processed %s and created %s", key, thumbnailKey)
	}

	return nil
}

func main() {
	lambda.Start(Handler)

}
