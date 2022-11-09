# API

We would like you to implement the following five HTTP endpoints.
If you want to, you can expand the capabilities of the API, but please ensure that the following endpoints will work as described below.
We will use a program to verify that these endpoints are working as expected.

| HTTP Verb | Path             | Request Content-Type | Request body | Response Content-Type | Example response body |
|-----------|------------------|----------------------|--------------|-----------------------|-----------------------|
| GET       | `/people`        | `application/json`   | -            | `application/json`    | `[ { "uuid": "49dc24bd-906d-4497-bcfc-ecc8c309ecfc", survived": true, "bookInfoClass": 3, "name": "Mr. Owen Harris Braund", "sex": "male", "age": 22, "siblingsOrSpousesAboard": 1, "parentsOrChildrenAboard":0, "fare":7.25}, ... ]` |
| POST      | `/people`        | `application/json`   | `{ "survived": true, "bookInfoClass": 3, "name": "Mr. Owen Harris Braund", "sex": "male", "age": 22, "siblingsOrSpousesAboard": 1, "parentsOrChildrenAboard":0, "fare":7.25}` | `application/json`    |  `{ "uuid": "49dc24bd-906d-4497-bcfc-ecc8c309ecfc", survived": true, "bookInfoClass": 3, "name": "Mr. Owen Harris Braund", "sex": "male", "age": 22, "siblingsOrSpousesAboard": 1, "parentsOrChildrenAboard":0, "fare":7.25}` |
| GET       | `/people/{uuid}` | `application/json`   | -            | `application/json`    | `{ "uuid": "49dc24bd-906d-4497-bcfc-ecc8c309ecfc", survived": true, "bookInfoClass": 3, "name": "Mr. Owen Harris Braund", "sex": "male", "age": 22, "siblingsOrSpousesAboard": 1, "parentsOrChildrenAboard":0, "fare":7.25}` |
| DELETE    | `/people/{uuid}` | `application/json`   | -            | `application/json`    | - |
| PUT       | `/people/{uuid}` | `application/json`   | `{ "survived": true, "bookInfoClass": 3, "name": "Mr. Owen Harris Braund", "sex": "male", "age": 22, "siblingsOrSpousesAboard": 1, "parentsOrChildrenAboard":0, "fare":7.25}` | `application/json`    | - |

We also have included a [swagger.yaml](./swagger.yml) file that describes this same API.

## Potential improvements
- Add a simple UI by responsing to `application/html` content types too on `GET /people` and `GET /people/{uuid}`.
