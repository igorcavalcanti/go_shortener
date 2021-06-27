package api

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	errs "github.com/pkg/errors"

	js "github.com/igorcavalcanti/go_shortener/serializer/json"
	ms "github.com/igorcavalcanti/go_shortener/serializer/msgpack"
	"github.com/igorcavalcanti/go_shortener/shortener"
)

type RedirectHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type handler struct {
	redirectService shortener.RedirectService
}

func NewRedirectHandler(redirectService shortener.RedirectService) RedirectHandler {
	return &handler{
		redirectService: redirectService,
	}
}

func setupResponse(writer http.ResponseWriter, contentType string, body []byte, statusCode int) {
	writer.Header().Set("Content-Type", contentType)
	writer.WriteHeader(statusCode)

	_, err := writer.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func (this *handler) serializer(contentType string) shortener.RedirectSerializer {
	var ret shortener.RedirectSerializer

	if contentType == "application/x-msgpack" {
		ret = &ms.Redirect{}
	} else {
		ret = &js.Redirect{}
	}
	return ret
}

func (this *handler) Get(writer http.ResponseWriter, request *http.Request) {
	code := chi.URLParam(request, "code")

	redirect, err := this.redirectService.Find(code)
	if err != nil {
		if errs.Cause(err) == shortener.ErrRedirectNotFound {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, request, redirect.URL, http.StatusMovedPermanently)
}

func (this *handler) Post(writer http.ResponseWriter, request *http.Request) {
	contentType := request.Header.Get("Content-Type")

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	redirect, err := this.serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = this.redirectService.Store(redirect)
	if err != nil {
		if errs.Cause(err) == shortener.ErrRedirectInvalid {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responseBody, err := this.serializer(contentType).Encode(redirect)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setupResponse(writer, contentType, responseBody, http.StatusCreated)
}
