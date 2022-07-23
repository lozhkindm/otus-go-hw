# Собираем в гошке
FROM golang:1.16.2 as build

ENV BIN_FILE /opt/calendar/goose
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/goose/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="goose"
LABEL MAINTAINERS="lozhkindm@yandex.ru"

ENV BIN_FILE "/opt/calendar/goose"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ARG CONFIG_FILE_NAME
ENV CONFIG_FILE /etc/calendar/.env
COPY ./configs/${CONFIG_FILE_NAME} ${CONFIG_FILE}

ENV MIGRATIONS_DIR /migrations
COPY ./migrations ${MIGRATIONS_DIR}

CMD ${BIN_FILE} --dir=${MIGRATIONS_DIR} --config=${CONFIG_FILE} up
