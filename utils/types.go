package utils

type Consumer struct {
	Name string `json:"name"`
}

type Pact struct {
	Consumer     Consumer    `json:"consumer"`
	Interactions interface{} `json:"interactions"`
	MetaData     interface{} `json:"metadata"`
	Provider     interface{} `json:"provider"`
}

type ConsumerBody struct {
	Contract        Pact   `json:"contract"`
	ConsumerName    string `json:"consumerName"`
	ConsumerVersion string `json:"consumerVersion"`
	ConsumerBranch  string `json:"consumerBranch"`
}

type ProviderBody struct {
	Spec            interface{} `json:"spec"`
	ProviderName    string      `json:"providerName"`
	ProviderVersion string      `json:"providerVersion"`
	ProviderBranch  string      `json:"providerBranch"`
	SpecFormat      string      `json:"specFormat"`
}

type EnvBody struct {
	EnvironmentName string `json:"environmentName"`
}

type DeploymentBody struct {
	EnvironmentName 	 string `json:"environmentName"`
	ParticipantName 	 string `json:"participantName"`
	ParticipantVersion string `json:"participantVersion"`
	Deployed 					 bool 	`json:"deployed"`
}

type MbProxy struct {
	To string `json:"to"`
	Mode string `json:"mode"`
}

type MbResponse struct {
	Proxy MbProxy `json:"proxy"`
}

type MbStub struct {
	Responses []MbResponse `json:"responses"`
}

type ProxyConfig struct {
	Port     int          `json:"port"`
	Name     string       `json:"name"`
	Protocol string       `json:"protocol"`
	Stubs    []MbStub `json:"stubs"`
}