#!/bin/bash
source logger.bash
source graceful-stop.bash

log.notice "Writing out Docker config file"
/bin/bash <<SCRIPT
mkdir -p /home/runner/.config/docker/

if [ ! -f /home/runner/.config/docker/daemon.json ]; then
  echo "{}" > /home/runner/.config/docker/daemon.json
fi

if [ -n "${MTU}" ]; then
jq ".\"mtu\" = ${MTU}" /home/runner/.config/docker/daemon.json > /tmp/.daemon.json && mv /tmp/.daemon.json /home/runner/.config/docker/daemon.json
# See https://docs.docker.com/engine/security/rootless/
echo "environment=DOCKERD_ROOTLESS_ROOTLESSKIT_MTU=${MTU}" >> /etc/supervisor/conf.d/dockerd.conf
fi

if [ -n "${DOCKER_REGISTRY_MIRROR}" ]; then
jq ".\"registry-mirrors\"[0] = \"${DOCKER_REGISTRY_MIRROR}\"" /home/runner/.config/docker/daemon.json > /tmp/.daemon.json && mv /tmp/.daemon.json /home/runner/.config/docker/daemon.json
fi
SCRIPT

log.notice "Starting Docker (rootless)"

dumb-init bash <<'SCRIPT' &
# Note that we don't want dockerd to be terminated before the runner agent,
# because it defeats the goal of the runner agent graceful stop logic implemenbed above.
# We can't rely on e.g. `dumb-init --single-child` for that, because with `--single-child` we can't even trap SIGTERM
# for not only dockerd but also the runner agent.
/home/runner/bin/dockerd-rootless.sh --config-file /home/runner/.config/docker/daemon.json >> /dev/null 2>&1 &

entrypoint.sh
SCRIPT

runner_init_pid=$!
log.notice "Runner init started with pid $runner_init_pid"
wait $runner_init_pid
log.notice "Runner init exited. Exiting this process with code 0 so that the container and the pod is GC'ed Kubernetes soon."
