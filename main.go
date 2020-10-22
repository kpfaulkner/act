package main

import (
	"flag"
	"fmt"
	"github.com/kpfaulkner/act/pkg"
	"os"
)

func main() {
	fmt.Printf("so it begins...\n")

	tenantID := flag.String("tenantid", "", "Azure Tenant ID")
	subscriptionID := flag.String("subscriptionid", "", "Azure Subscription ID")
	clientID := flag.String("clientid", "", "Azure Client ID")
	clientSecret := flag.String("clientsecret", "", "Azure Client Secret")
	cloudServiceName := flag.String("cloudservicename", "", "Azure Cloud Service Name to delete")
	deploymentSlot := flag.String("depslot", "", "Deployment slot")
	resourceGroup := flag.String("rg", "", "Resource Group")

	flag.Parse()

	if *tenantID == "" || *subscriptionID == "" || *clientSecret == "" || *clientID == "" || *cloudServiceName == "" || *deploymentSlot == "" || *resourceGroup == "" {
		fmt.Printf("All parameters are required. See --help")
		os.Exit(1)
	}

	ac := pkg.NewAzureClassic(*tenantID, *subscriptionID, *clientID, *clientSecret)

	ac.DeleteCloudServiceDeployment(*resourceGroup, *cloudServiceName, *deploymentSlot)
}
