#!/bin/bash

set -e

APP_NAME="demojob"
DOCKER_IMAGE="halfcoconut/hello-egg/${APP_NAME}-live:v0.0.1"
DOCKERFILE_PATH="."
K8S_YAML="demojob.yaml"

echo "开始部署..."

echo "1. 编译 Go 程序..."
if ! GOOS=linux GOARCH=arm go build -o ${APP_NAME} .; then
  echo "编译 Go 程序失败"
  exit 1
fi

echo "2. 构建 Docker 镜像..."
if ! docker build -t ${DOCKER_IMAGE} ${DOCKERFILE_PATH}; then
  echo "构建 Docker 镜像失败"
  exit 1
fi

echo "3. 推送 Docker 镜像 (可选)..."
# docker push ${DOCKER_IMAGE}

echo "4. 应用 Kubernetes YAML 文件..."
if ! kubectl apply -f ${K8S_YAML}; then
  echo "应用 Kubernetes YAML 文件失败"
  exit 1
fi

echo "5. 监听 Kubernetes Job 状态..."
kubectl get jobs --watch

echo "部署完成!"

