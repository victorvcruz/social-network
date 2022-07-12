package utils

import (
	"database/sql"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStringEnvOrElse(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		os.Setenv("TEST", "value")
		value := GetStringEnvOrElse("TEST", "not used")
		assert.Equal(t, "value", value)
	})
	t.Run("not found", func(t *testing.T) {
		os.Unsetenv("TEST")
		value := GetStringEnvOrElse("TEST", "value")
		assert.Equal(t, "value", value)
	})
}

func TestGetIntEnvOrElse(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		os.Setenv("TEST", "1")
		value := GetIntEnvOrElse("TEST", 1)
		assert.Equal(t, 1, value)
	})
	t.Run("not found", func(t *testing.T) {
		os.Unsetenv("TEST")
		value := GetIntEnvOrElse("TEST", 1)
		assert.Equal(t, 1, value)
	})
}

func TestReadBodyAndReturnMapBody(t *testing.T) {

	body := `{"Message": "Hello World"}`
	stringReadCloser := io.NopCloser(strings.NewReader(body))

	mapBody, err := ReadBodyAndReturnMapBody(stringReadCloser)
	assert.Nil(t, err)

	var mapBodyExpected map[string]interface{}

	err = json.Unmarshal([]byte(body), &mapBodyExpected)
	assert.Nil(t, err)

	assert.Equal(t, mapBodyExpected, mapBody)
}

func TestNewNullString(t *testing.T) {
	t.Run("len = 0", func(t *testing.T) {
		assert.Equal(t, sql.NullString{
			String: "",
			Valid:  false,
		}, NewNullString(""))
	})
	t.Run("len > 0", func(t *testing.T) {
		assert.Equal(t, sql.NullString{
			String: "1",
			Valid:  true,
		}, NewNullString("1"))
	})
}

func TestStringNullable(t *testing.T) {
	t.Run("str nil", func(t *testing.T) {
		var str interface{} = nil
		assert.Equal(t, "", StringNullable(str))
	})
	t.Run("str not nil", func(t *testing.T) {
		var str interface{} = "1"
		assert.Equal(t, "1", StringNullable(str))
	})

}
