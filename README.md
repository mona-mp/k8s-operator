# k8s-Operator

This project is a simple Kubernetes operator to deploy an application (like API) and create every object that this app needs, like service, ingress, persistentvolumeclaim, and secret.
To do this project, first read about these concepts were explained briefly below:

### What is an Operator in Kubernetes?

Operators are software extensions to Kubernetes that use [custom resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) to manage applications and their components.

### What is CRD?

The [CustomResourceDefinition](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/) API resource allows you to define custom resources. Defining a CRD object creates a new custom resource with a name and schema you specify. The Kubernetes API serves and handles the storage of your custom resource.

## Create an Operator

I divided the project into the following parts :

### Part1:  Create a Kubernetes cluster

This cluster has one master and two workers. The cluster has been initialized via  `kubeadm` and used to test and deploy this operator.
This process has been automated via the following ansible:

[HA-KubernetesCluster-Ansible](https://github.com/mona-mp/HA-K8sCluster-ansible)

### Part 2: Creating the Operator Project
There are different ways to create an operator. I would choose the framework Operator-SDK because it is easier to use, and the documentation is easy to read. The Operator SDK is a framework that uses the [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) library to make writing operators.

#### Prerequisites
The following softwares are required for creating the operator in this way:
- go to version 1.18
```bash
sudo apt update && sudo apt upgrade
sudo apt install wget software-properties-common apt-transport-HTTPS -y
wget https://golang.org/dl/go1.18.linux-amd64.tar.gz
sudo  tar -zxvf go1.18.linux-amd64.tar.gz -C /usr/local/
echo  "export PATH=/usr/local/go/bin:${PATH}"  |  sudo  tee /etc/profile.d/go.sh
source /etc/profile.d/go.sh
```

- gpg‍‍‍‍
```bash
sudo apt install gpg
```
- operator-SDK

Set platform information:
```bash
export ARCH=$(case $(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(uname -m) ;; esac)
export OS=$(uname | awk '{print tolower($0)}')
```
Download the binary for your platform:
```bash
export OPERATOR_SDK_DL_URL=https://github.com/operator-framework/operator-sdk/releases/download/v1.24.0
curl -LO ${OPERATOR_SDK_DL_URL}/operator-sdk_${OS}_${ARCH}
```
Import the operator-sdk release GPG key from  `keyserver.ubuntu.com`:
```bash
gpg --keyserver keyserver.ubuntu.com --recv-keys 052996E2A20B5C7E
```
Download the checksums file and its signature, then verify the signature:
"`bash
curl -LO ${OPERATOR_SDK_DL_URL}/checksums.txt
curl -LO ${OPERATOR_SDK_DL_URL}/checksums.txt.asc
gpg -u "Operator SDK (release) <cncf-operator-sdk@cncf.io>" --verify checksums.txt.asc
```
Make sure the checksums match:
```bash
grep operator-sdk_${OS}_${ARCH} checksums.txt | sha256sum -c -
```
The output should be like this:
```console
operator-sdk_linux_amd64: OK
```
Install the binary in the PATH:
```bash
chmod +x operator-sdk_${OS}_${ARCH} && sudo mv operator-sdk_${OS}_${ARCH} /usr/local/bin/operator-SDK
```

#### Init the project
Now it is time to use the [Operator SDK](https://sdk.operatorframework.io/) to create the project structure.
```bash
cd go/src/
mkdir k8s-operator && cd k8s-operator
operator-SDK init
```
#### Create the API and the Controller
With the below command, the API and the controller are created:
```bash
operator-SDK create --version v1alpha1 --kind Myapp --resource --controller
```
With these commands, some files create, so what each of them does?

- Makefile: Contains all the necessary commands to generate the artifacts for the operator.
- main.go: The central point of entry to the operator contains the main function.
- controllers/myapp_controller.go: The main logic of the operator goes here.
- API/v1alpha1/myapp_types.go: Contains the structure for the custom resource.
