package msgpack

import (
	"github.com/igorcavalcanti/go_shortener/shortener"
	errs "github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
)

type Redirect struct{}

func (this *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	ret := &shortener.Redirect{}
	var err error

	if err = msgpack.Unmarshal(input, ret); err != nil {
		err = errs.Wrap(err, "serializer.Redirect.Decode")
	}
	return ret, err
}

func (this *Redirect) Encode(redirect *shortener.Redirect) ([]byte, error) {
	raw, err := msgpack.Marshal(redirect)

	if err != nil {
		err = errs.Wrap(err, "serializer.Redirect.Encode")
	}
	return raw, err
}
