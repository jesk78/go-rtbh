#!/bin/bash

source $(basename $(dirname ${0}))/common.sh

DIR="$(basename $(dirname ${0}))"

${DIR}/stop_all.sh
d_run ${ROUTER} "> /var/log/auth.log"
d_run ${ROUTER} "> /var/log/suricata/eve.json"
d_run ${ELK_SERVER} "rm -rf /srv/elasticsearch/data/protected-es-cluster/*"
${DIR}/start_all.sh
