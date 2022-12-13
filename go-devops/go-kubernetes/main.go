package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var (
		client           *kubernetes.Clientset
		deploymentLabels map[string]string
		err              error
	)
	ctx := context.Background()

	if client, err = getClient(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	deploymentLabels, err = deploy(ctx, client)
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

func deploy(ctx context.Context, client *kubernetes.Clientset) (map[string]string, error) {
	var deployment *v1.Deployment

	appFile, err := ioutil.ReadFile("app.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read app.yaml file: %v", err)
	}

	obj, groupVersionKind, err := scheme.Codecs.UniversalDeserializer().Decode(appFile, nil, nil)

	switch obj.(type) {
	case *v1.Deployment:
		deployment = obj.(*v1.Deployment)
	default:
		return nil, fmt.Errorf("unrecognized type %v\n", groupVersionKind)
	}

	_, err = client.AppsV1().Deployments("default").Get(ctx, deployment.Name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		deploymentResponse, err := client.AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("error: deployment create failed: %v", err)
		}
		return deploymentResponse.Spec.Template.Labels, nil
	} else if err != nil && !errors.IsNotFound(err) {
		return nil, fmt.Errorf("error: deployment get failed: %v", err)
	}

	deploymentResponse, err := client.AppsV1().Deployments("default").Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("error: deployment create failed: %v", err)
	}

	return deploymentResponse.Spec.Template.Labels, nil

}
