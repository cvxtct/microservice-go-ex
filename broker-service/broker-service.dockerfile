# base go image
# FROM golang:1.18-alpine as builder

# RUN mkdir /app

# COPY . /app

# WORKDIR /app

# RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# RUN chmod +x /app/brokerApp

# build tiny docker image

FROM alpine:latest

RUN mkdir /app

# COPY --from=builder /app/brokerApp /app
# With this the make build_up will copy the brokerApp binary into the docker instead to recompile it within
# the container again.
COPY brokerApp /app

CMD ["/app/brokerApp"]