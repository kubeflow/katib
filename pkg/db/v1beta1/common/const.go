/*
Copyright 2022 The Kubeflow Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import "time"

const (
	ConnectInterval = 5 * time.Second

	DBUserEnvName     = "DB_USER"
	DBNameEnvName     = "DB_NAME"
	DBPasswordEnvName = "DB_PASSWORD"

	MySqlDBNameEnvValue = "mysql"

	MySQLDBHostEnvName = "KATIB_MYSQL_DB_HOST"
	MySQLDBPortEnvName = "KATIB_MYSQL_DB_PORT"
	MySQLDatabase      = "KATIB_MYSQL_DB_DATABASE"

	DefaultMySQLUser     = "root"
	DefaultMySQLDatabase = "katib"
	DefaultMySQLHost     = "katib-mysql"
	DefaultMySQLPort     = "3306"

	PostgresSQLDBNameEnvValue = "postgres"

	PostgreSQLDBHostEnvName = "KATIB_POSTGRESQL_DB_HOST"
	PostgreSQLDBPortEnvName = "KATIB_POSTGRESQL_DB_PORT"
	PostgreSQLDatabase      = "KATIB_POSTGRESQL_DB_DATABASE"
	PostgreSSLMode          = "KATIB_POSTGRESQL_SSL_MODE"

	DefaultPostgreSQLUser     = "katib"
	DefaultPostgreSQLDatabase = "katib"
	DefaultPostgreSQLHost     = "katib-postgres"
	DefaultPostgreSQLPort     = "5432"
	DefaultPostgreSSLMode     = "disable"

	SkipDbInitializationEnvName = "SKIP_DB_INITIALIZATION"
)
