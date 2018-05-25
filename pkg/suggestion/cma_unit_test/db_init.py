import mysql.connector

TABLES = {}
TABLES['studies'] = (
    "CREATE TABLE IF NOT EXISTS `studies` ("
    "  `id` CHAR(16) PRIMARY KEY ,"
    "  `name` VARCHAR (255),"
    "  `owner` VARCHAR (255),"
    "  `optimization_type` TINYINT,"
    "  `optimization_goal` DOUBLE ,"
    "  `parameter_configs` TEXT,"
    "  `suggest_algo` VARCHAR (255),"
    "  `early_stop_algo` VARCHAR (255),"
    "  `tags` TEXT,"
    "  `objective_value_name` VARCHAR (255),"
    "  `metrics` TEXT"
    ") ENGINE=InnoDB")

TABLES['trials'] = (
    "CREATE TABLE IF NOT EXISTS `trials` ("
    "  `id` CHAR(16) PRIMARY KEY ,"
    "  `study_id` CHAR (16),"
    "  `parameters` TEXT,"
    "  `status` TINYINT,"
    "  `objective_value` VARCHAR (255),"
    "  `tag` TEXT,"
    "  FOREIGN KEY (`study_id`) REFERENCES studies(`id`)"
    ") ENGINE=InnoDB")

TABLES['suggestion_param'] = (
    "CREATE TABLE IF NOT EXISTS `suggestion_param` ("
    "  `id` CHAR(16) PRIMARY KEY ,"
    "  `suggestion_algo` TEXT,"
    "  `parameters` TEXT,"
    "  `study_id` CHAR (16)"
    ") ENGINE=InnoDB")


def connect_db():
    cnx = mysql.connector.connect(
        user="root",
        password="zhangyingbo",
        database="vizier",
        host="localhost",
        port=3306,
    )

    return cnx


# def DB_Init(cnx):
#     cursor = cnx.cursor()
#     for name, ddl in TABLES.items():
#         cursor.execute(ddl)
#
#     cursor.close()
#     cnx.close()


def run():
    cnx = connect_db()
    # DB_Init(cnx)


# if __name__ == "__main__":
#     run()
