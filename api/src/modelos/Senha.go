package modelos

// Senha representa o novo formato de requisição de alteração de senhas
type Senha struct {
	Nova  string `json:"nova"`
	Atual string `json:"atual"`
}
