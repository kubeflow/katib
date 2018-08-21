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
		"tags TEXT, " +
		"objective_value_name VARCHAR(255), " +
		"metrics TEXT, " +
		"job_id TEXT)")
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
		"objective_value VARCHAR(255), " +
		"tags TEXT, " +
		"FOREIGN KEY(study_id) REFERENCES studies(id))")
	if err != nil {
		log.Fatalf("Error creating trials table: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS workers" +
		"(id CHAR(16) PRIMARY KEY, " +
		"study_id CHAR(16), " +
		"trial_id CHAR(16), " +
		"type VARCHAR(255), " +
		"status TINYINT, " +
		"template_path TEXT, " +
		"tags TEXT, " +
		"FOREIGN KEY(study_id) REFERENCES studies(id), " +
		"FOREIGN KEY(trial_id) REFERENCES trials(id))")
	if err != nil {
		log.Fatalf("Error creating workers table: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS worker_metrics" +
		"(worker_id CHAR(16) NOT NULL, " +
		"id INT AUTO_INCREMENT PRIMARY KEY, " +
		"time DATETIME(6), " +
		"name VARCHAR(255), " +
		"value TEXT, " +
		"is_objective TINYINT)")
	if err != nil {
		log.Fatalf("Error creating worker_metrics table: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS worker_lastlogs" +
		"(worker_id CHAR(16) PRIMARY KEY, " +
		"time DATETIME(6), " +
		"value TEXT)")
	if err != nil {
		log.Fatalf("Error creating worker_lastlogs table: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS suggestion_param" +
		"(id CHAR(16) PRIMARY KEY," +
		"suggestion_algo TEXT, " +
		"study_id CHAR(16), " +
		"parameters TEXT, " +
		"FOREIGN KEY(study_id) REFERENCES studies(id))")
	if err != nil {
		log.Fatalf("Error creating suggestion_param table: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS earlystop_param" +
		"(id CHAR(16) PRIMARY KEY, " +
		"earlystop_argo TEXT, " +
		"study_id CHAR(16), " +
		"parameters TEXT, " +
		"FOREIGN KEY(study_id) REFERENCES studies(id))")
	if err != nil {
		log.Fatalf("Error creating earlystop_param table: %v", err)
	}
}
