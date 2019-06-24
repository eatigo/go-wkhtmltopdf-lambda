# go-wkhtmltopdf-lambda
Run wkhtmltopdf in AWS Lambda

#### Not maintained for now, but mainly an example of how to use go-wkhtmltopdf on Lambda.

Uses https://github.com/SebastiaanKlippert/go-wkhtmltopdf to generate PDF files from JSON input.

# Usage

- Create a Golang AWS Lambda function in the AWS console or using the AWS CLI
- As source, use the lambda.zip (using 0.12.4 wkhtmltopdf general linux library)
  (how to zip: # zip lambda.zip ./receipt_handler ./wkhtmltopdf ./wkhtmltoimage ./fonts/*)
- The Handler name is `receipt_handler`
- If you want to build your own version, make sure you build the go executable for Linux (GOOS=linux) and make it executable (chmod +x)
  (how to compile: # GOOS=linux go build receipt_handler.go)
- Create an S3 trigger for you Lambda function (using suffix `.json`)
- Make sure your IAM Lambda role has S3 Read and Write access to your bucket
- Create a JSON file from the PDF generator, following the instructions at https://github.com/SebastiaanKlippert/go-wkhtmltopdf#saving-to-and-loading-from-json
- Upload the JSON to the S3 bucket you chose in your Lambda trigger
- The PDF is saved in the same bucket with extension `.pdf` with same filename
- The image is saved in the same bucket with extension `.jpg` with same filename
- for the fonts, many researches need to be done for thai and traditional chinese
