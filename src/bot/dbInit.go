package bot

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var password = os.Getenv("PASSWORD")
var dbName = os.Getenv("DBNAME")

var dbInfo = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)
var initDbInfo = fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, password, host, port)

func InitDb() error {
	db, err := sql.Open("mysql", initDbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE DATABASE IF NOT EXISTS ` + dbName)
	if err != nil {
		return err
	}

	_, err = db.Exec(`USE ` + dbName)
	if err != nil {
		return err
	}

	return nil
}

func CreateTables() error {
	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(
    	user_id BIGINT UNSIGNED NOT NULL,
		username VARCHAR(255) NOT NULL,
		first_name VARCHAR(255) NULL,
		last_name VARCHAR(255) NULL,
		PRIMARY KEY (user_id),
		UNIQUE INDEX ID_UNIQUE (user_id ASC) VISIBLE,
		UNIQUE INDEX USERNAME_UNIQUE (username ASC) VISIBLE
        ) ENGINE = InnoDB;`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS admins(
		user_ID BIGINT UNSIGNED NOT NULL,
		PRIMARY KEY (user_ID), INDEX fk_admins_users1_idx (user_ID ASC) VISIBLE,
		UNIQUE INDEX users_ID_UNIQUE (user_ID ASC) VISIBLE,
		CONSTRAINT fk_admins_users1 FOREIGN KEY (user_ID) REFERENCES users (user_id)
		ON DELETE NO ACTION
		ON UPDATE NO ACTION
        ) ENGINE = InnoDB;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS subjects(
		subject_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
		alias VARCHAR(255) NOT NULL,
		name VARCHAR(255) NULL,
		schedule VARCHAR(255) NULL DEFAULT 'WEEKLY',
		UNIQUE INDEX ALIAS_UNIQUE (alias ASC) VISIBLE,
		PRIMARY KEY (subject_id),
		UNIQUE INDEX ID_UNIQUE (subject_id ASC) VISIBLE
        ) ENGINE = InnoDB;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS queues_list(
    	queue_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    	subject_id BIGINT UNSIGNED NOT NULL,
    	name VARCHAR(255) NOT NULL,
    	date DATE NOT NULL,
    	PRIMARY KEY (queue_id, subject_id),
    	UNIQUE INDEX ID_UNIQUE (queue_id ASC) VISIBLE,
    	INDEX fk_queues_subjects1_idx (subject_id ASC) VISIBLE,
    	UNIQUE INDEX name_UNIQUE (name ASC) VISIBLE,
    	CONSTRAINT fk_queues_subjects1
        	FOREIGN KEY (subject_id) REFERENCES subjects (subject_id)
           		ON DELETE NO ACTION
            	ON UPDATE NO ACTION
    	) ENGINE = InnoDB;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS queue(
		subject_id BIGINT UNSIGNED NOT NULL,
		queue_id BIGINT UNSIGNED NOT NULL,
		user_id BIGINT UNSIGNED NOT NULL,
		position INT UNSIGNED NOT NULL,
		time DATETIME NULL,
		PRIMARY KEY (subject_id, queue_id, user_id),
		INDEX fk_queue_queues_list1_idx (queue_id ASC, subject_id ASC) VISIBLE,
		CONSTRAINT fk_queue_users1 FOREIGN KEY (user_id)
			REFERENCES users (user_id)
		    ON DELETE NO ACTION
		    ON UPDATE NO ACTION,
		CONSTRAINT fk_queue_queues_list1 FOREIGN KEY (queue_id , subject_id)
		    REFERENCES queues_list (queue_id , subject_id)
		    ON DELETE NO ACTION
		    ON UPDATE NO ACTION,
		CONSTRAINT fk_queue_subjects1 FOREIGN KEY (subject_id)
		    REFERENCES subjects (subject_id)
		    ON DELETE NO ACTION
		    ON UPDATE NO ACTION
		) ENGINE = InnoDB;`)
	if err != nil {
		return err
	}

	return nil
}
