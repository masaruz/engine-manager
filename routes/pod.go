package routes

import (
	"bytes"
	"encoding/json"
	"engine-manager/lib/kubernetes"
	"engine-manager/model"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// GetPod get a pod in kube
func GetPod(kube *kubernetes.Kube) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		name := vars["name"]
		req := kube.GetPodLogs(name)
		readCloser, err := req.Stream()
		if err != nil {
			model.Reply(w, model.Response{
				StatusCode: model.StatusGetPodLogsFail,
				Message:    err.Error(),
			})
			return
		}
		defer readCloser.Close()
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, readCloser)
		if err != nil {
			model.Reply(w, model.Response{
				StatusCode: model.StatusGetPodLogsFail,
				Message:    err.Error(),
			})
			return
		}
		logs := strings.Split(string(buf.Bytes()), "\n")

		pod, err := kube.GetPod(name)
		if err != nil {
			model.Reply(w, model.Response{
				StatusCode: model.StatusGetPodFail,
				Message:    err.Error(),
			})
			return
		}
		model.Reply(w, model.Response{
			StatusCode: model.StatusOK,
			Logs:       logs,
			Message:    fmt.Sprintf("This pod run on node: %s", pod.Spec.NodeName),
		})
	}
}

// GetPods list all pods in kube
func GetPods(kube *kubernetes.Kube) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

// CreatePod create pod in kube
func CreatePod(kube *kubernetes.Kube) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		name := vars["name"]
		kube.CreatePod("0.0.7", name)
		resp := model.Response{
			StatusCode: model.StatusOK,
			Message:    "Success",
		}
		j, _ := json.Marshal(resp)
		w.Write(j)
	}
}

// DeletePod delete pod in kube
func DeletePod(kube *kubernetes.Kube) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		name := vars["name"]
		kube.DeletePod(name)
		resp := model.Response{
			StatusCode: model.StatusOK,
			Message:    "Success",
		}
		j, _ := json.Marshal(resp)
		w.Write(j)
	}
}
