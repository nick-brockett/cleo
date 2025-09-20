package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"cleo.com/internal/adapter/auth"
	"cleo.com/internal/adapter/handler/http"
	"cleo.com/internal/core/service"

	"github.com/sethvargo/go-envconfig"
	"github.com/sirupsen/logrus"
)

func main() {
	s := "Cleo Health Metric Extraction API"
	fmt.Printf("Hello and welcome to %s!\n", s)

	ctx := context.Background()

	logger := initLogger()

	authCfg := auth.Config{}
	if err := envconfig.Process(ctx, &authCfg); err != nil {
		log.Fatal("failed to load auth config", "error", err)
	}

	authService := auth.NewService(logger, authCfg)
	parserService := service.NewParserService(logger)
	healthMetricHandler := http.NewHealthMetricParserHandler(logger, parserService)

	router, err := http.NewRouter(authService, healthMetricHandler)
	if err != nil {
		log.Fatal("error initializing router", "error", err)
	}

	_ = router.SetTrustedProxies(nil)

	err = router.Serve(":8080")
	if err != nil {
		log.Fatal("server error", "error", err)
	}

}

func initLogger() *logrus.Logger {
	logger := logrus.New()
	level := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if lvl, err := logrus.ParseLevel(level); err == nil {
		logger.SetLevel(lvl)
	} else {
		logger.SetLevel(logrus.InfoLevel) // default
	}
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}

/*
example test Bearer token
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTc1ODM4MDIzOCwicm9sZSI6IkNMSU5JQ0FMLUVESVRPUiJ9.KEN8skt0EzoarE-Wj-LckMqjSAuqwaULW72IdAYPRlM
*/
