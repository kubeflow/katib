package db

import (
	crand "crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"time"

	v1alpha2 "github.com/kubeflow/katib/pkg/api/v1alpha2"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriver     = "mysql"
	dbNameTmpl   = "root:%s@tcp(vizier-db:3306)/vizier?timeout=5s"
	mysqlTimeFmt = "2006-01-02 15:04:05.999999"

	connectInterval = 5 * time.Second
	connectTimeout  = 60 * time.Second
)

type GetWorkerLogOpts struct {
	Name       string
	SinceTime  *time.Time
	Descending bool
	Limit      int32
	Objective  bool
}

type WorkerLog struct {
	Time  time.Time
	Name  string
	Value string
}

type KatibDBInterface interface {
	DBInit()
	SelectOne() error

	RegisterExperiment(experiment *v1alpha2.Experiment) error
	DeleteExperiment(experimentName string) error
	GetExperiment(experimentName string) (*v1alpha2.Experiment, error)
	GetExperimentList() ([]*v1alpha2.ExperimentSummary, error)
	UpdateExperimentStatus(experimentName string, newStatus *v1alpha2.ExperimentStatus) error

	RegisterTrial(trial *v1alpha2.Trial) error
	GetTrialList(experimentName string) ([]*v1alpha2.Trial, error)
	GetTrial(trialName string) (*v1alpha2.Trial, error)
	UpdateTrialStatus(trialName string, newStatus *v1alpha2.TrialStatus) error
	DeleteTrial(trialName string) error

	RegisterObservationLog(trialName string, obsercationLog *v1alpha2.ObservationLog) error
	GetObservationLog(trialName string, startTime time.Time, endTime time.Time) (*v1alpha2.ObservationLog, error)
}

type dbConn struct {
	db *sql.DB
}

var rs1Letters = []rune("abcdefghijklmnopqrstuvwxyz")

func getDbName() string {
	dbPass := os.Getenv("MYSQL_ROOT_PASSWORD")
	if dbPass == "" {
		log.Printf("WARN: Env var MYSQL_ROOT_PASSWORD is empty. Falling back to \"test\".")

		// For backward compatibility, e.g. in case that all but vizier-core
		// is older ones so we do not have Secret nor upgraded vizier-db.
		dbPass = "test"
	}

	return fmt.Sprintf(dbNameTmpl, dbPass)
}

func openSQLConn(driverName string, dataSourceName string, interval time.Duration,
	timeout time.Duration) (*sql.DB, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	timeoutC := time.After(timeout)
	for {
		select {
		case <-ticker.C:
			if db, err := sql.Open(driverName, dataSourceName); err == nil {
				if err = db.Ping(); err == nil {
					return db, nil
				}
			}
		case <-timeoutC:
			return nil, fmt.Errorf("Timeout waiting for DB conn successfully opened.")
		}
	}
}

func NewWithSQLConn(db *sql.DB) (KatibDBInterface, error) {
	d := new(dbConn)
	d.db = db
	seed, err := crand.Int(crand.Reader, big.NewInt(1<<63-1))
	if err != nil {
		return nil, fmt.Errorf("RNG initialization failed: %v", err)
	}
	// We can do the following instead, but it creates a locking issue
	//d.rng = rand.New(rand.NewSource(seed.Int64()))
	rand.Seed(seed.Int64())

	return d, nil
}

func (d *dbConn) RegisterExperiment(experiment *v1alpha2.Experiment) error {
	return nil
}
func (d *dbConn) DeleteExperiment(experimentName string) error {
	return nil
}
func (d *dbConn) GetExperiment(experimentName string) (*v1alpha2.Experiment, error) {
	return nil, nil
}
func (d *dbConn) GetExperimentList() ([]*v1alpha2.ExperimentSummary, error) {
	return nil, nil
}
func (d *dbConn) UpdateExperimentStatus(experimentName string, newStatus *v1alpha2.ExperimentStatus) error {
	return nil
}

func (d *dbConn) RegisterTrial(trial *v1alpha2.Trial) error {
	return nil
}
func (d *dbConn) GetTrialList(experimentName string) ([]*v1alpha2.Trial, error) {
	return nil, nil
}
func (d *dbConn) GetTrial(trialName string) (*v1alpha2.Trial, error) {
	return nil, nil
}
func (d *dbConn) UpdateTrialStatus(trialName string, newStatus *v1alpha2.TrialStatus) error {
	return nil
}
func (d *dbConn) DeleteTrial(trialName string) error {
	return nil
}

func (d *dbConn) RegisterObservationLog(trialName string, obsercationLog *v1alpha2.ObservationLog) error {
	return nil
}
func (d *dbConn) GetObservationLog(trialName string, startTime time.Time, endTime time.Time) (*v1alpha2.ObservationLog, error) {
	return nil, nil
}
