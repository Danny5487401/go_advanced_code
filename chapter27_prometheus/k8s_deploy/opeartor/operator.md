# prometheus Operator

## 背景
为了在Kubernetes能够方便的管理和部署Prometheus，我们使用ConfigMap了管理Prometheus配置文件。
每次对Prometheus配置文件进行升级时，我们需要手动移除已经运行的Pod实例，从而让Kubernetes可以使用最新的配置文件创建Prometheus。
而如果当应用实例的数量更多时，通过手动的方式部署和升级Prometheus过程繁琐并且效率低下。
```yaml
    scrape_configs:
      - job_name: 'prometheus'
        static_configs:
        - targets: ['localhost:9090']
```

从本质上来讲Prometheus属于是典型的有状态应用，而其有包含了一些自身特有的运维管理和配置管理方式。而这些都无法通过Kubernetes原生提供的应用管理概念实现自动化。
为了简化这类应用程序的管理复杂度，CoreOS率先引入了Operator的概念，并且首先推出了针对在Kubernetes下运行和管理Etcd的Etcd Operator。并随后推出了Prometheus Operator


## 工作原理
![](../../.operator_images/operator_discipline.png)

Prometheus的本职就是一组用户自定义的CRD资源以及Controller的实现，Prometheus Operator负责监听这些自定义资源的变化，并且根据这些资源的定义自动化的完成如Prometheus Server自身以及配置的自动化管理工作

## Operator能做什么

Prometheus Operator为我们提供了哪些自定义的Kubernetes资源，列出了Prometheus Operator目前提供的️资源：

- Prometheus：声明式创建和管理Prometheus Server实例；
- Alertmanager：声明式的创建和管理Alertmanager实例。
- ServiceMonitor：负责声明式的管理监控配置；
- PrometheusRule：负责声明式的管理告警配置；


还有thanosRuler,podMonitor,Probe等

## 安装
由于需要对Prometheus Operator进行RBAC授权，而默认的bundle.yaml中使用了default命名空间，因此，在安装Prometheus Operator之前需要先替换一下bundle.yaml文件中所有namespace定义，由default修改为monitoring。
```shell
$ kubectl -n monitoring apply -f bundle.yaml
clusterrolebinding.rbac.authorization.k8s.io/prometheus-operator created
clusterrole.rbac.authorization.k8s.io/prometheus-operator created
deployment.apps/prometheus-operator created
serviceaccount/prometheus-operator created
service/prometheus-operator created
```

```shell
$ kubectl -n monitoring get pods
NAME                                   READY     STATUS    RESTARTS   AGE
prometheus-operator-6db8dbb7dd-2hz55   1/1       Running   0          19s
```


##  操作
1. 部署正常的业务http 服务:deployment-app.yaml
```yaml
kind: Service
apiVersion: v1
metadata:
  name: danny-example-app
  labels:
    app: danny-example-app
spec:
  selector:
    app: danny-example-app
  ports:
    - name: web
      port: 8080
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: danny-example-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: danny-example-app
  template:
    metadata:
      labels:
        app: danny-example-app
    spec:
      containers:
        - name: danny-example-app
          image: fabxc/instrumented_app
          ports:
            - name: web
              containerPort: 8080
```

