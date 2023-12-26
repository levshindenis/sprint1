package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/levshindenis/sprint1/internal/app/storages"
	"github.com/levshindenis/sprint1/internal/app/tools"
	"io"
	"net/http"
	"net/url"
)

type HStorage struct {
	storages.ServerStorage
}

func (serv *HStorage) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "There is not true method", http.StatusBadRequest)
		return
	}
	var body []byte
	if r.Header.Get("Content-Type") == "application/x-gzip" {
		var err error
		body, err = tools.Compression(r.Body)
		if err != nil {
			http.Error(w, "Something bad with compression", http.StatusBadRequest)
			return
		}
	} else {
		body, _ = io.ReadAll(r.Body)
		if _, err := url.ParseRequestURI(string(body)); err != nil {
			http.Error(w, "There is not url", http.StatusBadRequest)
			return
		}
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	myAddress := serv.GetAddress(string(body))

	if _, err := w.Write([]byte(myAddress)); err != nil {
		http.Error(w, "Something bad with write address", http.StatusBadRequest)
		return
	}
}

func (serv *HStorage) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "There is not true method", http.StatusBadRequest)
	}
	if _, in := (serv.GetStorage())[r.URL.Path[1:]]; in {
		w.Header().Add("Location", (serv.GetStorage())[r.URL.Path[1:]])
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "There is no such shortUrl", http.StatusBadRequest)
	}
}

func (serv *HStorage) JSONPostHandler(w http.ResponseWriter, r *http.Request) {
	type Decoder struct {
		MyURL string `json:"url"`
	}
	type Encoder struct {
		MyURL string `json:"result"`
	}
	var enc Encoder
	var dec Decoder
	var buf bytes.Buffer

	if r.Method != http.MethodPost {
		http.Error(w, "There is not true method", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "There is incorrect data format", http.StatusBadRequest)
		return
	}

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, "Something bad with read body", http.StatusBadRequest)
		return
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(buf.Bytes(), &dec); err != nil {
		http.Error(w, "Something bad with decoding JSON", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc.MyURL = serv.GetAddress(dec.MyURL)

	resp, err := json.Marshal(enc)
	if err != nil {
		http.Error(w, "Something bad with encoding JSON", http.StatusBadRequest)
		return
	}
	if _, err = w.Write(resp); err != nil {
		http.Error(w, "Something bad with write address", http.StatusBadRequest)
		return
	}
}
