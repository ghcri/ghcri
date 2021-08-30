# ghcri

ghcri is the repo for Github Container Registry Images. Just like [docker-library](https://github.com/docker-library) for Docker Registry.

## Usage

Replace all docker library from

```shell
docker pull golang:1.16
```

to

```shell
docker pull ghcr.io/ghcri/golang:1.16
```

All images from docker library will be copied AS-IS under `ghcr.io/ghcri`.

## Benefits

Say goodbye to:

```shell
ERROR: toomanyrequests: Too Many Requests.
```

or

```shell
You have reached your pull rate limit. You may increase the limit by authenticating and upgrading: https://www.docker.com/increase-rate-limits.
```

## Acknowledgements

- Thanks for hard work of [docker-library](https://github.com/docker-library). Without their work, this project is impossible.
