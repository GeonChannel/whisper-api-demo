package api

import (
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (s *Server) S3Download(bucket string, key string) (string, error) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		// Credentials:
	})
	if err != nil {
		s.logger.Errorw("Error creating session:", err)
		return "", err
	}
	s.logger.Info("S3 new session: ", sess)

	downloader := s3manager.NewDownloader(sess)
	uid := uuid.New()
	newFileName := uid.String() + ".wav"
	newFilePath := newFileName + viper.GetString("audioFilePath") + newFileName
	file, err := os.Create(newFilePath)
	if err != nil {
		s.logger.Errorw("Error creating file:", err)
		return "", err
	}
	defer file.Close()

	s.logger.Info("create new file: ", file)
	numBytes, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		s.logger.Errorw("Error downloading file:", err)
		return "", err
	}
	numMB := float64(numBytes) / (10 << 20)
	s.logger.Info("downloaded size(MB): ", numMB)

	return newFilePath, nil
}
