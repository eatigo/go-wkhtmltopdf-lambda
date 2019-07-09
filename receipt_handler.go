package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/eatigo/go-wkhtmltopdf"
)

const contentTypeImage = "image/jpeg"
const contentTypePdf = "application/pdf"

func main() {
	lambda.Start(handler)
}

func handler(s3Event events.S3Event) error {
	if len(s3Event.Records) > 0 {
		os.Setenv("FONTCONFIG_PATH", "/var/task/fonts")
		err := createPDF(s3Event.Records[0])
		if err != nil {
			return fmt.Errorf("create pdf error, %s", err.Error())
		}
		err = createImage(s3Event.Records[0])
		if err != nil {
			return fmt.Errorf("create image error, %s", err.Error())
		}
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
	return putS3Object(record.S3.Bucket.Name, strings.Replace(record.S3.Object.Key, ".json", ".pdf", -1), contentTypePdf, pdfg.Bytes())
}

func createImage(record events.S3EventRecord) error {

	// get json file from S3
	object, err := getS3Object(record.S3.Bucket.Name, record.S3.Object.Key)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer object.Close()

	// make sure we look for the included wkhtmltoimage binary
	os.Setenv("WKHTMLTOIMAGE_PATH", os.Getenv("LAMBDA_TASK_ROOT"))

	// create image
	image, err := wkhtmltopdf.ImageFromJSON(object)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// write PDF to same filename with .png added
	return putS3Object(record.S3.Bucket.Name, strings.Replace(record.S3.Object.Key, ".json", ".jpg", -1), contentTypeImage, image)
}

func getS3Object(bucket, key string) (io.ReadCloser, error) {

	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("error creating AWS session (check access keys): %s", err)
	}

	in := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	obj, err := s3.New(sess).GetObject(in)
	if err != nil {
		return nil, fmt.Errorf("error getting S3 object: %s", err)
	}
	return obj.Body, nil
}

func putS3Object(bucket, key, contentType string, buf []byte) error {

	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("error creating AWS session (check access keys): %s", err)
	}

	in := &s3.PutObjectInput{
		Bucket:             aws.String(bucket),
		Key:                aws.String(key),
		Body:               bytes.NewReader(buf),
		ContentType:        aws.String(contentType),
		ContentDisposition: aws.String(fmt.Sprintf("attachment;filename=%s", key)),
	}

	_, err = s3.New(sess).PutObject(in)
	if err != nil {
		return fmt.Errorf("error putting S3 object: %s", err)
	}
	return nil
}
