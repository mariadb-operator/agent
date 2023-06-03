package recovery

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/mariadb-operator/agent/pkg/errors"
	"github.com/mariadb-operator/agent/pkg/filemanager"
	"github.com/mariadb-operator/agent/pkg/galera"
	"github.com/mariadb-operator/agent/pkg/mariadbd"
	"github.com/mariadb-operator/agent/pkg/responsewriter"
)

var (
	defaultMariadbdReloadOpts = mariadbd.ReloadOptions{
		Retries:   3,
		WaitRetry: 1 * time.Second,
	}
	defaultRecoveryOpts = RecoveryOptions{
		Retries:   10,
		WaitRetry: 3 * time.Second,
	}
)

type RecoveryOptions struct {
	Retries   int
	WaitRetry time.Duration
}

type Recovery struct {
	fileManager           *filemanager.FileManager
	responseWriter        *responsewriter.ResponseWriter
	locker                sync.Locker
	logger                *logr.Logger
	mariadbdReloadOptions *mariadbd.ReloadOptions
	recoveryOptions       *RecoveryOptions
}

type Option func(*Recovery)

func WithMariadbdReload(opts *mariadbd.ReloadOptions) Option {
	return func(b *Recovery) {
		b.mariadbdReloadOptions = opts
	}
}

func WithRecovery(opts *RecoveryOptions) Option {
	return func(r *Recovery) {
		r.recoveryOptions = opts
	}
}

func NewRecover(fileManager *filemanager.FileManager, responseWriter *responsewriter.ResponseWriter, locker sync.Locker,
	logger *logr.Logger, opts ...Option) *Recovery {
	recovery := &Recovery{
		fileManager:           fileManager,
		responseWriter:        responseWriter,
		locker:                locker,
		logger:                logger,
		mariadbdReloadOptions: &defaultMariadbdReloadOpts,
		recoveryOptions:       &defaultRecoveryOpts,
	}
	for _, setOpts := range opts {
		setOpts(recovery)
	}
	return recovery
}

func (r *Recovery) Put(w http.ResponseWriter, req *http.Request) {
	r.locker.Lock()
	defer r.locker.Unlock()
	r.logger.V(1).Info("starting recovery")

	if err := r.fileManager.DeleteConfigFile(galera.BootstrapFileName); err != nil && !os.IsNotExist(err) {
		r.responseWriter.WriteErrorf(w, "error deleting existing bootstrap config: %v", err)
		return
	}

	if err := r.fileManager.DeleteStateFile(galera.RecoveryLogFileName); err != nil && !os.IsNotExist(err) {
		r.responseWriter.WriteErrorf(w, "error deleting existing recovery log: %v", err)
		return
	}

	if err := r.fileManager.WriteConfigFile(galera.RecoveryFileName, []byte(galera.RecoveryFile)); err != nil {
		r.responseWriter.WriteErrorf(w, "error writing recovery config: %v", err)
		return
	}

	r.logger.Info("reloading mariadbd process")
	if err := mariadbd.ReloadWithOptions(r.mariadbdReloadOptions); err != nil {
		r.logger.Error(err, "error reloading mariadbd process")
	} else {
		r.logger.Info("mariadbd process reloaded")
	}

	bootstrap, err := r.recover()
	if err != nil {
		r.responseWriter.WriteErrorf(w, "error recovering galera: %v", err)
		return
	}

	if err := r.fileManager.DeleteConfigFile(galera.RecoveryFileName); err != nil {
		r.responseWriter.WriteErrorf(w, "error deleting recovery config: %v", err)
		return
	}

	r.logger.V(1).Info("finished recovery")
	r.responseWriter.WriteOK(w, bootstrap)
}

func (r *Recovery) Delete(w http.ResponseWriter, req *http.Request) {
	if err := r.fileManager.DeleteConfigFile(galera.RecoveryFileName); err != nil {
		if os.IsNotExist(err) {
			r.responseWriter.Write(w, errors.NewAPIError("recovery config not found"), http.StatusNotFound)
			return
		}
		r.responseWriter.WriteErrorf(w, "error deleting recovery config: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (r *Recovery) recover() (*galera.Bootstrap, error) {
	for i := 0; i < r.recoveryOptions.Retries; i++ {
		time.Sleep(r.recoveryOptions.WaitRetry)

		bytes, err := r.fileManager.ReadStateFile(galera.RecoveryLogFileName)
		if err != nil {
			r.logger.Error(err, "error recovering galera from recovery log", "retry", i, "max-retries", r.recoveryOptions.Retries)
			continue
		}

		var recover galera.Bootstrap
		err = recover.Unmarshal(bytes)
		if err == nil {
			return &recover, nil
		}

		r.logger.Error(err, "error recovering galera from recovery log", "retry", i, "max-retries", r.recoveryOptions.Retries)
	}
	return nil, fmt.Errorf("maximum retries (%d) reached attempting to recover galera from recovery log", r.recoveryOptions.Retries)
}
