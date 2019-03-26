package db

import (
	"fmt"
	"log"
)

func (d *dbConn) DBInit() {
	db := d.db

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS experiments
		(id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255),
		parameters TEXT,
		objective TEXT,
		algorithm TEXT,
		trial_template TEXT,
		parallel_trial_count INT,
		max_trial_count INT,
		condition TINYINT,
		start_time DATETIME(6),
		metrics_collector_type TEXT,
		completion_time DATETIME(6),
		last_reconcile_time DATETIME(6))`)
	//TODO add nas config(may be it will be included in algorithm)
	if err != nil {
		log.Fatalf("Error creating experiments table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS trials
		(id INT AUTO_INCREMENT PRIMARY KEY,
		trial_name VARCHAR(255),
		experiment_name TEXT,
		parameter_assignments TEXT,
		run_spec TEXT,
		observation TEXT,
		condition TINYINT,
		start_time DATETIME(6),
		completion_time DATETIME(6),
		last_reconcile_time DATETIME(6),
		FOREIGN KEY(experiment_name) REFERENCES experiments(name) ON DELETE CASCADE)`)
	if err != nil {
		log.Fatalf("Error creating trials table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS observation_logs
		(trial_name VARCHAR(255) NOT NULL,
		id INT AUTO_INCREMENT PRIMARY KEY,
		time DATETIME(6),
		metric_name VARCHAR(255),
		value TEXT,
		FOREIGN KEY (trial_name) REFERENCES trials(trial_name) ON DELETE CASCADE)`)
	if err != nil {
		log.Fatalf("Error creating observation_logs table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS algorithm_variables
		(experiment_name VARCHAR(255) NOT NULL,
		id INT AUTO_INCREMENT PRIMARY KEY,
		variable_name VARCHAR(255),
		value TEXT,
		FOREIGN KEY (experiment_name) REFERENCES experiments(name) ON DELETE CASCADE)`)
	if err != nil {
		log.Fatalf("Error creating observation_logs table: %v", err)
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
