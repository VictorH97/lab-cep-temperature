package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var weatherAPIKey = "2083c7dd8e734e46971234222242102"

func TestValidCEP(t *testing.T) {
	valid, err := VerifyValidCEP("12345678")
	assert.NoError(t, err)
	assert.True(t, valid)

	valid, err = VerifyValidCEP("12345-678")
	assert.NoError(t, err)
	assert.True(t, valid)

	valid, err = VerifyValidCEP("123456")
	assert.NoError(t, err)
	assert.False(t, valid)
}

func TestGetCEPInfo(t *testing.T) {
	cepInfo, err := GetCEPInfo("12345678")
	assert.Error(t, err, "can not find zipcode")
	assert.Nil(t, cepInfo)

	cepInfo, err = GetCEPInfo("02201001")
	assert.NoError(t, err)
	assert.NotNil(t, cepInfo)
}

func TestGetWeatherInfo(t *testing.T) {
	weatherInfo, err := GetWeatherInfo("notexist", weatherAPIKey)
	assert.Error(t, err, "No matching location found.")
	assert.Nil(t, weatherInfo)

	weatherInfo, err = GetWeatherInfo("São Paulo", weatherAPIKey)
	assert.NoError(t, err)
	assert.NotNil(t, weatherInfo)
}
