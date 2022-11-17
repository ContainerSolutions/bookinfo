# BookInfo API

⚠️ WARNING: This file explains steps to be taken to complete the lab. Do not read this file until you have attempted the lab yourself.

## Step 1: Try to find out what's wrong with bookInfoAPI
Let's start by checking the logs of the bookInfoAPI pod. We can do this by running the following command:

```bash
kubectl logs -l app=infoapi
```
```bash
kubectl logs -l app=infoapi
{"level":"info","time":"2022-11-15T10:52:37Z","message":"Log Level from config: Fatal"}
2022/11/15 10:52:47 Initializing logging reporter
```
Even the logs do not show us any errors, we can see that the log level is set to `Fatal`. This may be too high for our needs. Let's try to lower it to `Error` and see if we get any more information. Also the log message reveals another information, that the log level is set from the config file. Let's first find that config file.

## Step 2: Find the config file
```bash
kubectl get pod -l app=infoapi -o yaml
```
We can easily find that a config file is mounted to the pod at `/go/bin/configuration` folder.
```yaml
...
      volumeMounts:
      - mountPath: /go/bin/configuration
        name: api-livesettings
...
```
and a bit below we can see that volume is actually a ConfigMap named `api-livesettings`.
```yaml
...
    volumes:
    - configMap:
        defaultMode: 420
        name: api-livesettings
      name: api-livesettings
...
```
Let's check the content of that ConfigMap.
```bash
kubectl get configmap api-livesettings -o yaml
```
```yaml
apiVersion: v1
data:
  livesettings.json: |-
    {
        "Logging": {
            "LogLevel": {
            "Default": "Fatal"
            }
        }
    }
kind: ConfigMap
metadata:
  name: api-livesettings
  namespace: bookinfo
```
## Step 3: Edit the config file
Here, we can see that the log level is set to `Fatal`. Let's change it to `Error` or even better to `Info` by editing this resource and see if we get any more information.
```bash
kubectl edit configmap api-livesettings
```
Now if we do the same request to the API by running 
```bash
curl http://bookinfo.localdev.me/book
``` 
and check the logs by running, 
```bash
kubectl logs -l app=infoapi
`
we can now see that we get more information about the error.
```bash
2022/11/15 12:58:30 /book
{"level":"error","error":"server selection error: server selection timeout, current topology: { Type: Unknown, Servers: [{ Addr: mongodb.svc:27017, Type: Unknown, Average RTT: 0, Last error: connection() error occured during connection handshake: dial tcp: lookup mongodb.svc on 10.43.0.10:53: no such host }, ] }","time":"2022-11-15T12:59:00Z","message":"Error getting BookInfos"}
```
## Step 4: This shouldn't be happening on Kubernetes
There's obviously something wrong with the connection to the database. It might be the related the MongoDB itself, but it might also be related to the way the connection is configured. But before jumping to the solution, we first make sure this won't happen again. A pod that's not working as intended should be known by Kubernetes, and it should be treated as such. Let's add a liveness and readiness probes to the pod. We can do this by editing the deployment. The application logs is informing us that there are endpoints for liveness and readiness probes, so let's use those.
```bash
{"level":"info","time":"2022-11-16T14:05:46+01:00","message":"Health is available at /healthz/live and /healthz/ready"}
```
```bash
kubectl edit deployment infoapi
```
```yaml
        readinessProbe:
          httpGet:
            path: /healthz/ready
            port: 5550
            scheme: HTTP
          timeoutSeconds: 2
          periodSeconds: 3
          successThreshold: 1
          failureThreshold: 5
        livenessProbe:
            httpGet:
              path: /healthz/live
              port: 5550
              scheme: HTTP
            timeoutSeconds: 2
            periodSeconds: 3
            successThreshold: 1
            failureThreshold: 5
        startupProbe:
            httpGet:
              path: /healthz/live
              port: 5550
            failureThreshold: 30
            periodSeconds: 3
```
Now the pod should be restarted and the newly created pod looks like it's running but not ready.
```bash
NAME                             READY   STATUS    RESTARTS      AGE
stockapi-7f88f897db-z9t2q        1/1     Running   0             20m
mongodb-598cb69f49-wxq4m         1/1     Running   0             20m
redis-794b86ddbd-7mzsg           1/1     Running   0             20m
infoapi-58bf87f775-m2klw         1/1     Running   0             20m
mongo-express-5c854f97c7-79mdc   1/1     Running   0             20m
stockapi2-595b754b75-5kjm5       1/1     Running   0             10m
infoapi-bcb7b5cfb-fwzl7          0/1     Running   0             7m29s
```
This shows that according to the developer of this service the pod is running, but it's not ready, which means restarting this pod will not resolve the issue, but it's not ready to serve. Let's investigate further and try to resolve the real problem.
## Step 5: Fix the connection string error
There's obviously something wrong with the connection to the database. The address mentioned in the error message is `mongodb.svc` but it should be just the service name for services within the same namespace, which is `mongodb`. Let's try to fix that by editing the deployment configuration.
```bash
kubectl edit deployment infoapi
```
and change the `MONGODB_HOST` environment variable to `mongodb` instead of `mongodb.svc`. Now if we do the same request to the API by running 
```bash
curl http://bookinfo.localdev.me/book
```
```bash
[{"uuid":"5f1b9b9b9b9b9b9b9b9b9b9b","name":"Lord of the Rings: The Fellowship of the Ring","author":"J.R.R. Tolkien","currentStock":0},{"uuid":"636bd0e2f2c4780497f3ad6c","name":"Lord of the Rings: : The Two Towers","author":"J.R.R. Tolkien","currentStock":0},{"uuid":"636bd174330aaaf4cb0bc1f0","name":"Northanger Abbey","author":"Austen, Jane","currentStock":0},{"uuid":"636bd18020beabc2639f361d","name":"War and Peace","author":"Tolstoy, Leo","currentStock":0},{"uuid":"636bd189db7a3afe9ac84af5","name":"Anna Karenina","author":"Tolstoy, Leo","currentStock":0},{"uuid":"636bd1943ae1d6488e79d814","name":"Mrs. Dalloway","author":"Woolf, Virginia","currentStock":0},{"uuid":"636bd19d7c8fb05f50ae83ed","name":"The Hours","author":"Cunnningham, Michael","currentStock":0},{"uuid":"636bd1a6526fa1880cc835e2","name":"Huckleberry Finn","author":"Twain, Mark","currentStock":0},{"uuid":"636bd1b181e5c2587a4c7e55","name":"Bleak House","author":"Dickens, Charles","currentStock":0},{"uuid":"636bd1bb76e29382a46689c5","name":"Tom Sawyer","author":"Twain, Mark","currentStock":0},{"uuid":"636bd1c39cb0744ca2412814","name":"A Room of One's Own","author":"Woolf, Virginia","currentStock":0},{"uuid":"636bd1cf5ae5d596f3a56b25","name":"Harry Potter","author":"Rowling, J.K.","currentStock":0},{"uuid":"636bd1d70fe72a696fffa236","name":"One Hundred Years of Solitude","author":"Marquez","currentStock":0},{"uuid":"636bd1e04bb308fc94d1d042","name":"Hamlet, Prince of Denmark","author":"Shakespeare","currentStock":0}]
```

[To continue](SOLUTION-2.md)