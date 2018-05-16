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
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	podInterface := clientset.CoreV1().Pods(metav1.NamespaceDefault)

	r := mux.NewRouter()
	r.HandleFunc("/create/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := fmt.Sprintf("game-%s", vars["name"])
		podInterface.Create(&v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: v1.PodSpec{
				HostNetwork:   true,
				RestartPolicy: v1.RestartPolicyNever,
				Containers: []v1.Container{
					v1.Container{
						Name:  "engine",
						Image: "masaruz/engine", // TODO need to dynamic versioning
					},
				},
			},
		})
		w.Header().Set("Content-Type", "application/json")
		// In the future we could report back on the status of our DB, or our cache
		// (e.g. Redis) by performing a simple PING, and include them in the response.
		io.WriteString(w, `{"message": "Dedicated server is creating ..."}`)
	})
	go func() {
		for {
			pods, err := podInterface.List(metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

			for _, item := range pods.Items {
				fmt.Printf("Pod name: %s\n", item.Name)
			}

			time.Sleep(10 * time.Second)
		}
	}()
	log.Fatal(http.ListenAndServe(":8080", r))
}
