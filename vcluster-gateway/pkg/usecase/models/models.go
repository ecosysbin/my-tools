package models

// Values 代表 helm 命令执行时的 values.yaml，使用结构体嵌套表示 yaml 语法中的层级关系
// 需要额外注意的是 Values Struct 被 processor.HelmValuesProcessor wrap，processor.HelmValuesProcessor 封装了从结构体序列化到 YAML 文件的方法
// 这些 YAML 格式的数据会被传递到 forkvcluster.CreateOptions 中的 Values 字段中
type Values struct {
	DefaultImageRegistry string            `yaml:"defaultImageRegistry,omitempty"`
	GlobalAnnotations    map[string]string `yaml:"globalAnnotations,omitempty"`
	Isolation            Isolation         `yaml:"isolation,omitempty"`
	Zetyun               Zetyun            `yaml:"zetyun,omitempty"`
	MapServices          MapServices       `yaml:"mapServices,omitempty"`
	Plugin               Plugin            `yaml:"plugin,omitempty"`
	Syncer               Syncer            `yaml:"syncer,omitempty"`
	Etcd                 Etcd              `yaml:"etcd,omitempty"`
	Labels               map[string]string `yaml:"labels,omitempty"`
	PodLabels            map[string]string `yaml:"podLabels,omitempty"`
	Sync                 Sync              `yaml:"sync,omitempty"`

	// Pro                    bool                   `yaml:"pro"`
	// Headless               bool                   `yaml:"headless"`
	// Monitoring             Monitoring             `yaml:"monitoring"`
	// EnableHA               bool                   `yaml:"enableHA"`
	// FallbackHostDns        bool                   `yaml:"fallbackHostDns"`
	// Proxy                  Proxy                  `yaml:"proxy"`
	// Controller             Controller             `yaml:"controller"`
	// Scheduler              Scheduler              `yaml:"scheduler"`
	// Api                    Api                    `yaml:"api"`
	// ServiceAccount         ServiceAccount         `yaml:"serviceAccount"`
	// WorkloadServiceAccount WorkloadServiceAccount `yaml:"workloadServiceAccount"`
	// Rbac                   Rbac                   `yaml:"rbac"`
	// Service                Service                `yaml:"service"`
	// Ingress                Ingress                `yaml:"ingress"`
	// Openshift              Openshift              `yaml:"openshift"`
	// Coredns                Coredns                `yaml:"coredns"`
	// Init                   Init                   `yaml:"init"`
	// MultiNamespaceMode     MultiNamespaceMode     `yaml:"multiNamespaceMode"`
	// Admission              Admission              `yaml:"admission"`
	// Telemetry              Telemetry              `yaml:"telemetry"`
}

func NewHelmValues() *Values {
	return &Values{
		GlobalAnnotations: make(map[string]string),
		Isolation: Isolation{
			ResourceQuota: ResourceQuota{
				Quota: make(map[string]interface{}),
			},
		},
		Syncer: Syncer{
			NodeSelector:       make(map[string]interface{}),
			Affinity:           make(map[string]interface{}),
			Labels:             make(map[string]interface{}),
			Annotations:        make(map[string]interface{}),
			PodAnnotations:     make(map[string]interface{}),
			PodLabels:          make(map[string]interface{}),
			SecurityContext:    make(map[string]interface{}),
			PodSecurityContext: make(map[string]interface{}),
			ServiceAnnotations: make(map[string]interface{}),
		},
		Etcd: Etcd{
			NodeSelector:       make(map[string]interface{}),
			Affinity:           make(map[string]interface{}),
			Labels:             make(map[string]interface{}),
			Annotations:        make(map[string]interface{}),
			PodAnnotations:     make(map[string]interface{}),
			PodLabels:          make(map[string]interface{}),
			SecurityContext:    make(map[string]interface{}),
			ServiceAnnotations: make(map[string]interface{}),
		},
		Labels:    make(map[string]string),
		PodLabels: make(map[string]string),
	}
}

