ARG GOLANG_VERSION=1.13.3
FROM golang:${GOLANG_VERSION}

ENV GOLANG_VERSION=${GOLANG_VERSION}

RUN apt update && apt install -y git

RUN go get -u github.com/onsi/ginkgo/ginkgo github.com/onsi/gomega/...

ARG PROJECT_NAME=redis-operator
ARG REPO_PATH=github.com/xiangyue520/$PROJECT_NAME

RUN mkdir -p /go/src/${REPO_PATH}
COPY . /go/src/${REPO_PATH}
RUN chmod +x /go/src/${REPO_PATH}/hack/e2e.sh

CMD /go/src/github.com/xiangyue520/redis-operator/hack/e2e.sh
