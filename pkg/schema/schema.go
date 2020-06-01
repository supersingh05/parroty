package schema

type Parroty struct {
	Clusters     []Cluster       `yaml:"clusters,omitempty"`
	GlobalExpect []ClusterExpect `yaml:"globalExpect,omitempty"`
}

type Cluster struct {
	Name           string          `yaml:"name"`
	Context        string          `yaml:"context"`
	Cloud          string          `yaml:"cloud,omitempty"`
	ClusterExpect  []ClusterExpect `yaml:"clusterExpect,omitempty"`
	KubeconfigPath string          `yaml:"kubeconfigPath,omitempty"`
	AwsSecretKey   string          `yaml:"awsSecretKey,omitempty"`
	AwsAccessKey   string          `yaml:"awsAccessKey,omitempty"`
}

type ClusterExpect struct {
	ObjectName string `yaml:"objectName,omitempty"`
	Group      string `yaml:"group,omitempty"`
	Kind       string `yaml:"kind,omitempty"`
	Version    string `yaml:"version,omitempty"`
	Namespace  string `yaml:"namespace,omitempty"`
}
