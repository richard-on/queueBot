package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/richard-on/QueueBot/internal/bot/model"
	"github.com/richard-on/QueueBot/internal/logger"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

type QueueDB struct {
	db  *sql.DB
	log logger.Logger
}

func NewQueueDB(db *sql.DB) QueueDB {
	return QueueDB{
		db,
		logger.NewLogger(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123},
			zerolog.TraceLevel,
			"queueBot-db"),
	}
}

func (q *QueueDB) GetUser(id int64) (*model.User, error) {
	res, err := q.db.Query(`SELECT * FROM users where tg_user_id = ?;`, id)
	if err != nil {
		return nil, err
	}

	var user model.User
	var subgroupID sql.NullInt64
	var tgFirstName sql.NullString
	var tgLastName sql.NullString
	var firstName sql.NullString
	var lastName sql.NullString
	flag := false
	for res.Next() {
		flag = true
		err = res.Scan(
			&user.UserID,
			&user.TgUsername,
			&user.GroupID,
			&subgroupID,
			&tgFirstName,
			&tgLastName,
			&firstName,
			&lastName)
		if err != nil {
			return &model.User{}, err
		}

		if tgFirstName.Valid {
			user.TgFirstName = tgFirstName.String
		}
		if tgLastName.Valid {
			user.TgLastName = tgLastName.String
		}
		if firstName.Valid {
			user.FirstName = firstName.String
		}
		if lastName.Valid {
			user.LastName = lastName.String
		}
		if subgroupID.Valid {
			user.SubgroupID = subgroupID.Int64
		}

	}

	if (user.GroupID == 0 || user.LastName == "" || user.FirstName == "") && flag {
		return &user, ErrLackUserInfo
	}

	if !flag {
		return &model.User{}, ErrNoUserInfo
	}

	q.log.Dedugf("got user @%v, id: %v, isRegistered: %v", user.TgUsername, user.UserID, user.IsRegistered)

	return &user, nil
}

func (q *QueueDB) GetGroupName(user *model.User) (groupName string, err error) {
	row := q.db.QueryRow("SELECT group_name FROM `groups` where group_id = ?;", user.GroupID)

	err = row.Scan(&groupName)
	if err != nil {
		return "", err
	}

	return groupName, nil
}

func (q *QueueDB) GetSubGroupName(user *model.User) (subGroupName string, err error) {
	res, err := q.db.Query(`SELECT subgroup_name FROM subgroups where subgroup_id = ?;`, user.SubgroupID)
	if err != nil {
		return "", err
	}

	if !res.Next() {
		return "Нет кафедры", nil
	} else {
		err = res.Scan(&subGroupName)
		if err != nil {
			return "", err
		}
	}

	return subGroupName, nil
}

func (q *QueueDB) AddUser(user *model.User) error {
	if _, err := q.db.Exec(`INSERT INTO users(tg_user_id, tg_username, tg_first_name, tg_last_name, group_id)
			VALUES(?, ?, ?, ?, 1);`, user.UserID, user.TgUsername, user.FirstName, user.LastName); err != nil {
		return err
	}

	return nil
}

func (q *QueueDB) GetSubjectList(group int64) (list []Subject, err error) {
	res, err := q.db.Query(
		`SELECT * FROM subjects WHERE is_subgroup_subject = FALSE AND group_id = ?`, group)
	if err != nil {
		return nil, err
	}

	var row Subject
	for res.Next() {
		err = res.Scan(&row.ID, &row.SubjectName, &row.IsSubgroupSubject, &row.GroupID, &row.SubGroupID)
		if err != nil {
			return nil, err
		}

		list = append(list, row)
	}
	if list == nil {
		return nil, errors.New("предметы не найдены")
	}

	return list, nil
}

func (q *QueueDB) GetSubgroupSubjectList(subGroup int64) (list []Subject, err error) {
	res, err := q.db.Query(
		`SELECT * FROM subjects WHERE is_subgroup_subject = TRUE AND subgroup_id = ?`, subGroup)
	if err != nil {
		return nil, err
	}

	var row Subject
	for res.Next() {
		err = res.Scan(&row.ID, &row.SubjectName, &row.IsSubgroupSubject, &row.GroupID, &row.SubGroupID)
		if err != nil {
			return nil, err
		}

		list = append(list, row)
	}
	if list == nil {
		return nil, errors.New("предметы не найдены")
	}

	return list, nil
}

func (q *QueueDB) GetQueueList(subject Subject) (list []Queue, err error) {
	res, err := q.db.Query(`SELECT queue_id, s.subject_id, name FROM queues_list
    	JOIN subjects s ON queues_list.subject_id = ? WHERE s.subject_id = queues_list.subject_id`, subject.ID)
	if err != nil {
		return nil, err
	}
	i := 0

	var row Queue
	for res.Next() {
		err = res.Scan(&row.ID, &row.SubjectId, &row.Name)
		if err != nil {
			return nil, err
		}
		i++
		list = append(list, row)
	}
	if list == nil {
		return nil, errors.New("очереди не найдены")
	}

	return list, nil
}

func (q *QueueDB) JoinQueue(user *model.User, queue *Queue) error {
	res, err := q.db.Query("SELECT * FROM queue WHERE queue_id = ? AND user_id = ?;", queue.ID, user.UserID)
	if err != nil {
		return err
	}

	var position sql.NullInt64
	if !res.Next() {
		res, err = q.db.Query("SELECT MAX(position) FROM queue WHERE queue_id = ?;", queue.ID)
		if err != nil {
			return err
		}
		if res.Next() {
			err := res.Scan(&position)
			if err != nil {
				return err
			}
		}

		if !position.Valid {
			_, err = q.db.Exec("INSERT INTO queue(subject_id, queue_id, user_id, position) VALUES (?, ?, ?, ?)",
				queue.SubjectId, queue.ID, user.UserID, 1)
		} else {
			_, err = q.db.Exec("INSERT INTO queue(subject_id, queue_id, user_id, position) VALUES (?, ?, ?, ?)",
				queue.SubjectId, queue.ID, user.UserID, position.Int64+1)
		}

		if err != nil {
			return err
		}
	} else {
		return errors.New("already in queue")
	}

	return nil
}

func (q *QueueDB) LeaveQueue(user *model.User, queue *Queue) error {
	res, err := q.db.Query("SELECT position FROM queue WHERE queue_id = ? AND user_id = ?;", queue.ID, user.UserID)
	if err != nil {
		return err
	}
	if res.Next() {
		var queuePos int64
		err = res.Scan(&queuePos)
		if err != nil {
			return err
		}

		_, err = q.db.Exec(`UPDATE queue SET position = position - 1 WHERE queue_id = ? AND position > ?;`,
			queue.ID, queuePos)
		if err != nil {
			return err
		}

		_, err = q.db.Exec("DELETE FROM queue WHERE subject_id = ? AND queue_id = ? AND user_id = ?",
			queue.SubjectId, queue.ID, user.UserID)
		if err != nil {
			return err
		}
	} else {
		return errors.New("not in queue")
	}

	return nil
}

func (q *QueueDB) PrintQueue(queueId int64) (string, error) {
	var queuePrint QueuePrint
	var sb strings.Builder

	sb.WriteString("Current queue is:\n")
	flag := false

	res, err := q.db.Query(`SELECT u.tg_username, u.first_name, u.last_name, position FROM queue JOIN users u ON
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
