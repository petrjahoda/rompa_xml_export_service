# ROMPA XML Export Service


## Installation
* use docker image from https://hub.docker.com/r/petrjahoda/rompa_xml_export_service
* use linux, mac or windows (using nssm) version and make it run like a service

## Settings
Settings are read from config.json.<br>
This file is created with default values, when not found.
* DatabaseType: "mysql"
* IpAddress:    "zapsidatabase"
* DatabaseName: "zapsi2"
* Port:         "3306"
* Login:        "zapsi_uzivatel"
* Password:     "zapsi"

## Tables used
* READ Workplace (WorkplaceGroupID=2)
* READ WorkplacePort
* READ User
* READ Order
* READ Product
* READ Package
* READ/WRITE Fail
* READ/WRITE Package
* READ/WRITE Product
* READ/WRITE Order
* READ/WRITE Terminal_input_package
* READ/WRITE Terminal_input_order
* READ/WRITE Terminal_input_order_terminal_input_fail
* READ/WRITE Terminal_input_fail
* READ/WRITE Rompa_workplace_data
    * **NEW NON-STANDARD TABLE !!**

## Description
Go service that process data every 10 seconds.<br>
Terminal_input_package are processed as boxes, new boxes are saved to mapped network folder in specific XML structure.<br>
Terminal_input_fail are processed as terminal fails, new fails are saved to mapped network folder in specific XML structure.<br>
Digital data are processed as machine fails, sum of new data is saved to mapped network folder in specific XML structure.<br>




www.zapsi.eu © 2020
