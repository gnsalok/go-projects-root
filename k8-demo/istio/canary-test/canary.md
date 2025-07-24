Below is a lightweight walk-through you can follow end-to-end on a single laptop that only has **Docker Desktop** installed.
The sequence is intentionally opinionated (Istio **demo** profile, a tiny “Hello” service, traffic shifting with `VirtualService` weights) so you can prove the concept in < 30 min and reuse the YAML in CI later.

---

## 0 . Prerequisites

| What                                                                                  | Why / notes                                                                                                                                        |
| ------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Docker Desktop ≥ 4.x**<br>with the built-in Kubernetes backend (v1.29-1.33) enabled | Turn on *Settings → Kubernetes → Enable Kubernetes* and **bump resources to ≥ 8 GiB RAM & 4 vCPU** so Istio + add-ons fit comfortably ([Istio][1]) |
| `kubectl` in your `$PATH`                                                             | Ships with Docker Desktop; verify with `kubectl version --short`                                                                                   |
| **Istio CLI (`istioctl`) v1.26.2**                                                    | Latest LTS as of July 2025 ([Istio][2])                                                                                                            |

```bash
curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.26.2 sh -
export PATH="$PATH:$HOME/istio-1.26.2/bin"
```

---

## 1 . Install Istio (demo profile)

```bash
istioctl install --set profile=demo -y
kubectl get pods -n istio-system   # wait until all are RUNNING
```

The *demo* profile gives you everything you need to demo traffic-splitting (ingress-gateway, Prometheus, Grafana, Kiali, Jaeger).

---

## 2 . Prepare a namespace

```bash
kubectl create namespace mesh-demo
kubectl label namespace mesh-demo istio-injection=enabled
```

Sidecars will now be injected automatically.

---

## 3 . Deploy two versions of the same service

Save as **hello.yaml** and apply:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata: {name: hello-v1, labels: {app: hello, version: v1}}
spec:
  replicas: 1
  selector: {matchLabels: {app: hello, version: v1}}
  template:
    metadata: {labels: {app: hello, version: v1}}
    spec:
      containers:
      - image: hashicorp/http-echo:0.2.3
        args: ["-text=Hello from v1"]
---
apiVersion: apps/v1
kind: Deployment
metadata: {name: hello-v2, labels: {app: hello, version: v2}}
spec:
  replicas: 1
  selector: {matchLabels: {app: hello, version: v2}}
  template:
    metadata: {labels: {app: hello, version: v2}}
    spec:
      containers:
      - image: hashicorp/http-echo:0.2.3
        args: ["-text=Hello from v2"]
---
apiVersion: v1
kind: Service
metadata: {name: hello}
spec:
  selector: {app: hello}
  ports: [{port: 80, name: http}]
```

```bash
kubectl apply -n mesh-demo -f hello.yaml
```

---

## 4 . Define traffic-control objects

### 4.1 DestinationRule — declare subsets **before** you route to them

Following Istio’s *make-before-break* best-practice ([Istio][3]) :

```yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata: {name: hello}
spec:
  host: hello.mesh-demo.svc.cluster.local
  subsets:
  - name: v1
    labels: {version: v1}
  - name: v2
    labels: {version: v2}
```

```bash
kubectl apply -n mesh-demo -f dr-hello.yaml
```

### 4.2 Gateway + edge VirtualService

Expose the service at `/hello` on the default ingress-gateway:

```yaml
# hello-gateway.yaml
apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata: {name: hello-gw}
spec:
  selector: {istio: ingressgateway}
  servers:
  - port: {number: 80, name: http, protocol: HTTP}
    hosts: ["*"]

---
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata: {name: hello-edge}
spec:
  hosts: ["*"]
  gateways: [hello-gw]
  http:
  - match: [{uri: {prefix: "/hello"}}]
    rewrite: {uri: "/"}
    route:
    - destination: {host: hello.mesh-demo.svc.cluster.local, subset: v1}
      weight: 100
```

```bash
kubectl apply -n mesh-demo -f hello-gateway.yaml
```

Retrieve the node-port on Docker Desktop:

```bash
export INGRESS_PORT=$(kubectl -n istio-system get svc istio-ingressgateway \
  -o jsonpath='{.spec.ports[?(@.port==80)].nodePort}')
export URL="http://localhost:${INGRESS_PORT}/hello"
watch -n1 curl -s $URL   # you should only see “Hello from v1”
```

---

## 5 . Shift traffic — your local canary

Create **vs-80-20.yaml** and apply it to start a canary at 20 %:

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata: {name: hello-edge}
spec:
  hosts: ["*"]
  gateways: [hello-gw]
  http:
  - match: [{uri: {prefix: "/hello"}}]
    rewrite: {uri: "/"}
    route:
    - destination: {host: hello.mesh-demo.svc.cluster.local, subset: v1}
      weight: 80
    - destination: {host: hello.mesh-demo.svc.cluster.local, subset: v2}
      weight: 20
```

```bash
kubectl apply -n mesh-demo -f vs-80-20.yaml
```

Run the `watch` loop again; roughly 1 in 5 responses should now say **v2**.

Increase or decrease weights at will (e.g., 50/50, 0/100) by patching the same `VirtualService`.
Because replica counts stay unchanged, you can test 1 % traffic with *one* Pod.

---

## 6 . Observe & validate

Istio’s add-ons are already running:

```bash
istioctl dashboard kiali     # topology + request distribution
istioctl dashboard prometheus
istioctl dashboard grafana
```

Watch error rates or latency for the `hello` workload; if the canary misbehaves, simply roll back:

```bash
kubectl apply -n mesh-demo -f hello-gateway.yaml   # 100 % back to v1
```

---

## 7 . Clean-up

```bash
kubectl delete namespace mesh-demo
istioctl uninstall -y
```

---

### What you just proved

* **Istio traffic splitting** lets you run *any* replica ratio (down to 1 %) without changing Deployment counts ([Istio][4]).
* The same pattern works identically in cloud clusters; swap the Gateway host from `*` to a DNS name and attach TLS.
* Following the *DestinationRule-first* order guarantees zero-downtime updates ([Istio][3]).

Happy shipping canaries 🎉

[1]: https://istio.io/latest/docs/setup/platform-setup/docker/?utm_source=chatgpt.com "Docker Desktop - Istio"
[2]: https://istio.io/latest/docs/releases/log/ "Istio / Website Content Changes"
[3]: https://istio.io/latest/docs/ops/best-practices/traffic-management/ "Istio / Traffic Management Best Practices"
[4]: https://istio.io/latest/blog/2017/0.1-canary/ "Istio / Canary Deployments using Istio"
