package k8s

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/Timozer/k8swt/common"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sWrapper struct {
	Client *kubernetes.Clientset
	Config *rest.Config
}

var (
	gK8sWrapper *K8sWrapper
	gOnce       sync.Once
)

func GetClient() *K8sWrapper {
	gOnce.Do(func() {
		gK8sWrapper = &K8sWrapper{}
		gK8sWrapper.init()
	})
	return gK8sWrapper
}

func (k *K8sWrapper) init() {
	var err error

	k.Config, err = rest.InClusterConfig()
	if err != nil {
		k.Config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(common.GetHomeDir(), ".kube", "config"))
		if err != nil {
			panic(err)
		}
	}

	k.Client, err = kubernetes.NewForConfig(k.Config)
	if err != nil {
		panic(err)
	}
}

func (k *K8sWrapper) ListPods(ctx context.Context, namespace string, opts metav1.ListOptions) ([]v1.Pod, error) {
	client := k.Client.CoreV1().Pods(namespace)
	pods, err := client.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

func (k *K8sWrapper) ListNamespaces(ctx context.Context, opts metav1.ListOptions) ([]v1.Namespace, error) {
	client := k.Client.CoreV1().Namespaces()
	namespaces, err := client.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return namespaces.Items, nil
}
