package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/sofiukl/io-service/models"
	"github.com/sofiukl/io-service/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

// UploadFileS3 - Upload file to s3
func UploadFileS3(DB *mongo.Database, config utils.Config, w http.ResponseWriter, r *http.Request) {
	uploadedFile, header, err := r.FormFile("file")
	if err != nil {
		utils.PrintErrorf("Unable to get file from request %q, %v", err)
		return
	}
	defer uploadedFile.Close()
	filename := header.Filename

	// Create session
	sess, err := awsSession.NewSessionWithOptions(awsSession.Options{
		Profile: config.AWSProfile,
		Config: aws.Config{
			Region: aws.String(config.AWSRegion),
		},
	})
	if err != nil {
		utils.PrintErrorf("Fail to create aws session %s", err)
		return
	}
	uuid := utils.GenUUID()
	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.AWSBucketName),
		Key:    aws.String(uuid),
		Body:   uploadedFile,
	})
	if err != nil {
		utils.PrintErrorf("Unable to upload %q to %q, %v", filename, config.AWSBucketName, err)
		return
	}

	// save data in db
	uh := models.UploadHistory{FileName: filename, FileKey: uuid}
	_, err = utils.SaveToDB(DB, utils.COLLECTION_UPLOAD_HISTORY, uh)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	fmt.Printf("Successfully uploaded %q to %q\n", filename, config.AWSBucketName)

	// Return response

	utils.RespondWithJSON(w, http.StatusOK, "File uploaded successfully", "", uh)
}

// GetDownloadURL - get download url
func GetDownloadURL(fileKey string, config utils.Config, w http.ResponseWriter, r *http.Request) {
	// Create session
	sess, err := awsSession.NewSessionWithOptions(awsSession.Options{
		Profile: config.AWSProfile,
		Config: aws.Config{
			Region: aws.String(config.AWSRegion),
		},
	})
	if err != nil {
		utils.PrintErrorf("Fail to create aws session %s", err)
		return
	}

	// Create S3 service client
	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(config.AWSBucketName),
		Key:    aws.String(fileKey),
	})
	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		log.Println("Failed to sign request", err)
	}

	log.Println("The URL is", urlStr)
	utils.RespondWithJSON(w, http.StatusOK, "Generated download url successfully", "", urlStr)
}
