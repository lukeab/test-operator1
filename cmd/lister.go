package lister

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}
	clientset := kubernetes.NewForConfigOrDie(config)
	ctxBg := context.Background()
	// nodeList, err := clientset.CoreV1().Nodes().List(ctxBg, metav1.ListOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Nodes:")
	// for _, n := range nodeList.Items {
	// 	fmt.Println(n.Name)
	// }

	podlist, err := clientset.CoreV1().Pods("default").List(ctxBg, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Pods:")
	for _, pod := range podlist.Items {
		fmt.Printf("%s", pod.Name)
	}
}
