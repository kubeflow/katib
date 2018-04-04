package db

import (
	"log"
)

func (d *db_conn) DB_Init() {
	db := d.db
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS studies" +
		"(id CHAR(16) PRIMARY KEY, " +
		"name VARCHAR(255), " +
		"owner VARCHAR(255), " +
		"optimization_type TINYINT, " +
		"optimization_goal DOUBLE, " +
		"parameter_configs TEXT, " +
		"suggest_algo VARCHAR(255), " +
		"autostop_algo VARCHAR(255), " +
		"study_task_name VARCHAR(255), " +
		"suggestion_parameters TEXT, " +
		"tags TEXT, " +
		"objective_value_name VARCHAR(255), " +
		"metrics TEXT, " +
		"image VARCHAR(255), " +
		"command TEXT, " +
		"gpu INT, " +
		"scheduler VARCHAR(255), " +
		"mount TEXT, " +
		"pull_secret TEXT)")
	if err != nil {
		log.Fatalf("Error creating studies table: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS study_permissions" +
		"(study_id CHAR(16) NOT NULL, " +
		"access_permission VARCHAR(255), " +
		"PRIMARY KEY (study_id, access_permission))")
	if err != nil {
		log.Fatalf("Error creating study_permissions table: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS trials" +
		"(id CHAR(16) PRIMARY KEY, " +
		"study_id CHAR(16), " +
		"parameters TEXT, " +
		"status TINYINT, " +
		"objective_value VARCHAR(255), " +
		"tags TEXT, " +
		"FOREIGN KEY(study_id) REFERENCES studies(id))")
	if err != nil {
		log.Fatalf("Error creating trials table: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS trial_logs" +
		"(trial_id CHAR(16) NOT NULL, " +
		"time DATETIME(6), " +
		"value TEXT, " +
		"PRIMARY KEY (trial_id, time))")
	// We can have "id INT AUTO_INCREMENT PRIMARY KEY" instead.
	if err != nil {
		log.Fatalf("Error creating trial_logs table: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS workers" +
		"(id CHAR(16) PRIMARY KEY, " +
		"trial_id CHAR(16), " +
		"status TINYINT, " +
		"FOREIGN KEY(trial_id) REFERENCES trials(id))")
	if err != nil {
		log.Fatalf("Error creating workers table: %v", err)
	}
}
