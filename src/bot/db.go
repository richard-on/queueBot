package bot

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

func CollectUserData(id int64, username string, firstName string, lastName string) error {
	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM users where user_id = ?;", id)
	if !res.Next() {
		if _, err = db.Exec("INSERT INTO users(user_id, username, first_name, last_name) VALUES(?, ?, ?, ?);", id, username, firstName, lastName); err != nil {
			return err
		}
	}

	return nil
}

func GetSubjects() ([]Subjects, error) {
	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var s []Subjects
	var row Subjects
	res, err := db.Query("SELECT * FROM subjects")
	for i := 0; res.Next(); i++ {
		err = res.Scan(&row.Id, &row.Alias, &row.Name, &row.Schedule)
		if err != nil {
			return nil, err
		}

		s = append(s, row)
	}
	if s == nil {
		return nil, errors.New("did not find subjects")
	}

	return s, nil
}

func GetQueues(subjectName string) ([]QueueInfo, error) {
	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var queueSlice QueueInfo
	var queueInfo []QueueInfo
	res, err := db.Query("SELECT s.subject_id, queue_id, queues_list.name FROM queues_list JOIN subjects s WHERE s.alias = ? OR s.name = ?", subjectName, subjectName)
	i := 0
	for res.Next() {
		err = res.Scan(&queueSlice.SubjectId, &queueSlice.QueueId, &queueSlice.Name)
		if err != nil {
			panic(err)
		}
		i++
		queueInfo = append(queueInfo, queueSlice)
	}
	if queueInfo == nil {
		return nil, errors.New("did not find queues")
	}

	return queueInfo, nil
}

func JoinQueue(subjectId int64, queueId int64, userId int64) error {
	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM queue WHERE queue_id = ? AND user_id = ?;", queueId, userId)
	var position int
	if !res.Next() {
		res, err = db.Query("SELECT MAX(position) FROM queue WHERE queue_id = ?;", queueId)
		if res.Next() {
			err := res.Scan(&position)
			if err != nil {
				return err
			}
		}
		_, err = db.Exec("INSERT INTO queue(subject_id, queue_id, user_id, position, time) VALUES (?, ?, ?, ?, NOW())", subjectId, queueId, userId, position+1)
		if err != nil {
			return err
		}
	} else {
		return errors.New("already in queue")
	}

	return nil
}

func LeaveQueue(subjectId int64, queueId int64, userId int64) error {
	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM queue WHERE queue_id = ? AND user_id = ?;", queueId, userId)
	if res.Next() {
		_, err = db.Exec("DELETE FROM queue WHERE subject_id = ? AND queue_id = ? AND user_id = ?", subjectId, queueId, userId)
		if err != nil {
			return err
		}
	} else {
		return errors.New("not in queue")
	}

	return nil
}

func PrintQueue(queueId int64, userId int64) (string, error) {
	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var queuePrint QueuePrint
	var sb strings.Builder
	sb.WriteString("Current queue is:\n")
	flag := false
	res, err := db.Query("SELECT u.username, u.first_name, u.last_name, position FROM queue JOIN users u ON u.user_id = queue.user_id WHERE queue_id = ? ORDER BY position;", queueId)
	for res.Next() {
		err = res.Scan(&queuePrint.Username, &queuePrint.FirstName, &queuePrint.LastName, &queuePrint.Position)
		if err != nil {
			return "", err
		}
		str := fmt.Sprintf("%d. %s %s (@%s)\n", queuePrint.Position, queuePrint.FirstName, queuePrint.LastName, queuePrint.Username)
		sb.WriteString(str)
		flag = true
	}
	if !flag {
		return "Queue is empty", nil
	}

	return sb.String(), nil
}
