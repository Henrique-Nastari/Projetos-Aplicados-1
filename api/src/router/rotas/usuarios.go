package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasUsuarios = []Rota{
	{
		URI:                "/usuarios",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarUsuario,
		RequerAutenticacao: false,
	},
	{
		URI:                "/usuarios",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarUsuarios,
		RequerAutenticacao: false,
	},
	{
		URI:                "/usuarios/{usuarioId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarUmUsuario,
		RequerAutenticacao: false,
	},
	{
		URI:                "/usuarios/{usuarioId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarUsuario,
		RequerAutenticacao: false,
	},
	{
		URI:                "/usuarios/{usuarioId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletarUsuario,
		RequerAutenticacao: false,
	},
	{
		URI:"/usuarios/{usuarioId}/seguir",
		Metodo: http.MethodPost,
		Funcao: controllers.SeguirUsuario,
		RequerAutenticacao: false,
	},
	{
		URI:"/usuarios/{usuarioId}/parar-de-seguir",
		Metodo: http.MethodPost,
		Funcao: controllers.PararDeSeguirUsuario,
		RequerAutenticacao: false,
	},
	{
		URI: "/usuarios/{usuarioId}/seguidores",
		Metodo: http.MethodGet,
		Funcao: controllers.BuscarSeguidores,
		RequerAutenticacao: false,
	},
	{
		URI: "/usuarios/{usuarioId}/seguindo",
		Metodo: http.MethodGet,
		Funcao: controllers.BuscarSeguindo,
		RequerAutenticacao: false,
	},
	{
		URI: "/usuarios/{usuarioId}/atualizar-senha",
		Metodo: http.MethodPost,
		Funcao: controllers.AtualizarSenha,
		RequerAutenticacao: false,
	},
}
