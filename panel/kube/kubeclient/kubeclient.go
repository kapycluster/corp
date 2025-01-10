package kubeclient

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"kapycluster.com/corp/panel/config"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KubeClient struct {
	clients    map[string]client.Client
	clientsets map[string]*kubernetes.Clientset
}

func New(c *config.Config) (*KubeClient, error) {
	k := &KubeClient{
		clients:    make(map[string]client.Client),
		clientsets: make(map[string]*kubernetes.Clientset),
	}

	files, err := os.ReadDir(c.Kubernetes.KubeconfigsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig directory: %w", err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		filename := f.Name()
		name := strings.TrimSuffix(filename, filepath.Ext(filename))
		kubeconfigPath := filepath.Join(c.Kubernetes.KubeconfigsDir, filename)

		restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to build rest config for region %s: %w", filename, err)
		}

		client, err := client.New(restConfig, client.Options{})
		if err != nil {
			return nil, fmt.Errorf("failed to create client for region %s: %w", filename, err)
		}

		clientset, err := kubernetes.NewForConfig(restConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create clientset for region %s: %w", filename, err)
		}

		k.clients[name] = client
		k.clientsets[name] = clientset
	}

	return k, nil
}

func (k *KubeClient) GetClient(region string) client.Client {
	return k.clients[region]
}

func (k *KubeClient) GetClientset(region string) *kubernetes.Clientset {
	return k.clientsets[region]
}

func (k *KubeClient) ValidateRegion(region string) bool {
	_, ok := k.clients[region]
	return ok
}

func (k *KubeClient) GetRegions() []string {
	regions := make([]string, 0, len(k.clients))
	for region := range k.clients {
		regions = append(regions, region)
	}
	return regions
}
