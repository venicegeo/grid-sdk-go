[![GoDoc](https://godoc.org/github.com/venicegeo/grid-sdk-go?status.svg)](https://godoc.org/github.com/venicegeo/grid-sdk-go)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/venicegeo/grid-sdk-go/blob/master/LICENSE.txt)
[![Build Status](https://travis-ci.org/venicegeo/grid-sdk-go.svg?branch=master)](https://travis-ci.org/venicegeo/grid-sdk-go)

# grid-sdk-go

`grid-sdk-go` provides two key elements:

* a Go package for accessing the GRiD API, and
* a Go CLI for interacting with GRiD.

## Installing

Use `go get` to install the latest version of both the CLI and the library.

    $ go get -v github.com/venicegeo/grid-sdk-go/...

To install.

    $ go install github.com/venicegeo/grid-sdk-go/...

To include it in your application.

    import "github.com/venicegeo/grid-sdk-go"

## Using the GRiD CLI

The first time we run `grid`, we must provide our GRiD credentials. These can be
updated at any time by running `grid configure`.

    $ grid configure
    GRiD Username: johnsmith
    GRiD Password:
    GRiD API Key: MyAPI-key
    GRiD Base URL: https://rsgis.erdc.dren.mil/te_ba/

This will create (or update) the configuration file in `$HOME/.grid/credentials`
on Linux/Mac OS X, or `%HOMEPATH%/.grid/credentials` on Windows. This
credentials file will be used each time GRiD authentication is required.

To get an overview of the available commands, just type `grid`.

    $ grid
    grid is a command-line interface to the GRiD database.

    Usage:
      grid [command]

    Available Commands:
      add         Add an AOI
      configure   Configure the CLI
      export      Initiate a GRiD Export
      lookup      Get suggested AOI name
      ls          List AOI/Export/File details
      pull        Download File
      task        Get task details
      version     Print the version number of the GRiD CLI

    Flags:
      -h, --help   help for grid

    Use "grid [command] --help" for more information about a command.

To view a complete listing of user AOIs:

    $ grid ls
    PRIMARY KEY    NAME    CREATED AT
    1              Foo     2015-06-22T08:15:33.513
    2              Bar     2013-12-17T14:08:53.316

To view details of an individual AOI:

    $ grid ls 1

    NAME: Foo
    CREATED AT: 2014-02-07T14:22:44.437

    RASTER COLLECTS
    PRIMARY KEY    NAME                   DATATYPE
    101            20091113_Foo           EO

    POINTCLOUD COLLECTS
    PRIMARY KEY    NAME                   DATATYPE
    201            20101106_Foo           LAS 1.2  

    EXPORTS
    PRIMARY KEY    NAME                   DATATYPE    STARTED AT
    301            Foo_2013-Sep-11.zip    N/A         2013-09-11T14:32:23.292031
    302            Foo_2013-Sep-11.zip    N/A         2013-09-11T11:43:38.729971

Or multiple AOIs:

    $ grid ls 1 2

You can also mix a match AOI and export primary keys (collect IDs are not
currently available):

    $ grid ls 1 301

To download an exported file:

    $ grid pull 7

To get a suggested AOI name:

    $ grid lookup "POLYGON ((30 10, 40 40, 20 40, 10 20, 30 10))"
    Great Sand Sea

To add an AOI (the AOI is automatically named using the name provided by
`grid lookup`):

    $ grid add "POLYGON ((30 10, 40 40, 20 40, 10 20, 30 10))"
    Successfully created AOI "Great Sand Sea" with primary key "2880" at 2016-04-01T15:59:00.587

To export a point cloud:

    $ grid export -h
    Export is used to initiate a GRiD export for the AOI and for each of the provided collects.

    Usage:
      grid export [AOI] [Collects]... [flags]

Currently, only point cloud exports are available via the CLI. Each export
requires specification of exactly one AOI primary key, and one or more point
cloud collect primary keys (these will be merged into a single file). The API
returns a task ID (for task status queries), and an export ID to later retrieve
export details (e.g., `grid ls <export ID>`).

    $ grid export 1 201
    TASK ID                               EXPORT ID
    c7def4ee-8b47-4434-b4f5-2eecf984c0a6  303

To get export task status:

    $ grid task c7def4ee-8b47-4434-b4f5-2eecf984c0a6
    ID                                    NAME                          STATE
    c7def4ee-8b47-4434-b4f5-2eecf984c0a6  export.tasks.generate_export  RUNNING

## Using the library

Simply create a GRiD client, and then begin making requests.

```go
package main

import (
  "encoding/base64"

  "github.com/venicegeo/grid-sdk-go"
)

func main() {
  // GRiD uses basic authentication. An authorization string is formed by base64
  // encoding the string formed by concatenating the user's username, a colon,
  // and the user's password.
  auth := base64.StdEncoding.EncodeToString([]byte("johnsmith:password"))

  // A GRiD client is created by providing the authorization string, API key,
  // and base URL as strings. If the base URL is empty, the default Test &
  // Evaluation instance of GRiD will be targeted.
  g := grid.NewClient(auth, "MyAPI-key", "")

  // Get details of the AOI with primary key of 100. The GRiD client does not
  // panic or set any HTTP status codes on error. Errors are returned from each
  // request, along with the HTTP response, for consumers of the SDK to act on
  // these as they see fit.
  aoiListObject, resp, err := g.GetAOI(100)
}
```
