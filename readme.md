
## Setup and Run

#### Install dependencies
```bash
go mod tidy
```

#### Create .env file
```bash
cp .env.example .env
```

#### Setup and run
```bash
make run
```

### Updating Graphql and gRPC

#### To generate GraphQL updates after updating schema.graphqls
```bash
make gen-graphql
```

#### To generate pb files updates after updating protofiles
```bash
make gen-proto
```

## Requests examples

### 1. REST requests:

http://localhost:8000

```
POST http://localhost:8000/order HTTP/1.1
content-type: application/json

{
    "id": "rest",
    "tax": 0.3,
    "price": 10.1
}
```

```
GET http://localhost:8000/orders HTTP/1.1
content-type: application/json
```

### 2. GraphQL requests
playground: http://localhost:8080/

```graphql
mutation createOrder {
  createOrder(input: { 
    id: "graphql",
    Price: 10.0,
    Tax: 1
  }) {
    id,
    Price,
    Tax
  }
}
```

```graphql
query orders {
  orders {
    id,
    Price,
    Tax,
    FinalPrice
  }
}
```

### 3. GRPC requests

#### Connect to gRPC client
```bash
make evans
```

Evans gRPC client:

```bash
package pb
```

```bash
service OrderService
```

```bash
call CreateOrder
"grpc"
10.1
1.1
```

```bash
call GetOrders
```

