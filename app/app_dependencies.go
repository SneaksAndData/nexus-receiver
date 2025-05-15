package app

import (
	"context"
	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/request"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
)

type ApplicationServices struct {
	cqlStore        *request.CqlStore
	recorder        record.EventRecorder
	completionActor *CompletionActor
}

func (appServices *ApplicationServices) WithCqlStore(ctx context.Context, bundleConfig *request.AstraBundleConfig) *ApplicationServices {
	if appServices.cqlStore == nil {
		logger := klog.FromContext(ctx)
		appServices.cqlStore = request.NewAstraCqlStore(logger, bundleConfig)
	}

	return appServices
}

func (appServices *ApplicationServices) WithCompletionActor(config *ReceiverConfig) *ApplicationServices {
	if appServices.completionActor == nil {
		appServices.completionActor = NewCompletionActor(appServices.cqlStore, config)
	}

	return appServices
}

func (appServices *ApplicationServices) CqlStore() *request.CqlStore {
	return appServices.cqlStore
}

func (appServices *ApplicationServices) CompletionActor() *CompletionActor {
	return appServices.completionActor
}

func (appServices *ApplicationServices) Start(ctx context.Context) {
	go appServices.completionActor.Start(ctx)
}
