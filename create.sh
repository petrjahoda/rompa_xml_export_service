#!/usr/bin/env bash
#docker rmi -f petrjahoda/rompa_xml_export_service:"$1"
#docker build -t petrjahoda/rompa_xml_export_service:"$1" .
#docker push petrjahoda/rompa_xml_export_service:"$1"


./update
name=${PWD##*/}
go get -u all
GOOS=linux go build -ldflags="-s -w" -o linux/"$name"
cd linux
upx "$name"
cd ..

docker rmi -f petrjahoda/"$name":latest
docker  build -t petrjahoda/"$name":latest .
docker push petrjahoda/"$name":latest

docker rmi -f petrjahoda/"$name":2020.4.2
docker build -t petrjahoda/"$name":2020.4.2 .
docker push petrjahoda/"$name":2020.4.2