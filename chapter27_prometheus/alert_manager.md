<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [alert告警](#alert%E5%91%8A%E8%AD%A6)
  - [操作](#%E6%93%8D%E4%BD%9C)
  - [流程](#%E6%B5%81%E7%A8%8B)
  - [应用举例](#%E5%BA%94%E7%94%A8%E4%B8%BE%E4%BE%8B)
    - [1. 告警分组](#1-%E5%91%8A%E8%AD%A6%E5%88%86%E7%BB%84)
    - [2. 告警抑制Inhibition](#2-%E5%91%8A%E8%AD%A6%E6%8A%91%E5%88%B6inhibition)
    - [3. silences告警静默](#3-silences%E5%91%8A%E8%AD%A6%E9%9D%99%E9%BB%98)
  - [Alertmanager的配置](#alertmanager%E7%9A%84%E9%85%8D%E7%BD%AE)
    - [route](#route)
    - [接收人（receivers)](#%E6%8E%A5%E6%94%B6%E4%BA%BAreceivers)
      - [集成邮件系统](#%E9%9B%86%E6%88%90%E9%82%AE%E4%BB%B6%E7%B3%BB%E7%BB%9F)
  - [参考链接](#%E5%8F%82%E8%80%83%E9%93%BE%E6%8E%A5)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# alert告警
![](.alert_images/alert_property.png)

在Prometheus的架构中被划分成两个独立的部分。Prometheus负责产生告警，而AlertManager负责告警产生后的后续处理.

## 操作
1. 搭建配置 Alertmanager
2. Prometheus 与 Alertmanager进行交流
3. 在Prometheus创建 alerting rules

## 流程
![](.alert_images/alert_process.png)

通过在Prometheus中定义AlertRule（告警规则），Prometheus会周期性的对告警规则进行计算，如果满足告警触发条件就会向Alertmanager发送告警信息。

## 应用举例

### 1. 告警分组
分组机制可以将某一类型的告警信息合并成一个大的告警信息，避免发送太多的告警邮件。

我们有3台服务器都介入了Prometheus，这3台服务器同时宕机了，那么如果不分组可能会发送3个告警信息，如果分组了，那么会合并成一个大的告警信息


1. 定义告警规则:监控服务器宕机的时间超过1分钟就发送告警邮件。
```yaml
groups:
- name: Test-Group-001 # 组的名字，在这个文件中必须要唯一
  rules:
  - alert: InstanceDown # 告警的名字，在组中需要唯一
    expr: up == 0 # 表达式, 执行结果为true: 表示需要告警
    for: 1m # 超过多少时间才认为需要告警(即up==0需要持续的时间)
    labels:
      severity: warning # 定义标签
    annotations:
      summary: "服务 {{ $labels.instance }} 下线了"
      description: "{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 1 minutes."

```
在每一个group中我们可以定义多个告警规则(rule)。一条告警规则主要由以下几部分组成：
- alert：告警规则的名称。
- expr：基于PromQL表达式告警触发条件，用于计算是否有时间序列满足该条件。
- for：评估等待时间，可选参数。用于表示只有当触发条件持续一段时间后才发送告警。在等待期间新产生告警的状态为pending。
- labels：自定义标签，允许用户指定要附加到告警上的一组附加标签。
- annotations：用于指定一组附加信息，比如用于描述告警详细信息的文字等，annotations的内容在告警产生时会一同作为参数发送到Alertmanager。
  在告警规则文件的annotations中使用summary描述告警的概要信息，description用于描述告警的详细信息。
  同时Alertmanager的UI也会根据这两个标签值，显示告警信息。为了让告警信息具有更好的可读性，Prometheus支持模板化label和annotations的中标签的值

通过$labels.<labelname>变量可以访问当前告警实例中指定标签的值。$value则可以获取当前PromQL表达式计算的样本值。


2. alertmanager.yml配置
```yaml
global:
  resolve_timeout: 5m
  # 整合qq邮件
  smtp_smarthost: 'smtp.qq.com:465'
  smtp_from: '1451578387@qq.com'
  smtp_auth_username: '1451578387@qq.com'
  smtp_auth_identity: 'xxxxxx'
  smtp_auth_password: 'xxxxxx'
  smtp_require_tls: false 
# 路由  
route:
  group_by: ['alertname'] # 根据什么分组，此处配置的是根据告警的名字分组,没有指定 group_by 貌似是根据规则文件的 groups[n].name 来分组的。
  group_wait: 10s # 当产生一个新分组时，告警信息需要等到 group_wait 才可以发送出去。
  group_interval: 10s # 如果上次告警信息发送成功，此时又来了一个新的告警数据，则需要等待 group_interval 才可以发送出去
  repeat_interval: 120s # 如果上次告警信息发送成功，且问题没有解决，则等待 repeat_interval 再次发送告警数据
  receiver: 'email' # 告警的接收者，需要和 receivers[n].name 的值一致。
receivers:
- name: 'email'
  email_configs:
  - to: '1451578387@qq.com'

```
- group_by :alertmanager可以对告警通知进行分组，将多条告警合合并为一个通知。这里我们可以使用group_by来定义分组规则。
基于告警中包含的标签，如果满足group_by中定义标签名称，那么这些告警将会合并为一个通知发送给接收器。

- group_wait: 有的时候为了能够一次性收集和发送更多的相关信息时，可以通过group_wait参数设置等待时间，如果在等待时间内当前group接收到了新的告警，这些告警将会合并为一个通知向receiver发送。
- group_interval :而group_interval配置，则用于定义相同的Group之间发送告警通知的时间间隔。

### 2. 告警抑制Inhibition

指的是当某类告警产生的时候，于此相关的别的告警就不用发送告警信息了。

我们对某台机器的CPU的使用率进行了监控，比如 使用到 80% 和 90% 都进行了监控，那么我们可能想如果CPU使用率达到了90%就不要发送80%的邮件了。

1. 告警规则
```yaml
groups:
- name: Cpu
  rules:
    - alert: Cpu01
      expr: "(1 - avg(irate(node_cpu_seconds_total{mode='idle'}[5m])) by (instance,job)) * 100 > 80"
      for: 1m
      labels:
        severity: info # 自定一个一个标签 info 级别
      annotations:
        summary: "服务 {{ $labels.instance }} cpu 使用率过高"
        description: "{{ $labels.instance }} of job {{ $labels.job }} 的 cpu 在过去5分钟内使用过高，cpu 使用率 {{humanize $value}}."
    - alert: Cpu02
      expr: "(1 - avg(irate(node_cpu_seconds_total{mode='idle'}[5m])) by (instance,job)) * 100 > 90"
      for: 1m
      labels:
        severity: warning # 自定一个一个标签 warning 级别
      annotations:
        summary: "服务 {{ $labels.instance }} cpu 使用率过高"
        description: "{{ $labels.instance }} of job {{ $labels.job }} 的 cpu 在过去5分钟内使用过高，cpu 使用率 {{humanize $value}}."

```

2. alertmanager.yml 配置抑制规则

如果 告警的名称 alertname = Cpu02 并且 告警级别 severity = warning ，那么抑制住 新的告警信息中 标签为 severity = info 的告警数据。并且源告警和目标告警数据的 instance 标签的值必须相等。
```yaml
# 抑制规则，减少告警数据
inhibit_rules:
- source_match: # 匹配当前告警规则后，抑制住target_match的告警规则
    alertname: Cpu02 # 标签的告警名称是 Cpu02
    severity: warning # 自定义的告警级别是 warning
  target_match: # 被抑制的告警规则
    severity: info # 抑制住的告警级别
  equal:
  - instance # source 和 target 告警数据中，instance的标签对应的值需要相等。

```

### 3. silences告警静默

指的是处于静默期，不发送告警信息。

我们系统某段时间进行停机维护，由此可能会产生一堆的告警信息，但是这个时候的告警信息是没有意义的，就可以配置静默规则过滤掉




## Alertmanager的配置

Alertmanager配置中一般会包含以下几个主要部分：
- 全局配置（global）：用于定义一些全局的公共参数，如全局的SMTP配置，Slack配置等内容；
- 模板（templates）：用于定义告警通知时的模板，如HTML模板，邮件模板等；
- 告警路由（route）：根据标签匹配，确定当前告警应该如何处理；
- 接收人（receivers）：接收人是一个抽象的概念，它可以是一个邮箱也可以是微信，Slack或者Webhook等，接收人一般配合告警路由使用；
- 抑制规则（inhibit_rules）：合理设置抑制规则可以减少垃圾告警的产生

主要介绍 路由(route)以及接收器(receivers)。所有的告警信息都会从配置中的顶级路由(route)进入路由树，根据路由规则将告警信息发送给相应的接收器。

### route
每一个告警都会从配置文件中顶级的route进入路由树，需要注意的是顶级的route必须匹配所有告警(即不能有任何的匹配设置match和match_re)，每一个路由都可以定义自己的接受人以及匹配规则。
默认情况下，告警进入到顶级route后会遍历所有的子节点，直到找到最深的匹配route，并将告警发送到该route定义的receiver中。
但如果route中设置continue的值为false，那么告警在匹配到第一个子节点之后就直接停止。如果continue为true，报警则会继续进行后续子节点的匹配。如果当前告警匹配不到任何的子节点，那该告警将会基于当前路由节点的接收器配置方式进行处理

```yaml
route:
  receiver: 'default-receiver'
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  group_by: [cluster, alertname]
  routes:
  - receiver: 'database-pager'
    group_wait: 10s
    match_re:
      service: mysql|cassandra
  - receiver: 'frontend-pager'
    group_by: [product, environment]
    match:
      team: frontend
```

默认情况下所有的告警都会发送给集群管理员default-receiver，因此在Alertmanager的配置文件的根路由中，对告警信息按照集群以及告警的名称对告警进行分组。


如果告警时来源于数据库服务如MySQL或者Cassandra，此时则需要将告警发送给相应的数据库管理员(database-pager)。
这里定义了一个单独子路由，如果告警中包含service标签，并且service为MySQL或者Cassandra,则向database-pager发送告警通知，由于这里没有定义group_by等属性，这些属性的配置信息将从上级路由继承，database-pager将会接收到按cluster和alertname进行分组的告警通知。


### 接收人（receivers)
每一个receiver具有一个全局唯一的名称，并且对应一个或者多个通知方式：

```yaml
name: <string>
email_configs:
  [ - <email_config>, ... ]
hipchat_configs:
  [ - <hipchat_config>, ... ]
pagerduty_configs:
  [ - <pagerduty_config>, ... ]
pushover_configs:
  [ - <pushover_config>, ... ]
slack_configs:
  [ - <slack_config>, ... ]
opsgenie_configs:
  [ - <opsgenie_config>, ... ]
webhook_configs:
  [ - <webhook_config>, ... ]
victorops_configs:
  [ - <victorops_config>, ... ]
```

目前官方内置的第三方通知集成包括：邮件、 即时通讯软件（如Slack、Hipchat）、移动应用消息推送(如Pushover)和自动化运维工具（例如：Pagerduty、Opsgenie、Victorops）。
Alertmanager的通知方式中还可以支持Webhook，通过这种方式开发者可以实现更多个性化的扩展支持。


#### 集成邮件系统
在Alertmanager中我们可以直接在配置文件的global中定义全局的SMTP配置：

```yaml
global:
  [ smtp_from: <tmpl_string> ]
  [ smtp_smarthost: <string> ]
  [ smtp_hello: <string> | default = "localhost" ]
  [ smtp_auth_username: <string> ]
  [ smtp_auth_password: <secret> ]
  [ smtp_auth_identity: <string> ]
  [ smtp_auth_secret: <secret> ]
  [ smtp_require_tls: <bool> | default = true ]
```

完成全局SMTP之后，我们只需要为receiver配置email_configs用于定义一组接收告警的邮箱地址即可，如下所示：
```yaml
name: <string>
email_configs:
  [ - <email_config>, ... ]
```





## 参考链接
1. https://prometheus.io/docs/alerting/latest/overview/