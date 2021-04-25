package main

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/dynamic/dynamiclister"
	"time"
)

type Factory interface {
	ForResource(gvr schema.GroupVersionResource) dynamiclister.Lister
}

func New(dc dynamic.Interface) Factory{
	return &directImpl{
		dc:      dc,
		listers: map[schema.GroupVersionResource]dynamiclister.Lister{},
	}
}

func NewCached(dc dynamic.Interface, defaultResync time.Duration, stopCh <-chan struct{}) Factory{
	return &cachedImpl{
		factory: dynamicinformer.NewDynamicSharedInformerFactory(dc, defaultResync),
		stopCh:  stopCh,
		listers: map[schema.GroupVersionResource]dynamiclister.Lister{},
	}
}

func NewFilteredCached(dc dynamic.Interface, defaultResync time.Duration, namespace string, tweakListOptions dynamicinformer.TweakListOptionsFunc, stopCh <-chan struct{}) Factory{
	return &cachedImpl{
		factory: dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, defaultResync, namespace, tweakListOptions),
		stopCh:  stopCh,
		listers: map[schema.GroupVersionResource]dynamiclister.Lister{},
	}
}
