FROM golang:1.14-alpine as deps

ADD go.mod /app/go.mod
WORKDIR /app

RUN ["go", "mod", "download", "-x"]

FROM deps as build

ADD . /app
RUN go build -o "/bin/hhgw" ./cmd/server && \
    chmod a+x /bin/hhgw

FROM alpine

COPY --from=build /bin/hhgw /bin/hhgw

RUN addgroup -g 9999 -S user && \
    adduser -u 9999 -G user -S -H user

USER user
EXPOSE 8080
CMD ["/bin/hhgw"]
