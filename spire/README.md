# SPIRE

## Installation

```
helm upgrade --install --create-namespace -n spire-mgmt spire-crds spire-crds --repo https://spiffe.github.io/helm-charts-hardened/
helm upgrade --install -n spire-mgmt spire spire --repo https://spiffe.github.io/helm-charts-hardened/ -f values.yaml
```

In case we want to disable container selectors to allow attestation to happen before containers are in a runnable state, run the following:
```
helm upgrade --install --create-namespace -n spire-mgmt spire-crds spire-crds --repo https://spiffe.github.io/helm-charts-hardened/
helm upgrade --install -n spire-mgmt spire spire --repo https://spiffe.github.io/helm-charts-hardened/ -f values.yaml --set spire-agent.workloadAttestors.k8s.disableContainerSelectors=true
```

## Add static workload entry

In case we want to avoid having to wait for the SPIRE controller manager to register workloads when their pods are starting, we can create static workload entries:

```
kubectl -n spire-server exec -it spire-server-0 -- spire-server entry create -spiffeID spiffe://spire-lifecycle-hook.ibm.com/ns/default/sa/default -parentID  spiffe://spire-lifecycle-hook.ibm.com/spire/agent/k8s_psat/spire-lifecycle-hook/226f3559-ffc0-45fb-bc5c-23d8e82c2a04 -selector k8s:ns:default -selector k8s:sa:default 
```

## Uninstallation

```
helm uninstall -n spire-mgmt spire
helm uninstall -n spire-mgmt spire-crds
kubectl delete namespace spire-mgmt

kubectl delete crd clusterfederatedtrustdomains.spire.spiffe.io
kubectl delete crd clusterspiffeids.spire.spiffe.io
kubectl delete crd clusterstaticentries.spire.spiffe.io
```
