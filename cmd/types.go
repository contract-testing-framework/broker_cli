package cmd

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

type HttpError struct {
	Error string `json:"error"`
}

type EnvBody struct {
	EnvironmentName string `json:"environmentName"`
}