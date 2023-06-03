package recovery

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/filemanager"
	"github.com/mariadb-operator/agent/pkg/galera"
	"github.com/mariadb-operator/agent/pkg/handler/jsonencoder"
	"github.com/mariadb-operator/agent/pkg/mariadbd"
)

var (
	defaultMariadbdRetryOpts = mariadbd.RetryOptions{
		Retries:   3,
		WaitRetry: 1 * time.Second,
	}
	defaultRecoverRetryOpts = RecoverRetryOptions{
		Retries:   10,
		WaitRetry: 3 * time.Second,
	}
)

type RecoverRetryOptions struct {
	Retries   int
	WaitRetry time.Duration
}

type Recovery struct {
	fileManager          *filemanager.FileManager
	jsonEncoder          *jsonencoder.JSONEncoder
	logger               *logr.Logger
	mariadbdRetryOptions *mariadbd.RetryOptions
	recoverRetryOptions  *RecoverRetryOptions
}

type Option func(*Recovery)

func WithMariadbdRetry(opts *mariadbd.RetryOptions) Option {
	return func(b *Recovery) {
		b.mariadbdRetryOptions = opts
	}
}

func WithRecoverRetry(opts *RecoverRetryOptions) Option {
	return func(r *Recovery) {
		r.recoverRetryOptions = opts
	}
}

func NewRecover(fileManager *filemanager.FileManager, jsonEncoder *jsonencoder.JSONEncoder, logger *logr.Logger,
	opts ...Option) *Recovery {
	recovery := &Recovery{
		fileManager:          fileManager,
		jsonEncoder:          jsonEncoder,
		logger:               logger,
		mariadbdRetryOptions: &defaultMariadbdRetryOpts,
		recoverRetryOptions:  &defaultRecoverRetryOpts,
	}
	for _, setOpts := range opts {
		setOpts(recovery)
	}
	return recovery
}

func (r *Recovery) Put(w http.ResponseWriter, req *http.Request) {
	if err := r.fileManager.DeleteConfigFile(galera.BootstrapFileName); err != nil && !os.IsNotExist(err) {
		r.logger.Error(err, "error deleting existing bootstrap config")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := r.fileManager.DeleteStateFile(galera.RecoveryLogFileName); err != nil && !os.IsNotExist(err) {
		r.logger.Error(err, "error deleting existing recovery log")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := r.fileManager.WriteConfigFile(galera.RecoveryFileName, []byte(galera.RecoveryFile)); err != nil {
		r.logger.Error(err, "error writing recovery config")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	r.logger.Info("reloading mariadbd process")
	if err := mariadbd.ReloadWithRetries(r.mariadbdRetryOptions); err != nil {
		r.logger.Error(err, "error reloading mariadbd process")
	} else {
		r.logger.Info("mariadbd process reloaded")
	}

	bootstrap, err := r.recover()
	if err != nil {
		r.logger.Error(err, "error recovering galera")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := r.fileManager.DeleteConfigFile(galera.RecoveryFileName); err != nil {
		r.logger.Error(err, "error deleting recovery file")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	r.jsonEncoder.Encode(w, bootstrap)
}

func (r *Recovery) Delete(w http.ResponseWriter, req *http.Request) {
	if err := r.fileManager.DeleteConfigFile(galera.RecoveryFileName); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (r *Recovery) recover() (*galera.Bootstrap, error) {
	for i := 0; i < r.recoverRetryOptions.Retries; i++ {
		time.Sleep(r.recoverRetryOptions.WaitRetry)

		bytes, err := r.fileManager.ReadStateFile(galera.RecoveryLogFileName)
		if err != nil {
			r.logger.Error(err, "error recovering galera from recovery log", "retry", i, "max-retries", r.recoverRetryOptions.Retries)
			continue
		}

		var recover galera.Bootstrap
		err = recover.Unmarshal(bytes)
		if err == nil {
			return &recover, nil
		}

		r.logger.Error(err, "error recovering galera from recovery log", "retry", i, "max-retries", r.recoverRetryOptions.Retries)
	}
	return nil, fmt.Errorf("maximum retries (%d) reached attempting to recover galera from recovery log", r.recoverRetryOptions.Retries)
}
