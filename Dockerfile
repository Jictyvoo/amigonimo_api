### STAGE 1: Build ###
FROM golang:1.26-alpine AS builder

ENV PATH=$GOPATH/bin:$PATH
ENV CGO_ENABLED=0
ENV GO1111MODULE=on

RUN apk add --no-cache git ca-certificates

ENV SOURCE_DIR=/go/src/amigonimo
WORKDIR $SOURCE_DIR

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
COPY go.* ./
RUN go mod download

# Copy all the Code and stuff to compile everything
COPY . .

# Builds the application as a statically linked binary
RUN GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -a -installsuffix cgo -o anonymigo_api ./cmd/api

########################################################################################################################
# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest AS runtime

RUN apk add --no-cache ca-certificates

WORKDIR /home/amigonimo

# Copy the generated binary (migrations SQL files are embedded inside it)
COPY --from=builder /go/src/amigonimo/anonymigo_api /bin/anonymigo_api

RUN chmod +x /bin/anonymigo_api

ENTRYPOINT ["/bin/anonymigo_api"]
