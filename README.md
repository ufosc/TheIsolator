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

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
