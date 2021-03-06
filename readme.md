
Service for managing configuration variables of IoT devices.

GRPC and HTTP interfaces are provided for setting configuration variables. The service will publish desired changes to an mqtt topic.

The service subscribes to another mqtt topic to receive values of configuration variables reported by devices in the field.

A LoRaWAN server such as [chirpstack.io](https://chirpstack.io) may be used to provide an mqtt interface.

The service checks intermittently for consistency between the desired and reported state of each device.

A database is required, currently there is support for [Couchbase](./internal/dbclient/nosql) and [PostgreSQL](./internal/dbclient/sql).

To run on Kubernetes,
- use the [Dockerfile](./Dockerfile) to build a docker container and push to a container registry. 
- Create a values.yaml file such as [this](./deployments/devicetwin-helm/values-example.yaml)
- run
```sh
helm upgrade --install \
    -f deployments/devicetwin-helm/values-example.yaml \
    devicetwin \
    devicetwin-helm
```


