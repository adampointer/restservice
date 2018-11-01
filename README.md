# restservice

Example of a RESTful service implemented in Golang with state persisted to BoltDB 

## Build, Test and Run

```
docker build . -t restservice
docker run -t -p 8080:8080 restservice
```

## API

`GET /payments`         |  Returns all payments
`GET /payments/{id}`    |  Returns payment by ID
`PUT /payments/{id}`    |  Create a new payment
`POST /payments/{id}`   |  Update a payment
`DELETE /payments/{id}` |  Delete a payment
