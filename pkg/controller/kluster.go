package controller

import (
	"log"
	"time"

	klientset "github.com/vitu1234/kluster/pkg/client/clientset/versioned"
	kinf "github.com/vitu1234/kluster/pkg/client/informers/externalversions/vitu.dev/v1alpha1"
	klister "github.com/vitu1234/kluster/pkg/client/listers/vitu.dev/v1alpha1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	// clientset for custom resource kluster
	klient klientset.Interface
	// check if kluster has synced
	klusterSynced cache.InformerSynced
	// lister
	kLister klister.KlusterLister
	// queue
	wq workqueue.RateLimitingInterface
}

func NewController(klient klientset.Interface, klusterInformer kinf.KlusterInformer) *Controller {
	c := &Controller{
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

	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		log.Printf("Error %s calling Namespace key funct on cache fro item\n", err.Error())
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		log.Printf("Splitting key into namespace and name, error %s\n", err.Error())
	}

	// get the created item from the lister
	kluster, err := c.kLister.Klusters(ns).Get(name)
	if err != nil {
		log.Printf("error %s, getting the cluster resource from the queue\n", err.Error())
	}

	log.Printf("Kluster spec that we have is %v\n", kluster.Spec)

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
