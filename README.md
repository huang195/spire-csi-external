# spire-csi

## pre-req

You will need a Kind cluster running locally.
You will need SPIRE installed in the Kind cluster.

## build image (optional)

```
make
```

## push image (optional)

```
make push
```

## deploy CSI

```
kubectl apply -f deploy/csi.yaml
```

## deploy workload

```
kubectl apply -f deploy/workload.yaml
```

## usage

Workloads need to include the following in the pod spec:

```
    volumeMounts:
    - name: csi-identity
      mountPath: /csi-identity
      readOnly: true
  volumes:
  - name: csi-identity
    csi:
      driver: "csi-identity.spiffe.io"
      readOnly: true
```

Once the pod is a running state, identity files can be found in `/csi-identity` directory, i.e.,

```
kubectl exec busybox -- ls /csi-identity
```

## clean up

Sometimes things are not cleaned up properly, e.g., a workload was force
deleted without its csi volumes fully cleaned up. This results in the apiserver
continuing to call the csi apis for these volumes that might no longer exists.

To properly clean these up, on the worker nodes, do

```
find /var/lib/kubelet/pods -name "csi-identity"
```

The pods that no longer exist, but its volume is still there will need to be
manually deleted so the apiserver will stop calling our csi webhook endpoint.
Note that sometimes that a `umount` is needed before the directory can be
manually deleted.
