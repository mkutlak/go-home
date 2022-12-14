package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const defaultNamespace string = "default"

func main() {
	var (
		client           *kubernetes.Clientset
		deploymentLabels map[string]string
		err              error
		expectedPods     int32
	)
	ctx := context.Background()

	if client, err = getClient(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	deploymentLabels, expectedPods, err = deploy(ctx, client)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	err = waitForPods(ctx, client, expectedPods, deploymentLabels)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	fmt.Printf("deployed finished! -- %v\n", deploymentLabels)
}

func getClient() (*kubernetes.Clientset, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func deploy(ctx context.Context, client *kubernetes.Clientset) (map[string]string, int32, error) {
	var deployment *v1.Deployment

	appFile, err := ioutil.ReadFile("app.yaml")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read app.yaml file: %v", err)
	}

	obj, groupVersionKind, _ := scheme.Codecs.UniversalDeserializer().Decode(appFile, nil, nil)

	switch obj.(type) {
	case *v1.Deployment:
		deployment = obj.(*v1.Deployment)
	default:
		return nil, 0, fmt.Errorf("unrecognized type %v", groupVersionKind)
	}

	_, err = client.AppsV1().Deployments(defaultNamespace).Get(ctx, deployment.Name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		deploymentResponse, err := client.AppsV1().Deployments(defaultNamespace).Create(ctx, deployment, metav1.CreateOptions{})

		if err != nil {
			return nil, 0, fmt.Errorf("error: deployment create failed: %v", err)
		}

		return deploymentResponse.Spec.Template.Labels, 0, nil

	} else if err != nil && !errors.IsNotFound(err) {
		return nil, 0, fmt.Errorf("error: deployment get failed: %v", err)
	}

	deploymentResponse, err := client.AppsV1().Deployments("default").Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return nil, 0, fmt.Errorf("error: deployment create failed: %v", err)
	}

	return deploymentResponse.Spec.Template.Labels, *deploymentResponse.Spec.Replicas, nil

}

func waitForPods(ctx context.Context, client *kubernetes.Clientset, expectedPods int32, deploymentLabels map[string]string) error {
	for {
		validatedLabels, err := labels.ValidatedSelectorFromSet(deploymentLabels)
		if err != nil {
			return fmt.Errorf("ValidatedSelectorFromSet error: %v", err)
		}

		podList, err := client.CoreV1().Pods(defaultNamespace).List(ctx, metav1.ListOptions{
			LabelSelector: validatedLabels.String(),
		})

		if err != nil {
			return fmt.Errorf("pods list error: %v", err)
		}

		podsRunning := 0
		for _, pod := range podList.Items {
			if pod.Status.Phase == "Running" {
				podsRunning++
			}
		}

		fmt.Printf("Waiting for pods to become ready (running %d / %d)\n", podsRunning, expectedPods)

		if podsRunning > 0 && podsRunning == len(podList.Items) && podsRunning == int(expectedPods) {
			break
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}
