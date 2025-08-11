package main

import (
	"context"
	"errors"
	"fmt"
	nexusconf "github.com/SneaksAndData/nexus-core/pkg/configurations"
	"github.com/SneaksAndData/nexus-core/pkg/signals"
	"github.com/SneaksAndData/nexus-core/pkg/telemetry"
	v1 "github.com/SneaksAndData/nexus-receiver/api/v1"
	"github.com/SneaksAndData/nexus-receiver/app"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	"os"
	"strconv"
)

func setupRouter(ctx context.Context, appConfig *app.ReceiverConfig) *gin.Engine {
	gin.DisableConsoleColor()
	router := gin.Default()
	router.Use(gin.Logger())
	// disable trusted proxies check
	_ = router.SetTrustedProxies(nil)
	// set runtime mode
	gin.SetMode(os.Getenv("GIN_MODE"))

	appServices := &app.ApplicationServices{}

	switch appConfig.CqlStoreType {
	case app.CqlStoreAstra:
		appServices = appServices.WithAstraCqlStore(ctx, &appConfig.AstraCqlStore)
	case app.CqlStoreScylla:
		appServices = appServices.WithScyllaCqlStore(ctx, &appConfig.ScyllaCqlStore)
	default:
		klog.FromContext(ctx).Error(errors.New("unknown store type "+appConfig.CqlStoreType), "failed to initialize a CqlStore")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	appServices = appServices.
		WithCompletionActor(appConfig)

	// version 1
	apiV1 := router.Group("algorithm/v1")

	apiV1.POST("complete/:algorithmName/requests/:requestId", v1.CompleteRun(appServices.CompletionActor()))

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

// @title           Nexus Receiver API
// @version         1.0
// @description     Nexus Receiver API specification. All Nexus supported clients conform to this spec.

// @contact.name   ESD Support
// @contact.email  esdsupport@ecco.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /algorithm/v1
func main() {
	ctx := signals.SetupSignalHandler()
	appConfig := nexusconf.LoadConfig[app.ReceiverConfig](ctx)
	appLogger, err := telemetry.ConfigureLogger(ctx, map[string]string{}, appConfig.LogLevel)
	ctx = telemetry.WithStatsd(ctx, "nexus_receiver")
	logger := klog.FromContext(ctx)

	if err != nil {
		logger.Error(err, "One of the logging handlers cannot be configured")
	}

	klog.SetSlogLogger(appLogger)

	r := setupRouter(ctx, &appConfig)
	// Configure webhost
	_ = r.Run(fmt.Sprintf(":%s", strconv.Itoa(appConfig.BindPort)))
}
