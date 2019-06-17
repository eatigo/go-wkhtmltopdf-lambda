package main

import (
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eatigo/go-wkhtmltopdf"
)

func main() {
	lambda.Start(handler)
}

func handler(s3Event events.S3Event) error {
	if len(s3Event.Records) > 0 {
		err := createImage(s3Event.Records[0])
		if err != nil {

		}
		err = createPDF(s3Event.Records[0])
		return err
	}
	return nil
}

func createPDF(record events.S3EventRecord) error {

	// get json file from S3
	object, err := getS3Object(record.S3.Bucket.Name, record.S3.Object.Key)
	if err != nil {
		return err
	}
	defer object.Close()

	// make sure we look for the included wkhtmltopdf binary
	os.Setenv("WKHTMLTOPDF_PATH", os.Getenv("LAMBDA_TASK_ROOT"))

	// create PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGeneratorFromJSON(object)
	if err != nil {
		return err
	}

	// create PDF
	err = pdfg.Create()
	if err != nil {
		return err
	}

	// write PDF to same filename with .pdf added
	return putS3Object(record.S3.Bucket.Name, strings.Replace(record.S3.Object.Key, ".json", ".pdf", -1), pdfg.Bytes())
}

func createImage(record events.S3EventRecord) error {

	// get json file from S3
	object, err := getS3Object(record.S3.Bucket.Name, record.S3.Object.Key)
	if err != nil {
		return err
	}
	defer object.Close()

	// make sure we look for the included wkhtmltoimage binary
	os.Setenv("WKHTMLTOIMAGE_PATH", os.Getenv("LAMBDA_TASK_ROOT"))

	// create image
	image, err := wkhtmltopdf.ImageFromJSON(object)
	if err != nil {
		return err
	}

	// write PDF to same filename with .png added
	return putS3Object(record.S3.Bucket.Name, strings.Replace(record.S3.Object.Key, ".json", ".png", -1), image)
}
