package bootstrap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/filemanager"
	"github.com/mariadb-operator/agent/pkg/galera"
	"github.com/mariadb-operator/agent/pkg/mariadbd"
)

var (
	defaultMariadbdRetryOpts = mariadbd.RetryOptions{
		Retries:   10,
		WaitRetry: 1 * time.Second,
	}
)

type Bootstrap struct {
	fileManager          *filemanager.FileManager
	logger               *logr.Logger
	mariadbdRetryOptions *mariadbd.RetryOptions
}

type Option func(*Bootstrap)

func WithMariadbdRetry(opts *mariadbd.RetryOptions) Option {
	return func(b *Bootstrap) {
		b.mariadbdRetryOptions = opts
	}
}

func NewBootstrap(fileManager *filemanager.FileManager, logger *logr.Logger, opts ...Option) *Bootstrap {
	bootstrap := &Bootstrap{
		fileManager:          fileManager,
		logger:               logger,
		mariadbdRetryOptions: &defaultMariadbdRetryOpts,
	}
	for _, setOpts := range opts {
		setOpts(bootstrap)
	}
	return bootstrap
}

func (b *Bootstrap) Put(w http.ResponseWriter, r *http.Request) {
	var bootstrap galera.Bootstrap
	if err := json.NewDecoder(r.Body).Decode(&bootstrap); err != nil {
		b.logger.Error(err, "error decoding bootstrap")
		http.Error(w, "invalid body: a valid bootstrap object must be provided", http.StatusBadRequest)
		return
	}
	if err := bootstrap.Validate(); err != nil {
		b.logger.Error(err, "invalid bootstrap")
		http.Error(w, fmt.Sprintf("invalid bootstrap: %v", err), http.StatusBadRequest)
		return
	}

	if err := b.fileManager.DeleteConfigFile(galera.RecoveryFileName); err != nil && !os.IsNotExist(err) {
		b.logger.Error(err, "error deleting existing recovery config")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := b.setSafeToBootstrap(&bootstrap); err != nil {
		b.logger.Error(err, "error setting safe to bootstrap")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := b.fileManager.WriteConfigFile(galera.BootstrapFileName, []byte(galera.BootstrapFile)); err != nil {
		b.logger.Error(err, "error writing bootstrap file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	b.logger.Info("reloading mariadbd process")
	if err := mariadbd.ReloadWithRetries(b.mariadbdRetryOptions); err != nil {
		b.logger.Error(err, "error reloading mariadbd process")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	b.logger.Info("mariadbd process reloaded")

	w.WriteHeader(http.StatusOK)
}

func (b *Bootstrap) Delete(w http.ResponseWriter, r *http.Request) {
	if err := b.fileManager.DeleteConfigFile(galera.BootstrapFileName); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		b.logger.Error(err, "error deleting bootstrap file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (b *Bootstrap) setSafeToBootstrap(bootstrap *galera.Bootstrap) error {
	bytes, err := b.fileManager.ReadStateFile(galera.GaleraStateFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error reading galera state: %v", err)
	}

	var galeraState galera.GaleraState
	if err := galeraState.Unmarshal(bytes); err != nil {
		return fmt.Errorf("error unmarshaling galera state: %v", err)
	}

	galeraState.UUID = bootstrap.UUID
	galeraState.Seqno = bootstrap.Seqno
	galeraState.SafeToBootstrap = true
	bytes, err = galeraState.Marshal()
	if err != nil {
		return fmt.Errorf("error marshallng galera state: %v", err)
	}

	if err := b.fileManager.WriteStateFile(galera.GaleraStateFileName, bytes); err != nil {
		return fmt.Errorf("error writing galera state: %v", err)
	}
	return nil
}