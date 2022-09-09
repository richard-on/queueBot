package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/richard-on/QueueBot/cmd/queueBot/initEnv"
	"github.com/richard-on/QueueBot/pkg/queueBot"
	"strings"
)

func GetUserData(id int64, tgUsername string) (queueBot.User, error) {
	db, err := sql.Open("mysql", initEnv.DbInfo)
	if err != nil {
		return queueBot.User{}, err
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM users where tg_user_id = ? OR tg_username = ?;", id, tgUsername)
	if err != nil {
		return queueBot.User{}, err
	}

	var user queueBot.User
	for res.Next() {
		err = res.Scan(
			&user.ID,
			&user.TgUsername,
			&user.GroupID,
			&user.SubgroupID,
			&user.TgFirstName,
			&user.TgLastName,
			&user.FirstName,
			&user.LastName)
		if err != nil {
			panic(err)
		}

		return user, nil
	}

	return queueBot.User{}, err
}

func GetGroup(user queueBot.User) (string, error) {
	db, err := sql.Open("mysql", initEnv.DbInfo)
	if err != nil {
		return "", err
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM `groups` where group_id = ?;", user.GroupID)
	if err != nil {
		return "", err
	}

	var groupID string
	var groupName string
	for res.Next() {
		err = res.Scan(&groupID, &groupName)
		if err != nil {
			panic(err)
		}

		return groupName, nil
	}

	return "", err
}

func GetSubGroup(user queueBot.User) (string, error) {
	db, err := sql.Open("mysql", initEnv.DbInfo)
	if err != nil {
		return "", err
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM subgroups where subgroup_id = ?;", user.SubgroupID)
	if err != nil {
		return "", err
	}

	var subgroupID string
	var subgroupName string
	for res.Next() {
		err = res.Scan(&subgroupID, &subgroupName)
		if err != nil {
			panic(err)
		}

		return subgroupName, nil
	}

	return "", err
}

func CollectUserData(id int64, tgUsername string, tgFirstName string, tgLastName string) error {
	db, err := sql.Open("mysql", initEnv.DbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM users where tg_user_id = ?;", id)
	if err != nil {
		return err
	}
	if !res.Next() {
		if _, err = db.Exec("INSERT INTO users(tg_user_id, tg_username, tg_first_name, tg_last_name, group_id) VALUES(?, ?, ?, ?, 0);",
			id, tgUsername, tgFirstName, tgLastName); err != nil {
			return err
		}
	}

	return nil
}

func GetSubjects(user queueBot.User) ([]queueBot.Subjects, error) {
	db, err := sql.Open("mysql", initEnv.DbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var s []queueBot.Subjects
	var row queueBot.Subjects
	res, err := db.Query(
		"SELECT * FROM subjects WHERE group_id = ? OR is_subgroup_subject = TRUE AND subgroup_id = ?",
		user.GroupID, user.SubgroupID)
	if err != nil {
		return nil, err
	}
	for i := 0; res.Next(); i++ {
		err = res.Scan(&row.ID, &row.SubjectName, &row.IsSubgroupSubject, &row.GroupID, &row.IsSubgroupSubject)
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

func GetQueues(subjectName string) ([]queueBot.QueueInfo, error) {
	db, err := sql.Open("mysql", initEnv.DbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var queueSlice queueBot.QueueInfo
	var queueInfo []queueBot.QueueInfo
	res, err := db.Query(`SELECT queue_id, s.subject_id, name FROM queues_list JOIN subjects s ON
    	s.subject_id = queues_list.subject_id WHERE subject_name = ?`, subjectName)
	if err != nil {
		return nil, err
	}
	i := 0
	for res.Next() {
		err = res.Scan(&queueSlice.QueueId, &queueSlice.SubjectId, &queueSlice.Name)
		if err != nil {
			panic(err)
		}
		i++
		queueInfo = append(queueInfo, queueSlice)
	}
	if queueInfo == nil {
		return nil, errors.New("did not find queueSlice")
	}

	return queueInfo, nil
}

func JoinQueue(subjectId int64, queueId int64, userId int64) error {
	db, err := sql.Open("mysql", initEnv.DbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM queue WHERE queue_id = ? AND user_id = ?;", queueId, userId)
	if err != nil {
		return err
	}
	var position sql.NullInt64
	if !res.Next() {
		res, err = db.Query("SELECT MAX(position) FROM queue WHERE queue_id = ?;", queueId)
		if err != nil {
			return err
		}
		if res.Next() {
			err := res.Scan(&position)
			if err != nil {
				return err
			}
		}

		if position.Valid == false {
			_, err = db.Exec("INSERT INTO queue(subject_id, queue_id, user_id, position) VALUES (?, ?, ?, ?)",
				subjectId, queueId, userId, 1)
		} else {
			_, err = db.Exec("INSERT INTO queue(subject_id, queue_id, user_id, position) VALUES (?, ?, ?, ?)",
				subjectId, queueId, userId, position.Int64+1)
		}

		if err != nil {
			return err
		}
	} else {
		return errors.New("already in queue")
	}

	return nil
}

func LeaveQueue(subjectId int64, queueId int64, userId int64) error {
	db, err := sql.Open("mysql", initEnv.DbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Query("SELECT position FROM queue WHERE queue_id = ? AND user_id = ?;", queueId, userId)
	if err != nil {
		return err
	}
	if res.Next() {
		var queuePos int64
		err = res.Scan(&queuePos)
		nextPos, err := db.Query(`SELECT position FROM queue WHERE queue_id = ? AND position = ?;`, queueId, queuePos-1)

		_, err = db.Exec("DELETE FROM queue WHERE subject_id = ? AND queue_id = ? AND user_id = ?", subjectId, queueId, userId)
		if err != nil {
			return err
		}
		_, err = db.Exec(``)
	} else {
		return errors.New("not in queue")
	}

	return nil
}

func PrintQueue(queueId int64, userId int64) (string, error) {
	db, err := sql.Open("mysql", initEnv.DbInfo)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var queuePrint queueBot.QueuePrint
	var sb strings.Builder
	sb.WriteString("Current queue is:\n")
	flag := false
	res, err := db.Query(`SELECT u.tg_username, u.first_name, u.last_name, position FROM queue JOIN users u ON
        u.tg_user_id = queue.user_id WHERE queue_id = ? ORDER BY position;`, queueId)
	if err != nil {
		return "", err
	}
	for res.Next() {
		err = res.Scan(&queuePrint.Username, &queuePrint.FirstName, &queuePrint.LastName, &queuePrint.Position)
		if err != nil {
			return "", err
		}
		str := fmt.Sprintf("%d. %s %s (@%s)\n", queuePrint.Position.Int64, queuePrint.FirstName, queuePrint.LastName, queuePrint.Username)
		sb.WriteString(str)
		flag = true
	}
	if !flag {
		return "Queue is empty", nil
	}

	return sb.String(), nil
}
