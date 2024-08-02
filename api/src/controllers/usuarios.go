package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// CriarUsuario cria um usuário no banco de dados.
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	if erro = usuario.Preparar("cadastro"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario.ID, erro = repositorio.Criar(usuario)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusCreated, usuario)
}

// BuscarUsuarios busca todos os usuários no banco de dados.
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
	}
	defer db.Close()
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.Buscar(nomeOuNick)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
	}
	respostas.JSON(w, http.StatusOK, usuarios)
}

// BuscarUmUsuario busca um único usuário no banco de dados.
func BuscarUmUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario, erro := repositorio.BuscarPorID(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
	}
	respostas.JSON(w, http.StatusOK, usuario)
}

// AtualizarUsuario altera as informações de um usuário no banco de dados.
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível atualizar um usuário que não seja o seu"))
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	if erro = usuario.Preparar("edicao"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Atualizar(usuarioID, usuario); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DeletarUsuario deleta um usuário do banco de dados.
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível deletar um usuário que não seja o seu"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Deletar(usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// SeguirUsuarios permite que um usuários siga outros.
func SeguirUsuario(w http.ResponseWriter, r *http.Request){
	seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro!= nil {
        respostas.Erro(w, http.StatusUnauthorized, erro)
        return
    }

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
        return
	}

	if seguidorID == usuarioID{
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível seguir a si mesmo"))
        return
	}
	db, erro := banco.Conectar()
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Seguir(usuarioID, seguidorID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
        return
	}

	respostas.JSON(w, http.StatusNoContent, nil)

}

// PararDeSeguirUsuario permite que um usuário deixe de seguir outro.
func PararDeSeguirUsuario(w http.ResponseWriter, r *http.Request){
	seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro!= nil {
        respostas.Erro(w, http.StatusUnauthorized, erro)
        return
    }

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro!= nil {
        respostas.Erro(w, http.StatusBadRequest, erro)
        return
    }

	if seguidorID == usuarioID{
        respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível parar de seguir a si mesmo"))
        return
    }

	db, erro := banco.Conectar()
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.PararDeSeguir(usuarioID, seguidorID); erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }

	respostas.JSON(w, http.StatusNoContent, nil)
}

// BuscarSeguidores busca os seguidores de um usuário.
func BuscarSeguidores(w http.ResponseWriter, r *http.Request){
	parametros := mux.Vars(r)
    usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
    if erro!= nil {
        respostas.Erro(w, http.StatusBadRequest, erro)
        return
    }

    db, erro := banco.Conectar()
    if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
    defer db.Close()

    repositorio := repositorios.NovoRepositorioDeUsuarios(db)
    seguidores, erro := repositorio.BuscarSeguidores(usuarioID)
    if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }

    respostas.JSON(w, http.StatusOK, seguidores)
}

// BuscarSeguindo busca os usuários aos quais um usuário está seguindo.
func BuscarSeguindo(w http.ResponseWriter, r *http.Request){
	parametros := mux.Vars(r)
    usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
    if erro!= nil {
        respostas.Erro(w, http.StatusBadRequest, erro)
        return
    }

    db, erro := banco.Conectar()
    if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
    defer db.Close()

    repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.BuscarSeguindo(usuarioID)
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }

	respostas.JSON(w, http.StatusOK, usuarios)
}

// AtualizarSenha permite que um usuário altere sua senha.
func AtualizarSenha (w http.ResponseWriter, r *http.Request){
	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro!= nil {
        respostas.Erro(w, http.StatusUnauthorized, erro)
        return
    }

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro!= nil {
        respostas.Erro(w, http.StatusBadRequest, erro)
        return
    }

	if usuarioID!= usuarioIDNoToken{
        respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível atualizar a senha de um usuário que não seja o seu"))
        return
    }

	corpoRequisicao, erro := io.ReadAll(r.Body)

	var senha modelos.Senha
	if erro = json.Unmarshal(corpoRequisicao, &senha); erro!= nil {
        respostas.Erro(w, http.StatusBadRequest, erro)
        return
    }

	db, erro := banco.Conectar()
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	senhaSalvaNoBanco, erro := repositorio.BuscarSenha(usuarioID)
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }

	if erro = seguranca.VerificarSenha(senhaSalvaNoBanco, senha.Atual); erro!= nil {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("Senha atual incorreta"))
        return
    }

	senhaComHash, erro := seguranca.Hash(senha.Nova)
	if erro!= nil {
        respostas.Erro(w, http.StatusBadRequest, erro)
        return
    }

	if erro = repositorio.AtualizarSenha(usuarioID, string(senhaComHash)); erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }

	respostas.JSON(w, http.StatusNoContent, nil)
}