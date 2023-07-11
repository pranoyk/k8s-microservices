# k8s-microservices

### Issues faced

While deleting the cluster components often pv and pvc were stuck at terminating state.
On closer inspection realised that kubernetes adds a `finalisers: - kubernetes.io/pvc-protection` annotation along with a `deletionTimestamp`

This is added in order to protect the component from preventing leaks when deleting out of order (for ex - trying to delete parent before child will leave a dangling child)

in order to fix this we can edit the yaml
```bash
kubectl -n <namespace> edit pvc <pvc-name>
```
and delete the `finalisers: - kubernetes.io/pvc-protection`

This will delete the pvc. Same can be done if a pv is stuck at terminating state.


I also faced a scenario where namespace was stuck in the terminating state.
There are various ways to fix this
1. force delete the namespace - the problem with this is that it may leave certain components dangling without namespace, they may or may not get attached later when the same namespace is again created
2. on closer inspection we can see that the namespace yaml has 
```bash
  finalizers:
  - kubernetes
```
in it. In order to fix this we can delete it and reapply this yaml. In my case this did not work because there were some dangling resources were present in the namespace
3. run 
```bash
kubectl api-resources --verbs=list --namespaced -o name \
  | xargs -n 1 kubectl get --show-kind --ignore-not-found -n <stuck-namespace> 
```
this will result in all the dangling resources.
output : 
```
NAME                         DATA   AGE
configmap/kube-root-ca.crt   1      3d23h
LAST SEEN   TYPE      REASON    OBJECT                             MESSAGE
60m         Warning   BackOff   pod/k8s-demo-db-7679f76556-84brp   Back-off restarting failed container
60m         Warning   BackOff   pod/k8s-demo-db-7679f76556-gjxbx   Back-off restarting failed container
60m         Warning   BackOff   pod/k8s-demo-db-7679f76556-lmtl7   Back-off restarting failed container
NAME                     SECRETS   AGE
serviceaccount/default   0         3d23h
NAME                                                       ADDRESSTYPE   PORTS     ENDPOINTS    AGE
endpointslice.discovery.k8s.io/k8s-demo-db-service-f9cm2   IPv4          <unset>   <unset>      3d23h
endpointslice.discovery.k8s.io/postgresdb-z6znz            IPv4          5432      10.244.1.2   3d23h
LAST SEEN   TYPE      REASON    OBJECT                             MESSAGE
60m         Warning   BackOff   pod/k8s-demo-db-7679f76556-84brp   Back-off restarting failed container
60m         Warning   BackOff   pod/k8s-demo-db-7679f76556-gjxbx   Back-off restarting failed container
60m         Warning   BackOff   pod/k8s-demo-db-7679f76556-lmtl7   Back-off restarting failed container
```
deletion of the following resource and reapplying the yaml fixed the problem.
In order to reapply the namespace yaml without the finalisers start the `kubectl proxy` then run 
```bash
kubectl get namespace <stuck-namespace> -o json \
  | tr -d "\n" | sed "s/\"finalizers\": \[[^]]\+\]/\"finalizers\": []/" \
  | kubectl replace --raw /api/v1/namespaces/<stuck-namespace>/finalize -f -
```
this will delete the namespace