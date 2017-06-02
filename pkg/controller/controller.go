/*
Copyright 2016 Skippbox, Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jpisaac/alert-watch/config"
	"github.com/jpisaac/alert-watch/pkg/handlers"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var log = logrus.New()

func Controller(conf *config.Config, eventHandler handlers.Handler) {

	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)

	if conf.Resource.Pod {
		var podsStore cache.Store
		podsStore = watchPods(kubeClient.Core().RESTClient(), podsStore, eventHandler)
	}

	if conf.Resource.Services {
		var servicesStore cache.Store
		servicesStore = watchServices(kubeClient.Core().RESTClient(), servicesStore, eventHandler)
	}

	if conf.Resource.ReplicationController {
		var rcStore cache.Store
		rcStore = watchReplicationControllers(kubeClient.Core().RESTClient(), rcStore, eventHandler)
	}

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func watchPods(client rest.Interface, store cache.Store, eventHandler handlers.Handler) cache.Store {
	//Define what we want to look for (Pods)
	watchlist := cache.NewListWatchFromClient(client, "pods", v1.NamespaceAll, fields.Everything())

	resyncPeriod := 30 * time.Minute

	//Setup an informer to call functions when the watchlist changes
	eStore, eController := cache.NewInformer(
		watchlist,
		&v1.Pod{},
		resyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    eventHandler.ObjectCreated,
			DeleteFunc: eventHandler.ObjectDeleted,
		},
	)

	//Run the controller as a goroutine
	go eController.Run(wait.NeverStop)

	return eStore
}

func watchServices(client rest.Interface, store cache.Store, eventHandler handlers.Handler) cache.Store {
	//Define what we want to look for (Services)
	watchlist := cache.NewListWatchFromClient(client, "services", v1.NamespaceAll, fields.Everything())

	resyncPeriod := 30 * time.Minute

	//Setup an informer to call functions when the watchlist changes
	eStore, eController := cache.NewInformer(
		watchlist,
		&v1.Service{},
		resyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    eventHandler.ObjectCreated,
			DeleteFunc: eventHandler.ObjectDeleted,
			UpdateFunc: eventHandler.ObjectUpdated,
		},
	)

	//Run the controller as a goroutine
	go eController.Run(wait.NeverStop)

	return eStore
}

func watchReplicationControllers(client rest.Interface, store cache.Store, eventHandler handlers.Handler) cache.Store {
	//Define what we want to look for (ReplicationControllers)
	watchlist := cache.NewListWatchFromClient(client, "replicationcontrollers", api.NamespaceAll, fields.Everything())

	resyncPeriod := 30 * time.Minute

	//Setup an informer to call functions when the watchlist changes
	eStore, eController := cache.NewInformer(
		watchlist,
		&api.ReplicationController{},
		resyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    eventHandler.ObjectCreated,
			DeleteFunc: eventHandler.ObjectDeleted,
		},
	)

	//Run the controller as a goroutine
	go eController.Run(wait.NeverStop)

	return eStore
}
