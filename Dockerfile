FROM golang:alpine AS build
WORKDIR /go/src
COPY ./app .

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

RUN apk add git gcc tzdata ca-certificates
RUN go build -a -installsuffix cgo -o app .


FROM scratch

COPY --from=build /go/src/app .
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENV TZ=Asia/Tokyo

EXPOSE 1323/tcp
ENTRYPOINT ["./app"]
