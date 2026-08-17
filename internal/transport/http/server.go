package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/wyw14/cry039/internal/application"
	"github.com/wyw14/cry039/internal/domain"
	"github.com/wyw14/cry039/internal/middleware"
	"go.uber.org/zap"
	"net/http"
)

type moveRequest struct {
	Key         string      `json:"idempotency_key" validate:"required"`
	Digest      string      `json:"digest" validate:"required"`
	MigrationID string      `json:"migration_id" validate:"required"`
	Source      domain.Area `json:"source"`
	Target      domain.Area `json:"target"`
	Actor       string      `json:"actor" validate:"required"`
	Reviewers   []string    `json:"reviewers" validate:"len=2"`
}

func Server(mover *application.Mover, logger *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.RequestContext())
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"component": "feedback-correction"}) })
	r.GET("/readyz", func(c *gin.Context) { c.JSON(200, gin.H{"ready": true}) })
	check := validator.New()
	r.POST("/api/v1/migrations/execute", func(c *gin.Context) {
		var in moveRequest
		if c.ShouldBindJSON(&in) != nil || check.Struct(in) != nil {
			c.JSON(422, gin.H{"code": "INVALID_MIGRATION", "message": "迁移参数不完整", "request_id": c.GetString("request_id")})
			return
		}
		job := &domain.Migration{ID: in.MigrationID, SourceArea: in.Source.ID, TargetArea: in.Target.ID}
		for _, reviewer := range in.Reviewers {
			if err := job.Approve(reviewer); err != nil {
				c.JSON(409, gin.H{"code": "REVIEW_CONFLICT", "message": err.Error(), "request_id": c.GetString("request_id")})
				return
			}
		}
		items, err := mover.Execute(c, in.Key, in.Digest, job, in.Source, in.Target, in.Actor)
		if err != nil {
			logger.Info("migration rejected", zap.Error(err))
			c.JSON(409, gin.H{"code": "MIGRATION_REJECTED", "message": err.Error(), "request_id": c.GetString("request_id")})
			return
		}
		c.JSON(200, gin.H{"migrated": len(items) + 1})
	})
	return r
}
