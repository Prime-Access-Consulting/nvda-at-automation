package client

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var validCapabilities = Capabilities{
	Name:     "testNVDA",
	Version:  "1.2",
	Platform: "windows",
}

var validSettings = Settings{}

type requestAssertions func(t *testing.T, r *http.Request)

func getSettingsHandlerFunc(t *testing.T, ra requestAssertions, returned *Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "client should issue a GET request")
		assert.Contains(t, r.URL.String(), "/settings?q=", "client should call the /settings endpoint")

		ra(t, r)

		response, _ := json.Marshal(validSettings)

		if returned != nil {
			response, _ = json.Marshal(returned)
		}
		_, err := w.Write(response)

		if err != nil {
			panic(err)
		}
	}
}

func getInfoHandlerFunc(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "client should issue a GET request")
		assert.Equal(t, "/info", r.URL.String(), "client should call the /info endpoint")

		response, _ := json.Marshal(validCapabilities)

		_, err := w.Write(response)

		if err != nil {
			panic(err)
		}
	}
}

func assertIsNvdaClientError(t *testing.T, err error, substr string) {
	assert.ErrorContains(t, err, "unexpected response from NVDA")
	assert.ErrorContains(t, err, substr)
}

func TestUnavailablePluginAddonFails(t *testing.T) {
	nvda, err := New("https://some-nvda.dev:9000")
	assert.Nil(t, nvda)
	assertIsNvdaClientError(t, err, "no such host")
}

func TestGetNvdaInfoCallsPluginAddon(t *testing.T) {
	ts := httptest.NewServer(getInfoHandlerFunc(t))
	defer ts.Close()

	nvda, err := New(ts.URL)

	assert.Nil(t, err)
	assert.IsType(t, &NVDA{}, nvda)
}

func TestInvalidGetInfoResponseFails(t *testing.T) {
	getInfoHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "some invalid response")

		if err != nil {
			panic(err)
		}
	}

	ts := httptest.NewServer(http.HandlerFunc(getInfoHandlerFunc))
	defer ts.Close()

	nvda, err := New(ts.URL)

	assert.Nil(t, nvda)
	assertIsNvdaClientError(t, err, "invalid character")
}

func TestValidGetInfoResponseProvidesCapabilities(t *testing.T) {
	ts := httptest.NewServer(getInfoHandlerFunc(t))
	defer ts.Close()

	nvda, err := New(ts.URL)

	assert.Nil(t, err)
	assert.Equal(t, validCapabilities.Name, nvda.Capabilities.Name)
	assert.Equal(t, validCapabilities.Version, nvda.Capabilities.Version)
	assert.Equal(t, validCapabilities.Platform, nvda.Capabilities.Platform)
}

func runGetSettingsTest(t *testing.T, requestedSettings []string, r requestAssertions, returnedSettings *Settings) Settings {
	mux := http.NewServeMux()
	mux.HandleFunc("/info", getInfoHandlerFunc(t))
	mux.HandleFunc("/settings", getSettingsHandlerFunc(t, r, returnedSettings))

	ts := httptest.NewServer(mux)
	defer ts.Close()

	nvda, err := New(ts.URL)

	assert.Nil(t, err)

	settings, err := nvda.GetSettings(requestedSettings)
	assert.Nil(t, err)
	assert.NotNil(t, settings)

	return *settings
}

func TestGetZeroSettingsCallsPluginEndpoint(t *testing.T) {
	requestAssertions := func(t *testing.T, r *http.Request) {
		query := r.URL.Query()
		assert.Contains(t, query, "q")
		assert.Equal(t, []string{""}, query["q"])
	}

	runGetSettingsTest(t, []string{}, requestAssertions, nil)
}

func TestGetSettingsBuildsQueryString(t *testing.T) {
	requestAssertions := func(t *testing.T, r *http.Request) {
		query := r.URL.Query()
		assert.Contains(t, query, "q")
		assert.Equal(t, []string{"first,second,third"}, query["q"])
	}

	runGetSettingsTest(t, []string{"first", "second", "third"}, requestAssertions, nil)
}

func TestGetSettingsProvidesSettings(t *testing.T) {
	requestAssertions := func(t *testing.T, r *http.Request) {}
	settings := Settings{
		"first":  0.0,
		"second": "foo",
		"third":  "some_other_value",
	}

	response := runGetSettingsTest(t, []string{}, requestAssertions, &settings)

	assert.EqualValues(t, settings, response)
}
