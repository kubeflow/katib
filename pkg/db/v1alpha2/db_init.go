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
		pallarel_trial_count INT,
		max_trial_count INT,
		condition TINYINT,
		start_time DATETIME(6),
		metrics_collector_type TEXT,
		completion_time DATETIME(6),
		last_reconcile_time DATETIME(6))`)

	if err != nil {
		log.Fatalf("Error creating experiments table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS trials
		(id INT AUTO_INCREMENT PRIMARY KEY,
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
		(trial_id INT NOT NULL,
		id INT AUTO_INCREMENT PRIMARY KEY,
		time DATETIME(6),
		name VARCHAR(255),
		value TEXT,
		FOREIGN KEY (trial_id) REFERENCES trials(id) ON DELETE CASCADE)`)
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
