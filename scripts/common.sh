ROUTER='router.protected'
RTBH_SERVER='rtbh-server.protected'
ELK_SERVER='elk-server.protected'

function info {
    echo "[+] ${@}"
}

function error {
    echo "[E] ${@}"
    exit 1
}

function d_run {
    HOST="${1}"
    CMD="${2}"

    if [ -z "${HOST}" ]; then
        error "Usage: d_run SRC is empty"
    fi

    if [ -z "${CMD}" ]; then
        error "Usage: d_run DST is empty"
    fi

    info "run ${HOST} ${CMD}"
    ssh -oControlMaster=no ${HOST} ${CMD}
}

function d_copy {
    SRC="${1}"
    DST="${2}"

    if [ -z "${SRC}" ]; then
        error "Usage: d_copy SRC is empty"
    fi

    if [ -z "${DST}" ]; then
        error "Usage: d_copy DST is empty"
    fi

    info "rsync ${SRC} ${DST}"
    rsync -apl ${SRC} ${DST}
}