type Syncer struct {
	Image                 string                 `yaml:"image,omitempty"`
	ImagePullPolicy       string                 `yaml:"imagePullPolicy,omitempty"`
	ExtraArgs             []interface{}          `yaml:"extraArgs,omitempty"`
	VolumeMounts          []interface{}          `yaml:"volumeMounts,omitempty"`
	ExtraVolumeMounts     []interface{}          `yaml:"extraVolumeMounts,omitempty"`
	Env                   []EnvConfig            `yaml:"env,omitempty"`
	LivenessProbe         ProbeConfig            `yaml:"livenessProbe,omitempty"`
	ReadinessProbe        ProbeConfig            `yaml:"readinessProbe,omitempty"`
	Resources             Resources              `yaml:"resources,omitempty"`
	Volumes               []interface{}          `yaml:"volumes,omitempty"`
	Replicas              int                    `yaml:"replicas,omitempty"`
	NodeSelector          map[string]interface{} `yaml:"nodeSelector,omitempty"`
	Affinity              map[string]interface{} `yaml:"affinity,omitempty"`
	Tolerations           []interface{}          `yaml:"tolerations,omitempty"`
	Labels                map[string]interface{} `yaml:"labels,omitempty"`
	Annotations           map[string]interface{} `yaml:"annotations,omitempty"`
	PodAnnotations        map[string]interface{} `yaml:"podAnnotations,omitempty"`
	PodLabels             map[string]interface{} `yaml:"podLabels,omitempty"`
	PriorityClassName     string                 `yaml:"priorityClassName,omitempty"`
	KubeConfigContextName string                 `yaml:"kubeConfigContextName,omitempty"`
	SecurityContext       map[string]interface{} `yaml:"securityContext,omitempty"`
	PodSecurityContext    map[string]interface{} `yaml:"podSecurityContext,omitempty"`
	ServiceAnnotations    map[string]interface{} `yaml:"serviceAnnotations,omitempty"`
}

