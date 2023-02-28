package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/homedir"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/vitu1234/kluster/pkg/apis/vitu.dev/v1alpha1"
	klient "github.com/vitu1234/kluster/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	k := v1alpha1.Kluster{}
	fmt.Print(k)

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

	clientset, err := klient.NewForConfig(config) // clientset because it is used to interact with clients from different API versions
	if err != nil {
		fmt.Printf("Error getting clientset: %s\n", err.Error())
	}

	// fmt.Println(clientset)

	klusters, err := clientset.VituV1alpha1().Klusters("").List(context.Background(), metav1.ListOptions{})

	if err != nil {
		log.Printf("Error getting klusters: %s \n", err)
	}

	fmt.Println(klusters)
}
