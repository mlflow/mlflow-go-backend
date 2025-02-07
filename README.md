[![PyPI - Version](https://img.shields.io/pypi/v/mlflow-go-backend)](https://pypi.org/project/mlflow-go-backend/)

# Go backend for MLflow

In order to increase the performance of the tracking server and the various stores, we propose to rewrite the server and store implementation in Go.

## Getting started

### Installation

```shell
pip install mlflow-go-backend
```

### CLI Usage

`mlflow-go-backend` is meant to be a drop-in-replacement for `mlflow`.

You can update your existing `mlflow` command with `mlflow-go`:

```diff
- mlflow server --backend-store-uri postgresql://postgres:postgres@localhost:5432/postgres
+ mlflow-go server --backend-store-uri postgresql://postgres:postgres@localhost:5432/postgres
```

Every existing setting of [mlflow server](https://mlflow.org/docs/latest/cli.html#mlflow-server) can be passed to `mlflow-go`.

> [!CAUTION]
> The Go implementation currently does not support file storage as a backend store. You must provide the --backend-store-uri argument pointing to a database.

### Python Usage

```py
import mlflow
import mlflow_go

# Enable the Go client implementation (disabled by default)
mlflow_go.enable_go()

# Set the tracking URI (you can also set it via the environment variable MLFLOW_TRACKING_URI)
# Currently only database URIs are supported
mlflow.set_tracking_uri("sqlite:///mlflow.db")

# Use MLflow as usual
mlflow.set_experiment("my-experiment")

with mlflow.start_run():
    mlflow.log_param("param", 1)
    mlflow.log_metric("metric", 2)
```

### Python Store Usage

```py
import logging
import mlflow
import mlflow_go

# Enable debug logging
logging.basicConfig()
logging.getLogger('mlflow_go').setLevel(logging.DEBUG)

# Enable the Go client implementation (disabled by default)
mlflow_go.enable_go()

# Instantiate the tracking store with a database URI
tracking_store = mlflow.tracking._tracking_service.utils._get_store('sqlite:///mlflow.db')

# Call any tracking store method
tracking_store.get_experiment(0)

# Instantiate the model registry store with a database URI
model_registry_store = mlflow.tracking._model_registry.utils._get_store('sqlite:///mlflow.db')

# Call any model registry store method
model_registry_store.get_latest_versions("model")
```

## Why bother?

For that sweet performance my friend!

Initial benchmarks show us that critical API calls are faster and Go can handle more concurrent requests.

![Duration Results](./benchmarks/results_duration.png)
![Iteration Results](./benchmarks/results_iterations.png)

> [!NOTE]
> These were initial results and can perhaps still be more optimized.

See [benchmarks](./benchmarks/README.md) for more information.


## Contribution

### Try it out!

We advocate this project as a drop-in replacement for `mlflow`.
Please give it a try, and let us know if anything isn't working for you!
We look forward to your feedback!

### Missing Endpoints

Hey there! We need your help! This project is not yet a full port of the current tracking server.
Not every [endpoint in the REST API](https://mlflow.org/docs/latest/rest-api.html) is implemented in Go.
Those missing are redirected to the existing Python implementation.
Please consider helping us out and implement [a missing endpoint](./docs/porting-a-new-endpoint.md).

Not yet a Go wizard? We've collected some [resources to learn Go](./docs/learning-go.md). Trust us, it will be an enriching and fun experience!

## Community

There is a dedicated Slack channel called `#mlflow-go` on the official [mlflow Slack](https://mlflow.org/slack).
Shoot us a message if you have any questions!
