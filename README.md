# csv_builder

## Development And Running

### Intall Docker
Ensure you have the correct docker version installed
If you are using linux you may have to use sudo for the following commands
```sh
$ docker --version
Docker version 17.06.0-ce
$
```

### Retreive the docker image

```sh
$ docker pull golang:alpine
```

### Build the docker image

You should Note that if you want to include additional input files that you can edit this file to include them in the build
```sh
$ docker build -f Dockerfile.builder -t mydockerimage:latest .

Sending build context to Docker daemon    130kB
Step 1/5 : FROM golang:alpine
.
.
.
Successfully built d82732716399
Successfully tagged mydockerimage:latest
$
```

### Run ash on the docker image

```sh
$ sudo docker run -it mydockerimage:latest /bin/ash
```

you should be greeted with the shell prompt
```sh
/go/src/csv_builder #
```

### Run the application
From the ash shell, execute the wp_merge function

```sh
/go/src/csv_builder # ./wpe_merge input.csv output.csv
```