# Controller without using kubebuilder or sdk

- We need yaml file to create our CR (Custom Resource).
- Kuberenetes should be aware of its definiton so we need CRD (Custom Resource Definition).

> Without CRD if we apply the CR, kubectl will throw the error. "No match found for the Kind".


## Proxy k8s server to access it locally

This will proxy the server to 8000. localhost:8000 will list all the resourses. After we applied our crd we can able to see our **Group** in the list.

```sh
kubectl proxy --port 8000
```

To see the CRD defintion in the browser instead of terminal
```curl
http://localhost:8000/apis/k8s.4m.fyi/v1
http://localhost:8000/apis/k8s.4m.fyi/v1/namespaces/default/helloworld
```

In the second url, we can see the resources created with the CRD in the **items** array.

Apply the CR yaml file now. It will create the resource based on the CRD.

```sh
kubectl get helloworld
kubectl describe helloworld
```