# S3 FUSE Flex Volume Drivers

[![Docker Image](https://img.shields.io/badge/docker-ready-blue.svg)](https://hub.docker.com/r/informaticslab/s3-fuse-flex-volume/) [![Docker Layers](https://images.microbadger.com/badges/image/informaticslab/s3-fuse-flex-volume.svg)](https://microbadger.com/#/images/informaticslab/s3-fuse-flex-volume)

This helm chart adds S3 FUSE flex volume drivers to your kubernetes cluster.

The flex volume drivers require the `fuse` package and the S3 fuse libraries to be installed on the host nodes, the chart assumes the hosts are running ubuntu and uses a privileged container to install these. It then installs the flex volume drivers.

This chart requires Kubernetes 1.8+ as previous versions require the `kubelet` to be restarted to pick up new flex volume drivers.

Included S3 FUSE libraries:
 - [pysssix](https://github.com/met-office-lab/pysssix)
 - [goofys](https://github.com/kahing/goofys)

## Installation

```
cd helm-chart
helm install --namespace kube-system --name s3-fuse-deployer s3-fuse-flex-volume
```

This helm chart will create a `DaemonSet` which uses privileged containers to install the fuse dependancies and the flex drivers on the kubernetes nodes. You are then able to use the drivers in your pod definitions.
## Usage examples

### pysssix

Pysssix will mount "all" of S3 which is accessible to the authenticating user. A mount point is created which referrs to all of S3 and then you access objects at `/<mount>/<bucket>/<object>`.

With this driver you are limited to read only.

```yaml
volumes:
  - name: pysssix
    flexVolume:
      driver: "informaticslab/pysssix-flex-volume"
      options:
        # Optional
        subPath: "key/prefix"
containers:
  - name: mycontainer
    ...
    volumeMounts:
      - name: pysssix
        mountPath: /s3
```

### goofys

Goofys will only mount a specific bucket so you must provide the `bucket` option.

```yaml
volumes:
  - name: goofys-mybucket
    flexVolume:
      driver: "informaticslab/goofys-flex-volume"
      options:
        # Required
        bucket: "mybucket"
        # Optional
        dirMode: "0755"
        fileMode: "0644"
        subPath: "key/prefix"
        access-key: "XXXXXXXXXXXXXXXXXXXX"
        secret-key: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
containers:
  - name: mycontainer
    ...
    volumeMounts:
      - name: goofys-mybucket
        mountPath: /s3/mybucket
```
