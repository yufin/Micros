FROM golang:1.21 AS builder

COPY . /src
WORKDIR /src
RUN go env -w GO111MODULE=on
RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 8000
#EXPOSE 9000
VOLUME /data/conf
VOLUME /data/certs

CMD ["./server", "-conf", "/data/conf/config_prod.yaml"]

#docker image build -t blmicros-dw:1.5 .

# docker run --name blmicros-api -p 8009:8000 -p 9000:9000 -v /Users/yufei/DockerVolume/blmicros/conf:/data/conf blmicros-api
# prod:
# docker run --name micros-dw -p 50052:50052 -v /www/wwwroot/BrillinkDataAdmin/brillinkMicros/micros-dw/configs:/data/conf blmicros-dw:1.5