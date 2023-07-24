FROM golang:1.18 AS builder

COPY . /src
WORKDIR /src
RUN go env -w GO111MODULE=on
RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

#RUN echo "deb https://mirrors.aliyun.com/debian/ bullseye main non-free contrib" > /etc/apt/sources.list && \
#    echo "deb-src https://mirrors.aliyun.com/debian/ bullseye main non-free contrib" >> /etc/apt/sources.list && \
#    echo "deb https://mirrors.aliyun.com/debian-security/ bullseye-security main" >> /etc/apt/sources.list && \
#    echo "deb-src https://mirrors.aliyun.com/debian-security/ bullseye-security main" >> /etc/apt/sources.list && \
#    echo "deb https://mirrors.aliyun.com/debian/ bullseye-updates main non-free contrib" >> /etc/apt/sources.list && \
#    echo "deb-src https://mirrors.aliyun.com/debian/ bullseye-updates main non-free contrib" >> /etc/apt/sources.list && \
#    echo "deb https://mirrors.aliyun.com/debian/ bullseye-backports main non-free contrib" >> /etc/apt/sources.list && \
#    echo "deb-src https://mirrors.aliyun.com/debian/ bullseye-backports main non-free contrib" >> /etc/apt/sources.list \

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

CMD ["./server", "-conf", "/data/conf"]
#docker image build -t blmicros-api .

#docker run --name blmicros-api -p 8009:8000 -p 9000:9000 -v /Users/yufei/DockerVolume/blmicros/conf:/data/conf blmicros-api
# prod:
# docker run --name blmicros-api -p 8009:8000 -v /www/wwwroot/BrillinkDataAdmin/micros/brillinkmicros/configs:/data/conf blmicros-api
#