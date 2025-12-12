package main

import (
	"flag"
	"log"
	"time"

	sriovv1 "github.com/k8snetworkplumbingwg/sriov-network-operator/api/v1"
	sriovNamespaces "github.com/k8snetworkplumbingwg/sriov-network-operator/test/util/namespaces"
	"github.com/openshift-kni/k8sreporter"
	"k8s.io/apimachinery/pkg/runtime"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "the kubeconfig path")
	testName := flag.String("testname", "test", "the file name used for the report")
	reportPath := flag.String("path", "./report", "the root path used for the report")

	flag.Parse()

	// When using custom crds, we need to add them to the scheme
	addToScheme := func(s *runtime.Scheme) error {
		err := sriovv1.AddToScheme(s)
		if err != nil {
			return err
		}
		return nil
	}

	namespacesToDump := map[string]bool{
		sriovNamespaces.Test: true,
	}

	// The list of CRDs we want to dump
	crds := []k8sreporter.CRData{
		{Cr: &sriovv1.SriovNetworkNodePolicyList{}},
		{Cr: &sriovv1.SriovNetworkList{}},
		{Cr: &sriovv1.SriovNetworkNodePolicyList{}},
		{Cr: &sriovv1.SriovOperatorConfigList{}},
	}

	// The namespaces we want to dump resources for (including pods and pod logs)
	dumpNamespace := func(ns string) bool {
		_, found := namespacesToDump[ns]
		return found
	}

	reporter, err := k8sreporter.New(*kubeconfig, addToScheme, dumpNamespace, *reportPath, crds...)
	if err != nil {
		log.Fatalf("Failed to initialize the reporter %s", err)
	}
	reporter.Dump(10*time.Minute, *testName)
}
