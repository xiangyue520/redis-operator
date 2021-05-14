# redis-operator

## 说明

- 该operator是基于https://github.com/ucloud/redis-cluster-operator,不过我看有段时间没有维护了,自己拿到本地,发现运行不了,稍加修复和整改,并且发现原git是redis
-cluster-operator,自己的想法是改成redis-operator,让其支持单机,哨兵和集群模式
- 现在目前只支持集群模式,所以节点嘛,就最少3个了
- 现在很多都是原来的github上的,operator镜像以及其他的,会逐渐进行更换
- 将redis的镜像修改为6.0
- 平时吧,偶尔时间可能不太够,欢迎有想法的朋友联系我,一起努力,嗯,就这么多

## 环境准备
- 下载当前git项目到本地
- k8s中需要有helm,来进行安装operator
- 需要集群中有redis-cluster-operator的镜像(dockerhub可以拉取): fishu/redis-cluster-operator:latest

## 安装redis-cluster-operator并查看
1. 进入到当前项目,然后执行命令:`cd charts/redis-cluster-operator && helm install redis-cluster-operator
 .`,当出现一个安装完成后显示的版本信息,就说明创建好了![crd_install_success](/static/crd_install.png)
2. 完成安装后,就可以跳过接下来的步骤(Deploy redis cluster operator)了,
3. 执行`kubectl get po`,可以找到一个叫`redis-cluster-operator`的po,并且已经ready,那么说明就部署好了
4. 直接进入下面的[使用步骤](#usage),如果想要自己亲手搭建crd的话,参照[Deploy redis cluster operator](#deploy-redis-cluster-operator)

## Overview

Redis Cluster Operator manages [Redis Cluster](https://redis.io/topics/cluster-spec) atop Kubernetes.

The operator itself is built with the [Operator framework](https://github.com/operator-framework/operator-sdk).

![Redis Cluster atop Kubernetes](/static/redis-cluster.png)

Each master node and its slave nodes is managed by a statefulSet, create a headless svc for each statefulSet,
and create a clusterIP service for all nodes.

Each statefulset uses PodAntiAffinity to ensure that the master and slaves are dispersed on different nodes.
At the same time, when the operator selects the master in each statefulset, it preferentially select the pod
with different k8s nodes as master.

Table of Contents
=================

   * [redis-cluster-operator](#redis-cluster-operator)
      * [Overview](#overview)
   * [Table of Contents](#table-of-contents)
      * [Prerequisites](#prerequisites)
      * [Features](#features)
      * [Quick Start](#quick-start)
         * [Deploy redis cluster operator](#deploy-redis-cluster-operator)
            * [Install Step by step](#install-step-by-step)
            * [Install using helm chart](#install-using-helm-chart)
         * [Usage](#usage)
            * [Deploy a sample Redis Cluster](#deploy-a-sample-redis-cluster)
            * [Scaling Up the Redis Cluster](#scaling-up-the-redis-cluster)
            * [Scaling Down the Redis Cluster](#scaling-down-the-redis-cluster)
            * [Backup and Restore](#backup-and-restore)
            * [Prometheus Discovery](#prometheus-discovery)
            * [Create Redis Cluster with password](#create-redis-cluster-with-password)
            * [Persistent Volume](#persistent-volume)
            * [Custom Configuration](#custom-configuration)
            * [Custom Service](#custom-service)
            * [Custom Resource](#custom-resource)
      * [ValidatingWebhook](#validatingwebhook)
      * [End to end tests](#end-to-end-tests)

## Prerequisites

* go version v1.13+.
* Access to a Kubernetes v1.13.10 cluster.

## Features

- __Customize the number of master nodes and the number of replica nodes per master__

- __Password__

- __Safely Scaling the Redis Cluster__

- __Backup and Restore__

- __Persistent Volume__

- __Custom Configuration__

- __Prometheus Discovery__

## Quick Start

### Deploy redis cluster operator

#### Install Step by step

Register the DistributedRedisCluster and RedisClusterBackup custom resource definition (CRD).
```
$ kubectl create -f deploy/crds/redis.kun_distributedredisclusters_crd.yaml
$ kubectl create -f deploy/crds/redis.kun_redisclusterbackups_crd.yaml
```

A namespace-scoped operator watches and manages resources in a single namespace, whereas a cluster-scoped operator watches and manages resources cluster-wide.
You can chose run your operator as namespace-scoped or cluster-scoped.
```
// cluster-scoped
$ kubectl create -f deploy/service_account.yaml
$ kubectl create -f deploy/cluster/cluster_role.yaml
$ kubectl create -f deploy/cluster/cluster_role_binding.yaml
$ kubectl create -f deploy/cluster/operator.yaml

// namespace-scoped
$ kubectl create -f deploy/service_account.yaml
$ kubectl create -f deploy/namespace/role.yaml
$ kubectl create -f deploy/namespace/role_binding.yaml
$ kubectl create -f deploy/namespace/operator.yaml
```

#### Install using helm chart

Add Helm repository
```
helm repo add ucloud-operator https://ucloud.github.io/redis-cluster-operator/
helm repo update
```

Install chart
```
helm install --generate-name ucloud-operator/redis-cluster-operator
```

Verify that the redis-cluster-operator is up and running:
```
$ kubectl get deployment
NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
redis-cluster-operator   1/1     1            1           1d
```

### Usage
#### Deploy a sample Redis Cluster

>注意:创建的时候需要有redis的镜像,默认为redis:6.0

NOTE: **Only the redis cluster that use persistent storage(pvc) can recover after accidental deletion or rolling update.Even if you do not use persistence(like rdb or aof), you need to set pvc for redis.**

```
$ kubectl apply -f deploy/example/redis.kun_v1alpha1_distributedrediscluster_cr.yaml
```

Verify that the cluster instances and its components are running.
```
$ kubectl get distributedrediscluster或者使用kubectl get drc
NAME                              MASTERSIZE   STATUS    AGE
example-distributedrediscluster   3            Scaling   11s

$ kubectl get all -l redis.kun/name=example-distributedrediscluster
NAME                                          READY   STATUS    RESTARTS   AGE
pod/drc-example-distributedrediscluster-0-0   1/1     Running   0          2m48s
pod/drc-example-distributedrediscluster-0-1   1/1     Running   0          2m8s
pod/drc-example-distributedrediscluster-1-0   1/1     Running   0          2m48s
pod/drc-example-distributedrediscluster-1-1   1/1     Running   0          2m13s
pod/drc-example-distributedrediscluster-2-0   1/1     Running   0          2m48s
pod/drc-example-distributedrediscluster-2-1   1/1     Running   0          2m15s

NAME                                        TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)              AGE
service/example-distributedrediscluster     ClusterIP   172.17.132.71   <none>        6379/TCP,16379/TCP   2m48s
service/example-distributedrediscluster-0   ClusterIP   None            <none>        6379/TCP,16379/TCP   2m48s
service/example-distributedrediscluster-1   ClusterIP   None            <none>        6379/TCP,16379/TCP   2m48s
service/example-distributedrediscluster-2   ClusterIP   None            <none>        6379/TCP,16379/TCP   2m48s

NAME                                                     READY   AGE
statefulset.apps/drc-example-distributedrediscluster-0   2/2     2m48s
statefulset.apps/drc-example-distributedrediscluster-1   2/2     2m48s
statefulset.apps/drc-example-distributedrediscluster-2   2/2     2m48s

$ kubectl get distributedrediscluster
NAME                              MASTERSIZE   STATUS    AGE
example-distributedrediscluster   3            Healthy   4m
```

#### Scaling Up the Redis Cluster

Increase the masterSize to trigger the scaling up.

```
apiVersion: redis.kun/v1alpha1
kind: DistributedRedisCluster
metadata:
  annotations:
    # if your operator run as cluster-scoped, add this annotations
    redis.kun/scope: cluster-scoped
  name: example-distributedrediscluster
spec:
  # Increase the masterSize to trigger the scaling.
  masterSize: 4
  ClusterReplicas: 1
  image: redis:5.0.4-alpine
```

#### Scaling Down the Redis Cluster

Decrease the masterSize to trigger the scaling down.

```
apiVersion: redis.kun/v1alpha1
kind: DistributedRedisCluster
metadata:
  annotations:
    # if your operator run as cluster-scoped, add this annotations
    redis.kun/scope: cluster-scoped
  name: example-distributedrediscluster
spec:
  # Increase the masterSize to trigger the scaling.
  masterSize: 3
  ClusterReplicas: 1
  image: redis:5.0.4-alpine
```

#### Backup and Restore

NOTE: **Only Ceph S3 object storage and PVC is supported now**

Backup
```
$ kubectl create -f deploy/example/backup-restore/redisclusterbackup_cr.yaml
```

Restore from backup
```
$ kubectl create -f deploy/example/backup-restore/restore.yaml
```

#### Prometheus Discovery

```
$ kubectl create -f deploy/example/prometheus-exporter.yaml
```

#### Create Redis Cluster with password

```
$ kubectl create -f deploy/example/custom-password.yaml
```

#### Persistent Volume

```
$ kubectl create -f deploy/example/persistent.yaml
```

#### Custom Configuration

```
$ kubectl create -f deploy/example/custom-config.yaml
```

#### Custom Service

```
$ kubectl create -f deploy/example/custom-service.yaml
```

#### Custom Resource

```
$ kubectl create -f deploy/example/custom-resources.yaml
```

## ValidatingWebhook

see [ValidatingWebhook](/hack/webhook/README.md)

## End to end tests

see [e2e](/test/e2e/README.md)
