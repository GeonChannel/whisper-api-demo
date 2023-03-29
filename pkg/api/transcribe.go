package api

import (
	"ch-whisper/internal/core"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-audio/wav"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"net/http"
	"os"
	"path/filepath"
)

func (s *Server) transcribeS3File(context *gin.Context) error {
	reqS3 := core.ReqS3DTO{}
	if err := context.ShouldBindJSON(&reqS3); err != nil {
		s.logger.Errorw("should bind JSON", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	wavFilePath, err := s.S3Download(reqS3.Bucket, reqS3.Key)
	if err != nil {
		return err
	}
	return s.decodeWavAndGetResult(context, wavFilePath)
}

func (s *Server) transcribeFile(context *gin.Context) error {
	wavFilePath, err := s.getConvertedWavFilePath(context)
	if err != nil {
		return err
	}

	return s.decodeWavAndGetResult(context, wavFilePath)
}

func (s *Server) decodeWavAndGetResult(context *gin.Context, wavFilePath string) error {
	decodedAudio, err := s.getDecodedAudio(wavFilePath)
	if err != nil {
		return err
	}

	resultDTO, err := s.model.Transcribe(decodedAudio)
	if err != nil {
		s.logger.Info("fail to transcribe decoded wav file: ", err)
		return err
	}

	context.JSON(http.StatusOK, resultDTO)

	return nil
}
func (s *Server) getDecodedAudio(wavFilePath string) ([]float32, error) {
	fh, err := os.Open(wavFilePath)
	if err != nil {
		s.logger.Errorw("fail to open converted wav file: ", err)
		return nil, err
	}
	defer fh.Close()

	var decodedAudio []float32

	dec := wav.NewDecoder(fh)
	if buf, err := dec.FullPCMBuffer(); err != nil {
		panic(err)
	} else {
		decodedAudio = buf.AsFloat32Buffer().Data
	}
	defer os.Remove(wavFilePath)
	return decodedAudio, nil
}

func (s *Server) getConvertedWavFilePath(context *gin.Context) (string, error) {
	file, _ := context.FormFile("file")
	s.logger.Info("filename: ", file.Filename)
	ext := filepath.Ext(file.Filename)
	s.logger.Info("ext: ", ext)

	uid := uuid.New()

	uploadedFile := fmt.Sprintf(viper.GetString("audioFilePath")+"%s%s", uid.String(), ext)
	if err := context.SaveUploadedFile(file, uploadedFile); err != nil {
		s.logger.Errorw("fail to upload file: ", err)
		return "", err
	}

	wavFile := fmt.Sprintf("tmp/converted_%s.wav", uid.String())

	if err := convertWav(uploadedFile, wavFile); err != nil {
		s.logger.Errorw("fail to convert uploaded file to wav file: ", err)
		return "", err
	}
	defer func() { // after converting audio, uploaded file must be deleted.
		os.Remove(uploadedFile)
	}()
	return wavFile, nil
}
func convertWav(inputFile, outputFile string) error {
	err := ffmpeg.Input(inputFile).
		Output(outputFile, ffmpeg.KwArgs{
			"ar": 16000, // change sample rate
			"ac": 1,     // make mono channel
		}).
		Run()
	return err
}
