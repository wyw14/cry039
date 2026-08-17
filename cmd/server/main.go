package main

import (
	"context"
	"errors"
	"github.com/wyw14/cry039/internal/application"
	"github.com/wyw14/cry039/internal/config"
	"github.com/wyw14/cry039/internal/domain"
	"github.com/wyw14/cry039/internal/repository"
	feedbackhttp "github.com/wyw14/cry039/internal/transport/http"
	"go.uber.org/zap"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.FromEnv()
	log, _ := zap.NewProduction()
	defer log.Sync()
	seed := []domain.Feedback{{ID: "f-demo", AreaID: "north-2", Tag: "noise", Status: "open", Level: 4, SubmittedAt: time.Now().Add(-time.Hour)}}
	mover := application.NewMover(repository.NewFeedbackMemory("demo", seed), time.Now)
	srv := http.Server{Addr: cfg.HTTP, Handler: feedbackhttp.Server(mover, log), ReadHeaderTimeout: 4 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("listen", zap.Error(err))
		}
	}()
	<-ctx.Done()
	closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(closeCtx)
}
