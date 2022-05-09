package main

import (
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	defaultTimeout := 30 * time.Second

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

	ch := make(chan struct{})
	informerFactory := informers.NewSharedInformerFactory(clientset, defaultTimeout)
	c := newController(clientset, informerFactory.Apps().V1().Deployments())
	informerFactory.Start(ch)
	c.run(ch)
}
