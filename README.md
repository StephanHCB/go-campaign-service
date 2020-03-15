# go-campaign-service

## Overview

This service is part of my example for [enterprise microservices](https://peter.bourgon.org/go-kit/) in 
[go](https://golang.org/).

The business scenario is rather contrived and not really the point here:
- [go-mailer-service](https://github.com/StephanHCB/go-mailer-service) 
  offers a REST API to send an email
  given an email address, a subject, and a body. When an email
  is sent, a Kafka message must be sent to inform some hypothetical
  downstream service, but only if a feature toggle is switched on
- [go-campaign-service](https://github.com/StephanHCB/go-campaign-service)
  (this service) offers a REST API to plan a campaign (really just a list of email addresses,
  plus a subject and a body) and execute it, 
  using the mailer service

See the README of [go-mailer-service](https://github.com/StephanHCB/go-mailer-service/README.md) for all further
details, including detailed discussions of the libraries used and the various nonfunctional requirements.

## Developer Instructions

### Development Project Setup

This service uses go modules to provide dependency management, see `go.mod`.

If you place this repository OUTSIDE of your gopath, go build and go test will clone
all required dependencies by default.

### Running on Localhost

On the command line, `go build main.go` will download all dependencies and build a standalone executable
for you.

The executable expects two configuration files `config.yaml` and `secrets.yaml` in the current directory.
You can override their path locations by passing command line options as follows:

```main --config-path=. --secrets-path=.``` 

Find configuration templates under docs, copy them to the main directory and edit them so they fit your
environment.

#### Database Configuration

For your convenience, this service includes an in-memory database which is enabled by default, and is also
used by the automated tests.

If you would like to instead use a mysql database, and thus retain your data between runs, add something
like this to your `config.yaml`:

```yaml
database:
  use: 'mysql' # defaults to 'inmemory'
  mysql:
    username: 'demouser'
    password: 'demopw'
    database: 'tcp(localhost:3306)/dbname'
    parameters:
      - 'charset=utf8mb4'
      - 'collation=utf8mb4_general_ci'
      - 'parseTime=True'
      - 'timeout=30s' # connection timeout
```

### Running the Automated Tests

This service comes with unit, acceptance, and consumer driven contract tests. 

You can run all of these on the command line:

```go test ./...```

In order for the **contract tests** to work, you will need to perform some additional installation:

#### Contract Tests

This service is an example for using pact-go for consumer driven contract tests.

This is the **consumer** side.

See [go-mailer-service](https://github.com/StephanHCB/go-mailer-service/) for the **producer** side.

##### Solution Concept

In order to automatically verify that any cross-service interaction will work as expected, we have 
implemented the consumer side of a consumer driven contract test (see `test/contract/consumer/main_ctr_test.go`).

```
TODO implement a real world example
```

When the test suite of this client runs, the consumer side tests are run, and pact json files are written out to
`test/contract/consumer/pacts`.

_Note how the consumer test calls into a very low level function, the one that uses a httpclient to make the call,
even below any circuit breaker etc. So we are not testing the business logic, only the actual technical client code._

When the test suite of the producer runs, it reads the pact json and uses it to replay the interaction.

_Again, we use a mock service underneath the web controller to only test the technical interaction,
not the business logic. This is easy to do using a httptest server._

For this to work as we have implemented it here, you must have cloned both 
[go-mailer-service](https://github.com/StephanHCB/go-mailer-service/) and
[go-campaign-service](https://github.com/StephanHCB/go-campaign-service) right next to each other, as the
producer uses a relative path to find the consumer generated pacts.

_In a more real world example you'd have some way to publish the generated pacts to a server and/or
check them into a repository. The producer side test can then use a URL on this server to download the current
pacts._

##### Installation of Pact

Download and install the pact command line tools and add them to your path as described in the
[pact-go manual](https://github.com/pact-foundation/pact-go#installation). This step is system
dependent.

##### Run The Contract Tests

Use

`go test -v github.com/StephanHCB/go-campaign-service/...`

to run the consumer side test and generate the pact json file.

Then do the same in the producer project.

`go test -v github.com/StephanHCB/go-mailer-service/...`

You should see output like this:

```
TODO
```
