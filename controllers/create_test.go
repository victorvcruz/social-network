package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCreate_readBodyAndReturnMapBody(t *testing.T) {

	body := `{"Message": "Hello World"}`

	r := &http.Request{
		Body: ioutil.NopCloser(bytes.NewBufferString(body)),
	}

	mapBody, err := readBodyAndReturnMapBody(r)
	assert.Nil(t, err)

	var mapBodyExpected map[string]interface{}

	err = json.Unmarshal([]byte(body), &mapBodyExpected)
	assert.Nil(t, err)

	assert.Equal(t, mapBodyExpected, mapBody)

}
