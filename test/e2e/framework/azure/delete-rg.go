package azure

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type DeleteResourceGroup struct {
	SubscriptionID    string
	ResourceGroupName string
	Location          string
}

func (d *DeleteResourceGroup) Run() error {
	log.Printf("deleting resource group \"%s\"...", d.ResourceGroupName)
	cred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		return fmt.Errorf("failed to obtain a credential: %w", err)
	}
	ctx := context.Background()
	clientFactory, err := armresources.NewClientFactory(d.SubscriptionID, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create resource group client: %w", err)
	}

	forceDeleteType := "Microsoft.Compute/virtualMachines,Microsoft.Compute/virtualMachineScaleSets"
	poller, err := clientFactory.NewResourceGroupsClient().BeginDelete(ctx, d.ResourceGroupName, &armresources.ResourceGroupsClientBeginDeleteOptions{ForceDeletionTypes: to.Ptr(forceDeleteType)})
	if err != nil {
		return fmt.Errorf("failed to finish the delete resource group request: %w", err)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to pull the result for delete resource group: %w", err)
	}

	log.Printf("resource group \"%s\" deleted successfully", d.ResourceGroupName)
	return nil
}

func (d *DeleteResourceGroup) Prevalidate() error {
	return nil
}

func (d *DeleteResourceGroup) Postvalidate() error {
	return nil
}
