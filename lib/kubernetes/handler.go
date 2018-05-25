package kubernetes

/*
Copyright 2016 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Kube handle client
type Kube struct {
	podInterface corev1.PodInterface
}

// Init read kubeconfig
// in case of inside cluster
func Init() *Kube {
	var (
		config *rest.Config
		err    error
	)
	if os.Getenv("SCOPE") == "localhost" {
		var kubeconfig *string
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	} else {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	kube := Kube{
		podInterface: clientset.CoreV1().Pods(metav1.NamespaceDefault),
	}

	return &kube
}

// GetPod get engine pod
func (k *Kube) GetPod(name string) (*v1.Pod, error) {
	return k.podInterface.Get(name, metav1.GetOptions{})
}

// GetPodLogs from pod
func (k *Kube) GetPodLogs(name string) *rest.Request {
	return k.podInterface.GetLogs(name, &v1.PodLogOptions{})
}

// CreatePod create engine pod with specific version
func (k *Kube) CreatePod(version string, name string) (*v1.Pod, error) {
	return k.podInterface.Create(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.PodSpec{
			HostNetwork:   true,
			RestartPolicy: v1.RestartPolicyNever,
			Containers: []v1.Container{
				v1.Container{
					Name:  "engine",
					Image: fmt.Sprintf("masaruz/engine:%s", version),
				},
			},
		},
	})
}

// DeletePod delete engine pod with specific name
func (k *Kube) DeletePod(name string) error {
	return k.podInterface.Delete(name, &metav1.DeleteOptions{})
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
