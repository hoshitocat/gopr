# gopr
Draw inspiration from git-pr-release.

## Installation

```shell
go get github.com/hoshitocat/gopr
```

## Usage

```shell
gopr -token={{ Your Github API Token }} -base={{ Base repository }} -head={{ Head repository }}
```

-token you must pass this argument

-base args default value is `master`

-head args default value is `develop`

## Configuration

### Pull Request Title

You can customize generated pull request title.
You should create configuration file on your own repository path `{{ your own repository path}}/.github/TITLE_TEMPLATE.md`

You can use a format string. The following table, title format specifiers:

| Specifier | Description |
|:---|:---|
| `%a` | Abbreviated weekday name( `Sun` ... `Sat` ) |
| `%d` | Number of day( `01` ... `31` ) |
| `%m` | Number of month( `01` ... `12` ) |
| `%Y` | Number of year, four digits |
