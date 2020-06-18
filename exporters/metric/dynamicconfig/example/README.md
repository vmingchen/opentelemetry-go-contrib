# Dynamic config service example

Right now the generated protobuf .go file we're using is in our fork of the repo. Eventually we want to import in from opentelemetry-proto.

Additionally, right now (June 18, 2020), the collector received an update that makes it incompatible with metric exporting from the SDK. We will need to export to an older version of the collector.

### Expected behaviour for this example
After about 10 seconds, you'll see metrics get outputted to stdout once every second.

You can change expected behaviour by editing the config served by the dummy configuration service in `./server/server.go`.

### Setup and run collector

```sh
// Clone repo
git clone https://github.com/open-telemetry/opentelemetry-collector.git

// Checkout proper commit
git branch compatible-master 746db761d19ed12ac2278cdfe7f30826a5ba6257 && git checkout compatible-master
```

### Run server

```sh
go run ./server
```

### Run sdk

```sh
go run ./
```
