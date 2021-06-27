package shortener

import (
	"errors"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrRedirectNotFound = errors.New("Redirect Not Found")
	ErrRedirectInvalid  = errors.New("Redirect Invalid")
)

type redirectService struct {
	redirectRepository RedirectRepository
}

func NewRedirectService(repository RedirectRepository) RedirectService {
	return &redirectService{
		redirectRepository: repository,
	}
}

func (this *redirectService) Find(code string) (*Redirect, error) {
	return this.redirectRepository.Find(code)
}

func (this *redirectService) Store(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errs.Wrap(ErrRedirectInvalid, "service.Redirect.Store")
	}
	redirect.Code = shortid.MustGenerate()
	redirect.CreateAt = time.Now().UTC().Unix()

	return this.redirectRepository.Store(redirect)
}
