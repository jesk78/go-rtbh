#!/bin/bash

source $(basename $(dirname ${0}))/common.sh

d_run_${ROUTER}

d_run ${ELK_SERVER} "systemctl stop elasticsearch"
d_run ${ELK_SERVER} "rm -rf /srv/elasticsearch/data/protected-es-cluster/*"
d_run ${ELK_SERVER} "systemctl start elasticsearch"
