FROM golang:1.24.6-alpine3.22
WORKDIR /app
RUN apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*
RUN go install github.com/air-verse/air@latest
CMD ["air"]