package app

import (
	"context"

	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/store"
	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/store/cassandra"
	"k8s.io/klog/v2"
)

type ApplicationServices struct {
	cqlStore        *store.CheckpointStore
	completionActor *CompletionActor
}

func (appServices *ApplicationServices) WithAstraCqlStore(ctx context.Context, bundleConfig *cassandra.AstraBundleConfig) *ApplicationServices {
	if appServices.cqlStore == nil {
		logger := klog.FromContext(ctx)
		appServices.cqlStore = new(cassandra.NewAstraStore(logger, bundleConfig))
	}

	return appServices
}

func (appServices *ApplicationServices) WithScyllaCqlStore(ctx context.Context, config *cassandra.ScyllaConfig) *ApplicationServices {
	if appServices.cqlStore == nil {
		logger := klog.FromContext(ctx)
		appServices.cqlStore = new(cassandra.NewScyllaStore(logger, config))
	}

	return appServices
}

func (appServices *ApplicationServices) WithKeyspacesCqlStore(ctx context.Context, config *cassandra.KeyspacesConfig) *ApplicationServices {
	if appServices.cqlStore == nil {
		logger := klog.FromContext(ctx)
		appServices.cqlStore = new(cassandra.NewKeyspacesStore(logger, config))
	}

	return appServices
}

func (appServices *ApplicationServices) WithCompletionActor(ctx context.Context, config *ReceiverConfig) *ApplicationServices {
	if appServices.completionActor == nil {
		appServices.completionActor = NewCompletionActor(ctx, appServices.cqlStore, config)
	}

	return appServices
}

func (appServices *ApplicationServices) CheckpointStore() *store.CheckpointStore {
	return appServices.cqlStore
}

func (appServices *ApplicationServices) CompletionActor() *CompletionActor {
	return appServices.completionActor
}

func (appServices *ApplicationServices) Start(ctx context.Context) {
	appServices.completionActor.Start(ctx, nil)
}