type EnvConfig struct {
	Name  string `yaml:"name,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type ProbeConfig struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

type Resources struct {
	Requests ResourceRequest `yaml:"requests,omitempty"`
	Limits   ResourceLimit   `yaml:"limits,omitempty"`
}

type ResourceRequest struct {
	EphemeralStorage string `yaml:"ephemeral-storage,omitempty"`
	Cpu              string `yaml:"cpu,omitempty"`
	Memory           string `yaml:"memory,omitempty"`
}

type ResourceLimit struct {
	EphemeralStorage string `yaml:"ephemeral-storage,omitempty"`
	Cpu              string `yaml:"cpu,omitempty"`
	Memory           string `yaml:"memory,omitempty"`
}

type Etcd struct {
	Image                            string                 `yaml:"image,omitempty"`
	ImagePullPolicy                  string                 `yaml:"imagePullPolicy,omitempty"`
	Replicas                         int                    `yaml:"replicas,omitempty"`
	NodeSelector                     map[string]interface{} `yaml:"nodeSelector,omitempty"`
	Affinity                         map[string]interface{} `yaml:"affinity,omitempty"`
	Tolerations                      []interface{}          `yaml:"tolerations,omitempty"`
	Labels                           map[string]interface{} `yaml:"labels,omitempty"`
	Annotations                      map[string]interface{} `yaml:"annotations,omitempty"`
	PodAnnotations                   map[string]interface{} `yaml:"podAnnotations,omitempty"`
	PodLabels                        map[string]interface{} `yaml:"podLabels,omitempty"`
	Resources                        EtcdResources          `yaml:"resources,omitempty"`
	Storage                          EtcdStorage            `yaml:"storage,omitempty"`
	PriorityClassName                string                 `yaml:"priorityClassName,omitempty"`
	SecurityContext                  map[string]interface{} `yaml:"securityContext,omitempty"`
	ServiceAnnotations               map[string]interface{} `yaml:"serviceAnnotations,omitempty"`
	AutoDeletePersistentVolumeClaims bool                   `yaml:"autoDeletePersistentVolumeClaims,omitempty"`
}

type EtcdResources struct {
	Requests EtcdResourceRequest `yaml:"requests,omitempty"`
	Limits   EtcdResourceLimit   `yaml:"limits,omitempty"`
}

type EtcdResourceRequest struct {
	Cpu    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

type EtcdResourceLimit struct {
	Cpu    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

type EtcdStorage struct {
	Persistence bool   `yaml:"persistence,omitempty"`
	Size        string `yaml:"size,omitempty"`
	ClassName   string `yaml:"className,omitempty"`
}

type Isolation struct {
	Enabled             bool                `yaml:"enabled,omitempty"`
	Namespace           string              `yaml:"namespace,omitempty"`
	PodSecurityStandard string              `yaml:"podSecurityStandard,omitempty"`
	NodeProxyPermission NodeProxyPermission `yaml:"nodeProxyPermission,omitempty"`
	ResourceQuota       ResourceQuota       `yaml:"resourceQuota,omitempty"`
	LimitRange          LimitRange          `yaml:"limitRange,omitempty"`
	NetworkPolicy       NetworkPolicy       `yaml:"networkPolicy,omitempty"`
}

type NodeProxyPermission struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

type ResourceQuota struct {
	Enabled       bool                   `yaml:"enabled,omitempty"`
	Quota         map[string]interface{} `yaml:"quota,omitempty"`
	ScopeSelector ScopeSelector          `yaml:"scopeSelector,omitempty"`
}

type ScopeSelector struct {
	MatchExpressions []interface{} `yaml:"matchExpressions,omitempty"`
	Scopes           []interface{} `yaml:"scopes,omitempty"`
}

type LimitRange struct {
	Enabled        bool           `yaml:"enabled,omitempty"`
	Default        Default        `yaml:"default,omitempty"`
	DefaultRequest DefaultRequest `yaml:"defaultRequest,omitempty"`
}

type Default struct {
	EphemeralStorage string `yaml:"ephemeral-storage,omitempty"`
	Memory           string `yaml:"memory,omitempty"`
	Cpu              string `yaml:"cpu,omitempty"`
}

type DefaultRequest struct {
	EphemeralStorage string `yaml:"ephemeral-storage,omitempty"`
	Memory           string `yaml:"memory,omitempty"`
	Cpu              string `yaml:"cpu,omitempty"`
}

type NetworkPolicy struct {
	Enabled     bool   `yaml:"enabled,omitempty"`
	FallbackDns string `yaml:"fallbackDns,omitempty"`
}

type OutgoingConnections struct {
	IpBlock IpBlock `yaml:"ipBlock,omitempty"`
}

type IpBlock struct {
	Cidr   string   `yaml:"cidr,omitempty"`
	Except []string `yaml:"except,omitempty"`
}

type Zetyun struct {
	Type         string       `yaml:"type,omitempty"`
	StorageClass StorageClass `yaml:"storageclass,omitempty"`
}

type StorageClass struct {
	ClusterId string       `yaml:"clusterId,omitempty"`
	Enabled   bool         `yaml:"enabled,omitempty"`
	List      []NameTypeNs `yaml:"list,omitempty"`
}

type NameTypeNs struct {
	Name string `yaml:"name,omitempty"`
	Type string `yaml:"type,omitempty"`
	Ns   string `yaml:"ns,omitempty"`
}

type MapServices struct {
	FromHost    []FromTo `yaml:"fromHost,omitempty"`
	FromVirtual []FromTo `yaml:"fromVirtual,omitempty"`
}

type FromTo struct {
	From string `yaml:"from,omitempty"`
	To   string `yaml:"to,omitempty"`
}

type Plugin struct {
	Hooks                 Hooks                 `yaml:"hooks,omitempty"`
	MultiStorageRequester MultiStorageRequester `yaml:"multi-storage-requester,omitempty"`
}

type Hooks struct {
	Image           string      `yaml:"image,omitempty"`
	ImagePullPolicy string      `yaml:"imagePullPolicy,omitempty"`
	Env             []NameValue `yaml:"env,omitempty"`
}

type NameValue struct {
	Name  string `yaml:"name,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type MultiStorageRequester struct {
	Optional        bool            `yaml:"optional,omitempty"`
	Image           string          `yaml:"image,omitempty"`
	ImagePullPolicy string          `yaml:"imagePullPolicy,omitempty"`
	Command         []string        `yaml:"command,omitempty"`
	VolumeMounts    []NameMountPath `yaml:"volumeMounts,omitempty"`
	Env             []NameValue     `yaml:"env,omitempty"`
	Args            []string        `yaml:"args,omitempty"`
}

type NameMountPath struct {
	Name      string `yaml:"name,omitempty"`
	MountPath string `yaml:"mountPath,omitempty"`
}

