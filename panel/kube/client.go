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
	"k8s.io/apimachinery/pkg/labels"
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

func (k *Kube) CreateControlPlane(ctx context.Context, cp ControlPlane) error {
	kcp := cp.ToKubeObject()
	namespaceObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: kcp.Namespace,
		},
	}

	kcp.Spec.Server.Image = "ghcr.io/kapycluster/kapyserver:master"
	kcp.Spec.Server.Persistence = "sqlite"

	var err error

	err = k.client.Create(ctx, namespaceObj)
	if err != nil {
		return fmt.Errorf("failed to create namespace: %w", err)
	}

	err = k.client.Create(ctx, kcp)
	if err != nil {
		go k.cleanup(ctx, *kcp)
		return fmt.Errorf("failed to create ControlPlane: %w", err)
	}
	return nil
}

func (k *Kube) WatchControlPlane(ctx context.Context, cp ControlPlane) (<-chan bool, error) {
	kcp := cp.ToKubeObject()
	watcher := cache.NewListWatchFromClient(
		k.clientset.CoreV1().RESTClient(),
		"controlplanes",
		kcp.Namespace,
		fields.OneTermEqualSelector(metav1.ObjectNameField, kcp.Name),
	)

	isReady := make(chan bool)

	_, informer := cache.NewInformer(watcher, kcp, time.Second*0,
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

func (k *Kube) UpdateControlPlane(ctx context.Context, cp ControlPlane) error {
	return nil
}

func (k *Kube) DeleteControlPlane(ctx context.Context, cp ControlPlane) error {
	return nil
}

func (k *Kube) GetControlPlane(ctx context.Context, cp ControlPlane) (*ControlPlane, error) {
	kcp := &kapyv1.ControlPlane{}
	err := k.client.Get(ctx, client.ObjectKey{Namespace: cp.ID, Name: cp.Name}, kcp)
	if err != nil {
		return nil, fmt.Errorf("failed to get ControlPlane: %w", err)
	}

	return FromKubeObject(kcp), nil
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

func (k *Kube) ListControlPlanes(ctx context.Context, userID string) ([]*ControlPlane, error) {
	listOpts := client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{
			labelUserID: userID,
		}),
	}
	list := &kapyv1.ControlPlaneList{}
	err := k.client.List(ctx, list, &listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list ControlPlanes: %w", err)
	}

	cps := make([]*ControlPlane, 0, len(list.Items))

	for _, cp := range list.Items {
		cps = append(cps, FromKubeObject(&cp))
	}

	return cps, nil
}
