FROM golang:alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN go build -o /out/mqtt_lorawan_consumer .

FROM scratch
COPY --from=build /out/mqtt_lorawan_consumer /
ENTRYPOINT ["./mqtt_lorawan_consumer"]