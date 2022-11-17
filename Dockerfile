FROM golang:1.19-alpine3.16 AS build

RUN apk --no-cache add tzdata

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o app

FROM alpine:3.16

RUN apk --update add ca-certificates

WORKDIR /app

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /app/app .
COPY --from=build /app/config.yml config.yml

EXPOSE 8000
CMD [ "/app/app" ]
