package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Kubectl struct {
	ctx       context.Context
	clientset *kubernetes.Clientset
}

var kubeconfig *string

func NewKubectl(ctx context.Context) *Kubectl {
	kubeconfig = flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal("fail to build config from flags:", err) // TODO
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error()) // TODO
	}

	return &Kubectl{
		ctx:       ctx,
		clientset: clientset,
	}
}

func (k *Kubectl) GetPodsFromNamespace(ns string) (podsnames []string, err error) {

	pods, err := k.clientset.CoreV1().Pods(ns).List(k.ctx, metav1.ListOptions{})
	if err != nil {
		log.Println("get pods list by namespace error:", err)
		return nil, err
	}
	for _, v := range pods.Items {
		podsnames = append(podsnames, v.Name)
	}

	return
}

func (k *Kubectl) GetServiceFromNamespace(ns string) (service []string, err error) {

	services, err := k.clientset.CoreV1().Services(ns).List(k.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, v := range services.Items {
		service = append(service, v.Name)
	}

	return
}

func (k *Kubectl) GetNamespaces() (ns []string, err error) {

	namespace, err := k.clientset.CoreV1().Namespaces().List(k.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, v := range namespace.Items {
		ns = append(ns, v.Name)
	}

	return
}

func initKubeconfig() {
	kubeconfig = flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal("fail to build config from flags:", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods("backend").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, v := range pods.Items {
		fmt.Println(v.Name)
	}

}
