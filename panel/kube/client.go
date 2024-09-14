package kube

import (
	"context"
	"fmt"
	"os"
	"time"

	kapyv1 "github.com/kapycluster/corpy/controller/api/v1"
	"github.com/kapycluster/corpy/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Kube struct {
	client    client.Client
	clientset kubernetes.Interface
	dynamic   dynamic.Interface
}

// NewKube creates a new Kube client
func NewKube() (*Kube, error) {
	if err := kapyv1.AddToScheme(scheme.Scheme); err != nil {
		return nil, fmt.Errorf("failed to add ControlPlane to scheme: %w", err)
	}

	restConfig, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return nil, fmt.Errorf("failed to create rest config: %w", err)
	}
	dynamic, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s dynamic client: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s clientset: %w", err)
	}

	client, err := client.New(restConfig, client.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client: %w", err)
	}

	return &Kube{clientset: clientset, client: client, dynamic: dynamic}, nil
}

func (k *Kube) CreateControlPlane(ctx context.Context, cp kapyv1.ControlPlane) error {
	namespaceObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: cp.Namespace,
		},
	}

	cp.Spec.Server.Image = "ghcr.io/kapycluster/kapyserver:master"
	cp.Spec.Server.Persistence = "sqlite"

	var err error

	err = k.client.Create(ctx, namespaceObj)
	if err != nil {
		return fmt.Errorf("failed to create namespace: %w", err)
	}

	err = k.client.Create(ctx, &cp)
	if err != nil {
		go k.cleanup(ctx, cp)
		return fmt.Errorf("failed to create ControlPlane: %w", err)
	}
	return nil
}

func (k *Kube) WatchControlPlane(ctx context.Context, cp kapyv1.ControlPlane) (<-chan bool, error) {
	watcher := cache.NewListWatchFromClient(
		k.clientset.CoreV1().RESTClient(),
		"controlplanes",
		cp.Namespace,
		fields.OneTermEqualSelector(metav1.ObjectNameField, cp.Name),
	)

	isReady := make(chan bool)

	_, informer := cache.NewInformer(watcher, &cp, time.Second*0,
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldObj, newObj interface{}) {
				if newObj.(*kapyv1.ControlPlane).Status.Ready {
					close(isReady)
				}
			},
		},
	)

	stopCh := make(chan struct{})
	defer close(stopCh)

	go informer.Run(stopCh)

	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(isReady)
				if ctx.Err() == context.DeadlineExceeded {
					isReady <- false
				}
				return
			case <-isReady:
				isReady <- true
				return
			}
		}
	}()

	return isReady, nil
}

func (k *Kube) UpdateControlPlane(ctx context.Context, cp kapyv1.ControlPlane) error {
	return nil
}

func (k *Kube) DeleteControlPlane(ctx context.Context, cp kapyv1.ControlPlane) error {
	return nil
}

func (k *Kube) GetControlPlane(ctx context.Context, namespace, name string) (*kapyv1.ControlPlane, error) {
	cp := &kapyv1.ControlPlane{}
	err := k.client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, cp)
	if err != nil {
		return nil, fmt.Errorf("failed to get ControlPlane: %w", err)
	}

	return cp, nil
}

func (k *Kube) cleanup(ctx context.Context, cp kapyv1.ControlPlane) error {
	err := k.client.Delete(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: cp.Namespace,
		},
	})

	if err != nil {
		log.FromContext(ctx).Error("failed to delete namespace", "namespace", cp.Namespace)
	}

	return nil
}
