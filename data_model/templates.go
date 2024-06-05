package data_model

type Templates struct {
	Templates []Template `json:"templates"`
}

type Template struct {
	EnvTemplateName      string            `json:"envTemplateName"`
	Org                  string            `json:"org"`
	Pod                  string            `json:"pod"`
	EnvOwners            string            `json:"envOwners"`
	TTLInHours           string            `json:"ttlInHours"`
	AdditionalComponents []string          `json:"additionalComponents"`
	Services             []Service         `json:"services"`
	CreatedBy            string            `json:"createdBy"`
	CreatedAt            string            `json:"createdAt"`
	UpdatedBy            string            `json:"updatedBy"`
	UpdatedAt            string            `json:"updatedAt"`
	EnvType              string            `json:"envType"`
	EnvironmentTask      []EnvironmentTask `json:"environmentTask"`
	TTL                  int64             `json:"ttl"`
}

type EnvironmentTask struct {
	TaskTriggerEnvState string   `json:"taskTriggerEnvState"`
	TaskType            string   `json:"taskType"`
	TaskPlaneType       string   `json:"taskPlaneType"`
	TaskData            TaskData `json:"taskData"`
}

type TaskData struct {
	Endpoint string  `json:"endpoint"`
	Header   Header  `json:"header"`
	Method   string  `json:"method"`
	Payload  Payload `json:"payload"`
}

type Header struct {
	UserID string `json:"userId"`
}

type Payload struct {
	ArtifactID      string       `json:"artifact_id"`
	EnvInstanceName string       `json:"envInstanceName"`
	GroupID         string       `json:"group_id"`
	Repository      string       `json:"repository"`
	TestDetails     []TestDetail `json:"test_details"`
	TestType        string       `json:"test_type"`
	Version         string       `json:"version"`
}

type TestDetail struct {
	ResourceName string      `json:"resource_name"`
	ResourceType string      `json:"resource_type"`
	Tests        interface{} `json:"tests"`
}

type Service struct {
	Name              string   `json:"name"`
	Version           string   `json:"version"`
	DependentServices []string `json:"dependentServices"`
	IsMockService     bool     `json:"isMockService"`
	CommitSHA         string   `json:"commitSha"`
	File              string   `json:"file"`
	BranchName        string   `json:"branchName"`
	Repo              string   `json:"repo"`
}
