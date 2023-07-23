FROM golang:1.20-alpine as builder

COPY . /src

WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux go build -a -o kvs


FROM scratch

COPY --from=builder /src/kvs .

COPY --from=builder /src/*.pem .

EXPOSE 8080

CMD ["/kvs"]
