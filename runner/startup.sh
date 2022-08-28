#!/bin/bash
source logger.bash
source graceful-stop.bash

sudo /bin/bash <<SCRIPT
mkdir -p /etc/docker

if [ ! -f /etc/docker/daemon.json ]; then
  echo "{}" > /etc/docker/daemon.json
fi

if [ -n "${MTU}" ]; then
jq ".\"mtu\" = ${MTU}" /etc/docker/daemon.json > /tmp/.daemon.json && mv /tmp/.daemon.json /etc/docker/daemon.json
# See https://docs.docker.com/engine/security/rootless/
echo "environment=DOCKERD_ROOTLESS_ROOTLESSKIT_MTU=${MTU}" >> /etc/supervisor/conf.d/dockerd.conf
fi

if [ -n "${DOCKER_REGISTRY_MIRROR}" ]; then
jq ".\"registry-mirrors\"[0] = \"${DOCKER_REGISTRY_MIRROR}\"" /etc/docker/daemon.json > /tmp/.daemon.json && mv /tmp/.daemon.json /etc/docker/daemon.json
fi
SCRIPT

dumb-init bash <<'SCRIPT' &
source logger.bash
source wait.bash

dump() {
  local path=${1:?missing required <path> argument}
  shift
  printf -- "%s\n---\n" "${*//\{path\}/"$path"}" 1>&2
  cat "$path" 1>&2
  printf -- '---\n' 1>&2
}

for config in /etc/docker/daemon.json /etc/supervisor/conf.d/dockerd.conf; do
  dump "$config" 'Using {path} with the following content:'
done

log.debug 'Starting supervisor daemon'
sudo /usr/bin/supervisord -n >> /dev/null 2>&1 &

log.debug 'Waiting for processes to be running...'
processes=(dockerd)

for process in "${processes[@]}"; do
    wait_for_process "$process"
    if [ $? -ne 0 ]; then
        log.error "$process is not running after max time"
        dump /var/log/dockerd.err.log 'Dumping {path} to aid investigation'
        exit 1
    else
        log.debug "$process is running"
    fi
done

if [ -n "${MTU}" ]; then
  sudo ifconfig docker0 mtu ${MTU} up
fi

entrypoint.sh
SCRIPT

runner_init_pid=$!
log.notice "Runner init started with pid $runner_init_pid"
wait $runner_init_pid
log.notice "Runner init exited. Exiting this process with code 0 so that the container and the pod is GC'ed Kubernetes soon."
