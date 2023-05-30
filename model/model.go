package model

type Swagger struct {
	Paths       map[string]map[string]Path `json:"paths"`
	Definitions map[string]DefinitionWrap  `json:"definitions"`
}

// ================ path ======================

type Path struct {
	Tags       []string    `json:"tags"`
	Summary    string      `json:"summary"`
	Responses  Responses   `json:"responses"`
	Parameters []Parameter `json:"parameters"`
	Produces   []string    `json:"produces"`
}

type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	TypeGo      string `json:"type"`
	DefaultGo   string `json:"default"`
	EnumGo      string `json:"enum"`
	Schema      struct {
		Ref string `json:"$ref"`
	} `json:"schema"`
}

type Responses struct {
	HttpOk struct {
		Description string `json:"description"`
		Schema      struct {
			Ref string `json:"$ref"`
		} `json:"schema"`
	} `json:"200"`
}

// ================ 模型 ======================

// 用于保存操作生成和使用的数据类型的对象
type DefinitionWrap struct {
	TypeGo     string                `json:"type"`
	Properties map[string]Properties `json:"properties"`
}

type Properties struct {
	TypeGo      string `json:"type"`
	Format      string `json:"format"`
	Description string `json:"description"`
	Ref         string `json:"$ref"`
	Items       struct {
		TypeGo string `json:"type"`
		Ref    string `json:"$ref"`
	} `json:"items"`
	AdditionalProperties struct {
		Ref    string `json:"$ref"`
		TypeGo string `json:"type"`
	} `json:"additionalProperties"`
}
