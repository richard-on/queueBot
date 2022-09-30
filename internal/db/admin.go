package db

/*import (
	"database/sql"
	"errors"
	"github.com/richard-on/queueBot/internal"
)

/*func CheckAdmin(ID int64, username string) bool {
	db, err := sql.Open("mysql", internal.DbInfo)
	if err != nil {
		return false
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM admins WHERE user_id = ?;", ID)
	if err != nil {
		return false
	}
	if res.Next() {
		/*if _, err = db.Exec("INSERT INTO admins(is_logged) VALUES(1);"); err != nil {
			return false
		}

		return true
	}

	return false
}

func AddSubject(alias string, name string, schedule string) error {
	db, err := sql.Open("mysql", internal.DbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM subjects WHERE alias = ? OR name = ?;", alias, name)
	if err != nil {
		return err
	}
	if !res.Next() {
		if _, err = db.Exec("INSERT INTO subjects(subject_id, alias, name, schedule) VALUES(?, ?, ?, 'WEEKLY');", nil, alias, name); err != nil {
			return err
		}
	} else {
		return errors.New("this subject already exists")
	}

	return nil
}

func RmSubject(alias string, name string) error {
	db, err := sql.Open("mysql", internal.DbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM subjects WHERE alias = ? OR name = ?;", alias, name)
	if err != nil {
		return err
	}

	return nil
}

func AddQueue(subjectAlias string, queueName string, queueDate string) error {
	db, err := sql.Open("mysql", internal.DbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	var subjectId int64
	res, err := db.Query("SELECT subject_id FROM subjects WHERE alias = ?;", subjectAlias)
	if err != nil {
		return err
	}
	for res.Next() {
		if err = res.Scan(&subjectId); err != nil {
			return errors.New("this subject does not exist")
		}
	}

	res, err = db.Query("SELECT * FROM queues_list WHERE subject_id = ? AND name = ?;", subjectId, queueName)
	if err != nil {
		return err
	}
	if res.Next() {
		return errors.New("this queue already exists")
	} else {
		_, err = db.Exec("INSERT INTO queues_list(queue_id, subject_id, name, date) VALUES (?, ?, ?, NOW());", nil, subjectId, queueName)
		if err != nil {
			return err
		}
	}

	return nil
}

func RmQueue(subjectAlias string, queueName string, queueDate string) error {
	db, err := sql.Open("mysql", internal.DbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DROP TABLE " + internal.DbName + ".queue_" + subjectAlias + "_" + queueName)
	if err != nil {
		return err
	}

	return nil
}*/
