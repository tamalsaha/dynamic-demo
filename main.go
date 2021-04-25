package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	clientscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"sigs.k8s.io/cli-utils/pkg/kstatus/polling"
	"sigs.k8s.io/cli-utils/pkg/kstatus/status"
	"sigs.k8s.io/cli-utils/pkg/object"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func main() {
	masterURL := ""
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	dc := dynamic.NewForConfigOrDie(config)
	nodes, err := dc.Resource(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "nodes",
	}).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, obj := range nodes.Items {
		s, err := status.Compute(&obj)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", *s)
	}

	disco := discovery.NewDiscoveryClientForConfigOrDie(config)
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(disco))

	reader, err := client.New(config, client.Options{
		Scheme: clientscheme.Scheme,
		Mapper: mapper,
	})
	if err != nil {
		panic(err)
	}
	poller := polling.NewStatusPoller(reader, mapper)

	ids := []object.ObjMetadata{
		{
			Namespace: "default",
			Name:      "busybox",
			GroupKind: schema.GroupKind{
				Group: "",
				Kind:  "Pod",
			},
		},
	}

	fmt.Println("----")

	ch := poller.Poll(context.TODO(), ids, polling.Options{
		PollInterval: 2 * time.Second,
		UseCache:     false,
	})
	for e := range ch {
		fmt.Printf("%s, err: %v, rs:%+v\n", e.EventType, e.Error, *e.Resource)
	}
}