//type Controller struct {
//	Image             string                 `yaml:"image"`
//	ImagePullPolicy   string                 `yaml:"imagePullPolicy"`
//	Replicas          int                    `yaml:"replicas"`
//	NodeSelector      map[string]string      `yaml:"nodeSelector"`
//	Affinity          map[string]interface{} `yaml:"affinity"`
//	Tolerations       []interface{}          `yaml:"tolerations"`
//	Labels            map[string]string      `yaml:"labels"`
//	Annotations       map[string]string      `yaml:"annotations"`
//	PodAnnotations    map[string]string      `yaml:"podAnnotations"`
//	PodLabels         map[string]string      `yaml:"podLabels"`
//	Resources         ResourceRequests       `yaml:"resources"`
//	PriorityClassName string                 `yaml:"priorityClassName"`
//	SecurityContext   map[string]interface{} `yaml:"securityContext"`
//}
//
//type ResourceRequests struct {
//	Requests map[string]string `yaml:"requests"`
//}
//
//type Scheduler struct {
//	Image             string                 `yaml:"image"`
//	ImagePullPolicy   string                 `yaml:"imagePullPolicy"`
//	Replicas          int                    `yaml:"replicas"`
//	NodeSelector      map[string]string      `yaml:"nodeSelector"`
//	Affinity          map[string]interface{} `yaml:"affinity"`
//	Tolerations       []interface{}          `yaml:"tolerations"`
//	Labels            map[string]string      `yaml:"labels"`
//	Annotations       map[string]string      `yaml:"annotations"`
//	PodAnnotations    map[string]string      `yaml:"podAnnotations"`
//	PodLabels         map[string]string      `yaml:"podLabels"`
//	Resources         ResourceRequests       `yaml:"resources"`
//	PriorityClassName string                 `yaml:"priorityClassName"`
//}

//type Api struct {
//	Image              string                 `yaml:"image"`
//	ImagePullPolicy    string                 `yaml:"imagePullPolicy"`
//	ExtraArgs          []string               `yaml:"extraArgs"`
//	Replicas           int                    `yaml:"replicas"`
//	NodeSelector       map[string]string      `yaml:"nodeSelector"`
//	Affinity           map[string]interface{} `yaml:"affinity"`
//	Tolerations        []interface{}          `yaml:"tolerations"`
//	Labels             map[string]string      `yaml:"labels"`
//	Annotations        map[string]string      `yaml:"annotations"`
//	PodAnnotations     map[string]string      `yaml:"podAnnotations"`
//	PodLabels          map[string]string      `yaml:"podLabels"`
//	Resources          Resources              `yaml:"resources"`
//	PriorityClassName  string                 `yaml:"priorityClassName"`
//	SecurityContext    map[string]interface{} `yaml:"securityContext"`
//	ServiceAnnotations map[string]string      `yaml:"serviceAnnotations"`
//}
//
//type ServiceAccount struct {
//	Create           bool     `yaml:"create"`
//	Name             string   `yaml:"name"`
//	ImagePullSecrets []string `yaml:"imagePullSecrets"`
//}
//
//type WorkloadServiceAccount struct {
//	Annotations map[string]string `yaml:"annotations"`
//}
//
//type Rbac struct {
//	ClusterRole ClusterRole `yaml:"clusterRole"`
//	Role        Role        `yaml:"role"`
//}
//
//type ClusterRole struct {
//	Create bool `yaml:"create"`
//}
//
//type Role struct {
//	Create               bool     `yaml:"create"`
//	Extended             bool     `yaml:"extended"`
//	ExcludedApiResources []string `yaml:"excludedApiResources"`
//}
//
//type Service struct {
//	Type                     string   `yaml:"type"`
//	ExternalIPs              []string `yaml:"externalIPs"`
//	ExternalTrafficPolicy    string   `yaml:"externalTrafficPolicy"`
//	LoadBalancerIP           string   `yaml:"loadBalancerIP"`
//	LoadBalancerSourceRanges []string `yaml:"loadBalancerSourceRanges"`
//	LoadBalancerClass        string   `yaml:"loadBalancerClass"`
//}
//
//type Ingress struct {
//	Enabled          bool              `yaml:"enabled"`
//	PathType         string            `yaml:"pathType"`
//	IngressClassName string            `yaml:"ingressClassName"`
//	Host             string            `yaml:"host"`
//	Annotations      map[string]string `yaml:"annotations"`
//	Tls              []interface{}     `yaml:"tls"`
//}
//
//type Openshift struct {
//	Enable bool `yaml:"enable"`
//}
//
//type Coredns struct {
//	Integrated     bool              `yaml:"integrated"`
//	Enabled        bool              `yaml:"enabled"`
//	Plugin         PluginConfig      `yaml:"plugin"`
//	Replicas       int               `yaml:"replicas"`
//	NodeSelector   map[string]string `yaml:"nodeSelector"`
//	Image          string            `yaml:"image"`
//	Config         string            `yaml:"config"`
//	Service        Service           `yaml:"service"`
//	Resources      Resources         `yaml:"resources"`
//	PodAnnotations map[string]string `yaml:"podAnnotations"`
//	PodLabels      map[string]string `yaml:"podLabels"`
//}
//
//type PluginConfig struct {
//	Enabled bool          `yaml:"enabled"`
//	Config  []interface{} `yaml:"config"`
//}
//type Monitoring struct {
//	ServiceMonitor ServiceMonitor `yaml:"serviceMonitor"`
//}
//
//type ServiceMonitor struct {
//	Enabled bool `yaml:"enabled"`
//}

