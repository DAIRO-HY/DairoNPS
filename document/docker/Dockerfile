FROM frolvlad/alpine-glibc:glibc-2.34

RUN mkdir /app
WORKDIR /app

#下载程序
ADD https://github.com/DAIRO-HY/DairoNPS/releases/download/1.0.0/dairo-nps-linux-amd64 /app

#授权可执行
RUN chmod +x dairo-nps-linux-amd64

#启动脚本
ENTRYPOINT "/app/dairo-nps-linux-amd64"