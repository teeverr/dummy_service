FROM golang:1.16 AS build

ENV GOSUMDB="off"

WORKDIR /usr/local/src

COPY go.mod .
COPY go.sum .
COPY ./cmd/api/k8s/prod/configs/config.yml .
RUN go mod download

ADD . .

# Setup the compilation environment
ENV CGO_CPPFLAGS="-I/usr/include"
ENV CGO_LDFLAGS="-L/usr/lib -lpthread -lrt -lstdc++ -lm -lc -lgcc"
ENV CC="/usr/bin/gcc"
ENV CFLAGS="-march=x86-64"
ENV PKG_CONFIG_PATH="/usr/local/lib/pkgconfig"
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -v -mod readonly -o /go/bin/ ./cmd/...

FROM alpine:3.13

RUN apk add --no-cache curl

WORKDIR /

COPY --from=build /go/bin/api /opt/app/api
COPY --from=build /usr/local/src/config.yml /opt/app/config/config.yml

EXPOSE 8080

CMD ["/opt/app/api"]
