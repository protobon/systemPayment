# System Payment

## Build Project


```console
$ make build stage=test
$ make build stage=prod
```

</br>

## Compose Up
```console
$ make up stage=test  # dup for detached
$ make dup stage=prod
```

</br>

## Docker Logs (view)
```console
$ make logs stage=test
$ make logs stage=prod
```

</br>

## Run locally
```console
$ make build stage=local  # only the first time
$ make start stage=local
$ make run
```

</br>

# [Swagger](http://localhost:8080/swagger/index.html)
