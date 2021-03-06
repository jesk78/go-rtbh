#!/bin/bash

source $(basename $(dirname ${0}))/common.sh

PRE_DEPLOY_SCRIPT='/tmp/.pre_deploy_script'
PREFIX="/opt/go-rtbh"

## Script to run on the remote system to perform an installation
cat > ${PRE_DEPLOY_SCRIPT} <<EOF
for DIR in "${PREFIX}" "${PREFIX}/bin" "${PREFIX}/etc" "${PREFIX}/share"; do
    if [ ! -d "\${DIR}" ]; then
        mkdir -p "\${DIR}"
        echo "[+] mkdir \${DIR}"
    fi
done
EOF

## Perform the pre deployment
cat ${PRE_DEPLOY_SCRIPT} | ssh ${RTBH_SERVER} /bin/bash
rm -f ${PRE_DEPLOY_SCRIPT}

## Perform the software deployment
d_copy ./go-rtbh ${RTBH_SERVER}:${PREFIX}/bin
d_copy ./go-rtbh.yml ${RTBH_SERVER}:${PREFIX}/etc
d_copy ./templates/* ${RTBH_SERVER}:${PREFIX}/share

## Cleanup redis database
d_run ${RTBH_SERVER} "systemctl stop redis-server"
d_run ${RTBH_SERVER} "rm -f /var/lib/redis/dump.rdb"
d_run ${RTBH_SERVER} "systemctl start redis-server"

## Run the software
d_run ${RTBH_SERVER} "pkill go-rtbh"
d_run ${RTBH_SERVER} "${PREFIX}/bin/go-rtbh -f ${PREFIX}/etc/go-rtbh.yml -D"
