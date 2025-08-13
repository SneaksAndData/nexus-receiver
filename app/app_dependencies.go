package app

import (
	"context"
	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/request"
	"k8s.io/klog/v2"
)

type ApplicationServices struct {
	cqlStore        *request.CqlStore
	completionActor *CompletionActor
}

func (appServices *ApplicationServices) WithAstraCqlStore(ctx context.Context, bundleConfig *request.AstraBundleConfig) *ApplicationServices {
	if appServices.cqlStore == nil {
		logger := klog.FromContext(ctx)
		appServices.cqlStore = request.NewAstraCqlStore(logger, bundleConfig)
	}

	return appServices
}

func (appServices *ApplicationServices) WithScyllaCqlStore(ctx context.Context, config *request.ScyllaCqlStoreConfig) *ApplicationServices {
	if appServices.cqlStore == nil {
		logger := klog.FromContext(ctx)
		appServices.cqlStore = request.NewScyllaCqlStore(logger, config)
	}

	return appServices
}

func (appServices *ApplicationServices) WithCompletionActor(ctx context.Context, config *ReceiverConfig) *ApplicationServices {
	if appServices.completionActor == nil {
		appServices.completionActor = NewCompletionActor(ctx, appServices.cqlStore, config)
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
	appServices.completionActor.Start(ctx, nil)
}
