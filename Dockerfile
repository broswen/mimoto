FROM golang:alpine AS build
# run build process in /app directory 
WORKDIR /app
# copy dependencies and get them
COPY ./go.mod ./go.sum ./
RUN go get -d -v ./...
# copy go src file(s)
COPY ./cmd ./cmd
COPY ./internal ./internal
# build the binary
RUN GOOS=linux GOARCH=amd64 go build -o mimoto ./cmd/main.go

FROM alpine
# copy the binary from build stage
COPY --from=build /app/mimoto /bin/mimoto
# use non root
USER 1000:1000
# start server
CMD ["./bin/mimoto"]