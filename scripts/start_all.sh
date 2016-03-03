#!/bin/bash

source $(basename $(dirname ${0}))/common.sh

d_run ${ELK_SERVER} "systemctl start elasticsearch"
d_run ${ELK_SERVER} "systemctl start rabbitmq-server"
d_run ${ELK_SERVER} "systemctl start logstash"
d_run ${ELK_SERVER} "systemctl start kibana"
d_run ${ROUTER} "systemctl start suricata"
d_run ${ROUTER} "systemctl start logshipper"

KIBANA_SCRIPT="/tmp/.kibana_thingies"
cat > ${KIBANA_SCRIPT} <<EOF
echo '[+] Waiting for Kibana startup'
i=0
while :; do
    if [ \${i} -ge 30 ]; then
        echo "[E] Kibana startup failed!"
        exit 1
    else
        i=\$(expr \${i} + 1)
        sleep 1
    fi
    nc -z localhost 5601 && break
done
EOF
cat ${KIBANA_SCRIPT} | ssh ${ELK_SERVER} /bin/bash
rm -f ${KIBANA_SCRIPT}

ELK_SCRIPT="/tmp/.elk_thingies"
cat > ${ELK_SCRIPT} <<EOF

RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
NC="\033[0m"

echo '[+] Waiting for ES startup'
i=0
while :; do
    if [ \${i} -ge 30 ]; then
        echo "[E] Elasticsearch startup failed!"
        exit 1
    else
        i=\$(expr \${i} + 1)
        sleep 1
    fi
    nc -z localhost 9200 && break
done

echo '[+] Waiting for Kibana index creation'
NUM_FOUND=0
i=0
while :; do
    if [ \${i} -ge 30 ]; then
        echo "${RED}[E] Kibana index creation failed!${NC}"
        exit 1
    else
        i=\$(expr \${i} + 1)
        sleep 1
    fi

    NUM_FOUND=\$(curl 'localhost:9200/_cat/shards?pretty' 2>/dev/null | wc -l)
    if [ \${NUM_FOUND} -eq 2 ]; then
        break
    fi
done

echo '[+] Setting index.number_of_replicas: 0'
curl -XPUT 'localhost:9200/_settings' -d '{
    "index": {
        "number_of_replicas": 0
    }
}' &>/dev/null

echo -n '[+] Elasticsearch cluster is: '
STATUS="\$(curl -XGET 'localhost:9200/_cluster/health?pretty' 2>/dev/null| grep status | cut -d\" -f4)"
case "\${STATUS}" in
    "green")
        echo -e "\${GREEN}green\${NC}" ;;
    "yellow")
        echo -e "\${YELLOW}yellow\${NC}" ;;
    "red")
        echo -e "\${RED}red\${NC}" ;;
esac

EOF

cat ${ELK_SCRIPT} | ssh ${ELK_SERVER} /bin/bash
rm -f ${ELK_SCRIPT}
