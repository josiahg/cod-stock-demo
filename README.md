# Cloudera Operational Database Demo

This app demonstrates modern web apps can be easily built using Cloudera Operational Database as a data store.

# Quick Start

## Prerequisites

[Install golang](https://golang.org/doc/install)

Deploy a Cloudera Operational Database on [CDP Public Cloud](https://docs.cloudera.com/runtime/7.1.0/howto-operational-database.html)

Ensure you have set your [CDP Workload Password](https://docs.cloudera.com/management-console/cloud/user-management/topics/mc-setting-the-ipa-password.html)

Sign up for a free [AlphaVantage API Key](https://www.alphavantage.co/support/#api-key)

## Go!
Clone this repository

```
$ git clone https://github.com/josiahg/go-cod-demo
$ cd go-cod-demo
```

Set the following environment variables:

```
$ export AV_API_KEY=<Your AlphaVantage API Key>
$ export COD_USER=<Your CDP Workload Username>
$ export COD_PASS=<Your CDP Workload Password>
$ export COD_URL=<Your COD base URL>
```

Run the app

```
$ go run main.go
```

By default the app will run on port 8080 - you can access it at http://localhost:8080/


# Built With

[Cloudera Data Platform](https://www.cloudera.com/products/cloudera-data-platform.html)

[Golang](https://golang.org/)

[React](https://reactjs.org/)

[uPlot](https://github.com/leeoniya/uPlot)

[calcite-avatica-go](https://github.com/apache/calcite-avatica-go)

[Alpha Vantage](https://www.alphavantage.co/)

# Thanks 

Josh Elser - Original Java version of this concept
