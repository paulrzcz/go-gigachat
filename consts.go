package gigachat

const (
	AuthUrl    = "https://ngw.devices.sberbank.ru:9443/api/"
	BaseUrl    = "https://gigachat.devices.sberbank.ru/api/"
	OAuthPath  = "v2/oauth"
	ModelsPath = "v1/models"
	ChatPath   = "v1/chat/completions"
	Files      = "v1/files"
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

const (
	GIGACHAT_2_LITE = "GigaChat-2"
	GIGACHAT_2_PRO  = "GigaChat-2-Pro"
	GIGACHAT_2_MAX  = "GigaChat-2-Max"
	GIGACHAT_1_LITE = "GigaChat"
	GIGACHAT_1_PRO  = "GigaChat-Pro"
	GIGACHAT_1_MAX  = "GigaChat-Max"
)
