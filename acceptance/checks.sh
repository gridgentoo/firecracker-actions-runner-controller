#!/usr/bin/env bash

set +e

repo_runnerdeployment_passed="skipped"
repo_runnerset_passed="skipped"

echo "Checking if RunnerDeployment repo test is set"
if [ "${TEST_REPO}" ] && [ ! "${USE_RUNNERSET}" ]; then
  runner_name=
  count=0
  while [ $count -le 10 ]; do
    echo "Finding Runner ..."
    runner_name=$(kubectl get runner --output=jsonpath="{.items[*].metadata.name}")
    if [ "${runner_name}" ]; then
      while [ $count -le 10 ]; do
        runner_pod_name=
        echo "Found Runner \""${runner_name}"\""
        echo "Finding underlying pod ..."
        runner_pod_name=$(kubectl get pod --output=jsonpath="{.items[*].metadata.name}" | grep ${runner_name})
        if [ "${runner_pod_name}" ]; then
          echo "Found underlying pod \""${runner_pod_name}"\""
          echo "Waiting for pod \""${runner_pod_name}"\" to become ready..."
          kubectl wait pod/${runner_pod_name} --for condition=ready --timeout 270s
          break 2
        fi
        sleep 5
        let "count=count+1"
      done
    fi
    sleep 5
    let "count=count+1"
  done
  if [ $count -ge 10 ]; then
    echo "Timeout reached, failed to find the RunnerDeployment or underlying pod"
    repo_runnerdeployment_passed=false
  else
    repo_runnerdeployment_passed=true
  fi
fi

echo "Checking if RunnerSet repo test is set"
if [ "${TEST_REPO}" ] && [ "${USE_RUNNERSET}" ]; then
  runnerset_name=
  count=0
  while [ $count -le 10 ]; do
    echo "Finding RunnerSet ..."
    runnerset_name=$(kubectl get runnerset --output=jsonpath="{.items[*].metadata.name}")
    if [ "${runnerset_name}" ]; then
      while [ $count -le 10 ]; do
        runnerset_pod_name=
        echo "Found RunnerSet \""${runnerset_name}"\""
        echo "Finding underlying pod ..."
        runnerset_pod_name=$(kubectl get pod --output=jsonpath="{.items[*].metadata.name}" | grep ${runnerset_name})
        if [ "${runnerset_pod_name}" ]; then
          echo "Found underlying pod \""${runnerset_pod_name}"\""
          echo "Waiting for pod \""${runnerset_pod_name}"\" to become ready..."
          kubectl wait pod/${runnerset_pod_name} --for condition=ready --timeout 270s
          break 2
        fi
      sleep 5
      let "count=count+1"
      done
    fi
    sleep 5
    let "count=count+1"
  done
  if [ $count -ge 10 ]; then
    echo "Timeout reached, failed to find the RunnerSet or underlying pod"
    repo_runnerset_passed=false
  else
    repo_runnerset_passed=true
  fi
fi

if [ "${repo_runnerset_passed}" == "true" ] || [ "${repo_runnerset_passed}" == "skipped" ] && \
   [ "${repo_runnerdeployment_passed}" == "true" ] || [ "${repo_runnerdeployment_passed}" == "skipped" ]; then
  echo "INFO : All tests passed or skipped"
  echo "RunnerSet Repo Test Status : ${repo_runnerset_passed}"
  echo "RunnerDeployment Repo Test Status : ${repo_runnerdeployment_passed}"
else
  echo "ERROR : Some tests failed"
  echo "RunnerSet Repo Test Status : ${repo_runnerset_passed}"
  echo "RunnerDeployment Repo Test Status : ${repo_runnerdeployment_passed}"
  exit 1
fi