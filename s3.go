package judgego

import (
	"bytes"
	"errors"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	// soundClipRegex is the regex used to pull sound clip names from the bucket
	soundClipRegex  = audioFilePrefix + "(.+)"
	audioFilePrefix = "sound-clips/"
)

var bucketName string = os.Getenv("BUCKET_NAME")

// TODO: This entire file can be genericized a bit, there is a mix of app specific code and generic functionality.

func listSoundsS3() ([]string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return nil, errors.New("Unable to access AWS")
	}
	svc := s3.New(sess)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucketName)})
	if err != nil {
		return nil, errors.New("Unable to access sound bucket")
	}

	sounds := make([]string, 0)
	re := regexp.MustCompile(soundClipRegex)
	for _, item := range resp.Contents {
		matches := re.FindSubmatch([]byte(*item.Key))
		if matches != nil {
			sounds = append(sounds, string(matches[1]))
		}
	}
	return sounds, nil
}

func putSoundS3(sound *bytes.Buffer, name string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(audioFilePrefix + name),
		Body:   sound,
	})
	if err != nil {
		return err
	}

	return nil
}

func getSoundS3(name string) []byte {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	downloader := s3manager.NewDownloader(sess)

	buf := aws.NewWriteAtBuffer([]byte{})
	_, _ = downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(audioFilePrefix + name),
		})
	return buf.Bytes()
}

func writeToS3(b *bytes.Buffer, name string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(name),
		Body:   b,
	})
	if err != nil {
		return err
	}

	return nil

}

func getFromS3(name string) ([]byte, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return nil, err
	}

	downloader := s3manager.NewDownloader(sess)

	buf := aws.NewWriteAtBuffer([]byte{})
	_, err = downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(name),
		})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
