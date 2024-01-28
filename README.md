# Cobble

A simple tool for creating a Minecraft Bedrock add-on from the CLI. Cobble provides a quick and easy way to initialize the default add-on files and dependencies.

## Installation

Since cobble is still not fully implemented, there is currently no build pipeline. To build cobble from source, clone this repository and run `go build` in the cobble directory.

## Usage

Run `cobble help` for a list of commands.
To initialize a new add-on, run `cobble new [name]`
Cobble will then prompt you with a list of options to configure the add-on's dependencies and capabilities. If `name` is not provided, cobble will prompt for a name. Otherwise, the name passed will be used and the name question will be skipped.
