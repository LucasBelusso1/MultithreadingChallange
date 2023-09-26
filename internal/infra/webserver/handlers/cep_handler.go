package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/LucasBelusso1/MultithreadingChallange/internal/dto"
	"github.com/go-chi/chi"
)

type Response struct {
	Error   bool        `json:"erro"`
	Message string      `json:"mensagem,omitempty"`
	Origin  string      `json:"origem,omitempty"`
	Cep     interface{} `json:"cep,omitempty"`
}

func GetAddress(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	ch := make(chan Response)

	go requestToApiCEP(cep, ch)
	go requestToViaCEP(cep, ch)
	go requestToOpenCEP(cep, ch)

	for {
		select {
		case result := <-ch:
			w.Header().Add("Content-Type", "application/json")
			if result.Error {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(result)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Println("Origem:", result.Origin)
			fmt.Println("Resposta:", result)
			json.NewEncoder(w).Encode(result)
			return
		case <-time.After(time.Second):
			fmt.Println("Timeout!")
			w.WriteHeader(http.StatusRequestTimeout)
			return
		}
	}
}

func requestToApiCEP(cep string, response chan<- Response) {
	var apiCepDto dto.ApiCepOutput
	var res Response
	body, err := requestCepFromUrl("https://cdn.apicep.com/file/apicep/" + cep[:5] + "-" + cep[5:] + ".json")

	if err != nil {
		res.Error = true
		res.Message = err.Error()
	} else {
		json.Unmarshal(body, &apiCepDto)
		if apiCepDto.Status != 200 {
			res.Error = true
			res.Message = apiCepDto.Message
		} else {
			res.Cep = apiCepDto
		}
	}

	res.Origin = "API CEP"
	response <- res
}

func requestToViaCEP(cep string, response chan<- Response) {
	var viaCepDto dto.ViaCepOutput
	var res Response
	body, err := requestCepFromUrl("http://viacep.com.br/ws/" + cep + "/json/")

	if err != nil {
		res.Error = true
		res.Message = err.Error()
	} else {
		json.Unmarshal(body, &viaCepDto)
		res.Cep = viaCepDto
	}

	res.Origin = "Via CEP"
	response <- res
}

func requestToOpenCEP(cep string, response chan<- Response) {
	var openCepDto dto.OpenCepOutput
	var res Response
	body, err := requestCepFromUrl("http://opencep.com/v1/" + cep)

	if err != nil {
		res.Error = true
		res.Message = err.Error()
	} else {
		json.Unmarshal(body, &openCepDto)
		res.Cep = openCepDto
	}

	res.Origin = "Open CEP"
	response <- res
}

func requestCepFromUrl(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
