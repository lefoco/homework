FROM golang:1.17 AS builder
ENV 	GO111MODULE=off \	
	CGO_ENABLED=0 \	
	GOOS=linux \	
	GOARCH=amd64 \
	GOPROXY=https://goproxy.cn
WORKDIR /build
COPY . .
RUN go build -o httpserver .

FROM scratch
COPY --from=builder /build/httpserver /
EXPOSE 80
ENTRYPOINT ["/httpserver"]
