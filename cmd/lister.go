package main

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
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

	informerFactory := informers.NewSharedInformerFactory(clientset, 30*time.Second)
	podinformer := informerFactory.Core().V1().Pods()

	podinformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(new interface{}) {
			fmt.Println("add was called")
		},
		UpdateFunc: func(old, new interface{}) {
			fmt.Println("update was called")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("delte was called")
		},
	})

	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
	pod, err := podinformer.Lister().Pods("default").Get("default")
	if err != nil {
		panic(err)
	}
	fmt.Println(pod)

	// ctxBg := context.Background()
	// nodeList, err := clientset.CoreV1().Nodes().List(ctxBg, metav1.ListOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Nodes:")
	// for _, n := range nodeList.Items {
	// 	fmt.Println(n.Name)
	// }

	// podlist, err := clientset.CoreV1().Pods("default").List(ctxBg, metav1.ListOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Pods:")
	// for _, pod := range podlist.Items {
	// 	fmt.Printf("%s", pod.Name)
	// }
}
