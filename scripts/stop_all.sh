#!/bin/bash

source $(basename $(dirname ${0}))/common.sh

d_run ${ROUTER} "systemctl stop logshipper"
d_run ${ROUTER} "systemctl stop suricata"


d_run ${ELK_SERVER} "systemctl stop kibana"
d_run ${ELK_SERVER} "systemctl stop logstash"
d_run ${ELK_SERVER} "systemctl stop rabbitmq-server"
d_run ${ELK_SERVER} "systemctl stop elasticsearch"
