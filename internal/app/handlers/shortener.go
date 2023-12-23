package handlers

import (
	"io"
	"net/http"
	"net/url"

	"github.com/levshindenis/sprint1/internal/app/storages"
	"github.com/levshindenis/sprint1/internal/app/tools"
)

type HStorage struct {
	storages.ServerStorage
}

func (serv *HStorage) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "There is not true method", http.StatusBadRequest)
		return
	}
	body, _ := io.ReadAll(r.Body)
	if err := r.Body.Close(); err != nil {
		return
	}

	if _, err := url.ParseRequestURI(string(body)); err != nil {
		http.Error(w, "There is not url", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	myAddress := serv.GetBaseSA() + "/"
	if value, ok := (serv.GetStorage())[string(body)]; ok {
		if _, err := w.Write([]byte(myAddress + value)); err != nil {
			return
		}
	} else {
		shortKey := tools.GenerateShortKey()
		for {
			if _, in := (serv.GetStorage())[shortKey]; !in {
				serv.SetStorage(shortKey, string(body))
				break
			}
			shortKey = tools.GenerateShortKey()
		}
		if _, err := w.Write([]byte(myAddress + shortKey)); err != nil {
			return
		}
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
