package config

import (
	"github.com/spf13/viper"
	"time"
)

type Environment struct {
	HTTP, DatabaseURL, AttachmentRoot string
	UndoWindow                        time.Duration
}

func FromEnv() Environment {
	v := viper.New()
	v.SetEnvPrefix("OFFICE")
	v.AutomaticEnv()
	v.SetDefault("HTTP", ":8080")
	v.SetDefault("UNDO_WINDOW", "30m")
	v.SetDefault("ATTACHMENT_ROOT", "./feedback-files")
	return Environment{v.GetString("HTTP"), v.GetString("DATABASE_URL"), v.GetString("ATTACHMENT_ROOT"), v.GetDuration("UNDO_WINDOW")}
}
