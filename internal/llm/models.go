package llm

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role    Role
	Content string
}

type OutputFormat string

const (
	OutputFormatText OutputFormat = "text"
	OutputFormatJSON OutputFormat = "json"
	OutputFormatEnum OutputFormat = "enum"
)

// ResponseSchema specifies exactly what the caller expects back from the LLM.
type ResponseSchema struct {
	Format OutputFormat

	// Used when Format == OutputFormatJSON to enforce a schema (can be empty string for generic JSON)
	JSONSchema string

	// Used when Format == OutputFormatEnum to restrict to specific choices (e.g., ["Yes", "No"])
	AllowedChoices []string
}

type GenerateRequest struct {
	Messages       []Message
	Model          string   // Optional: to override default model
	Temperature    *float32 // Optional: parameter to control randomness
	ExpectedOutput ResponseSchema
}

type GenerateResponse struct {
	Text             string
	PromptTokens     int
	CompletionTokens int
}
