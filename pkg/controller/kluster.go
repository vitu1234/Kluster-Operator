package controller

import (
	"log"
	"time"

	klientset "github.com/vitu1234/kluster/pkg/client/clientset/versioned"
	kinf "github.com/vitu1234/kluster/pkg/client/informers/externalversions/vitu.dev/v1alpha1"
	klister "github.com/vitu1234/kluster/pkg/client/listers/vitu.dev/v1alpha1"
	"github.com/vitu1234/kluster/pkg/do"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {

	//k8s client object
	client kubernetes.Interface

	// clientset for custom resource kluster
	klient klientset.Interface
	// check if kluster has synced
	klusterSynced cache.InformerSynced
	// lister
	kLister klister.KlusterLister
	// queue
	wq workqueue.RateLimitingInterface
}

func NewController(client kubernetes.Interface, klient klientset.Interface, klusterInformer kinf.KlusterInformer) *Controller {
	c := &Controller{
		client:        client,
		klient:        klient,
		klusterSynced: klusterInformer.Informer().HasSynced,
		kLister:       klusterInformer.Lister(),
		wq:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "kluster"),
	}

	// these functions will be triggered
	klusterInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDel,
		},
	)

	return c
}

func (c *Controller) Run(ch chan struct{}) error {
	if ok := cache.WaitForCacheSync(ch, c.klusterSynced); !ok {
		log.Printf("Error on cache sync\n")
	}

	go wait.Until(c.worker, time.Second, ch)
	<-ch

	return nil
}

func (c *Controller) worker() {
	for c.processNextItem() {

	}
}

func (c *Controller) processNextItem() bool {
	item, shutDown := c.wq.Get()
	if shutDown {
		//logs as well
		return false
	}

	defer c.wq.Forget(item)

	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		log.Printf("Error %s calling Namespace key funct on cache fro item\n", err.Error())
		return false
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		log.Printf("Splitting key into namespace and name, error %s\n", err.Error())
		return false
	}

	// get the created item from the lister
	kluster, err := c.kLister.Klusters(ns).Get(name)
	if err != nil {
		log.Printf("error %s, getting the cluster resource from the queue\n", err.Error())
	}

	log.Printf("Kluster spec that we have is %v\n", kluster.Spec)

	//call digital ocean apis
	clusterID, err := do.Create(c.client, kluster.Spec)
	if err != nil {
		log.Printf("error %s, creating native k8s clientset\n", err.Error())
	}

	log.Printf("Cluster ID: %s\n", clusterID)

	return true
}

func (c *Controller) handleAdd(obj interface{}) {
	log.Println("handleAdd was called")
	c.wq.Add(obj)
}

func (c *Controller) handleDel(obj interface{}) {
	log.Println("handleDel was called")
	c.wq.Add(obj)
}
