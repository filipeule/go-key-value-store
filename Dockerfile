FROM golang:1.24 AS build

COPY . /src

WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./cmd/api -a -tags netgo -ldflags '-s -w -extldflags "-static"' -o kvs .

FROM scratch

COPY --from=build /src/cmd/api/kvs .

COPY --from=build /src/cert/* ./cert/

EXPOSE 8080

CMD [ "/kvs" ]