package main

import (
	"context"
	"flag"
	"github.com/SneaksAndData/nexus-core/pkg/signals"
	"github.com/SneaksAndData/nexus-core/pkg/telemetry"
	v1 "github.com/SneaksAndData/nexus-receiver/api/v1"
	"github.com/SneaksAndData/nexus-receiver/app"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

const (
	MaxBodySize = 512 * 1024 * 1024
)

var (
	logLevel string
)

func init() {
	flag.StringVar(&logLevel, "log-level", "INFO", "Log level for the application.")
}

func setupRouter(ctx context.Context) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.Default()
	router.MaxMultipartMemory = MaxBodySize
	router.Use(gin.Logger())
	appConfig := app.LoadConfig(ctx)

	appServices := (&app.ApplicationServices{}).
		WithCqlStore(ctx, &appConfig.CqlStore)

	// version 1.2
	apiV12 := router.Group("algorithm/v1.2")

	apiV12.POST("complete/:algorithmName/requests/:requestId", v1.CompleteRun(appServices.CqlStore()))

	go func() {
		appServices.Start(ctx)
		// handle exit
		logger := klog.FromContext(ctx)
		reason := ctx.Err()
		if reason.Error() == context.Canceled.Error() {
			logger.V(0).Info("Received SIGTERM, shutting down gracefully")
			klog.FlushAndExit(klog.ExitFlushTimeout, 0)
		}

		logger.V(0).Error(reason, "Fatal error occurred.")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}()

	return router
}

func main() {
	ctx := signals.SetupSignalHandler()
	appLogger, err := telemetry.ConfigureLogger(ctx, map[string]string{}, logLevel)
	ctx = telemetry.WithStatsd(ctx, "nexus_receiver")
	logger := klog.FromContext(ctx)

	if err != nil {
		logger.Error(err, "One of the logging handlers cannot be configured")
	}

	klog.SetSlogLogger(appLogger)

	r := setupRouter(ctx)
	// Configure webhost
	_ = r.Run(":8080")
}
