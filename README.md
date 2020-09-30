# Build

1. Build executable binaries: `docker-compose run --rm mqtt_builder make clean build`

# Run in Docker

1. Build docker image using `docker build -t <tag> .`
2. Run docker container with config-file as volume: `docker run -dt -v <path-to-config>:/mqtt_lorawan_consumer.toml <tag>`

# Configuration Template

```
[mqtt_broker]
    url=""
    username=""
    password=""
    topic=""
    client_id=""

[parser]
    measurement_key=""
    tag_keys= [""]
    values_key=""

# No persistence when commented out
# [influx_db]
#     url=""
#     database=""
#     username=""
#     password=""

```