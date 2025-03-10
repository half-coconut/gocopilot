package utils

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"strings"
)

type ClientGo struct {
	ClientSet       *kubernetes.Clientset
	DynamicClient   dynamic.Interface
	DiscoveryClient discovery.DiscoveryInterface
}

func NewClientGo(kubeconfig string) (*ClientGo, error) {
	// ~/.kube/config
	if strings.HasPrefix(kubeconfig, "~") {
		homeDir := homedir.HomeDir()
		kubeconfig = filepath.Join(homeDir, kubeconfig[1:])
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}
	return &ClientGo{
		ClientSet:       clientset,
		DynamicClient:   dynamicClient,
		DiscoveryClient: discoveryClient,
	}, nil
}
