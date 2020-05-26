package mysql

import (
	"fmt"

	"k8s.io/klog"
)

func (d *dbConn) DBInit() {
	db := d.db
	klog.Info("Initializing v1beta1 DB schema")

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS observation_logs
		(trial_name VARCHAR(255) NOT NULL,
		id INT AUTO_INCREMENT PRIMARY KEY,
		time DATETIME(6),
		metric_name VARCHAR(255) NOT NULL,
		value TEXT NOT NULL)`)
	if err != nil {
		klog.Fatalf("Error creating observation_logs table: %v", err)
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
