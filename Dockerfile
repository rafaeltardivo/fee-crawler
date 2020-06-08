FROM golang:1.14-buster as build

# build image
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN mkdir /build 
ADD . /build/
WORKDIR /build/main
RUN go build -o app .

# final image
FROM gcr.io/distroless/base-debian10
COPY --from=build /build/main /


