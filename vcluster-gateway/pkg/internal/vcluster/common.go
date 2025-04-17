package forkvcluster

// GlobalFlags is the flags that contains the global flags
type GlobalFlags struct {
	Silent    bool
	Debug     bool
	Config    string
	Context   string
	Namespace string
	LogOutput string
}

func NewDefaultGlobalFlags() *GlobalFlags {
	return &GlobalFlags{
		Silent:    false,
		Debug:     true,
		Config:    "",
		Context:   "",
		Namespace: "",
		LogOutput: "",
	}
}

func (gf *GlobalFlags) SetK8sContext(context string) {
	gf.Context = context
}
