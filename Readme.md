# Operator SDK Learning

In this project we are trying to achieve the following

1. Based on the manifest file, take the Image and Port. With those details create deployment and service
2. Based on the UTC time, do the replica scale up.

To have the details , please visit the [Wiki](https://github.com/mahboobmonnamd/operator-sdk-learning/wiki).

## How to use this in local

* Perform the sequence of commands below

```sh
make manifests
kubectl apply -f operator-src/config/crd/bases/web-app.4m.fyi_mywebapps.yaml
make run
```
or
```sh
make install run
```
> install will do the manifests, applying to cluster and run the operator.

With the above step we applied CRD and ran the controller.

* In new terminal, apply the manifest file

```sh
kubectl apply -f operator-src/config/samples/web-app_v1_mywebapp.yaml
```

* To access the application in browser using local host

```sh
kubectl port-forward svc/mywebapp-sample 8081:8080
```

This will do the portforwarding based on the service from 8080 to the host 8081. Application should alive in localhost:8081

* To list the apps under our CRD

```sh
kubectl get mywebapp.web-app.4m.fyi

kubectl describe mywebapp.web-app.4m.fyi <app-name>
```
