# BookInfo API

Besides the expected functional endpoints, the API will expose these endpoints for the following purposes:

| Path            | Purpose             |
|-----------------|---------------------|
| `/metrics`      | The metrics endpoint will be exposing metrics of the API to be consumed by Prometheus | 
| `/health/live`  | This endpoint will be the liveness probe's target. It returns 200 OK if the server is up and running |
| `/health/ready` | This endpoint will be the readiness probe's target. It checks the database connection and returns `200 OK` if connection is good |
| `/version`      | This endpoint simply reads the version.txt file inside the static folder. Normally version.txt will contain a token instead of a version number which can be replaced by the pipeline |
| `/docs`         | This endpoint uses redoc interface to show the OpenAPI documentation of the API |

## Initial state
There are two services named bookInfoAPI and bookStockAPI. bookInfoAPI has two endpoints namely `/book` and `/book/{id}`. The first loads all the book information from the backing MongoDB database. The latter loads the book information from the MongoDB with the given `{id}` and also loads its current stock count from `bookStockAPI` using its `/book{id}` endpoint. bookStockAPI loads stock information from its backing Redis database.
These databases come with their preloaded data, so the services should be working when they are first deployed.
