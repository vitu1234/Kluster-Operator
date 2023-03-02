package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/homedir"

	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	klient "github.com/vitu1234/kluster/pkg/client/clientset/versioned"
	kinfFac "github.com/vitu1234/kluster/pkg/client/informers/externalversions"
	controller "github.com/vitu1234/kluster/pkg/controller"
	// "sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	// k := v1alpha1.Kluster{}
	// fmt.Print(k)

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		//if failed to find kubeconfig location because the code is now running in cluster, do the following
		fmt.Println(err)
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("Error getting incluster config: %s\n", err.Error())
		}

	}

	klientset, err := klient.NewForConfig(config) // clientset because it is used to interact with clients from different API versions
	if err != nil {
		fmt.Printf("Error getting klientset: %s\n", err.Error())
	}

	//clientset for k8s native resources

	clientset, err := kubernetes.NewForConfig(config) // clientset because it is used to interact with clients from different API versions
	if err != nil {
		fmt.Printf("Error getting standard clientset: %s\n", err.Error())
	}

	// fmt.Println(clientset)

	// klusters, err := klientset.VituV1alpha1().Klusters("").List(context.Background(), metav1.ListOptions{})

	// if err != nil {
	// 	log.Printf("sError getting klusters: %s \n", err)
	// }

	// fmt.Println(klusters)

	infoFactory := kinfFac.NewSharedInformerFactory(klientset, 20*time.Minute)

	ch := make(chan struct{})
	c := controller.NewController(clientset, klientset, infoFactory.Vitu().V1alpha1().Klusters())

	infoFactory.Start(ch)
	if err := c.Run(ch); err != nil {
		log.Printf("Error running controller: %s\n", err.Error())
	}
}
