package main

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDel,
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
	for c.ProcessItem() {

	}
}

func (c *controller) ProcessItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		fmt.Print("Error getting key from cache %s\n", err.Error())
	}
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		fmt.Printf("Error splitting key into namespace and name $s\n", err.Error())
	}
	c.syncDeployment(ns, name)

	if err != nil {
		//potentiually retry here (re-queue)
		fmt.Printf("Error syncing deployment %s\n", err.Error())
		return false
	}

	return true
}

func (c *controller) syncDeployment(ns, name string) error {
	//create service
	ctx := context.Background()
	dep, err := c.deploymentLister.Deployments(ns).Get(name)
	if err != nil {
		fmt.Printf("Error geteting deployment from lister %s\n", err.Error())
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dep.Name,
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
			},
		},
	}
	_, err = c.clientset.CoreV1().Services(ns).Create(ctx, &svc, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error create service %s\n", err.Error())
		return err
	}
	fmt.Printf("Created service %s\n", svc.Name)
	//create ingress

	return nil
}

func (c *controller) handleAdd(obj interface{}) {
	fmt.Println("add was called")
	c.queue.Add(obj)
}

func (c *controller) handleDel(obj interface{}) {
	fmt.Println("del was called")
	c.queue.Add(obj)
}
