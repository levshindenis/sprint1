package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/levshindenis/sprint1/internal/app/models"
)

func (serv *HStorage) SetJsonLongURL(w http.ResponseWriter, r *http.Request) {
	var (
		enc  models.JsonEncoder
		dec  models.JsonDecoder
		buf  bytes.Buffer
		flag bool
		err  error
	)

	cookie, _ := r.Cookie("UserID")
	http.SetCookie(w, &http.Cookie{Name: "UserID", Value: cookie.Value})

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "There is incorrect data format", http.StatusBadRequest)
		return
	}

	if _, err = buf.ReadFrom(r.Body); err != nil {
		http.Error(w, "Something bad with read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(buf.Bytes(), &dec); err != nil {
		http.Error(w, "Something bad with decoding JSON", http.StatusBadRequest)
		return
	}

	if _, err = url.ParseRequestURI(dec.LongURL); err != nil {
		http.Error(w, "There is not url", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	enc.ShortURL, flag, err = serv.MakeShortURL(dec.LongURL)
	if err != nil {
		http.Error(w, "Something bad with MakeShortURL", http.StatusBadRequest)
		return
	}
	if flag {
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err = serv.GetStorage().SetData(enc.ShortURL, dec.LongURL, cookie.Value); err != nil {
		http.Error(w, "Something bad with Save", http.StatusBadRequest)
		return
	}

	enc.ShortURL = serv.GetServerConfig("baseURL") + "/" + enc.ShortURL

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
