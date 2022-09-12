package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/richard-on/QueueBot/internal/bot/model"
	"strings"
)

type QueueDB struct {
	db *sql.DB
}

func NewQueueDB(db *sql.DB) QueueDB {
	return QueueDB{db}
}

func (q *QueueDB) GetUser(id int64) (*model.User, error) {
	res, err := q.db.Query(`SELECT * FROM users where tg_user_id = ?;`, id)

	var user model.User
	for res.Next() {
		err = res.Scan(
			&user.UserID,
			&user.TgUsername,
			&user.GroupID,
			&user.SubgroupID,
			&user.TgFirstName,
			&user.TgLastName,
			&user.FirstName,
			&user.LastName)
		if err != nil {
			return &model.User{}, err
		}
	}

	if user.GroupID == 0 || !user.SubgroupID.Valid {
		return &user, errors.New("user not initialised")
	}

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
			VALUES(?, ?, ?, ?, 0);`, user.UserID, user.TgUsername, user.FirstName, user.LastName); err != nil {
		return err
	}

	return nil
}

func (q *QueueDB) GetSubjectList(group int64) (list []Subject, err error) {
	res, err := q.db.Query(
		`SELECT * FROM subjects WHERE group_id = ? OR is_subgroup_subject = TRUE AND subgroup_id = ?`,
		group, group)
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
		return nil, errors.New("did not find subjects")
	}

	return list, nil
}

func (q *QueueDB) GetQueueList(subject Subject) (list []Queue, err error) {
	res, err := q.db.Query(`SELECT queue_id, s.subject_id, name FROM queues_list
    	JOIN subjects s ON queues_list.subject_id = ?`, subject.ID)
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
		return nil, errors.New("did not find queues")
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

		if position.Valid == false {
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

/*func GetUserData(id int64, tgUsername string) (internal.User, error) {
	db, err := sql.Open("mysql", config.DbInfo)
	if err != nil {
		return internal.User{}, err
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM users where tg_user_id = ? OR tg_username = ?;", id, tgUsername)
	if err != nil {
		return internal.User{}, err
	}

	var user internal.User
	for res.Next() {
		err = res.Scan(
			&user.UserID,
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

	return internal.User{}, err
}

func GetGroup(user internal.User) (string, error) {
	db, err := sql.Open("mysql", config.DbInfo)
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

func GetSubGroup(user internal.User) (string, error) {
	db, err := sql.Open("mysql", config.DbInfo)
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

func AddUser(user *internal.User) error {
	db, err := sql.Open("mysql", config.DbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err = db.Exec(`INSERT INTO users(tg_user_id, tg_username, tg_first_name, tg_last_name, group_id)
			VALUES(?, ?, ?, ?, 0);`, user.UserID, user.TgUsername, user.FirstName, user.LastName); err != nil {
		return err
	}

	return nil
}

func CheckUserData(tgUser *tgbotapi.User) (internal.User, error) {
	db, err := sql.Open("mysql", config.DbInfo)
	if err != nil {
		return internal.User{}, err
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM users where tg_user_id = ?;", tgUser.ID)
	if err != nil {
		return internal.User{}, err
	}

	var user internal.User
	if !res.Next() {
		return internal.User{}, errors.New("no such user")
		if _, err = db.Exec(`INSERT INTO users(tg_user_id, tg_username, tg_first_name, tg_last_name, group_id)
			VALUES(?, ?, ?, ?, 0);`, user.ID, user.UserName, user.FirstName, user.LastName); err != nil {
			return err
		}
	} else {
		err = res.Scan(
			&user.UserID,
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

		if user.GroupID == 0 || !user.SubgroupID.Valid {
			return internal.User{}, errors.New("unreg")
		}
	}

	return user, nil
}

func GetSubjects(user internal.User) ([]Subjects, error) {
	db, err := sql.Open("mysql", config.DbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var s []Subjects
	var row Subjects
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

func GetQueues(subjectName string) ([]QueueInfo, error) {
	db, err := sql.Open("mysql", config.DbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var queueSlice QueueInfo
	var queueInfo []QueueInfo
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
	db, err := sql.Open("mysql", config.DbInfo)
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
	db, err := sql.Open("mysql", config.DbInfo)
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

		_, err = db.Exec(`UPDATE queue SET position = ? WHERE queue_id = ? AND position > ?;`,
			queuePos, queueId, queuePos)
		if err != nil {
			return err
		}

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
	db, err := sql.Open("mysql", config.DbInfo)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var queuePrint QueuePrint
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
}*/
