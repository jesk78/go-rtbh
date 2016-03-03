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

    if [ -z "${SRC}" ]; then
        error "Usage: d_copy SRC DST"
    fi

    if [ -z "${DST}" ]; then
        error "Usage: d_copy SRC DST"
    fi

    info "run ${HOST} ${CMD}"
    ssh -oControlMaster=no ${HOST} ${CMD}
}

function d_copy {
    SRC="${1}"
    DST="${2}"

    if [ -z "${SRC}" ]; then
        error "Usage: d_copy SRC DST"
    fi

    if [ -z "${DST}" ]; then
        error "Usage: d_copy SRC DST"
    fi

    info "rsync ${SRC} ${DST}"
    rsync -apl ${SRC} ${DST}
}
