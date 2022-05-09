package main

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	informersappsv1 "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	listersappsv1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type controller struct {
	clientset        kubernetes.Interface
	deploymentLister listersappsv1.DeploymentLister
	depCacheSynced   cache.InformerSynced
	queue            workqueue.RateLimitingInterface
}

func newController(clientset kubernetes.Interface, deploymentInformer informersappsv1.DeploymentInformer) *controller {
	c := &controller{
		clientset:        clientset,
		deploymentLister: deploymentInformer.Lister(),
		depCacheSynced:   deploymentInformer.Informer().HasSynced,
		queue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ekspose"),
	}
	deploymentInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    handleAdd,
			DeleteFunc: handleDel,
		},
	)
	return c

}

func (c *controller) run(ch <-chan struct{}) {
	if !cache.WaitForCacheSync(ch, c.depCacheSynced) {
		fmt.Println("waiting for cache to be synced")
	}
	go wait.Until(c.worker, 1*time.Second, ch)
	<-ch
}

func (c *controller) worker() {

}

func handleAdd(obj interface{}) {
	fmt.Println("add was called")
}

func handleDel(obj interface{}) {
	fmt.Println("del was called")
}
