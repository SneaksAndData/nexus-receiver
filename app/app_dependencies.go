package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/SneaksAndData/nexus-core/pkg/buildmeta"
	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/request"
	nexuscore "github.com/SneaksAndData/nexus-core/pkg/generated/clientset/versioned"
	"github.com/SneaksAndData/nexus-core/pkg/generated/clientset/versioned/scheme"
	nexusscheme "github.com/SneaksAndData/nexus-core/pkg/generated/clientset/versioned/scheme"
	"github.com/SneaksAndData/nexus-core/pkg/pipeline"
	"github.com/SneaksAndData/nexus-core/pkg/shards"
	//"github.com/SneaksAndData/nexus/services"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"time"
)

type ApplicationServices struct {
	cqlStore *request.CqlStore
	recorder record.EventRecorder
}

func (appServices *ApplicationServices) WithCqlStore(ctx context.Context, bundleConfig *request.AstraBundleConfig) *ApplicationServices {
	if appServices.cqlStore == nil {
		logger := klog.FromContext(ctx)
		appServices.cqlStore = request.NewAstraCqlStore(logger, bundleConfig)
	}

	return appServices
}

func (appServices *ApplicationServices) CqlStore() *request.CqlStore {
	return appServices.cqlStore
}

func (appServices *ApplicationServices) schedule(output *request.BufferOutput) (types.UID, error) {
	if output == nil {
		return types.UID(""), fmt.Errorf("buffer is nil")
	}

	var job = output.Checkpoint.ToV1Job("kubernetes.sneaksanddata.com/service-node-group", output.Checkpoint.AppliedConfiguration.Workgroup, fmt.Sprintf("%s-%s", buildmeta.AppVersion, buildmeta.BuildNumber))
	var submitted *batchv1.Job
	var submitErr error

	// submit to controller cluster if workgroup host is not provided
	if output.Checkpoint.AppliedConfiguration.WorkgroupHost == "" {
		submitted, submitErr = appServices.kubeClient.BatchV1().Jobs(appServices.defaultNamespace).Create(context.TODO(), &job, v1.CreateOptions{})
	}

	if shard := appServices.getShardByName(output.Checkpoint.AppliedConfiguration.WorkgroupHost); shard != nil {
		submitted, submitErr = shard.SendJob(shard.Namespace, &job)
	} else {
		return "", errors.New(fmt.Sprintf("Shard API server %s not configured", output.Checkpoint.AppliedConfiguration.WorkgroupHost))
	}

	if submitErr != nil {
		return "", submitErr
	}

	return submitted.UID, nil
}

func (appServices *ApplicationServices) Start(ctx context.Context) {
	logger := klog.FromContext(ctx)
	err := appServices.configCache.Init(ctx)
	if err != nil {
		logger.Error(err, "Error building in-cluster kubeconfig for the scheduler")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}
	submissionActor := pipeline.NewDefaultPipelineStageActor[*request.BufferOutput, types.UID](
		"checkpoint_buffer",
		map[string]string{},
		time.Second*1,
		time.Second*5,
		10,
		100,
		10,
		appServices.schedule,
		nil,
	)

	go submissionActor.Start(ctx)
	appServices.checkpointBuffer.Start(submissionActor)
}
