package seguranca

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash recebe uma string e coloca uma hash nela.
func Hash (senha string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
}

// Verificar compara duas senhas.
func VerificarSenha(senhaString, senhaComHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(senhaComHash), []byte(senhaString))
}