# Build

1. Build executable binaries for Windows and Linux with the docker-compose Command: 
<br>
`docker-compose run --rm mqtt_builder make clean build`

# Run in Docker

1. Build docker image using `docker build -t <tag> .`
2. Run docker container with config-file as volume: `docker run -dt -v <path>:/mqtt_lorawan_consumer.toml`