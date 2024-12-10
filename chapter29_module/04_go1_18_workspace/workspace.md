<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [workspace 工作区模式](#workspace-%E5%B7%A5%E4%BD%9C%E5%8C%BA%E6%A8%A1%E5%BC%8F)
  - [使用](#%E4%BD%BF%E7%94%A8)
  - [应用-->kubernetes](#%E5%BA%94%E7%94%A8--kubernetes)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


# workspace 工作区模式

多 Module 问题的解决方式

## 使用

```shell
➜ go help work
Work provides access to operations on workspaces.

Note that support for workspaces is built into many other commands, not
just 'go work'.

See 'go help modules' for information about Go's module system of which
workspaces are a part.

See https://go.dev/ref/mod#workspaces for an in-depth reference on
workspaces.

See https://go.dev/doc/tutorial/workspaces for an introductory
tutorial on workspaces.

A workspace is specified by a go.work file that specifies a set of
module directories with the "use" directive. These modules are used as
root modules by the go command for builds and related operations.  A
workspace that does not specify modules to be used cannot be used to do
builds from local modules.

go.work files are line-oriented. Each line holds a single directive,
made up of a keyword followed by arguments. For example:

	go 1.18

	use ../foo/bar
	use ./baz

	replace example.com/foo v1.2.3 => example.com/bar v1.4.5

The leading keyword can be factored out of adjacent lines to create a block,
like in Go imports.

	use (
	  ../foo/bar
	  ./baz
	)

The use directive specifies a module to be included in the workspace's
set of main modules. The argument to the use directive is the directory
containing the module's go.mod file.

The go directive specifies the version of Go the file was written at. It
is possible there may be future changes in the semantics of workspaces
that could be controlled by this version, but for now the version
specified has no effect.

The replace directive has the same syntax as the replace directive in a
go.mod file and takes precedence over replaces in go.mod files.  It is
primarily intended to override conflicting replaces in different workspace
modules.

To determine whether the go command is operating in workspace mode, use
the "go env GOWORK" command. This will specify the workspace file being
used.

Usage:

	go work <command> [arguments]

The commands are:

	edit        edit go.work from tools or scripts
	init        initialize workspace file
	sync        sync workspace build list to modules
	use         add modules to workspace file

Use "go help work <command>" for more information about a command.
```

## 应用-->kubernetes

```go
// https://github.com/kubernetes/kubernetes/blob/v1.31.0/go.work

// This is a generated file. Do not edit directly.

go 1.22.0

use (
	.
	./staging/src/k8s.io/api
	./staging/src/k8s.io/apiextensions-apiserver
	./staging/src/k8s.io/apimachinery
	./staging/src/k8s.io/apiserver
	./staging/src/k8s.io/cli-runtime
	./staging/src/k8s.io/client-go
	./staging/src/k8s.io/cloud-provider
	./staging/src/k8s.io/cluster-bootstrap
	./staging/src/k8s.io/code-generator
	./staging/src/k8s.io/component-base
	./staging/src/k8s.io/component-helpers
	./staging/src/k8s.io/controller-manager
	./staging/src/k8s.io/cri-api
	./staging/src/k8s.io/cri-client
	./staging/src/k8s.io/csi-translation-lib
	./staging/src/k8s.io/dynamic-resource-allocation
	./staging/src/k8s.io/endpointslice
	./staging/src/k8s.io/kms
	./staging/src/k8s.io/kube-aggregator
	./staging/src/k8s.io/kube-controller-manager
	./staging/src/k8s.io/kube-proxy
	./staging/src/k8s.io/kube-scheduler
	./staging/src/k8s.io/kubectl
	./staging/src/k8s.io/kubelet
	./staging/src/k8s.io/metrics
	./staging/src/k8s.io/mount-utils
	./staging/src/k8s.io/pod-security-admission
	./staging/src/k8s.io/sample-apiserver
	./staging/src/k8s.io/sample-cli-plugin
	./staging/src/k8s.io/sample-controller
)
```


## 参考
- [通过一个例子让你彻底掌握 Go 工作区模式](https://polarisxu.studygolang.com/posts/go/workspace/)