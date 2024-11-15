package app

import (
	"fizzbuzz/pkg/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	t.Run("WithUnsetEnvs", func(t *testing.T) {
		assert.Eq(t, getEnvValOrDefault(EnvLstnAddrName, EnvLstnAddrDefault), EnvLstnAddrDefault)
		assert.Eq(t, getEnvValOrDefault(EnvDBFileDefault, EnvDBFileDefault), EnvDBFileDefault)
	})
	t.Run("WithSetEnvs", func(t *testing.T) {
		addr := "127.0.0.1:8080"
		dbFile := "test.db"
		os.Setenv(EnvLstnAddrName, addr)
		os.Setenv(EnvDBFileName, dbFile)
		assert.Eq(t, getEnvValOrDefault(EnvLstnAddrName, EnvLstnAddrDefault), addr)
		assert.Eq(t, getEnvValOrDefault(EnvDBFileName, EnvDBFileDefault), dbFile)
		os.Unsetenv(EnvLstnAddrName)
		os.Unsetenv(EnvDBFileName)
	})
}

func TestNew(t *testing.T) {
	app := New()
	assert.Eq(t, app.listenAddr, EnvLstnAddrDefault)
	assert.Eq(t, app.dbFile, EnvDBFileDefault)
}

func TestHealth(t *testing.T) {
	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	Health(w, req)
	assert.Eq(t, w.Code, http.StatusOK)
	assert.Eq(t, strings.Contains(w.Body.String(), "OK"), true)
}
