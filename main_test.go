package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guilhermeonrails/api-go-gin/controllers"
	"github.com/guilhermeonrails/api-go-gin/database"
	"github.com/guilhermeonrails/api-go-gin/models"
	"github.com/stretchr/testify/assert"
)

var ID int

func SetupRotasTeste() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	rotas := gin.Default()
	return rotas
}

func CriaAlunoMock() {
	aluno := models.Aluno{Nome: "Aluno Teste", CPF: "12345678901", RG: "123456789"}

	database.DB.Create(&aluno)
	ID = int(aluno.ID)
}

func DeletaAlunoMock() {
	var aluno models.Aluno
	database.DB.Delete(&aluno, ID)
}

func TestVerificaStatusCodeSaudacao(t *testing.T) {
	r := SetupRotasTeste()
	r.GET("/:nome", controllers.Saudacao)
	req, _ := http.NewRequest("GET", "/gui", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)

	assert.Equal(t, http.StatusOK, resposta.Code, "Status diferente de 200")

	mockResposta := `{"API diz:":"E ai gui, tudo beleza?"}`

	respostaBody, _ := ioutil.ReadAll(resposta.Body)

	assert.Equal(t, mockResposta, string(respostaBody))
}

func TestListandoAlunos(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupRotasTeste()
	path := "/alunos"
	r.GET(path, controllers.ExibeTodosAlunos)
	req, _ := http.NewRequest("GET", path, nil)
	resposta := httptest.NewRecorder()

	r.ServeHTTP(resposta, req)

	assert.Equal(t, http.StatusOK, resposta.Code)
}

func TestBuscaAlunoCPF(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()

	r := SetupRotasTeste()
	r.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCPF)

	req, _ := http.NewRequest("GET", "/alunos/cpf/12345678901", nil)
	resposta := httptest.NewRecorder()

	r.ServeHTTP(resposta, req)

	assert.Equal(t, http.StatusOK, resposta.Code)
}

func TestBuscaAlunoID(t *testing.T) {
	database.ConectaComBancoDeDados()

	CriaAlunoMock()
	defer DeletaAlunoMock()

	r := SetupRotasTeste()

	r.GET("/alunos/:id", controllers.BuscaAlunoPorID)

	pathBusca := "/alunos/" + strconv.Itoa(ID)

	req, _ := http.NewRequest("GET", pathBusca, nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	var alunoMock models.Aluno

	json.Unmarshal(res.Body.Bytes(), &alunoMock)

	assert.Equal(t, "Aluno Teste", alunoMock.Nome)
	assert.Equal(t, "12345678901", alunoMock.CPF)
	assert.Equal(t, "123456789", alunoMock.RG)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestDeletaAlunoHandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()

	r := SetupRotasTeste()

	r.DELETE("/alunos/:id", controllers.DeletaAluno)

	pathBusca := "/alunos/" + strconv.Itoa(ID)

	req, _ := http.NewRequest("DELETE", pathBusca, nil)

	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestEditaAlunoHandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupRotasTeste()
	r.PATCH("/alunos/:id", controllers.EditaAluno)

	aluno := models.Aluno{Nome: "Aluno Teste", CPF: "4734567890", RG: "123456700"}

	valorJson, _ := json.Marshal(aluno)

	pathEditar := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("PATCH", pathEditar, bytes.NewBuffer(valorJson))
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	var alunoMockAtualizado models.Aluno
	json.Unmarshal(res.Body.Bytes(), &alunoMockAtualizado)

	assert.Equal(t, http.StatusOK, res.Code)

}
