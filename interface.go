package main

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamiclister"
)

type Factory interface {
	ForResource(gvr schema.GroupVersionResource) dynamiclister.Lister
}
