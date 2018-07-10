# _gen Code Generator

This command line tool is part of Godash - an efficient [Lodash](https://lodash.com) port for golang. Since golang doesn't currently support generics, the closest workaround is to work with code generation. The `_gen` command line tool generates code based on your project requirements and adds it to the `go-dash/slice` library.

## Installation

```
brew install go-dash/tools/gen
``` 

Verify the installation by running `_gen --version`.

## Usage

1. Write your project code normally according to the instructions in the `go-dash/slices` library.

2. Using terminal, open your project root and run the command line tool once by typing `_gen`.

3. That's it, the new code for the types you need will appear in `$GOPATH/src/github.com/go-dash/slice`.

## How it works

The tool first analyzes your source code and looks for imports of the format:

```
"github.com/go-dash/slice/_TYPE"
```  

For example:

```go
import "github.com/go-dash/slice/_string"
import "github.com/go-dash/slice/_int"
import "github.com/go-dash/slice/_Person"
```

It then generates new subsets of the library for your specific types. The generated code will be placed in:

```
$GOPATH/src/github.com/go-dash/slice/_string
$GOPATH/src/github.com/go-dash/slice/_int
$GOPATH/src/github.com/go-dash/slice/_Person
```

## License

MIT

## Who made this

Godash is developed by the [Orbs.com](https://orbs.com) engineering team. Orbs is a public blockchain infrastructure for decentralized consumer applications with millions of users. Orbs core has an open source [implementation](https://github.com/orbs-network/orbs-network-go) in golang.
