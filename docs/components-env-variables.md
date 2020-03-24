# Katib components environment variables

Here you can find information about environment variables in each Katib component. If you want to deploy Katib not in default way, you can modify some of this variables.

## Katib Controller

This is the list of environment variables for the [Katib controller](https://github.com/kubeflow/katib/blob/master/manifests/v1alpha3/katib-controller/katib-controller.yaml) deployment.

| Variable                             | Description                                                       | Default Value                       |
| ------------------------------------ | ----------------------------------------------------------------- | ----------------------------------- |
| `KATIB_CORE_NAMESPACE`               | Base Namespace for all Katib components and default Experiment    | kubeflow                            |
| `KATIB_SUGGESTION_COMPOSER`          | Composer for the Katib Suggestions. You can use your own Composer | general                             |
| `KATIB_DB_MANAGER_SERVICE_NAMESPACE` | Katib DB Manager Namespace                                        | `KATIB_CORE_NAMESPACE` env variable |
| `KATIB_DB_MANAGER_SERVICE_IP`        | Katib DB Manager IP                                               | katib-db-manager                    |
| `KATIB_DB_MANAGER_SERVICE_PORT`      | Katib DB Manager Port                                             | 6789                                |

Katib DB Manager address is created by this expression:

`KATIB_DB_MANAGER_SERVICE_IP` . `KATIB_DB_MANAGER_SERVICE_NAMESPACE` : `KATIB_DB_MANAGER_SERVICE_PORT`

If you set `KATIB_DB_MANAGER_SERVICE_NAMESPACE=""`, Katib DB Manager address will be:

`KATIB_DB_MANAGER_SERVICE_IP` : `KATIB_DB_MANAGER_SERVICE_PORT`

## Katib UI

This is the list of environment variables for the [Katib UI](https://github.com/kubeflow/katib/blob/master/manifests/v1alpha3/ui/deployment.yaml) deployment.

| Variable                             | Description                                                    | Default Value                       |
| ------------------------------------ | -------------------------------------------------------------- | ----------------------------------- |
| `KATIB_CORE_NAMESPACE`               | Base Namespace for all Katib components and default Experiment | kubeflow                            |
| `KATIB_DB_MANAGER_SERVICE_NAMESPACE` | Katib DB Manager Namespace                                     | `KATIB_CORE_NAMESPACE` env variable |
| `KATIB_DB_MANAGER_SERVICE_IP`        | Katib DB Manager IP                                            | katib-db-manager                    |
| `KATIB_DB_MANAGER_SERVICE_PORT`      | Katib DB Manager Port                                          | 6789                                |

Katib DB Manager address is created by above expression.

## Katib DB Manager

This is the list of environment variables for the [Katib DB manager](https://github.com/andreyvelich/katib/blob/doc-katib-config/manifests/v1alpha3/db-manager/deployment.yaml) deployment.

| Variable                  | Description                    | Default Value    |
| ------------------------- | ------------------------------ | ---------------- |
| `DB_NAME`                 | Katib DB Name, must be `mysql` | No default value |
| `DB_PASSWORD`             | Katib DB Password, must be set | No default value |
| `DB_USER`                 | Katib DB user                  | root             |
| `KATIB_MYSQL_DB_HOST`     | Katib MYSQL host               | katib-mysql      |
| `KATIB_MYSQL_DB_PORT`     | Katib MYSQL port               | 3306             |
| `KATIB_MYSQL_DB_DATABASE` | Katib MYSQL database name      | katib            |

Katib DB manager creates connection to the DB, using `mysql` driver and this data source name:

`DB_USER` : `DB_PASSWORD`@tcp(`KATIB_MYSQL_DB_HOST` : `KATIB_MYSQL_DB_PORT`)/`KATIB_MYSQL_DB_DATABASE`?timeout=5s

## Katib MySQL DB

For [Katib MySQL](https://github.com/kubeflow/katib/blob/master/manifests/v1alpha3/mysql-db/deployment.yaml) we set `MYSQL_ROOT_PASSWORD` as value from [katib-mysql-secrets](https://github.com/kubeflow/katib/blob/master/manifests/v1alpha3/mysql-db/secret.yaml), `MYSQL_ALLOW_EMPTY_PASSWORD` as `true`, `MYSQL_DATABASE` as `katib`.

Check [here](https://github.com/docker-library/docs/tree/master/mysql#environment-variables) about all environment variable for MySQL docker image.

Katib MySQL env variables must be matched with Katib DB Manager env variables.
