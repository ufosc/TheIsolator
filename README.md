# TheIsolator

An ActivityPub implementation for static websites. Built during SwampHacks 2019.

## Installation

Install [Ruby and Jekyll](https://jekyllrb.com/docs/installation/).

Install [Go](https://golang.org/dl/) and [Dep](https://github.com/golang/dep). Note: Please make sure to have the project with-in your [Go src path](https://golang.org/doc/code.html).

<!-- Add Go Info -->

## Running

For local development, first setup Jekyll:

```bash
cd jekyll
bundle exec jekyll serve
```

To run Go:

```bash
cd go
go run
```

to create and run a local Go executable:

```bash
cd go
go build
./build
```

To get all dependencies (Assuming Dep is already installed):
```bash
cd go
dep install
#there should be no error outputs, you should only see "Fetching sources"
#followed by download status indicators. 
```

To add dependency (assuming dep is already installed)
If youre using VS code, the GOlang extension will prevent you from testing if
anything is actually compiling because it will automatically remove unused imports by default
but, you can manually add and track the dependency using:
```bash
cd go
dep ensure -add [GITHUB URL/PACKAGE SRC]
# no feedback on this aside from "fetching packages",
#if there's errors it'll note them
```


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
