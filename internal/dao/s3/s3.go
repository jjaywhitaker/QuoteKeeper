package s3

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	//"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jjaywhitaker/QuoteKeeper/internal/model"
)

type s3DAO interface {
}

type s3Dao struct {
}

var (
	region = os.Getenv("AWS_REGION")
	bucket = os.Getenv("S3_BUCKET")
	key    = os.Getenv("S3_KEY")
)

func NewS3Dao() s3Dao {
	dao := s3Dao{}
	return dao
}

func (s s3Dao) ReadQuoteFile() ([]model.QuoteResponse, error) {
	//TODO how to login to account?
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		panic(err)
	}

	s3Client := s3.New(sess)
	requestInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := s3Client.GetObject(requestInput)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Body.Close()
	body1, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyString1 := fmt.Sprintf("%s", body1)

	var s3data []model.QuoteResponse
	decoder := json.NewDecoder(strings.NewReader(bodyString1))
	err = decoder.Decode(&s3data)
	if err != nil {
		fmt.Println("Yikes! There an error decoding")
	}

	fmt.Println(s3data)

	return s3data, err
}
