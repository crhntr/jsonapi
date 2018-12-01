#!/usr/bin/env bash

mockgen -source="${PROJECT_DIR}/response_writers.go" -destination="${PROJECT_DIR}/mocks/response_writers_mocks.go"
