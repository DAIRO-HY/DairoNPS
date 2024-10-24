#!/bin/bash

#删除上次编译文件
rm DairoNPS.zip
rm -rf DairoNPS-main
rm /app/DairoNPS/datadairo-nps

curl -L -o DairoNPS.zip https://github.com/DAIRO-HY/DairoNPS/archive/refs/heads/main.zip
unzip DairoNPS.zip
cd DairoNPS-main

#开始编译 由于使用了sqlite插件，编译时需要指定参数CGO_ENABLED=1
CGO_ENABLED=1 go build -o /app/DairoNPS/dairo-nps-linux-amd64

/app/DairoNPS/dairo-nps-linux-amd64