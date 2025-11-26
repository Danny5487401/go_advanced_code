package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"net"
)

func main() {
	domain := "www.baidu.com"
	// 返回 ip 类型
	fmt.Println(net.LookupIP(domain))            // [183.2.172.177 183.2.172.17 240e:ff:e020:98c:0:ff:b061:c306 240e:ff:e020:99b:0:ff:b099:cff1] <nil>
	fmt.Println(net.LookupHost(domain))          // [183.2.172.177 183.2.172.17 240e:ff:e020:98c:0:ff:b061:c306 240e:ff:e020:99b:0:ff:b099:cff1] <nil>
	fmt.Println(net.LookupAddr("142.250.190.4")) // [ord37s32-in-f4.1e100.net.] <nil>

	fmt.Println(net.LookupCNAME(domain)) // www.a.shifen.com. <nil>

	/*
		k8s集群内执行
		(⎈|kubeasz-test:monitoring)➜  ~ kubectl get svc -n default kubernetes -o yaml
		apiVersion: v1
		kind: Service
		metadata:
		  creationTimestamp: "2025-08-20T06:48:08Z"
		  labels:
		    component: apiserver
		    provider: kubernetes
		  name: kubernetes
		  namespace: default
		  resourceVersion: "227"
		  uid: 437e4917-b093-45ec-8664-ac3ae843fb04
		spec:
		  clusterIP: 10.233.0.1
		  clusterIPs:
		  - 10.233.0.1
		  internalTrafficPolicy: Cluster
		  ipFamilies:
		  - IPv4
		  ipFamilyPolicy: SingleStack
		  ports:
		  - name: https
		    port: 443
		    protocol: TCP
		    targetPort: 6443
		  sessionAffinity: None
		  type: ClusterIP
	*/
	cname, addresses, err := net.LookupSRV("", "", "_https._tcp.kubernetes.default.svc.cluster.local")
	if err != nil {
		fmt.Printf("failed: %s\n", err)
		return
	}
	fmt.Printf("cname: %s\n", cname) //cname: _https._tcp.kubernetes.default.svc.cluster.local.
	spew.Dump(addresses)
	/*
		([]*net.SRV) (len=1 cap=1) {
		 (*net.SRV)(0xc0001181e0)({
		  Target: (string) (len=37) "kubernetes.default.svc.cluster.local.",
		  Port: (uint16) 443,
		  Priority: (uint16) 0,
		  Weight: (uint16) 100
		 })
		}

	*/

}