type Sync struct {
	Services               SyncConfig             `yaml:"services,omitempty"`
	Configmaps             SyncConfig             `yaml:"configmaps,omitempty"`
	Secrets                SyncConfig             `yaml:"secrets,omitempty"`
	Endpoints              SyncConfig             `yaml:"endpoints,omitempty"`
	Pods                   PodsSync               `yaml:"pods,omitempty"`
	Events                 SyncConfig             `yaml:"events,omitempty"`
	PersistentVolumeClaims SyncConfig             `yaml:"persistentvolumeclaims,omitempty"`
	Ingresses              SyncConfig             `yaml:"ingresses,omitempty"`
	IngressClasses         map[string]interface{} `yaml:"ingressclasses,omitempty"`
	FakeNodes              SyncConfig             `yaml:"fake-nodes,omitempty"`
	FakePersistentVolumes  SyncConfig             `yaml:"fake-persistentvolumes,omitempty"`
	Nodes                  NodesSync              `yaml:"nodes,omitempty"`
	PersistentVolumes      SyncConfig             `yaml:"persistentvolumes,omitempty"`
	StorageClasses         SyncConfig             `yaml:"storageclasses,omitempty"`
	HostStorageClasses     SyncConfig             `yaml:"hoststorageclasses,omitempty"`
	PriorityClasses        SyncConfig             `yaml:"priorityclasses,omitempty"`
	NetworkPolicies        SyncConfig             `yaml:"networkpolicies,omitempty"`
	VolumeSnapshots        SyncConfig             `yaml:"volumesnapshots,omitempty"`
	PodDisruptionBudgets   SyncConfig             `yaml:"poddisruptionbudgets,omitempty"`
	ServiceAccounts        SyncConfig             `yaml:"serviceaccounts,omitempty"`
	Generic                GenericConfig          `yaml:"generic,omitempty"`
}

type SyncConfig struct {
	Enabled bool `yaml:"enabled"`
}
type PodsSync struct {
	Enabled             bool `yaml:"enabled"`
	EphemeralContainers bool `yaml:"ephemeralContainers"`
	Status              bool `yaml:"status"`
}

type NodesSync struct {
	FakeKubeletIPs  bool   `yaml:"fakeKubeletIPs"`
	Enabled         bool   `yaml:"enabled"`
	SyncAllNodes    bool   `yaml:"syncAllNodes"`
	NodeSelector    string `yaml:"nodeSelector"`
	EnableScheduler bool   `yaml:"enableScheduler"`
	SyncNodeChanges bool   `yaml:"syncNodeChanges"`
}

type GenericConfig struct {
	Config string `yaml:"config"`
}

//type Proxy struct {
//	MetricsServer MetricsServer `yaml:"metricsServer"`
//}
//
//type MetricsServer struct {
//	Nodes MetricsServerConfig `yaml:"nodes"`
//	Pods  MetricsServerConfig `yaml:"pods"`
//}
//
//type MetricsServerConfig struct {
//	Enabled bool `yaml:"enabled"`
//}

//type Init struct {
//	Manifests         string        `yaml:"manifests"`
//	ManifestsTemplate string        `yaml:"manifestsTemplate"`
//	Helm              []interface{} `yaml:"helm"`
//}

//type MultiNamespaceMode struct {
//	Enabled bool `yaml:"enabled"`
//}

//type Admission struct {
//	ValidatingWebhooks []interface{} `yaml:"validatingWebhooks"`
//	MutatingWebhooks   []interface{} `yaml:"mutatingWebhooks"`
//}

//type Telemetry struct {
//	Disabled           bool   `yaml:"disabled"`
//	InstanceCreator    string `yaml:"instanceCreator"`
//	PlatformUserID     string `yaml:"platformUserID"`
//	PlatformInstanceID string `yaml:"platformInstanceID"`
//	MachineID          string `yaml:"machineID"`
//}
