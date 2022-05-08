package main

import (
	informersappsv1 "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	listersappsv1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type controller struct {
	clientset        kubernetes.Interface
	deploymentLister listersappsv1.DeploymentLister
	deepCacheSynced  cache.InformerSynced
	queue            workqueue.RateLimitingInterface
}

func newController(clientset kubernetes.Interface, deploymentInformer informersappsv1.DeploymentInformer) *controller {
	return &controller{
		clientset:        clientset,
		deploymentLister: deploymentInformer.Lister(),
		deepCacheSynced:  deploymentInformer.Informer().HasSynced,
		queue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ekspose"),
	}

}
