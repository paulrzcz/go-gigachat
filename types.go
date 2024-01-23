package gigachat

const (
	AuthUrl    = "https://ngw.devices.sberbank.ru:9443/api/"
	BaseUrl    = "https://gigachat.devices.sberbank.ru/api/"
	OAuthPath  = "v2/oauth"
	ModelsPath = "v1/models"
	ChatPath   = "v1/chat/completions"
)

const (
	ScopeApiIndividual = "GIGACHAT_API_PERS"
	ScopeApiBusiness   = "GIGACHAT_API_CORP"
)

const (
	ModelLatest = "GigaChat:latest"
)

const (
	UserRole      = "user"
	AssistantRole = "assistant"
	SystemRole    = "system"
)

type Config struct {
	AuthUrl      string
	BaseUrl      string
	ClientId     string
	ClientSecret string
	Scope        string
	Insecure     bool
}

type OAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

type ChatRequest struct {
	Model             string    `json:"model"`
	Messages          []Message `json:"messages"`
	Temperature       *float64  `json:"temperature"`
	TopP              *float64  `json:"top_p"`
	N                 *int64    `json:"n"`
	Stream            *bool     `json:"stream"`
	MaxTokens         *int64    `json:"max_tokens"`
	RepetitionPenalty *float64  `json:"repetition_penalty"`
	UpdateInterval    *int64    `json:"update_interval"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Model   string   `json:"model"`
	Created int64    `json:"created"`
	Method  string   `json:"object"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int64  `json:"index"`
	FinishReason string `json:"finish_reason"`
	Message      Message
}

type Usage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}
