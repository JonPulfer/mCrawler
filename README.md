# mCrawler

## Overview

This is a simple website crawler with configurable site and worker parameters from command line options. The design considers that the crawler should be well behaved and only request pages once. It should also limit the number of concurrent requests made to the target site.

The implementation includes a naive stack to handle making requests to discovered pages and a b-tree to store the site structure.

## Installation

A working Go development environment is required with a correctly configured GOPATH. The repos can then be placed in the GOPATH and built with `go install`

## Usage

Once built, a command line tool will be available which is controlled by flags. Running `mCrawler -h` will display help information about the options including example input.
