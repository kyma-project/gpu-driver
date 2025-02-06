# Find Gardenlinux version


### Option 1 - Using host filesystem

Mount host file `/etc/os-release` in the pod and read the Gardenlinux version from value GARDENLINUX_VERSION

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test
spec:
  os: { name: linux }
  nodeSelector:
    kubernetes.io/os: linux
  containers:
    - name: test
      image: ubuntu
      imagePullPolicy: IfNotPresent
      command:
        - "/bin/bash"
        - "-c"
        - "--"
      args:
        - "source /mnt/host/etc/os-release && echo $GARDENLINUX_VERSION"
      volumeMounts:
        - name: host
          mountPath: /mnt/host
          readOnly: true
  restartPolicy: Never
  volumes:
    - name: host
      hostPath:
        path: /
        type: Directory
```

Example for the `/etc/os-release` file

```
ID=gardenlinux
NAME="Garden Linux"
PRETTY_NAME="Garden Linux 1592.4"
HOME_URL="https://gardenlinux.io"
SUPPORT_URL="https://github.com/gardenlinux/gardenlinux"
BUG_REPORT_URL="https://github.com/gardenlinux/gardenlinux/issues"
GARDENLINUX_CNAME=aws-gardener_prod-amd64-1592.4
GARDENLINUX_FEATURES=log,sap,ssh,_boot,_nopkg,_prod,_slim,base,server,cloud,aws,gardener
GARDENLINUX_VERSION=1592.4
GARDENLINUX_COMMIT_ID=00942a0a
GARDENLINUX_COMMIT_ID_LONG=00942a0a5c2f25e89651e3f780b913619ffe5212
```

## Option 2 - read it from k8s node resource

```bash
kubectl get nodes -o=jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.nodeInfo.osImage}{"\n"}{end}'

ip-10-250-0-46.ec2.internal	Garden Linux 1592.4
ip-10-250-1-12.ec2.internal	Garden Linux 1592.4
ip-10-250-2-64.ec2.internal	Garden Linux 1592.4
```

## Options 3 - read it from the shoot resource

Take care about different node pools.

```
kind: Shoot
apiVersion: core.gardener.cloud/v1beta1
spec:
  provider:
    workers:
      - machine:
        type: n2-standard-2
        image:
          name: gardenlinux
          version: 1592.4.0
```