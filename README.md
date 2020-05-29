# Validating Admission Webhook for Service Binding

A simple admission webhook which rejects every `ServiceBindingRequest`

## Deploying it!

1. Run `make` to build the image. Run `make push` to push it your registry. By default, a pre-built image `quay.io/shbose/service-binding-admission-controller:v0.2` will be used. 

2. Run `./deploy.sh`. This will create a CA, a certificate and private key for the webhook server,
and deploy the resources in the newly created `service-binding-webhook` namespace in your Kubernetes cluster.


## Verify

1. The `webhook-server` pod in the `service-binding-webhook` namespace should be running:
```
$ kubectl -n service-binding-webhook get pods
NAME                             READY     STATUS    RESTARTS   AGE
webhook-server-6f976f7bf-hssc9   1/1       Running   0          35m
```

2. A `ValidatingWebhookConfiguration` named `sbr-webhook` should exist:
```
$ kubectl get validatingwebhookconfigurations
NAME           AGE
sbr-webhook   36m
```
3. Try creating any `ServiceBindingRequest` and find the request rejected!

4. IMPORTANT: Ensure you clean up after you are done, else `ServiceBindingRequest` creation will be blocked for eternity!

```
$ oc delete project service-binding-webhook
$ oc delete validatingwebhookconfiguration sbr-webhook
```

