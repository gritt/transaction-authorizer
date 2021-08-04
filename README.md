# Transaction Authorizer

I've followed a very standard and flat Go pkg structure in this project:

- `/cmd`: contains the only entrypoint of the app, which handles CLI input and output.
- `/internal`: because is not a library, all code lives inside the `internal` directory, in there, it is separated in
  simple pkgs:
- `/internal/core`: encapsulates the business logic for authorizing transactions and creating accounts as well as domain
  types, such as `Transaction` and `Account`, it does not know details like JSON parsing, CLI input / output, it talks
  to dependencies through contracts, and it's protected of external changes.

Other pkgs inside `internal` are responsible for implementing the interfaces needed by the core business logic, in this
case, there is a single `repository/memory_repository` that stores the state in memory.

In general, it's very simple implementation
of [hexagonal architecture](https://netflixtechblog.com/ready-for-changes-with-hexagonal-architecture-b315ec967749).
Because of the interfaces and dependency injection, this architecture is very easy to test and also very flexible.

If we needed to fetch the account or the transactions from an external APIs or even a database, it would only take a new
implementation of the core repository interface like `FindAccount` or `FindTransactionsAfter`, for the business logic it
doesn't matter where it came from, neither the format, the business logic only refers to domain types.

### Libraries used

- `github.com/stretchr/testify`

> A widely used toolkit with common assertions and mocks to improve the development of tests.

## Setup

The `Makefile` provides all commands needed to install, test and build the app.

If you simply `make` it will display all the options:

```shell
Available Commands:
-> install                 install dependencies
-> test                    run all tests with coverage
-> build                   generates a build
-> build-linux             generates a build for linux
```

## Running

* Required `go 1.16+`

To compile a new binary run:

- `make build`

It will create a binary called `authorizer` in the root of the project.

To execute the app, you can call it with a single argument:

```shell
./authorizer < path/to/input/file
```

Example Response:

```text
{"account":{"active-card":true,"available-limit":100},"violations":[]}
{"account":{"active-card":true,"available-limit":90},"violations":[]}
{"account":{"active-card":true,"available-limit":70},"violations":[]}
{"account":{"active-card":true,"available-limit":65},"violations":[]}
{"account":{"active-card":true,"available-limit":65},"violations":["high-frequency-small-interval","double-transaction"]}
{"account":{"active-card":true,"available-limit":65},"violations":["insufficient-limit","high-frequency-small-interval"]}
{"account":{"active-card":true,"available-limit":65},"violations":["insufficient-limit","high-frequency-small-interval"]}
{"account":{"active-card":true,"available-limit":50},"violations":[]}
```