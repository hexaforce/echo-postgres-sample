FROM golang:alpine AS build

RUN apk add git gcc tzdata ca-certificates
RUN go install github.com/swaggo/swag/cmd/swag@latest

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src
COPY ./app app/
COPY ./go.mod .
COPY ./go.sum .

WORKDIR /go/src/app 
RUN swag init

WORKDIR /go/src
RUN go build -a -installsuffix cgo -o app app/main.go 

FROM scratch

COPY --from=build /go/src/app/main .
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENV TZ=Asia/Tokyo

EXPOSE 1323/tcp
ENTRYPOINT ["./main"]
