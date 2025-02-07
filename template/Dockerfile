FROM golang:1.23 AS build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN apt-get update && apt-get install -y gcc libc-dev unzip
RUN curl -fsSL https://deno.land/install.sh | sh
WORKDIR /build
COPY go.* ./
RUN go mod download
COPY . .
RUN cd frontend && /root/.deno/bin/deno install && /root/.deno/bin/deno task build && cd ..
RUN go build main.go
RUN chmod +x /build/{{SOSO_APPNAME}}

FROM golang:1.23 AS dev
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN apt-get update && apt-get install -y gcc libc-dev unzip
RUN apt-get install -y ca-certificates && update-ca-certificates
RUN curl -fsSL https://deno.land/install.sh | sh
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
EXPOSE 8080

FROM alpine:latest AS production
WORKDIR /app
COPY --from=build /build/{{SOSO_APPNAME}} /app
EXPOSE 8080
ENTRYPOINT [ "/app/{{SOSO_APPNAME}}" ]

