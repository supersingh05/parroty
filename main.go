/*
Copyright © 2020 Asavir Kalla kalla.asavir@gmail.com

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
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	parrotyschema "github.com/supersingh05/parroty/pkg/schema"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// cmd.Execute()

	p := parrotyschema.Parroty{}
	yamlFile, err := ioutil.ReadFile("./example.yaml")
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return
	}
	// fmt.Println(string(yamlFile))
	err = yaml.Unmarshal(yamlFile, &p)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	modifyExpects(&p)
	for _, s := range p.Clusters {
		config, err := setupKubeClient(s.KubeconfigPath, s.Context)
		if err != nil {
			fmt.Println("failing in config")
			// log.Fatal(err)
		}

		dc, err := discovery.NewDiscoveryClientForConfig(config)
		mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))
		dynClient, errClient := dynamic.NewForConfig(config)
		if errClient != nil {
			fmt.Println("error in dyrnamic client")
		}
		for _, t := range s.ClusterExpect {
			gvk := schema.GroupKind{}
			gvk.Group = t.Group
			gvk.Kind = t.Kind
			mapping, err := mapper.RESTMapping(gvk, t.Version)
			if err != nil {
				fmt.Println("mapping didnt work")
				log.Fatal(err)
			}
			var x dynamic.ResourceInterface
			if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
				x = dynClient.Resource(mapping.Resource).Namespace(t.Namespace)
			} else {
				x = dynClient.Resource(mapping.Resource)
			}
			_, err = x.Get(context.TODO(), t.ObjectName, metav1.GetOptions{})
			if err != nil {
				fmt.Println("Cluster: " + s.Name + " does NOT have resource: " + t.ObjectName + " " + t.Group + "/" + t.Version + " " + t.Kind + " in namespace: " + t.Namespace)
			} else {
				fmt.Println("Cluster: " + s.Name + " has resource: " + t.ObjectName + " " + t.Group + "/" + t.Version + " " + t.Kind + " in namespace: " + t.Namespace)
			}
		}
	}

}

func modifyExpects(parroty *parrotyschema.Parroty) {
	for i, s := range parroty.Clusters {
		parroty.Clusters[i].ClusterExpect = append(s.ClusterExpect, parroty.GlobalExpect...)
	}
}

func setupKubeClient(kubeConfig, context string) (*rest.Config, error) {
	rules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig}
	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
	overrides.CurrentContext = context

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()

	if err != nil {
		return nil, err
	}
	return config, nil
}
