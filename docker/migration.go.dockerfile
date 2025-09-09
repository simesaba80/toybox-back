FROM golang:1.24.6-alpine3.22
WORKDIR /app
RUN apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN apk --no-cache add ca-certificates postgresql-client bash curl tzdata

RUN cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone

CMD ["go", "run", "cmd/movedata/main.go"]