package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/richard-on/QueueBot/config"
)

func InitDb() error {
	db, err := sql.Open("mysql", config.InitDbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE SCHEMA IF NOT EXISTS ` + config.DbName)
	if err != nil {
		return err
	}

	_, err = db.Exec(`USE ` + config.DbName)
	if err != nil {
		return err
	}

	return nil
}

func CreateTables() error {
	db, err := sql.Open("mysql", config.DbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS queue_db.groups (
  		group_id BIGINT UNSIGNED NOT NULL,
  		group_name VARCHAR(45) NOT NULL,
  		PRIMARY KEY (group_id))
		ENGINE = InnoDB;`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS queue_db.subgroups (
  		subgroup_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  		subgroup_name VARCHAR(45) NOT NULL,
  		PRIMARY KEY (subgroup_id))
		ENGINE = InnoDB;`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS queue_db.subjects (
  	subject_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  	subject_name VARCHAR(255) NOT NULL,
  	is_subgroup_subject TINYINT UNSIGNED ZEROFILL NOT NULL,
  	group_id BIGINT UNSIGNED NOT NULL,
  	subgroup_id BIGINT UNSIGNED NULL,
  	PRIMARY KEY (subject_id),
  	CONSTRAINT fk_subjects_groups1
  		FOREIGN KEY (group_id)
  		REFERENCES queue_db.groups (group_id)
		ON DELETE NO ACTION
    	ON UPDATE NO ACTION,
  	CONSTRAINT fk_subjects_subgroups1
    	FOREIGN KEY (subgroup_id)
    	REFERENCES queue_db.subgroups (subgroup_id)
		ON DELETE NO ACTION
    	ON UPDATE NO ACTION)
	ENGINE = InnoDB;`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS queue_db.users (
  	tg_user_id BIGINT UNSIGNED NOT NULL,
  	tg_username VARCHAR(255) NOT NULL,
  	group_id BIGINT UNSIGNED NOT NULL,
  	subgroup_id BIGINT UNSIGNED NULL,
  	tg_first_name VARCHAR(255) NULL,
  	tg_last_name VARCHAR(255) NULL,
  	first_name VARCHAR(255) NULL,
  	last_name VARCHAR(255) NULL,
  	PRIMARY KEY (tg_user_id, group_id),
  	CONSTRAINT fk_users_groups1
    	FOREIGN KEY (group_id)
    	REFERENCES queue_db.groups (group_id)
    	ON DELETE NO ACTION
    	ON UPDATE NO ACTION,
  	CONSTRAINT fk_users_subgroups1
    	FOREIGN KEY (subgroup_id)
    	REFERENCES queue_db.subgroups (subgroup_id)
    	ON DELETE NO ACTION
    	ON UPDATE NO ACTION)
	ENGINE = InnoDB;`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS queue_db.admins (
  	user_ID BIGINT UNSIGNED NOT NULL,
  	PRIMARY KEY (user_ID),
  	CONSTRAINT fk_admins_users1
    	FOREIGN KEY (user_ID)
    	REFERENCES queue_db.users (tg_user_id)
    	ON DELETE NO ACTION
    	ON UPDATE NO ACTION)
	ENGINE = InnoDB;`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS queue_db.admins (
  	user_ID BIGINT UNSIGNED NOT NULL,
  	PRIMARY KEY (user_ID),
  	CONSTRAINT fk_admins_users1
    	FOREIGN KEY (user_ID)
    	REFERENCES queue_db.users (tg_user_id)
    	ON DELETE NO ACTION
    	ON UPDATE NO ACTION)
	ENGINE = InnoDB;`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS queue_db.queues_list (
  	queue_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  	subject_id BIGINT UNSIGNED NOT NULL,
  	name VARCHAR(255) NOT NULL,
  	PRIMARY KEY (queue_id, subject_id),
  	CONSTRAINT fk_queues_subjects1
    	FOREIGN KEY (subject_id)
    	REFERENCES queue_db.subjects (subject_id)
    	ON DELETE NO ACTION
    	ON UPDATE NO ACTION)
	ENGINE = InnoDB;`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS queue_db.queue (
  	subject_id BIGINT UNSIGNED NOT NULL,
  	queue_id BIGINT UNSIGNED NOT NULL,
  	user_id BIGINT UNSIGNED NOT NULL,
  	position INT UNSIGNED NOT NULL,
  	PRIMARY KEY (subject_id, queue_id, user_id),
  	CONSTRAINT fk_queue_users1
    	FOREIGN KEY (user_id)
    	REFERENCES queue_db.users (tg_user_id)
    	ON DELETE NO ACTION
    	ON UPDATE NO ACTION,
  	CONSTRAINT fk_queue_queues_list1
    	FOREIGN KEY (queue_id , subject_id)
    	REFERENCES queue_db.queues_list (queue_id , subject_id)
    	ON DELETE NO ACTION
    	ON UPDATE NO ACTION,
  	CONSTRAINT fk_queue_subjects1
    	FOREIGN KEY (subject_id)
    	REFERENCES queue_db.subjects (subject_id)
    	ON DELETE NO ACTION
    	ON UPDATE NO ACTION)
	ENGINE = InnoDB;`)

	if err != nil {
		return err
	}

	return nil
}
