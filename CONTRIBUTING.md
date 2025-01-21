# Contributing

## Table of Contents
- [Contributing](#contributing)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
    - [Python](#python)
    - [Go](#go)
    - [Protocol Buffer Compiler](#protocol-buffer-compiler)
  - [Installation](#installation)
  - [Mage](#mage)
  - [Common information](#common-information)
    - [MLflow source code](#mlflow-source-code)
    - [Protos](#protos)
    - [Request validation](#request-validation)
    - [Data access](#data-access)
    - [Linting](#linting)
    - [Building the Go binary](#building-the-go-binary)
  - [Go MLflow server](#go-mlflow-server)
    - [Approach](#approach)
    - [Run with Mage](#run-with-mage)
  - [Supported endpoints](#supported-endpoints)
  - [Porting an Endpoint](#porting-an-endpoint)
  - [Run tests](#run-tests)
    - [Debug Failing Tests](#debug-failing-tests)
    - [Targeting Local Postgres in Python Tests](#targeting-local-postgres-in-python-tests)



## Prerequisites

To contribute to this project, you need the following:

### Python

- [UV](https://docs.astral.sh/uv/getting-started/installation/)
- [pre-commit](https://pre-commit.com/) (via `uv tool install pre-commit`)
- [ruff](https://astral.sh/ruff) (via `uv tool install ruff`)

### Go

- [Go 1.23](https://go.dev/doc/install)
- [Mage](https://magefile.org/) (via `go install github.com/magefile/mage@v1.15.0`)
- [protoc-gen-go](https://pkg.go.dev/github.com/golang/protobuf/protoc-gen-go) (via `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0`)
- [Mockery](https://vektra.github.io/mockery/latest/) (via `go install github.com/vektra/mockery/v2@v2.43.2`)
- [Golangci-lint](https://golangci-lint.run/) (via `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1`)

### Protocol Buffer Compiler

See [Protocol Buffer Compiler Installation](https://grpc.io/docs/protoc-installation/)

Alternatively, you can use the [development container](https://containers.dev/) that includes all the required tools.

## Installation

to configure all the development environment just run `mage` target:

```bash
mage configure
```

it will configure MLflow and all the Python dependencies required by the project or run each step manually:

```bash
# Install our Python package and its dependencies
pip install -e .

# Install the dreaded psycho
pip install psycopg2-binary

# Archive the MLflow pre-built UI
tar -C /usr/local/python/current/lib/python3.8/site-packages/mlflow -czvf ./ui.tgz ./server/js/build

# Clone the MLflow repo
git clone https://github.com/jgiannuzzi/mlflow.git -b master .mlflow.repo

# Add the UI back to it
tar -C .mlflow.repo/mlflow -xzvf ./ui.tgz

# Install it in editable mode
pip install -e .mlflow.repo
```

## Mage

This repository uses [mage](https://magefile.org/) to streamline some utility functions.

```bash
# Install mage (already done in the dev container)
go install github.com/magefile/mage@v1.15.0
```

```bash
# See all targets
mage
```

```shell
# Output
Targets:
  build                  a Python wheel.
  configure              development environment.
  dev                    Start the mlflow-go-backend dev server connecting to postgres.
  endpoints              Print an overview of implementated API endpoints.
  generate               Go files based on proto files and other configuration.
  repo:init              Clone or reset the .mlflow.repo fork.
  repo:update            Forcefully update the .mlflow.repo according to the .mlflow.ref.
  test:all               Run all tests.
  test:python            Run mlflow Python tests against the Go backend.
  test:pythonSpecific    Run specific Python test against the Go backend.
  test:unit              Run the Go unit tests.
```

```bash
# Execute single target
mage generate
```

The beauty of Mage is that we can use regular Go code for our scripting. That being said, we are not married to this tool.

## Common information

### MLflow source code

To integrate with MLflow, you need to include the source code. The [mlflow/mlflow](https://github.com/mlflow/mlflow/) repository contains proto files that define the tracking API. It also includes Python tests that we use to verify our Go implementation produces identical behaviour.

We use a `.mlflow.ref` file to specify the exact location from which to pull our sources. The format should be `remote#reference`, where `remote` is a git remote and `reference` is a branch, tag, or commit SHA.

If the `.mlflow.ref` file is modified and becomes out of sync with the current source files, the mage target will automatically detect this. To manually force a sync, you can run `mage repo:update`.

### Protos

To ensure we stay compatible with the Python implementation, we aim to generate as much as possible based on the `.proto` files.

By running

```bash
mage generate
```

Go code will be generated. Use the protos files from `.mlflow.repo` repository.

This includes the generation of:

- Structs for each endpoint. ([pkg/protos](./pkg/protos/service.pb.go))
- Go interfaces for each service. ([pkg/contract/service/*.g.go](./pkg/contract/service/tracking.g.go))
- [fiber](https://gofiber.io/) routes for each endpoint. ([pkg/server/routes/*.g.go](./pkg/server/routes/tracking.g.go))

If there is any change in the proto files, this should ripple into the Go code.

### Request validation

We use [Go validator](https://github.com/go-playground/validator) to validate all incoming request structs.
As the proto files don't specify any validation rules, we map them manually in [pkg/cmd/generate/validations.go](./cmd/generate/validations.go).

Once the mapping has been done, validation will be invoked automatically in the generated fiber code.

When the need arises, we can write custom validation function in [pkg/validation/validation.go](./validation/validation.go).

### Data access

Initially, we want to focus on supporting Postgres SQL. We chose [Gorm](https://gorm.io/) as ORM to interact with the database.

We do not generate any Go code based on the database schema. Gorm has generation capabilities but they didn't fit our needs. The plan would be to eventually assert the current code still matches the database schema via an integration test.

All the models use pointers for their fields. We do this for performance reasons and to distinguish between zero values and null values.

### Linting

We have enabled various linters from [golangci-lint](https://golangci-lint.run/), you can run these via:

```bash
pre-commit run golangci-lint --all-files
```

Sometimes `golangci-lint` can complain about unrelated files, run `golangci-lint cache clean` to clear the cache.

### Building the Go binary

To ensure everything still compiles:

```bash
go build -o /dev/null ./pkg/cmd/server
```

or

```bash
python -m mlflow_go.lib . /tmp
```

## Go MLflow server

### Approach

To enable use of the Go server, users can run the `mlflow-go server` command.

```bash
# Start the Go server with a database URI
# Other databases are supported as well: sqlite, mysql and mssql
mlflow-go server --backend-store-uri postgresql://postgres:postgres@localhost:5432/postgres
```

This will launch the python process as usual. Within Python, a random port is chosen to start the existing server and a Go child process is spawned. The Go server will use the user specified port (5000 by default) and spawn the actual Python server as its own child process (`gunicorn` or `waitress`).
Any incoming requests the Go server cannot process will be proxied to the existing Python server.

Any Go-specific options can be passed with `--go-opts`, which takes a comma-separated list of key-value pairs.

```bash
mlflow-go server --backend-store-uri postgresql://postgres:postgres@localhost:5432/postgres --go-opts log_level=debug,shutdown_timeout=5s
```

MLflow client could be pointed the Go server:

```python
import mlflow

# Use the Go server
mlflow.set_tracking_uri("http://localhost:5000")

# Use MLflow as usual
mlflow.set_experiment("my-experiment")

with mlflow.start_run():
    mlflow.log_param("param", 1)
    mlflow.log_metric("metric", 2)
```

### Run with Mage

To start the mlflow-go-backend dev server connecting to postgres just run next `mage` target:

```bash
mage dev
```

The postgres database should already be running prior to this command. By default service uses next connection string:

```
postgresql://postgres:postgres@localhost:5432/postgres
```

but it could be configured in [mage](./magefiles/run.go)

## Supported endpoints

The currently supported endpoints can be found by running mage command:

```bash
mage endpoints
```

## Porting an Endpoint

If you wish to contribute to the porting of an existing Python endpoint, you can read our [dedicated guide](./docs/porting-a-new-endpoint.md).

## Run tests

The Python integration tests have been adapted to also run against the Go implementation.
Next `mage` targets are available to run different types of tests:

```bash
# Run all the available tests
mage test:all
```

```bash
# Run just MLflow Python tests
mage test:python
```

```bash
# Run specific MLflow Python tests (matches all tests containing the argument)
mage test:pythonSpecific <test_file::test_name>

#Example
mage test:pythonSpecific ".mlflow.repo/tests/tracking/test_rest_tracking.py::test_rename_experiment"
```

```bash
# Run just unit tests 
mage test:unit
```

Additionally, there is always an option to run, specific test\tests if it is necessary:

```bash
pytest tests/tracking/test_rest_tracking.py
```

To run only the tests targeting the Go implementation, you can use the `-k` flag:

```bash
pytest tests/tracking/test_rest_tracking.py -k '[go-'
```

If you'd like to run a specific test and see its output 'live', you can use the `-s` flag:

```bash
pytest -s "tests/tracking/test_rest_tracking.py::test_create_experiment_validation[go-postgresql]"
```

See the [pytest documentation](https://docs.pytest.org/en/8.2.x/how-to/usage.html#specifying-which-tests-to-run) for more details.

```bash
# Build the Go binary in a temporary directory
libpath=$(mktemp -d)
python -m mlflow_go.lib . $libpath

# Run the tests (currently just the server ones)
MLFLOW_GO_LIBRARY_PATH=$libpath pytest --confcutdir=. \
  .mlflow.repo/tests/tracking/test_rest_tracking.py \
  .mlflow.repo/tests/tracking/test_model_registry.py \
  .mlflow.repo/tests/store/tracking/test_sqlalchemy_store.py \
  .mlflow.repo/tests/store/model_registry/test_sqlalchemy_store.py \
  -k 'not [file'

# Remove the Go binary
rm -rf $libpath

# If you want to run a specific test with more verbosity
# -s for live output
# --log-level=debug for more verbosity (passed down to the Go server/stores)
MLFLOW_GO_LIBRARY_PATH=$libpath pytest --confcutdir=. \
  .mlflow.repo/tests/tracking/test_rest_tracking.py::test_create_experiment_validation \
  -k 'not [file' \
  -s --log-level=debug
```

### Debug Failing Tests

Sometimes, it can be very useful to modify failing tests and use `print` statements to display the current state or differences between objects from Python or Go services.

Adding `"-vv"` to the `pytest` command in `magefiles/tests.go` can also provide more information when assertions are not met.

### Targeting Local Postgres in Python Tests

At times, you might want to apply store calls to your local database to investigate certain read operations via the local tracking server.

You can achieve this by changing:

```python
def test_search_runs_datasets(store: SqlAlchemyStore):
```

to:

```python
def test_search_runs_datasets():
    db_uri = "postgresql://postgres:postgres@localhost:5432/postgres"
    artifact_uri = Path("/tmp/artifacts")
    artifact_uri.mkdir(exist_ok=True)
    store = SqlAlchemyStore(db_uri, artifact_uri.as_uri())
```

in the test file located in `.mlflow.repo`.