package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

type Storage struct {
	db *sql.DB
}

type Comix struct {
	Description string
	UploadDate  string
	Views       int
}

type ComixFromAllComix struct {
	ComixName   string
	ComixTag    string
	Description string
	ComixDate   string
	Views       int
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.postgres.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return &Storage{db: db}, nil
}

/*
*
  - Функция добавляет комикс в таблицу с тэгом
    @param
  - tagName - название тэга
  - name - название комикса
  - description - описание комикса
  - currentDate - дата добавления
    @return
  - err - ошибка
    *
*/
func (s *Storage) AddComixByTagName(tagName string, name string, description string, currentDate string) error {
	const fn = "storage.postgres.addComixByTagName"

	query := fmt.Sprintf("INSERT into %s (name,description,upload_date,views) VALUES('%s','%s','%s',1);", tagName, name, description, currentDate)

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

/*
*
  - Функция проверяет существование таблицы по тэгу
    @param
  - tagName - название тэга
  - name - название комикса
    @return
  - err - ошибка
    *
*/
func (s *Storage) CheckComixExists(tagName string, name string) (bool, error) {
	const fn = "storage.postgres.CheckComixExists"

	query := fmt.Sprintf("SELECT EXISTS(SELECT * FROM %s WHERE name='%s')", tagName, name)

	var rowExist bool

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	err = stmt.QueryRow().Scan(&rowExist)
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return rowExist, nil
}

/*
*
  - Функция возвращает комикс, находя его по его названию
    @param
  - tagName - название тэга
  - name - название комикса
    @return
    -Comix - возвращает структуру комикса
  - err - ошибка
    *
*/
func (s *Storage) GetComixByName(tagName string, name string) (Comix, error) {
	const fn = "storage.postgres.GetComixByName"

	query := fmt.Sprintf("SELECT description,upload_date,views FROM %s WHERE name='%s'", tagName, name)

	comix := Comix{}

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return Comix{}, fmt.Errorf("%s: %w", fn, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return Comix{}, fmt.Errorf("%s: %w", fn, err)
	}

	err = stmt.QueryRow().Scan(&comix.Description, &comix.UploadDate, &comix.Views)
	if err != nil {
		return Comix{}, fmt.Errorf("%s: %w", fn, err)
	}

	return comix, nil
}

/*
*
  - Добавляет комикс в таблицу всех комиксов
    @param
  - tagName - название тэга
  - name - название комикса
  - description - описание комикса
  - currentDate - дата добавления
    @return
  - err - ошибка
    *
*/
func (s *Storage) AddComixToAllComixTable(tagName string, name string, description string, currentDate string) error {
	const fn = "storage.postgres.AddComixToAllComixTable"

	query := fmt.Sprintf("INSERT INTO all_comix (comix_name, comix_tag, description, comix_date, views)  VALUES('%s','%s','%s','%s',1);", name, strings.ToLower(tagName), description, currentDate)

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

/*
*
  - Добавляет новый тэг в таблицу тэгов
    @param
  - tagName - название тэга
    @return
  - err - ошибка
    *
*/
func (s *Storage) AddComixTagToAllTags(tagName string) error {
	const fn = "storage.postgres.AddComixTagToAllTags"

	query := fmt.Sprintf("INSERT INTO all_tags VALUES('%s');", tagName)

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

/*
*
  - Возвращает количество существующих комиксов
    @return
  - err - ошибка
  - int - количество
    *
*/
func (s *Storage) GetComixQuantity() (int, error) {
	const fn = "storage.postgres.GetComixQuantity"

	var quantity int

	query := fmt.Sprintf("SELECT COUNT(*) FROM all_comix")

	err := s.db.QueryRow(query).Scan(&quantity)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return quantity, nil
}

/*
*
  - Возвращает количество существующих тэгов
    @param
  - tagName - название тэга
    @return
  - err - ошибка
  - int - количество
    *
*/
func (s *Storage) GetComixQuantityFromTag(tagName string) (int, error) {
	const fn = "storage.postgres.GetComixQuantityFromTag"

	var quantity int

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tagName)

	err := s.db.QueryRow(query).Scan(&quantity)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return quantity, nil
}

/*
*
  - Возвращает комиксы в названии которых есть входная переменная "name"
    @param
  - name - название тэга
    @return
  - err - ошибка
  - int - количество
    *
*/
func (s *Storage) GetComixQuantityFromName(name string) (int, error) {
	const fn = "storage.postgres.GetComixQuantityFromTag"

	var quantity int

	query := fmt.Sprintf("SELECT COUNT(*) FROM all_comix WHERE comix_name LIKE '%%%s%%';", name)

	err := s.db.QueryRow(query).Scan(&quantity)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return quantity, nil
}

/*
*
  - Возвращает описание тэга
    @param
  - tagName - название тэга
    @return
  - err - ошибка
  - string - описание
    *
*/
func (s *Storage) GetTagDescription(tagName string) (string, error) {
	const fn = "storage.postgres.GetTagDescription"

	var description string

	query := fmt.Sprintf("SELECT description FROM tags_description WHERE tag='%s'", tagName)

	err := s.db.QueryRow(query).Scan(&description)
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return description, nil
}

/*
*
  - Возвращает 16 комиксов из таблицы всех комиксов по параметру "pageToDisplay"
    @param
  - pageToDisplay - номер страницы для отображения
    @return
  - err - ошибка
  - []ComixFromAllComix - 16 комиксов, каждый в виде структуры ComixFromAllComix
    *
*/
func (s *Storage) GetComixForMainPage(pageToDisplay int) ([]ComixFromAllComix, error) {
	const fn = "storage.postgres.GetComixForMainPage"
	const numberComicsPerPage = 16

	foo := 0

	var comixList []ComixFromAllComix

	offset := (pageToDisplay - 1) * numberComicsPerPage

	query := fmt.Sprintf("SELECT * FROM all_comix ORDER BY id DESC LIMIT %d OFFSET %d", numberComicsPerPage, offset)

	rows, err := s.db.Query(query)
	if err != nil {
		return []ComixFromAllComix{}, fmt.Errorf("%s: %w", fn, err)
	}

	for rows.Next() {
		var comix ComixFromAllComix
		err := rows.Scan(&foo, &comix.ComixName, &comix.ComixTag, &comix.Description, &comix.ComixDate, &comix.Views)
		if err != nil {
			return []ComixFromAllComix{}, fmt.Errorf("%s: %w", fn, err)
		}
		comixList = append(comixList, comix)
	}

	return comixList, nil
}

/*
*
  - Возвращает 16 комиксов из таблицы по тэгу
    @param
  - pageToDisplay - номер страницы для отображения
    -tagName - название тэга
    @return
  - err - ошибка
  - []ComixFromAllComix - 16 комиксов, каждый в виде структуры ComixFromAllComix
    *
*/
func (s *Storage) GetAllTagComix(pageToDisplay int, tagName string) ([]ComixFromAllComix, error) {
	const fn = "storage.postgres.GetAllTagComix"
	const numberComicsPerPage = 16

	foo := 0

	var comixList []ComixFromAllComix

	offset := (pageToDisplay - 1) * numberComicsPerPage

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id DESC LIMIT %d OFFSET %d", tagName, numberComicsPerPage, offset)

	rows, err := s.db.Query(query)
	if err != nil {
		return []ComixFromAllComix{}, fmt.Errorf("%s: %w", fn, err)
	}

	for rows.Next() {
		var comix ComixFromAllComix
		err := rows.Scan(&foo, &comix.ComixName, &comix.Description, &comix.ComixDate, &comix.Views)
		if err != nil {
			return []ComixFromAllComix{}, fmt.Errorf("%s: %w", fn, err)
		}
		comix.ComixTag = tagName
		comixList = append(comixList, comix)
	}

	return comixList, nil
}

/*
*
  - Добавляет просмотры к комиксу
    @param
  - tag - название тэга
    -name - название комикса
    @return
  - err - ошибка
    *
*/
func (s *Storage) AddViews(tag string, name string) error {
	const fn = "storage.postgres.addViews"

	query := fmt.Sprintf("UPDATE %s SET views = views + 1 WHERE name='%s';", tag, name)

	_, err := s.db.Query(query)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

/*
*
  - Добавляет просмотры к комиксу
    @param
  - tag - название тэга
    -description - описание комикса
    @return
  - err - ошибка
    *
*/
func (s *Storage) AddTagDescription(tag string, description string) error {
	const fn = "storage.postgres.AddTagDescription"

	query := fmt.Sprintf("INSERT INTO tags_description values('%s','%s')", tag, description)

	_, err := s.db.Query(query)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

/*
*
  - Возвращает список всех тэгов
    @return
  - err - ошибка
    -[]string - список тэгов
    *
*/
func (s *Storage) GetAllTags() ([]string, error) {
	const fn = "storage.postgres.GetAllTags"

	var TagsList []string

	query := fmt.Sprintf("SELECT * FROM all_tags")

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		TagsList = append(TagsList, tag)
	}

	return TagsList, nil
}

/*
*
  - Удаляет комикс из таблицы по тэгу
    @param
  - tag - название тэга
    -name - название комикса
    @return
  - err - ошибка
    *
*/
func (s *Storage) DeleteComixFromTagTable(tag string, name string) error {
	const fn = "storage.postgres.DeleteComixFromTagTable"

	query := fmt.Sprintf("DELETE FROM %s WHERE name = '%s'", tag, name)

	_, err := s.db.Query(query)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

/*
*
  - Удаляет комикс из таблицы всех комиксов по имени
    @param
    -name - название комикса
    @return
  - err - ошибка
    *
*/
func (s *Storage) DeleteComixFromAllComixTable(name string) error {
	const fn = "storage.postgres.DeleteComixFromAllComixTable"

	query := fmt.Sprintf("DELETE FROM all_comix WHERE comix_name = '%s'", name)

	_, err := s.db.Query(query)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

/*
*
  - Изменяет параметры комикса из таблицы всех комиксов: название, описание
    @param
  - param - параметр
    -name - название комикса
    -newValue - новое значение параметра комикса
    @return
  - err - ошибка
    *
*/
func (s *Storage) EditComixFromAllComixTable(name string, param string, newValue string) error {
	const fn = "storage.postgres.EditComixFromAllComixTable"

	query := fmt.Sprintf("UPDATE all_comix SET %s = '%s' WHERE comix_name = '%s'", param, newValue, name)

	_, err := s.db.Query(query)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

/*
*
  - Изменяет тэг комикса
    @param
  - tag - название тэга
    -name - название комикса
    -newValue - новое значение тэга комикса
    @return
  - err - ошибка
    *
*/
func (s *Storage) EditComixTag(tag string, name string, newValue string) error {
	const fn = "storage.postgres.EditComixTag"

	comix := Comix{}

	query := fmt.Sprintf("SELECT description,upload_date,views FROM %s WHERE name = '%s'", tag, name)

	err := s.db.QueryRow(query).Scan(&comix.Description, &comix.UploadDate, &comix.Views)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE name = '%s'", tag, name)

	_, err = s.db.Query(query)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	query = fmt.Sprintf("INSERT into %s (name,description,upload_date,views) VALUES('%s','%s','%s',1);", newValue, name, comix.Description, comix.UploadDate)

	_, err = s.db.Query(query)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

/*
*
  - Изменяет параметры комикса из таблицы тэга: название, описание
    @param
  - param - параметр
  - tag - сам тэг
    -name - название комикса
    -newValue - новое значение параметра комикса
    @return
  - err - ошибка
    *
*/
func (s *Storage) EditComixFromTagTable(tag string, name string, param string, newValue string) error {
	const fn = "storage.postgres.EditComixFromAllComixTable"
	if param == "tag" {
		comix := Comix{}

		query := fmt.Sprintf("SELECT (description,upload_date,views)FROM %s WHERE comix_name = '%s'", tag, name)

		err := s.db.QueryRow(query).Scan(&comix.Description, &comix.UploadDate, &comix.Views)
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}

		query = fmt.Sprintf("DELETE FROM %s WHERE name = '%s'", tag, name)

		_, err = s.db.Query(query)
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}

		query = fmt.Sprintf("INSERT into %s (name,description,upload_date,views) VALUES('%s','%s','%s',1);", tag, name, comix.Description, comix.UploadDate)

		_, err = s.db.Query(query)
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}

		return nil
	}
	query := fmt.Sprintf("UPDATE %s SET %s = '%s' WHERE name = '%s'", tag, param, newValue, name)

	_, err := s.db.Query(query)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

/*
*
  - Функция создаёт таблицу тэга
  - Она содержит в себе столбцы: id комикса, название, описани, дату загрузки, кол-во просмотров
    @param
  - tagName - название тэга
    @return
  - err - ошибка
    *
*/
func (s *Storage) CreateNewTag(tagName string) error {
	const fn = "storage.postgres.CreateNewTag"

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id SERIAL PRIMARY KEY,name TEXT NOT NULL,description TEXT NOT NULL,upload_date DATE,views INTEGER NOT NULL);", tagName)

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

/*
*
  - Проверяет существование таблицы тэна
    @param
  - tagName - название тэга
    @return
  - err - ошибка
  - bool - существут или нет
    *
*/
func (s *Storage) TagExist(tagName string) (bool, error) {

	const fn = "storage.postgres.TagExist"

	var tagExists bool

	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = '%s')", strings.ToLower(tagName))

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	err = stmt.QueryRow().Scan(&tagExists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return tagExists, nil
}

/*
*
  - Находит комикс в таблице всех комиксов, по имени
    @param
  - name - название комикса
  - pageToDisplay - страница для отображения
    @return
  - err - ошибка
  - []ComixFromAllComix - список комиксов
    *
*/
func (s *Storage) FindComixFromAllComix(name string, pageToDisplay int) ([]ComixFromAllComix, error) {

	const fn = "storage.postgres.FindComixFromAllComix"
	const numberComicsPerPage = 16

	foo := 0

	var comixList []ComixFromAllComix

	offset := (pageToDisplay - 1) * numberComicsPerPage

	query := fmt.Sprintf("SELECT * FROM all_comix WHERE comix_name LIKE '%%%s%%' ORDER BY id DESC LIMIT %d OFFSET %d;", name, numberComicsPerPage, offset)

	rows, err := s.db.Query(query)
	if err != nil {
		return []ComixFromAllComix{}, fmt.Errorf("%s: %w", fn, err)
	}

	for rows.Next() {
		var comix ComixFromAllComix
		err := rows.Scan(&foo, &comix.ComixName, &comix.ComixTag, &comix.Description, &comix.ComixDate, &comix.Views)
		if err != nil {
			return []ComixFromAllComix{}, fmt.Errorf("%s: %w", fn, err)
		}
		comixList = append(comixList, comix)
	}

	return comixList, nil
}

/*
*
  - Проверяет правильность пароля
    @param
  - inputPass - пароль вводимый пользователем
    @return
  - err - ошибка
  - bool - правльный или нет
    *
*/
func (s *Storage) CheckPass(inputPass string) (bool, error) {

	const fn = "storage.postgres.TagExist"

	var passIsCorrect bool
	var pass string

	query := fmt.Sprintf("SELECT pass FROM pass")

	err := s.db.QueryRow(query).Scan(&pass)
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)

	}

	if pass == inputPass {
		passIsCorrect = true

	} else {
		passIsCorrect = false

	}

	return passIsCorrect, nil
}
