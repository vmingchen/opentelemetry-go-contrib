# Dynamic config service example

Right now the generated protobuf .go file we're using is in our fork of the repo. Eventually we want to import in from opentelemetry-proto.

### Expected behaviour for this example
After about 10 seconds, you'll see metrics get outputted to stdout once every second.

### Run server

```sh
go run ./server
```

### Run sdk

```sh
go run ./
```
