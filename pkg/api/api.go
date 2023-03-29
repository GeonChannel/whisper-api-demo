package api

import (
	"ch-whisper/internal/core"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	router *gin.Engine
	model  *core.Model
	logger *zap.SugaredLogger
}

func NewServer(model *core.Model) *Server {
	return &Server{
		router: gin.Default(),
		model:  model,
		logger: zap.L().Named("ApiServer").Sugar(),
	}
}

func (s *Server) Serve() {
	s.router.Use(gin.Logger())

	s.router.MaxMultipartMemory = viper.GetInt64("maximumAudioFileMiBSize") << 20 // 80 MiB

	s.router.POST("/api/transcribe/audioFile", func(context *gin.Context) {
		if err := s.transcribeFile(context); err != nil {
			s.logger.Errorw("fail to transcribeFile", err)
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	s.router.POST("/api/transcribe/S3", func(context *gin.Context) {
		if err := s.transcribeS3File(context); err != nil {
			s.logger.Errorw("fail to transcribeS3File", err)
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	go func() {
		port := viper.GetString("api.port")
		s.logger = s.logger.With("port", port)
		s.logger.Infow("starting api server")
		if err := s.router.Run(port); err != nil {
			s.logger.Errorw("failed to serve api server", err)
		}
	}()
}
