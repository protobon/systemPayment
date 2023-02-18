# Awesome API Server

## Build Project


```console
$ make build
$ make build stage=local
```
Production: 'make build stage=prod'
</br>
</br>

## Run Project

### Local
```console
$ make run
```

### Compose Test
```console
$ make up
```

### Detached 
```console
$ make dup
```
Production: 'make dup stage=prod'
</br>
</br>

## Docker Logs (view)
```console
$ make logs
```
Production: 'make logs stage=prod'
</br>


### [Swagger](http://localhost:8080/swagger/index.html)