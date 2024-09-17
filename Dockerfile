# Используем Go 1.23.0
FROM golang:1.23.0-alpine as builder

RUN set -x \
    && mkdir /src \
    && apk add --no-cache ca-certificates \
    && update-ca-certificates

WORKDIR /src

COPY ./go.* ./

# Загружаем и проверяем модули
RUN set -x \
    && go version \
    && go mod download \
    && go mod verify

COPY . /src

# Сборка приложения
RUN set -x \
    && go version \
    && GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /tmp/app .

# Образ для выполнения
FROM alpine:latest as runtime

RUN set -x \
    && adduser --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" --no-create-home --uid "10001" "appuser"

USER appuser:appuser

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /tmp/app /bin/app

ENTRYPOINT ["/bin/app"]
