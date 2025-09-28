# Initializing Terraform with Remote State

This folder contains Terraform code for creating an AKS instance in Azure.  
Terraform state will be stored remotely in an Azure Storage Account.

## Prerequisites
- [Terraform](https://developer.hashicorp.com/terraform/downloads) installed  
- [Azure CLI](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli) installed  
- Logged in with `az login` and connected to a subscription 

## Setup

Run the bootstrap script to create the resource group and storage account:

```sh
sh bootstrap_tfstate.sh
```

## Retrieve the Storage Account Name

Check the script logs or run:

```sh
az storage account list --resource-group tfstate --output table
```

## Export the Access Key

Replace `<storage-account-name>` with the created storage account name:

```sh
ACCOUNT_KEY=$(az storage account keys list \
  --resource-group tfstate \
  --account-name <storage-account-name> \
  --query '[0].value' -o tsv)

export ARM_ACCESS_KEY=$ACCOUNT_KEY
```

The ARM_ACCESS_KEY will be used by terraform. Wherever you run the terraform plan or apply, you need this variable.

Finally make sure that the correct storage account name is added in providers.tf file. (backend, storage_account_name)

## Terraform apply
If everything else is configured, you can run

sh´´´
terraform init
terraform plan
terraform apply
´´´

When all has successfully been created, you can verify the results

sh```
resource_group_name=$(terraform output -raw resource_group_name)
```

Display the name of the created k8s cluster
sh```
az aks list --resource-group $resource_group_name --query "[].{\"K8s cluster name\":name}" --output table
```

And get the kubeconfig for kubectl
sh```
echo "$(terraform output kube_config)" > kubeconfig
```