2. 部署监听服务serviceMonitor.yaml：使用matchlabels中的app: danny-example-app去监听业务服务
```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: danny-service-monitor-example-app
  namespace: monitoring
  labels:
    team: frontend
spec:
  # namespaceSelector定义让其可以跨命名空间
  namespaceSelector:
    matchNames:
      - danny-xia
  selector:
    matchLabels:
      app: danny-example-app
  endpoints:
    - port: web
```
3. 配置serviceAccount方便promethues拉取数据:主要是nonResourceURLs: ["/metrics"]获取资源
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: danny-prometheus
  namespace: monitoring
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: danny-prometheus-cluster-role
rules:
  - apiGroups: [""]
    resources:
      - nodes
      - services
      - endpoints
      - pods
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources:
      - configmaps
    verbs: ["get"]
  - nonResourceURLs: ["/metrics"]
    verbs: ["get"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: danny-prometheus-cluster-role-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: danny-prometheus-cluster-role
subjects:
  - kind: ServiceAccount
    name: danny-prometheus
    namespace: monitoring
```
4. 配置prometheus rules:alert_rules.yaml,其中设置标签prometheus: example，role: alert-rules
```
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  labels:
    prometheus: example
    role: alert-rules
  name: danny-prometheus-example-rules
  namespace: monitoring
spec:
  groups:
    - name: ./example.rules
      rules:
        - alert: ExampleAlert
          expr: vector(1)
```
5. 定义alert manager的全局配置secret:alertmanager.yaml,注意alertmanager-danny-alert-instance这名字是固定的，在默认情况下，会通过alertmanager-{ALERTMANAGER_NAME}的命名规则去查找Secret配置并以文件挂载的方式
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: alertmanager-danny-alert-instance
  namespace: monitoring
type: Opaque
stringData:
  alertmanager.yaml: |-
    global:
      resolve_timeout: 5m
    route:
      group_by: ['job']
      group_wait: 30s
      group_interval: 5m
      repeat_interval: 12h
      receiver: 'webhook'
    receivers:
      - name: 'webhook'
        webhook_configs:
          - url: 'http://alertmanagerwh:30500/'
```
6. 开启alert manager实例:alert-manager-instance.yaml
```yaml
apiVersion: monitoring.coreos.com/v1
kind: Alertmanager
metadata:
  name: danny-alert-instance
  namespace: monitoring
spec:
  replicas: 3
```
7. 暴露alert manager 的service:alertmanager-danny-alert-instance-srv.yaml
```yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"labels":{"alertmanager":"danny-alert-instance"},"name":"alertmanager-danny-alert-instance","namespace":"monitoring"},"spec":{"ports":[{"name":"web","port":9093,"targetPort":"web"}],"selector":{"alertmanager":"danny-alert-instance","app":"alertmanager"},"sessionAffinity":"ClientIP"}}
  labels:
    alertmanager: danny-alert-instance
  name: alertmanager-danny-alert-instance
  namespace: monitoring

spec:
  ports:
    - name: web
      port: 9093
      protocol: TCP
      targetPort: web
  selector:
    alertmanager: danny-alert-instance
    app: alertmanager


```


8. 定义prometheus实例：serviceAccountName指定步骤3的danny-prometheus账号，ruleSelector指定步骤4的规则，使用serviceMonitorSelector中的team: frontend去关联步骤2的monitor实例，alerting找步骤7暴露的endpoint
```yaml
apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: danny-instance
  namespace: monitoring
spec:
  serviceAccountName: danny-prometheus
  serviceMonitorSelector:
    matchLabels:
      team: frontend
  ruleSelector:
    matchLabels:
      role: alert-rules
      prometheus: example
  alerting:
    alertmanagers:
      - name: alertmanager-danny-alert-instance
        namespace: monitoring
        port: web
  resources:
    requests:
      memory: 400Mi

```

### 使用后效果
1. 配置变化
![](../../.operator_images/operator_effect_on_config.png)
2. 服务发现变化
![](../../.operator_images/operator_effect_on_discovery.png)
3. rules变化
![](../../.operator_images/operator_effect_on_rule.png)



Note:不使用账号会报没有权限
```shell
(⎈ |teleport.gllue.com-test:danny-xia)➜  go_advanced_code git:(feature/monitor) ✗ kubectl logs prometheus-danny-instance-0 prometheus -n monitoring --tail 2 
level=error ts=2022-05-05T08:17:55.023Z caller=klog.go:94 component=k8s_client_runtime func=ErrorDepth msg="/app/discovery/kubernetes/kubernetes.go:263: Failed to list *v1.Pod: pods is forbidden: User \"system:serviceaccount:monitoring:default\" cannot list resource \"pods\" in API group \"\" in the namespace \"danny-xia\""
level=error ts=2022-05-05T08:17:55.030Z caller=klog.go:94 component=k8s_client_runtime func=ErrorDepth msg="/app/discovery/kubernetes/kubernetes.go:262: Failed to list *v1.Service: services is forbidden: User \"system:serviceaccount:monitoring:default\" cannot list resource \"services\" in API group \"\" in the namespace \"danny-xia\""

```