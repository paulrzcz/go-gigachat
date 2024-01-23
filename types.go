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
