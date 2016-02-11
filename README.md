[![GoDoc](https://godoc.org/github.com/venicegeo/grid-sdk-go?status.svg)](https://godoc.org/github.com/venicegeo/grid-sdk-go)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/venicegeo/grid-sdk-go/blob/master/LICENSE.txt)

# grid-sdk-go

`grid-sdk-go` provides two key elements:

* a Go package for accessing the GRiD API, and
* a Go CLI for interacting with GRiD.

## Installing

Use `go get` to install the latest version of both the CLI and the library.

    $ go get -v github.com/venicegeo/grid-sdk-go

To include it in your application.

    import "github.com/venicegeo/grid-sdk-go"

## Using the GRiD CLI

To get an overview of the available commands, type `grid`.

    $ grid
    grid is the main command.

    grid provide CLI access to GRiD.

    Usage:
      grid [command]

    Available Commands:
      add         Add an AOI
      configure   Configure the CLI
      lookup      Get suggested AOI name
      ls          List AOI/Export/File details
      pull        Download File
      version     Print the version number of the GRiD CLI
      help        Help about any command

    Flags:
      -h, --help=false: help for grid


    Use "grid help [command]" for more information about a command.

First, we must provide our GRiD credentials via `grid configure`.

    $ grid configure
    GRiD Username: johnsmith
    GRiD Password:

This will create a file in `$HOME/.grid/credentials` on Linux/Mac OS X, or `%HOMEPATH%/.grid/credentials` on Windows. This credentials file will be used each time GRiD authentication is required.

To view a complete listing of user AOIs:

    $ grid ls
    PRIMARY KEY    NAME    CREATED AT
    1              Foo     2015-06-22T08:15:33.513
    2              Bar     2013-12-17T14:08:53.316

To view details of an individual AOI:

    $ grid ls 1
    PRIMARY KEY    NAME                   DATATYPE    STARTED AT
    3              Foo_2013-Sep-11.zip    N/A         2013-09-11T14:32:23.292031
    4              Foo_2013-Sep-11.zip    N/A         2013-09-11T11:43:38.729971

Or multiple AOIs:

    $ grid ls 1 2
    PRIMARY KEY    NAME                   DATATYPE    STARTED AT
    3              Foo_2013-Sep-11.zip    N/A         2013-09-11T14:32:23.292031
    4              Foo_2013-Sep-11.zip    N/A         2013-09-11T11:43:38.729971
    PRIMARY KEY    NAME                   DATATYPE    STARTED AT
    5              Bar_2013-Sep-11.zip    N/A         2013-09-11T14:32:23.292031
    6              Bar_2013-Sep-11.zip    N/A         2013-09-11T11:43:38.729971

Likewise, one or more exports within an AOI:

    $ grid ls 3 5
    PRIMARY KEY    NAME
    7              Foo_2013-Sep-11.las
    PRIMARY KEY    NAME
    8              Bar_2013-Sep-11.las

To download an exported file:

    $ grid pull 7

To add an AOI:

To get a suggested AOI name:

## Using the library

GRiD currently uses Basic Authentication for all API calls. We begin by creating a `BasicAuthTransport`, supplying just username and password for our GRiD account. From this, we generate a new GRiD client and begin making API calls. We then indicate our desire to invoke the `List` function provided as part of the `AOI` service, finally iterating over the returned AOIs and printing the AOI names.

```go
package main

import "github.com/venicegeo/grid-sdk-go"

func main() {
  tp := grid.BasicAuthTransport{
    Username:  strings.TrimSpace(/*user's GRiD username*/),
    Password:  strings.TrimSpace(/*user's GRiD password*/),
  }

  client := grid.NewClient(tp.Client())

  a, _, err := client.AOI.List("")

  for _, pk := range a {
    for _, aoi := range pk.AOIs {
      fmt.Println(aoi.Fields.Name)
    }
  }
}
```
