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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
	parrotyschema "github.com/supersingh05/parroty/pkg/schema"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// cmd.Execute()
	var yamlPath string
	flag.StringVar(&yamlPath, "config", "", "path to spec config")
	flag.Parse()
	p := parrotyschema.Parroty{}
	yamlFile, err := ioutil.ReadFile(yamlPath)
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
	resp := parrotyschema.Response{}

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
			fmt.Println(err)
			fmt.Println("error in dyrnamic client")
		}
		if len(s.AwsAccessKey) != 0 {
			os.Setenv("AWS_SESSION_TOKEN", s.AwsSessionToken)
			os.Setenv("AWS_ACCESS_KEY_ID", s.AwsAccessKey)
			os.Setenv("AWS_SECRET_ACCESS_KEY", s.AwsSecretKey)
			os.Setenv("AWS_SECURITY_TOKEN", s.AwsSecurityToken)
		}

		cr := parrotyschema.ClusterResponse{}
		cr.Name = s.Name
		cr.Type = s.Cloud
		for _, t := range s.ClusterExpect {
			gvk := schema.GroupKind{}
			gvk.Group = t.Group
			gvk.Kind = t.Kind
			mapping, err := mapper.RESTMapping(gvk, t.Version)
			if err != nil {
				fmt.Println("mapping didnt work")
				continue
			}
			var x dynamic.ResourceInterface
			if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
				x = dynClient.Resource(mapping.Resource).Namespace(t.Namespace)
			} else {
				x = dynClient.Resource(mapping.Resource)
			}
			_, err = x.Get(context.TODO(), t.ObjectName, metav1.GetOptions{})
			if err != nil {
				cr.AddCheck(t, false)
			} else {
				cr.AddCheck(t, true)
			}
		}
		sort.Slice(cr.Checks, func(i, j int) bool { return cr.Checks[i].Namespace < cr.Checks[j].Namespace })
		resp.AddClusterResponse(cr)
	}
	printResponse(resp)
}

func printResponse(resp parrotyschema.Response) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Cluster Name", "Resource", "Namespace", "Name", "Exists"})
	for _, cr := range resp.ClusterResponses {
		fmt.Println("For cluster: " + cr.Name)
		for _, check := range cr.Checks {
			t.AppendRow(table.Row{cr.Name, check.Group + "/" + check.Version + " " + check.Kind, check.Namespace, check.ObjectName, check.Passed})
		}
		t.AppendSeparator()
	}
	t.Render()
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
