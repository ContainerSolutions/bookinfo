# BookInfo API

Besides the expected functional endpoints, the API will expose these endpoints for the following purposes:

| Path            | Purpose             |
|-----------------|---------------------|
| `/metrics`      | The metrics endpoint will be exposing metrics of the API to be consumed by Prometheus | 
| `/health/live`  | This endpoint will be the liveness probe's target. It returns 200 OK if the server is up and running |
| `/health/ready` | This endpoint will be the readiness probe's target. It checks the database connection and returns `200 OK` if connection is good |
| `/version`      | This endpoint simply reads the version.txt file inside the static folder. Normally version.txt will contain a token instead of a version number which can be replaced by the pipeline |
| `/docs`         | This endpoint uses redoc interface to show the OpenAPI documentation of the API |

## Pipeline
Application's pipeline is built on Azure DevOps. It checks the gitlab repository every three seconds for changes, and triggers a new build if it detects a change in the codebase.
Unit tests are automatically run by the pipeline.
When everything goes as expected (unit tests pass and build succeeded), the binary, livesettings configuration and version file (with its token replaced with build number) are copied into a new `scratch` image and pushed to [repository](https://hub.docker.com/repository/docker/serdarkalayci/bookInfoapi).
![Azure DevOps Pipeline1](img/pipeline-1.jpg) 
![Azure DevOps Pipeline2](img/pipeline-2.jpg)
## Installation - docker compose
Running the command below on a computer with docker running will run all the components, import the initial data to the database and everything will be ready to go!
```
docker compose up --build --force-recreate
```
### Mongo Express
`http://localhost:8081` address will let you access to the Mongo Express interface. The `bookInfo` database would be populated from the initial data.
![Mongo Express](img/mongo-express.jpg)

### Prometheus data on Grafana
`http://localhost:3000` address will let you access to Grafana interface. The default username and password is `admin`. It'll require that you change the password on the first run. Later on you have to define a Data Source, adding `http://prometheus:9090` will give access to the data on Prometheus, and you can define dashboards using these data.
![Prometheus data on Grafana 1](img/prometheus-1.jpg)
![Prometheus data on Grafana 2](img/prometheus-2.jpg)

### Tracing with Jaeger
`http://localhost:16686` address will let you access the Jaeger interface. After sending some requests to the API, Jaeger interface will show these requests as traces.
![Jaeger 1](img/jaeger-ui-1.jpg)
![Jaeger 2](img/jaeger-ui-2.jpg)

## Installation - kubernetes

The whole system can be run on a local [minikube](https://minikube.sigs.k8s.io/docs/start/) cluster. 
All components within the system communicate each other over the service network (eg `mongodb.bookInfo.svc.cluster.local`) within cluster. There are four endpoints which are accessible by the client.
- api
- jaeger ui
- mongo express ui
- grafana

In order to access these public endpoints from outside the cluster, [ingress plugin](https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/) must also been enabled. All the ingress definitions included are programmed to respond to on `api.bookInfo.io`, `jaeger.bookInfo.io`, `me.bookInfo.io` and `grf.bookInfo.io` addresses respectively. So these subdomains must be forwarded to minikube ip. These records can be added to `/etc/host` file by running the following commands.
```sh
echo -e "$(minikube ip)\tme.bookInfo.io" | sudo tee -a /etc/hosts
echo -e "$(minikube ip)\tjaeger.bookInfo.io" | sudo tee -a /etc/hosts
echo -e "$(minikube ip)\tapi.bookInfo.io" | sudo tee -a /etc/hosts
echo -e "$(minikube ip)\tgrf.bookInfo.io" | sudo tee -a /etc/hosts
```

All the files under k8s folder (excluding the statefulset subfolder) must be applied in order to deploy whole system on the K8s cluster. The files can be applied using the following command.
```sh
kubectl apply -f <filename>
```
### 01 - Namespace
Creates the `bookInfo` namespace. The whole system following will be deployed under `bookInfo` namespace. The default namespace to which the `kubectl` commands are sent to can be changed by using
```sh
kubectl config set-context --current --namespace=bookInfo
```

### 02 - Mongo PVC
Creates the PersistentVolumeClaim of a 1 GiB to be used as data disk of mongodb. This pvc will be mounted to the `/var/lib/mongodb/data` folder of the mongodb.

### 03 - Mongo Init
Creates a configmap consisting three files. This configmap will be mounted to the `/docker-entrypoint-initdb.d` folder, which is a special folder whose contents is run when mongodb is run for the first time.
- `bookInfo.csv` includes the initial data
- `import.sh` imports the csv file into the mongodb database
- `import.js` creates uuid for each record and updates them

### 04 - Mongo Credentials
Creates a secret of the type basic auth (which by design contains two keys, `username` and `password`). This secret will be used as environment settings by mongodb, mongo express and api deployments. This way mongodb's own password and the two other service that need to connect it will get their username and password info from one place.

### 05 - MongoDB Deployment
Creates the mongodb instance. [LADR 007 - MongoDB Deployment](LADR-007-MongoDB-Deployment.md) explains what a deployment yaml is preferred. StatefulSet yaml files can also be found under `k8s/statefulset` folder.

### 06 - MongoDB Service
Creates a service to allow the other services to access MongoDB pods using 27017 port.

### 07 - Mongo Express Deployment
Creates the Mongo Express instance. This deployment uses `mongo-credentials` secret created earlier to access the MongoDB instance throug MongoDB Service.

### 08 - Mongo Express Service
Creates a service to allow the other services to access Mongo Express pods using 8081 port.

### 09 - Mongo Express Ingress
Creates an ingress to let clients outside the cluster to access Mongo Express Service using `me.bookInfo.io` domain name.

### 10 - Jaeger All In One
The standard Jaeger deployment contains multiple components like agent, collector, query. The all-in-one image is a spesific image which contains all these components in one image for development purposes. This yaml file also contains the service records and it's prepered by Jaeger team, so I didn't see a point to disect them.
This yaml creates the jaeger-all-in-one instance and services for agent, collector, query and zipkin.

### 11 - Jaeger Ingress
Creates an ingress to let clients outside the cluster to access Jaeger Query Service using `jaeger.bookInfo.io` domain name. Sample Jaeger UI screenshot can be seen using this endpoint below:
![Jaeger-k8s](img/jaeger-k8s.jpg)

### 12 - API Livesettings
The api reads its log level from a file named `livesettings.json` under `configuration` folder. This yaml creates a configmap containing this file as the data and this configmap will be mounted to api pods' `configuration` folder to enable changing the log level by updating the configmap.
Kubernetes reflects the changes to the configmaps to the mounted volumes but it does not recycle pods for configmap changes. But as its name suggests, the application will track this spesific file by its own means and change its logging level when this configmap changes, without any downtime or a need of recycling the pods. For example when changed from `Info` to `Error` the pod's log displays the change and stops reporting `Connection to MongoDB checked successfuly!` as it does before when ready endpoint is queried.

```
2021/06/08 10:07:58 /health/live
2021/06/08 10:07:59 /health/ready
{"level":"info","time":"2021-06-08T10:07:59Z","message":"Connection to MongoDB checked successfuly!"}
2021/06/08 10:08:01 /health/live
{"level":"info","time":"2021-06-08T10:08:01Z","message":"Log Level from config: Error"}
2021/06/08 10:08:02 /health/ready
2021/06/08 10:08:04 /health/live
2021/06/08 10:08:05 /health/ready
```

### 13 - API Deployment
Creates the api instance. This deployment takes Jaeger Agent Service address `jaeger-agent.bookInfo.svc.cluster.local` and MongoDB Service address `mongodb://{username}:{password}@mongodb.bookInfo.svc.cluster.local:27017` as environment variables. The username and password needed to access MongoDB is also fed through environment variables but they are bound to mongo-credentials secret.
API Livesettings configmap is mounted to configuration folder to allow administrators change the level the api writes its logs by changing the configmap without recycling the pod.

### 14 - API Service
Creates a service to allow the other services to access API pods using 5550 port.

### 15 - API Ingress
Creates an ingress to let clients outside the cluster to access API Service using `api.bookInfo.io` domain name. Sample queries using this endpoint can be seen below:
![API-k8s-1](img/postman-k8s-1.jpg)
![API-k8s-2](img/postman-k8s-2.jpg)

### 16 - Prometheus ScrapeConfig
Prometheus does not waits for data to be written to it. Instead it queries compatible endpoints to collect metrics data. This operation is called scraping. Prometheus needs a configuration to point it to where to scrape. This configmap contains the `prometheus.yaml` file prepared for that purpose.

### 17 - Prometheus Deployment
Creates Prometheus instance. Prometheus ScrapeConfig is mounted to `/etc/prometheus/` folder to enable the pods read their configuration from configmap.

### 18 - Prometheus Service
Creates a service to allow the other services to access Prometheus pods using 9090 port.

### 19 - Grafana Deployment
Creates the Grafana instance. 

### 14 - API Service
Creates a service to allow the other services to access Grafana pods using 3000 port.

### 15 - API Ingress
Creates an ingress to let clients outside the cluster to access Grafana Service using `grf.bookInfo.io` domain name. After accessing the Grafana UI and logging in using default credentials `username: admin, password: admin`, the Prometheus instance can be added as a data source by addressing it using Prometheus Service address `http://prometheus.bookInfo.svc.cluster.local:9090` and various metrics can be displayed on Dashboards as seen below:
![Grafana-k8s-1](img/prometheus-k8s.jpg)
![Grafana-k8s-2](img/prometheus-k8s-2.jpg)