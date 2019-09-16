package mysql

import (
	"fmt"

	"k8s.io/klog"
)

func (d *dbConn) DBInit() {
	db := d.db
	klog.Info("Initializing v1alpha3 DB schema")

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS experiments
		(id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		parameters TEXT,
		objective TEXT,
		algorithm TEXT,
		trial_template TEXT,
		metrics_collector_spec TEXT,
		parallel_trial_count INT,
		max_trial_count INT,
		status TINYINT,
		start_time DATETIME(6),
		completion_time DATETIME(6),
		nas_config TEXT)`)
	//TODO add nas config(may be it will be included in algorithm)
	if err != nil {
		klog.Fatalf("Error creating experiments table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS trials
		(id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		experiment_name VARCHAR(255) NOT NULL,
		objective TEXT,		
		parameter_assignments TEXT,		
		run_spec TEXT,
		metrics_collector_spec TEXT,
		observation TEXT,
		status TINYINT,
		start_time DATETIME(6),
		completion_time DATETIME(6),
		FOREIGN KEY(experiment_name) REFERENCES experiments(name) ON DELETE CASCADE)`)
	if err != nil {
		klog.Fatalf("Error creating trials table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS observation_logs
		(trial_name VARCHAR(255) NOT NULL,
		id INT AUTO_INCREMENT PRIMARY KEY,
		time DATETIME(6),
		metric_name VARCHAR(255) NOT NULL,
		value TEXT NOT NULL)`)
	if err != nil {
		klog.Fatalf("Error creating observation_logs table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS extra_algorithm_settings
		(experiment_name VARCHAR(255) NOT NULL,
		id INT AUTO_INCREMENT PRIMARY KEY,
		setting_name VARCHAR(255) NOT NULL,
		value TEXT NOT NULL,
		FOREIGN KEY (experiment_name) REFERENCES experiments(name) ON DELETE CASCADE)`)
	if err != nil {
		klog.Fatalf("Error creating extra_algorithm_settings table: %v", err)
	}

}

func (d *dbConn) SelectOne() error {
	db := d.db
	_, err := db.Exec(`SELECT 1`)
	if err != nil {
		return fmt.Errorf("Error `SELECT 1` probing: %v", err)
	}
	return nil
}
