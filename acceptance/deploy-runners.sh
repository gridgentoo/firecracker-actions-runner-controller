#!/usr/bin/env bash

if [ "${TEST_REPO}" ]; then
  if [ "${USE_RUNNERSET}" ]; then
    echo "Deploying acceptance/testdata/repo.runnerset.yaml to ${TEST_REPO}"
    cat acceptance/testdata/repo.runnerset.yaml | envsubst | kubectl apply -f -
  else
    echo "Deploying acceptance/testdata/repo.runnerdeploy.yaml to ${TEST_REPO}"
    cat acceptance/testdata/repo.runnerdeploy.yaml | envsubst | kubectl apply -f -
    echo "Deploying acceptance/testdata/repo.hra.yaml to ${TEST_REPO}"
    cat acceptance/testdata/repo.hra.yaml | envsubst | kubectl apply -f -
  fi
else
  echo 'Skipped deploying Runnerdeployment / RunnerSet and HorizontalRunnerAutoscaler, set TEST_REPO to "your-org/your-repo" to deploy.'
fi

if [ "${TEST_ORG}" ]; then
  echo "Deploying acceptance/testdata/org.runnerdeploy.yaml to ${TEST_ORG}"
  cat acceptance/testdata/org.runnerdeploy.yaml | envsubst | kubectl apply -f -
  if [ "${TEST_ORG_REPO}" ]; then
    echo "Deploying acceptance/testdata/org.hra.yaml to ${TEST_ORG_REPO}"
    cat acceptance/testdata/org.hra.yaml | envsubst | kubectl apply -f -
  else
    echo 'Skipped deploying organizational HorizontalRunnerAutoscaler, set TEST_ORG_REPO to "yourorg/yourrepo" to deploy.'
  fi
else
  echo 'Skipped deploying organizational Runnerdeployment, set TEST_ORG to deploy.'
fi
