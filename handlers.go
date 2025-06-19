package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type StringInt int

func (st *StringInt) UnmarshalJSON(b []byte) error {
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case int:
		*st = StringInt(v)
	case float64:
		*st = StringInt(int(v))
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return err

		}
		*st = StringInt(i)

	}
	return nil
}

// @Summury Создать новую запись в учебном плане
// @Accept json
// @Produce json
// @Param item body Item "Создаем план"
// @Success 201 {object} Item
// @Router /home/create_item [post]
// HandlerCreateItem обрабатывает POST-запрос для сохранения учебного плана в базе данных
func HandlerCreateItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("Новый запрос: %s %s %s", r.URL, r.Method, r.Proto)
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	if r.Method != http.MethodPost {
		log.Printf("Метод не разрешен: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		msg := fmt.Sprintf("получен метод %s, ожидался %s", r.Method, http.MethodPost)
		w.Write([]byte(msg))
		return
	}

	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = saveItem(item)
	if err != nil {
		log.Printf("Ошибка сохранения элемента: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf("Не удалось добавить запись: %v", err)
		w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandlerUpdateItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("Новый запрос: %s %s %s", r.URL, r.Method, r.Proto)
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	if r.Method != http.MethodPatch {
		log.Printf("Метод не разрешен: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		msg := fmt.Sprintf("получен метод %s, ожидался %s", r.Method, http.MethodPatch)
		w.Write([]byte(msg))
		return
	}

	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = saveItem(item)
	if err != nil {
		log.Printf("Ошибка сохранения элемента: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf("Не удалось обновить запись: %v", err)
		w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// getPlan выполняет SELECT запрос к базе данных и возвращает список планов
func getPlan(w http.ResponseWriter, r *http.Request) {
	// Параметры подключения к базе данных
	db, err := sql.Open("postgres", "user=bmo password=exp dbname=bmo port=5433 sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Подготавливаем SQL-запрос для выбора данных
	rows, err := db.Query("SELECT id, week_number, week_day, time, groups, teachers, disciplines, lesson_option, classroom FROM study_plan.plans order by week_number, week_day, time, groups")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Создаем срез для хранения данных
	var plans []Item

	// Проходим по всем строкам и заполняем срез
	for rows.Next() {
		var plan Item
		err = rows.Scan(&plan.Id, &plan.WeekNumber, &plan.WeekDay, &plan.Time, &plan.Group, &plan.Teacher, &plan.Subject, &plan.Subject_lvl2, &plan.ClassRoom)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		plans = append(plans, plan)
	}

	// Конвертируем данные в JSON и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plans)
}

func saveItem(item Item) error {
	// Параметры подключения к базе данных
	db, err := sql.Open("postgres", "user=bmo password=exp dbname=bmo port=5433 sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	// Подготавливаем SQL-запрос
	stmt, err := db.Prepare("INSERT INTO study_plan.plans (week_number, week_day, time, groups, teachers, disciplines, lesson_option, classroom) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполняем запрос с передачей параметров
	_, err = stmt.Exec(item.WeekNumber, item.WeekDay, item.Time, item.Group, item.Teacher, item.Subject, item.Subject_lvl2, item.ClassRoom)
	if err != nil {
		return err
	}

	log.Printf("Запись данных: %s", msgSaveSuccess)
	log.Printf("Значения: WeekNumber=%d, WeekDay=%s, Time=%s", item.WeekNumber, item.WeekDay, item.Time) // Добавляем логирование

	return nil
}

func updatePlan(w http.ResponseWriter, r *http.Request) {
	// Параметры подключения к базе данных
	db, err := sql.Open("postgres", "user=bmo password=exp dbname=bmo port=5433 sslmode=disable")
	if err != nil {
		log.Println("Ошибка подключения к базе данных:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Проверяем, что метод запроса — POST
	if r.Method != http.MethodPost {
		log.Println("Неподдерживаемый метод запроса:", r.Method)
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Читаем тело запроса
	var plan Item
	err = json.NewDecoder(r.Body).Decode(&plan)
	if err != nil {
		if err == io.EOF {
			log.Println("Пустое тело запроса")
			http.Error(w, "Пустое тело запроса", http.StatusBadRequest)
			return
		}
		log.Println("Ошибка при чтении данных:", err)
		http.Error(w, "Ошибка при чтении данных", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Формируем SQL-запрос с учётом возможных NULL значений
	query := "UPDATE study_plan.plans SET "
	params := []interface{}{}
	setClauses := []string{}

	if plan.WeekNumber != 0 {
		setClauses = append(setClauses, "week_number = $1")
		params = append(params, plan.WeekNumber)
	}
	if plan.WeekDay != "" {
		setClauses = append(setClauses, "week_day = $2")
		params = append(params, plan.WeekDay)
	}
	if plan.Time != "" {
		setClauses = append(setClauses, "time = $3")
		params = append(params, plan.Time)
	}
	if plan.Group != "" {
		setClauses = append(setClauses, "groups = $4")
		params = append(params, plan.Group)
	}
	if plan.Teacher != "" {
		setClauses = append(setClauses, "teachers = $5")
		params = append(params, plan.Teacher)
	}
	if plan.Subject != "" {
		setClauses = append(setClauses, "disciplines = $6")
		params = append(params, plan.Subject)
	}
	if plan.Subject_lvl2 != "" {
		setClauses = append(setClauses, "lesson_option = $7")
		params = append(params, plan.Subject_lvl2)
	}
	if plan.ClassRoom != "" {
		setClauses = append(setClauses, "classroom = $8")
		params = append(params, plan.ClassRoom)
	}

	query += strings.Join(setClauses, ", ") + " WHERE id = $9"
	params = append(params, plan.Id)

	// Подготавливаем SQL-запрос для обновления данных
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("Ошибка подготовки запроса:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Выполняем запрос с передачей параметров
	_, err = stmt.Exec(params...)
	if err != nil {
		log.Println("Ошибка выполнения запроса:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ об успешном обновлении
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "План успешно обновлен"})
}

func deletePlan(w http.ResponseWriter, r *http.Request) {
	// Параметры подключения к базе данных
	db, err := sql.Open("postgres", "user=bmo password=exp dbname=bmo port=5433 sslmode=disable")
	if err != nil {
		log.Println("Ошибка подключения к базе данных:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Проверяем, что метод запроса — DELETE
	if r.Method != http.MethodDelete {
		log.Println("Неподдерживаемый метод запроса:", r.Method)
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Читаем тело запроса
	var plan Item
	err = json.NewDecoder(r.Body).Decode(&plan)
	if err != nil {
		log.Println("Ошибка при чтении данных:", err)
		http.Error(w, "Ошибка при чтении данных", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Подготавливаем SQL-запрос для удаления данных
	stmt, err := db.Prepare("DELETE FROM study_plan.plans WHERE id = $1")
	if err != nil {
		log.Println("Ошибка подготовки запроса:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Выполняем запрос с передачей параметров
	_, err = stmt.Exec(plan.Id)
	if err != nil {
		log.Println("Ошибка выполнения запроса:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ об успешном удалении
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "План успешно удален"})
	log.Printf("Значения: WeekNumber=%d, WeekDay=%s, Time=%s", plan.WeekNumber, plan.WeekDay, plan.Teacher)
}
