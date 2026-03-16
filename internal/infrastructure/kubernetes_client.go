package infrastructure

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubernetesClient wraps the Kubernetes client
type KubernetesClient struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

// NewKubernetesClient creates a new Kubernetes client
func NewKubernetesClient(config *Config) (*KubernetesClient, error) {
	var restConfig *rest.Config
	var err error

	if config.InCluster {
		// Use in-cluster configuration
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create in-cluster config: %w", err)
		}
	} else {
		// Use kubeconfig file
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		loadingRules.ExplicitPath = config.KubeConfig

		configOverrides := &clientcmd.ConfigOverrides{}
		if config.KubeContext != "" {
			configOverrides.CurrentContext = config.KubeContext
		}

		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			loadingRules,
			configOverrides,
		)

		restConfig, err = kubeConfig.ClientConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create kubeconfig: %w", err)
		}
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	return &KubernetesClient{
		clientset: clientset,
		config:    restConfig,
	}, nil
}

// GetClientset returns the Kubernetes clientset
func (k *KubernetesClient) GetClientset() *kubernetes.Clientset {
	return k.clientset
}

// GetConfig returns the REST config
func (k *KubernetesClient) GetConfig() *rest.Config {
	return k.config
}

// Made with Bob
