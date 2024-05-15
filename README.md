## Description

On startup application declares `input-A` and `output-A` queues and starts listen on `input-A`.
Upon delivery service processes it (basically does `{"result": " + string(input) + "}`).
Processed message is being published to default exchange with routing key `output-A`.

### Black box
Client publishes 'hello' to `input-A` and receives '{"result": "hello"}' from `output-A`.


## How to

### Lint
Format the code, fix dependencies and validate the project.
```bash
make lint
```

### Run
Bootstrap RabbitMQ, run the service and start processing messages received from the queue.
```bash
make
```

## Objective

Add support for the second stream of data using `input-B` and `output-B` queues.

### Black box
1. Client publishes 'hello' to `input-A` and receives '{"result": "hello"}' from `output-A`.
2. Client publishes 'hello' to `input-B` and receives '{"result": "hello"}' from `output-B`.
