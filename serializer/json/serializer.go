package json

import (
	"encoding/json"

	"github.com/igorcavalcanti/go_shortener/shortener"
	errs "github.com/pkg/errors"
)

type Redirect struct{}

func (this *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	ret := &shortener.Redirect{}
	var err error

	if err = json.Unmarshal(input, ret); err != nil {
		err = errs.Wrap(err, "serializer.Redirect.Decode")
	}
	return ret, err
}

func (this *Redirect) Encode(redirect *shortener.Redirect) ([]byte, error) {
	raw, err := json.Marshal(redirect)

	if err != nil {
		err = errs.Wrap(err, "serializer.Redirect.Encode")
	}
	return raw, err
}
