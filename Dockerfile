FROM golang:1.10.2 as build
WORKDIR /go/src/github.com/the-little-pedestrians/tlp-proxy
COPY . .

ENV DEP_VERSION=v0.4.1

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/${DEP_VERSION}/dep-linux-amd64
RUN chmod +x /usr/local/bin/dep
RUN dep ensure

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o tlp-proxy

FROM alpine:3.7
LABEL maintainer="Thomas Sauvajon <thomas.sauvajon.dev@gmail.com>"

ENV PORT=80

COPY --from=build /go/src/github.com/the-little-pedestrians/tlp-proxy/tlp-proxy /bin/tlp-proxy

RUN chmod +x /bin/tlp-proxy

CMD ["tlp-proxy"]

EXPOSE 80
