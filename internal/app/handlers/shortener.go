package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/levshindenis/sprint1/internal/app/storages"
	"github.com/levshindenis/sprint1/internal/app/tools"
	"io"
	"net/http"
	"net/url"
	"time"
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
	var err error
	if r.Header.Get("Content-Type") == "application/x-gzip" {
		body, err = tools.Compression(r.Body)
		if err != nil {
			http.Error(w, "Something bad with compression", http.StatusBadRequest)
			return
		}
	} else {
		body, _ = io.ReadAll(r.Body)
		if _, err = url.ParseRequestURI(string(body)); err != nil {
			http.Error(w, "There is not url", http.StatusBadRequest)
			return
		}
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusCreated)

	address, flag, err := serv.MakeShortURL(string(body))
	if err != nil {
		http.Error(w, "Something bad with MakeShortURL", http.StatusBadRequest)
		return
	}
	if !flag {
		if err = serv.Save(address, string(body)); err != nil {
			http.Error(w, "Something bad with Save", http.StatusBadRequest)
			return
		}
	}

	address = serv.GetConfigParameter("baseURL") + "/" + address
	if _, err := w.Write([]byte(address)); err != nil {
		http.Error(w, "Something bad with write address", http.StatusBadRequest)
		return
	}
}

func (serv *HStorage) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "There is not true method", http.StatusBadRequest)
	}

	if result, err := serv.Get(r.URL.Path[1:], "key"); err == nil && result != "" {
		w.Header().Add("Location", result)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "There is no such shortUrl", http.StatusBadRequest)
	}
}

func (serv *HStorage) JSONPostHandler(w http.ResponseWriter, r *http.Request) {
	type Decoder struct {
		LongURL string `json:"url"`
	}
	type Encoder struct {
		ShortURL string `json:"result"`
	}
	var enc Encoder
	var dec Decoder
	var buf bytes.Buffer
	var err error

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
	defer r.Body.Close()

	if err := json.Unmarshal(buf.Bytes(), &dec); err != nil {
		http.Error(w, "Something bad with decoding JSON", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var flag bool
	enc.ShortURL, flag, err = serv.MakeShortURL(dec.LongURL)
	if err != nil {
		http.Error(w, "Something bad with MakeShortURL", http.StatusBadRequest)
		return
	}

	if !flag {
		if err = serv.Save(enc.ShortURL, dec.LongURL); err != nil {
			http.Error(w, "Something bad with Save", http.StatusBadRequest)
			return
		}
	}

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

func (serv *HStorage) GetPingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "There is not true method", http.StatusBadRequest)
	}

	db, err := sql.Open("pgx", serv.GetConfigParameter("db"))
	if err != nil {
		http.Error(w, "Something bad with open db", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		http.Error(w, "Something bad with ping", http.StatusInternalServerError)
		return
	}
}
