### STAGE 1: Build ###
FROM golang:1.25-bullseye AS builder

ENV PATH $GOPATH/bin:$PATH
ENV CGO_ENABLED=1
ENV GO1111MODULE=on


RUN apt update && \
    apt install -y --no-install-recommends \
    gcc \
    git \
    libc6-dev \
    libcurl4-openssl-dev && \
    pkg-config \
    rm -rf /var/lib/apt/lists/*


ENV SOURCE_DIR=/go/src/amigonimo
WORKDIR $SOURCE_DIR

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
COPY go.* ./
RUN go mod download

# Copy all the Code and stuff to compile everything
COPY . .

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -a -installsuffix cgo -o anonymigo_api ./cmd/discordbot
#
#
########################################################################################################################
# Moving the binary to the 'final Image' to make it smaller
FROM debian:buster-slim

RUN apt update && \
        apt install -y --no-install-recommends \
        ca-certificates curl \
        && rm -rf /var/lib/apt/lists/*

RUN curl -sSf https://atlasgo.sh | sh

WORKDIR /home/amigonimo

# Copy the generated binary from builder image to execution image
COPY --from=builder /go/src/amigonimo/anonymigo_api /bin/anonymigo_api

COPY build/migrations ./build/migrations
COPY build/container-entrypoint.sh ./entrypoint.sh

# Ensure the binary is executable
RUN chmod +x /bin/anonymigo_api && chmod +x ./entrypoint.sh && \
  mkdir -p /app/data

# Run the binary program
ENTRYPOINT ["/home/amigonimo/entrypoint.sh"]
CMD ["/home/amigonimo/entrypoint.sh"]
