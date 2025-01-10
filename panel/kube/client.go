package kube

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	kapyv1 "kapycluster.com/corp/controller/api/v1"
	"kapycluster.com/corp/log"
	"kapycluster.com/corp/panel/config"
	"kapycluster.com/corp/panel/kube/kubeclient"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Kube struct {
	c  *config.Config
	kc *kubeclient.KubeClient
}

// NewKube creates a new Kube client
func NewKube(ctx context.Context, c *config.Config) (*Kube, error) {
	if err := kapyv1.AddToScheme(scheme.Scheme); err != nil {
		return nil, fmt.Errorf("failed to add ControlPlane to scheme: %w", err)
	}

	kc, err := kubeclient.New(c)
	if err != nil {
		return nil, fmt.Errorf("failed to setup kubeclients: %w", err)
	}

	for _, r := range kc.GetRegions() {
		log.FromContext(ctx).Info("loaded kubeclient", "region", r)
	}

	return &Kube{
		c:  c,
		kc: kc,
	}, nil
}

// +kubebuilder:rbac:groups=kapy.sh,resources=controlplanes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kapy.sh,resources=controlplanes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="core",resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (k *Kube) CreateControlPlane(ctx context.Context, cp ControlPlane) error {
	var err error
	kcp := cp.ToKubeObject()

	namespaceObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: kcp.Namespace,
		},
	}

	if k.c.Server.ListenHost == "localhost" {
		kcp.Spec.Network.LoadBalancerAddress = "0.0.0.0"
	}
	kcp.Spec.Server.Image = "ghcr.io/kapycluster/kapyserver@sha256:594cf0bdc606804088b1da78bcb4fb2b1869aeee56feaa31da0140f42498cb54"

	cl := k.kc.GetClient(cp.Region)

	err = cl.Create(ctx, namespaceObj)
	if err != nil {
		return fmt.Errorf("failed to create namespace: %w", err)
	}

	dockerRegistrySecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "regcred",
			Namespace: kcp.Namespace,
		},
		Type: corev1.SecretTypeDockerConfigJson,
		StringData: map[string]string{
			corev1.DockerConfigJsonKey: k.c.Server.PullToken,
		},
	}

	err = cl.Create(ctx, dockerRegistrySecret)
	if err != nil {
		return fmt.Errorf("failed to create docker registry secret: %w", err)
	}

	err = cl.Create(ctx, kcp)
	if err != nil {
		go func() {
			err := cl.Delete(ctx, &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: kcp.Namespace,
				},
			})

			if err != nil {
				log.FromContext(ctx).Error("failed to delete namespace", "namespace", kcp.Namespace)
			}
		}()
		return fmt.Errorf("failed to create ControlPlane: %w", err)
	}
	return nil
}

func (k *Kube) WatchControlPlane(ctx context.Context, cp ControlPlane) (<-chan bool, error) {
	kcp := cp.ToKubeObject()

	cls := k.kc.GetClientset(cp.Region)

	watcher := cache.NewListWatchFromClient(
		cls.CoreV1().RESTClient(),
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
	cl := k.kc.GetClient(cp.Region)
	err := cl.Get(ctx, client.ObjectKey{Namespace: cp.ID, Name: cp.Name}, kcp)
	if err != nil {
		return nil, fmt.Errorf("failed to get ControlPlane: %w", err)
	}

	return FromKubeObject(kcp), nil
}

func (k *Kube) ListControlPlanes(ctx context.Context, userID string, regions []string) ([]*ControlPlane, error) {
	cps := make([]*ControlPlane, 0)

	if len(regions) == 0 {
		regions = k.kc.GetRegions()
	}

	for _, region := range regions {
		listOpts := client.ListOptions{
			LabelSelector: labels.SelectorFromSet(labels.Set{
				labelUserID: userID,
			}),
		}

		list := &kapyv1.ControlPlaneList{}
		err := k.kc.GetClient(region).List(ctx, list, &listOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to list ControlPlanes in region %s: %w", region, err)
		}

		for _, kcp := range list.Items {
			cp := FromKubeObject(&kcp)
			cp.Region = region
			cps = append(cps, cp)
		}
	}

	return cps, nil
}

func (k *Kube) GetKubeconfig(ctx context.Context, cpID string, region string) ([]byte, error) {
	secret := &corev1.Secret{}
	err := k.kc.GetClient(region).Get(ctx, client.ObjectKey{Namespace: cpID, Name: "kubeconfig"}, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig secret: %w", err)
	}

	return secret.Data["value"], nil
}

func (k *Kube) GetRegions() []string {
	return k.kc.GetRegions()
}
