<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/testcontainers/testcontainers-go](#githubcomtestcontainerstestcontainers-go)
  - [使用](#%E4%BD%BF%E7%94%A8)
  - [第三方使用--grafana clickhouse 插件](#%E7%AC%AC%E4%B8%89%E6%96%B9%E4%BD%BF%E7%94%A8--grafana-clickhouse-%E6%8F%92%E4%BB%B6)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/testcontainers/testcontainers-go

在做集成测试的时候，每次测试前，如果通过docker重启一个干净的容器是不是免去了数据清理的苦恼。



其他方案:
- https://github.com/orlangure/gnomock, 使用 docker 作为临时容器,已经预制一些database, service or other tools.
- https://github.com/ory/dockertest

## 使用
使用技巧: 由于每次拉取镜像和启动docker代价比较大，比较耗时，我们一般在单测的入口TestMain方法里做初始化，也就是一个模块进行一次容器初始化。
由于单测case之间没有数据的清理，因此我们每个单测结束后都需要注意清理和还原数据。


容器配置
```go
type GenericContainerRequest struct {
	ContainerRequest              // embedded request for provider
	Started          bool         // whether to auto-start the container
	ProviderType     ProviderType // 默认docker ,可以是 podman 等
	Logger           Logging      // provide a container specific Logging - use default global logger if empty
	Reuse            bool         // 是否复用已经根据name过滤存在的容器
}


type ContainerRequest struct {
	FromDockerfile
	HostAccessPorts         []int
	Image                   string
	ImageSubstitutors       []ImageSubstitutor
	Entrypoint              []string
	Env                     map[string]string
	ExposedPorts            []string // allow specifying protocol info
	Cmd                     []string
	Labels                  map[string]string
	Mounts                  ContainerMounts
	Tmpfs                   map[string]string
	RegistryCred            string // Deprecated: Testcontainers will detect registry credentials automatically
	WaitingFor              wait.Strategy // readiness 等待策略
	Name                    string // for specifying container name
	Hostname                string
	WorkingDir              string                                     // specify the working directory of the container
	ExtraHosts              []string                                   // Deprecated: Use HostConfigModifier instead
	Privileged              bool                                       // For starting privileged container
	Networks                []string                                   // for specifying network names
	NetworkAliases          map[string][]string                        // for specifying network aliases
	NetworkMode             container.NetworkMode                      // Deprecated: Use HostConfigModifier instead
	Resources               container.Resources                        // Deprecated: Use HostConfigModifier instead
	Files                   []ContainerFile                            // files which will be copied when container starts
	User                    string                                     // for specifying uid:gid
	SkipReaper              bool                                       // Deprecated: The reaper is globally controlled by the .testcontainers.properties file or the TESTCONTAINERS_RYUK_DISABLED environment variable
	ReaperImage             string                                     // Deprecated: use WithImageName ContainerOption instead. Alternative reaper image
	ReaperOptions           []ContainerOption                          // Deprecated: the reaper is configured at the properties level, for an entire test session
	AutoRemove              bool                                       // Deprecated: Use HostConfigModifier instead. If set to true, the container will be removed from the host when stopped
	AlwaysPullImage         bool                                       // Always pull image
	ImagePlatform           string                                     // ImagePlatform describes the platform which the image runs on.
	Binds                   []string                                   // Deprecated: Use HostConfigModifier instead
	ShmSize                 int64                                      // Amount of memory shared with the host (in bytes)
	CapAdd                  []string                                   // Deprecated: Use HostConfigModifier instead. Add Linux capabilities
	CapDrop                 []string                                   // Deprecated: Use HostConfigModifier instead. Drop Linux capabilities
	ConfigModifier          func(*container.Config)                    // Modifier for the config before container creation
	HostConfigModifier      func(*container.HostConfig)                // Modifier for the host config before container creation
	EnpointSettingsModifier func(map[string]*network.EndpointSettings) // Modifier for the network settings before container creation
	LifecycleHooks          []ContainerLifecycleHooks                  // define hooks to be executed during container lifecycle
	LogConsumerCfg          *LogConsumerConfig                         // define the configuration for the log producer and its log consumers to follow the logs
}
```

发起启动
```go
func GenericContainer(ctx context.Context, req GenericContainerRequest) (Container, error) {
	if req.Reuse && req.Name == "" {
		return nil, ErrReuseEmptyName
	}

	logging := req.Logger
	if logging == nil {
		logging = Logger
	}
	provider, err := req.ProviderType.GetProvider(WithLogger(logging))
	if err != nil {
		return nil, fmt.Errorf("get provider: %w", err)
	}
	defer provider.Close()

	var c Container
	if req.Reuse {
		// we must protect the reusability of the container in the case it's invoked
		// in a parallel execution, via ParallelContainers or t.Parallel()
		reuseContainerMx.Lock()
		defer reuseContainerMx.Unlock()

		c, err = provider.ReuseOrCreateContainer(ctx, req.ContainerRequest)
	} else {
		c, err = provider.CreateContainer(ctx, req.ContainerRequest)
	}
	if err != nil {
        // ...
	}

	if req.Started && !c.IsRunning() {
		if err := c.Start(ctx); err != nil {
			return c, fmt.Errorf("start container: %w", err)
		}
	}
	return c, nil
}
```

## 第三方使用--grafana clickhouse 插件

```go
// https://github.com/grafana/clickhouse-datasource/blob/28f86d02d120e38a11fff363fac846224580550b/pkg/plugin/driver_test.go

func TestMain(m *testing.M) {
	useDocker := strings.ToLower(getEnv("CLICKHOUSE_USE_DOCKER", "true"))
	if useDocker == "false" {
		fmt.Printf("Using external ClickHouse for IT tests -  %s:%s\n",
			getEnv("CLICKHOUSE_PORT", "9000"), getEnv("CLICKHOUSE_HOST", "localhost"))
		os.Exit(m.Run())
	}
	// create a ClickHouse container
	ctx := context.Background()
	// attempt use docker for CI
	provider, err := testcontainers.ProviderDocker.GetProvider()
	if err != nil {
		fmt.Printf("Docker is not running and no clickhouse connections details were provided. Skipping IT tests: %s\n", err)
		os.Exit(0)
	}
	err = provider.Health(ctx)
	if err != nil {
		fmt.Printf("Docker is not running and no clickhouse connections details were provided. Skipping IT tests: %s\n", err)
		os.Exit(0)
	}
	chVersion := GetClickHouseTestVersion()
	fmt.Printf("Using Docker for IT tests with ClickHouse %s\n", chVersion)
	cwd, err := os.Getwd()
	if err != nil {
		// can't test without container
		panic(err)
	}

	customHostPath := "../../config/custom.xml"
	adminHostPath := "../../config/admin.xml"
	if chVersion == "21.8" {
		customHostPath = "../../config/custom.21.8.xml"
		adminHostPath = "../../config/admin.21.8.xml"
	}
	req := testcontainers.ContainerRequest{
		Env: map[string]string{
			"CLICKHOUSE_SKIP_USER_SETUP": "1", // added because of https://github.com/ClickHouse/ClickHouse/commit/65435a3db7214261b793acfb69388567b4c8c9b3
			"TZ":                         time.Local.String(),
		},
		ExposedPorts: []string{"9000/tcp", "8123/tcp"},
		Files: []testcontainers.ContainerFile{
			{
				ContainerFilePath: "/etc/clickhouse-server/config.d/custom.xml",
				FileMode:          0644,
				HostFilePath:      path.Join(cwd, customHostPath),
			},
			{
				ContainerFilePath: "/etc/clickhouse-server/users.d/admin.xml",
				FileMode:          0644,
				HostFilePath:      path.Join(cwd, adminHostPath),
			},
		},
		Image: fmt.Sprintf("clickhouse/clickhouse-server:%s", chVersion),
		Resources: container.Resources{
			Ulimits: []*units.Ulimit{
				{
					Name: "nofile",
					Hard: 262144,
					Soft: 262144,
				},
			},
		},
		WaitingFor: wait.ForLog("Ready for connections"),
	}
	clickhouseContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		// can't test without container
		panic(err)
	}
	p, _ := clickhouseContainer.MappedPort(ctx, "9000")
	os.Setenv("CLICKHOUSE_PORT", p.Port())
	hp, _ := clickhouseContainer.MappedPort(ctx, "8123")
	os.Setenv("CLICKHOUSE_HTTP_PORT", hp.Port())
	os.Setenv("CLICKHOUSE_HOST", "localhost")
	defer clickhouseContainer.Terminate(ctx) //nolint
	os.Exit(m.Run())
}
```


## 参考

- https://golang.testcontainers.org/