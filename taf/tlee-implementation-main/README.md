# TLEE implementation

This repository contains a Go library of the TLEE implementation with a preloaded trust model and numerical values for trust opinions.

## Usage

To run the `main.go` file, which executes the TLEE, follow these steps:

1. Download the Go dependencies:

    ```sh
    go mod tidy
    ```

2. Execute the `main.go` file:

    ```sh
    go run main.go
    ```

The `main` function will run `RunTLEE_TM0()`, which executes the TLEE with a preloaded trust model and numerical values for trust opinions.
