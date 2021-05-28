# THe MINARIA project
#### A project name which you can read but you never understand.

## Swagger Docs
Swagger docs are available at [Heroku free tier](https://minaria.herokuapp.com/docs).

## Pre-Push hook
The pre-push hook is configured to run following tasks:

- executes `tests.sh`
- `make swagger`

## How to use makefile?
To install swagger:

```
make check_install
```

To generate swagger.yml file:

```
make swagger
```

To generate clients (required for running end to end testing):

```
make swagger_gen_client
```