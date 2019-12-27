#!/usr/bin/env bash
docker rmi -f petrjahoda/rompa_xml_export_service:"$1"
docker build -t petrjahoda/rompa_xml_export_service:"$1" .
docker push petrjahoda/rompa_xml_export_service:"$1"