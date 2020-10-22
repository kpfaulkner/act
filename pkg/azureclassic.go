package pkg

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/kpfaulkner/azureauth"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type AzureClassic struct {
	subscriptionID string
	tenantID       string
	clientID       string
	clientSecret   string
	azureAuth      *azureauth.AzureAuth
}

// NewAzureClassic. Works against a single Subscription.
// Passing in ServicePrincipal for auth.
func NewAzureClassic(tenantID string, subscriptionID string, clientID string, clientSecret string) *AzureClassic {
	ac := AzureClassic{}
	ac.clientID = clientID
	ac.clientSecret = clientSecret
	ac.subscriptionID = subscriptionID
	ac.tenantID = tenantID
	ac.azureAuth = azureauth.NewAzureAuth(subscriptionID, tenantID, clientID, clientSecret)
	return &ac
}

func (ac *AzureClassic) DeleteCloudServiceDeployment(resourceGroup string, cloudServiceName string, deploymentSlot string) error {

	// refresh all the tokens!!!
	err := ac.azureAuth.RefreshToken()
	if err != nil {
		return err
	}

	url := generateDeleteCloudServiceURL(ac.subscriptionID, resourceGroup, cloudServiceName, deploymentSlot)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Renegotiation:      tls.RenegotiateOnceAsClient,
			InsecureSkipVerify: true},
	}

	timeout := time.Duration(20 * time.Second)

	client := &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}

	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", "Bearer "+ac.azureAuth.CurrentToken().AccessToken)
	req.Header.Add("Content-type", "application/json")
	//req.Header.Add("x-ms-version", "2017-03-30")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error on put %s\n", err.Error())
		panic(err)
	}

	fmt.Printf("status code is %d\n", resp.StatusCode)
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("body is %s\n", string(b))

	// if status begins with 4.... assume failure.
	if strings.HasPrefix(resp.Status, "4") {
		return errors.New("unable to delete cloud service instance")
	}

	return nil
}

func generateDeleteCloudServiceURL(subscriptionID string, resourceGroup string, cloudServiceName string, deploymentSlot string) string {
	//template := "https://management.core.windows.net/%s/services/hostedservices/%s/deploymentslots/%s"
	template := "https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ClassicCompute/domainNames/%s/slots/%s?api-version=2015-06-01"
	url := fmt.Sprintf(template, subscriptionID, resourceGroup, cloudServiceName, deploymentSlot)
	return url
}
