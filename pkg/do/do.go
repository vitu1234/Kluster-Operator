package do

import (
	"context"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/vitu1234/kluster/pkg/apis/vitu.dev/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Create(c kubernetes.Interface, spec v1alpha1.KlusterSpec) (string, error) {
	token, err := getToken(c, spec.TokenSecret)
	if err != nil {
		log.Printf("Error getting secret: %s\n", err.Error())
		return "", err
	}

	client := godo.NewFromToken(token)
	// fmt.Println(client)
	//query k8s apis and get the secret name

	createRequest := &godo.KubernetesClusterCreateRequest{
		Name:        spec.Name,
		RegionSlug:  spec.Region,
		VersionSlug: spec.Version,
		NodePools: []*godo.KubernetesNodePoolCreateRequest{
			&godo.KubernetesNodePoolCreateRequest{
				Name:  spec.NodePools[0].Name,
				Size:  spec.NodePools[0].Size,
				Count: spec.NodePools[0].Count,
			},
		},
	}

	cluster, _, err := client.Kubernetes.Create(context.Background(), createRequest)

	return cluster.ID, nil
}

func getToken(client kubernetes.Interface, sec string) (string, error) {

	namespace := strings.Split(sec, "/")[0]
	name := strings.Split(sec, "/")[1]

	//call k8s api server
	s, err := client.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})

	if err != nil {
		return "", nil
	}

	return string(s.Data["token"]), nil
}
