# BookInfo API

⚠️ WARNING: This file explains steps to be taken to complete the lab. Do not read this file until you have attempted the lab yourself.

## Step 6: Check the single book information
We have to make sure the other endpoints are working as well. Let's try to get the information of a single book.
```bash
curl http://bookinfo.localdev.me/book/5f1b9b9b9b9b9b9b9b9b9b9b
```
```bash
{"uuid":"5f1b9b9b9b9b9b9b9b9b9b9b","name":"Lord of the Rings: The Fellowship of the Ring","author":"J.R.R. Tolkien","currentStock":234}
```
It all works as expected, which means stockAPI also works. We can now move on to the next step. In which we will put the services under some load. For this we'll be using the `hey` tool. `hey` is a load testing tool written in Go. It is very easy to use and displays the results in a very readable format. 
```bash
hey -z 30s -c 50 -m GET http://boookinfo.localdev.me/book/636bd189db7a3afe9ac84af5
```
The status code distrubution shows that all oour requests returned successfully
```bash
Status code distribution:
  [200]	3808 responses
```
but the latency distribution shows that we have some problems with the response time.
```bash
Response time histogram:
  0.002 [1]	    |
  0.711 [3492]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  1.421 [0]	    |
  2.130 [0]   	|
  2.840 [0]	    |
  3.549 [59]   	|■
  4.259 [55]  	|■
  4.968 [0]   	|
  5.678 [76]  	|■
  6.387 [63]  	|■
  7.096 [62]    |■
```
We can clearly see that some requests took more than 3 seconds, up to 7 seconds to complete which is not acceptable. Let's try to find out what's going on.
First lets check the `/book` endpoint, because we know that this endpoint does not query the stockAPI, so we can rule out the stockAPI as the cause of the problem.
```bash
Response time histogram:
  0.001 [1]	    |
  0.034 [13319]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.067 [17224]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.100 [3412]	|■■■■■■■■
  0.133 [612]	  |■
  0.166 [160]	  |
  0.199 [58]	  |
  0.232 [37]	  |
  0.264 [15]	  |
  0.297 [5]	    |
  0.330 [1]	    |
```
We can clearly see that the response time is very low, which means that the problem is not in the infoAPI. We have to check the stockAPI, but it does not have an ingress on its own.
It's good to see these information via the load testing tool, but this also means that our customers are also suffering from this problem and we were not aware of it until we run our load testing tool. We need to find a way to monitor our services and get notified when something goes wrong.  

## Step 7: Add Prometheus
The startup logs of our application clearly shows that it has an endpoint called `/metrics` which exposes metrics compatible with [Open Metrics](https://openmetrics.io/) which is is a sandbox project of [Cloud Native Computing Foundation](https://www.cncf.io/). [Prometheus](https://prometheus.io/), which is a graduate project of CNCF can consume this endpoint and store the metrics in a time series database. [Grafana](https://grafana.com/) is a tool that can be used to visualize the metrics stored in Prometheus. We will use these tools to monitor our services and get notified when something goes wrong.
We can deploy Prometheus using Prometheus Operator with the following command.
```bash
kubectl create -f https://raw.githubusercontent.com/ContainerSolutions/bookinfo/main/k8s/final/prom-bundle/prom-bundle.yaml
```
It can take a few minutes for the operator to be up and running. We can check for completion with the following command:
```bash
kubectl wait --for=condition=Ready pods -l  app.kubernetes.io/name=prometheus-operator -n default
```
Prometheus needs some rights to access resources on our namespace. By default it uses the `ServiceAccount` named `prometheus`. If it's not created yet, we can create it with the following command:
```bash
kubectl create serviceaccount prometheus
```
If the needed `ClusterRole` is not created yet, we can create it by applyin the following manifest:
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: prometheus
rules:
- apiGroups: [""]
  resources:
  - nodes
  - nodes/metrics
  - services
  - endpoints
  - pods
  verbs: ["get", "list", "watch"]
- apiGroups: [""]
  resources:
  - configmaps
  verbs: ["get"]
- apiGroups:
  - networking.k8s.io
  resources:
  - ingresses
  verbs: ["get", "list", "watch"]
- nonResourceURLs: ["/metrics"]
  verbs: ["get"]
```
Then add a `ClusterRoleBinding` to bind this `ClusterRole` to the `ServiceAccount` for our namespace by applying the following manifest:
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: prometheus
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: prometheus
subjects:
- kind: ServiceAccount
  name: prometheus
  namespace: bookinfo
```
We have to create a Prometheus instance to scrape the metrics from our services. We can do this by applying the following manifest:
```yaml
apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: prometheus
spec:
  serviceAccountName: prometheus
  serviceMonitorSelector:
    matchLabels:
      layer: api
  resources:
    requests:
      memory: 100Mi
  enableAdminAPI: false
```
After it's up and running, we should add `ServiceMonitor` objects to tell Prometheus to scrape the metrics from our services. We can do this by applying the following manifest:
```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: svcmon-infoapi
  labels:
    app: infoapi
    layer: api
spec:
  selector:
    matchLabels:
      app: infoapi
      layer: api
  endpoints:
  - port: web
---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: svcmon-stockapi
  labels:
    app: stockapi
    layer: api
spec:
  selector:
    matchLabels:
      app: stockapi
      layer: api
  endpoints:
  - port: web
```
## Step 8: Add Grafana
We can deploy Grafana, create a service and and ingress using the following command:
```bash
kubectl create -f https://raw.githubusercontent.com/ContainerSolutions/bookinfo/main/k8s/final/20-grafana.yaml
```
Now we can access to the grafana interface by navigating to `http://grf.localdev.me` and login with the default credentials `admin:admin`. We can add a Prometheus data source by clicking on the `Add data source` button and filling the `URL` field as `http://prometheus-operated:9090` which ic the service address of Prometheus instance. We can then click on the `Save & Test` button to test the connection and use this connection to create our dashboards.


[To continue](SOLUTION-3.md)