package middleware

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ContainerSolutions/bookinfo/bookStockAPI/application"
	"github.com/rs/zerolog/log"
)

// ToJSON serializes the given interface into a string based JSON format
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

// FromJSON deserializes the object from JSON string
// in an io.Reader to the given interface
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}

func readPayload(r *http.Request) (payload []byte, e error) {
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		e = &application.ErrorReadPayload{}
		log.Error().Err(err)
		return
	}
	if len(payload) == 0 {
		e = &application.ErrorPayloadMissing{}
		log.Error().Err(err)
		return
	}
	return
}
