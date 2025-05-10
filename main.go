package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var db *sql.DB

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	var err error
	dsn := "root:password@tcp(localhost:3306)/event_management?parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("failed to connect db")
	}
	if err := db.Ping(); err != nil {
		log.Fatal("db is unreachable")
	}
	http.HandleFunc("/events", eventsHandler)
	http.HandleFunc("/events/", eventHandler)
	fmt.Println("Server running on port: 8080")
	http.ListenAndServe(":8080", nil)
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listEvent(w)
	case http.MethodPost:
		createEvent(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	fmt.Println("header: ", w.Header().Get("Content-Type"))
}

func eventHandler(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/events/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	switch r.Method {
	case http.MethodGet:
		getEvent(w, id)
	case http.MethodPut:
		updateEvent(w, id, r)
	case http.MethodDelete:
		deleteEvent(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusBadRequest)

	}
}

func deleteEvent(w http.ResponseWriter, id int) {
	result, err := db.Exec(`DELETE from events where id = ?`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		http.Error(w, "No event found with the given id", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Event deleted")
}

func updateEvent(w http.ResponseWriter, id int, r *http.Request) {
	event, err := getEventFromDb(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var updatedEvent Event
	if err := json.NewDecoder(r.Body).Decode(&updatedEvent); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	updatedEvent.ID = id
	updatedEvent.CreatedAt = event.CreatedAt
	updatedEvent.UpdatedAt = time.Now()
	_, err = db.Exec(`
	Update events SET title=?, description=?, location=?, start_time=?, end_time=?, created_by=?, updated_at=?
	WHERE id = ?`,
		updatedEvent.Title, updatedEvent.Description, updatedEvent.Location, updatedEvent.StartTime, updatedEvent.EndTime, updatedEvent.CreatedAt, updatedEvent.UpdatedAt, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := Response{
		Message: "event updated",
		Data:    updatedEvent,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getEvent(w http.ResponseWriter, id int) {
	e, err := getEventFromDb(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}


func getEventFromDb(id int) (*Event, error) {
	var e Event
	err := db.QueryRow(`SELECT * from events where id = ?`, id).Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.EndTime, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	result, err := db.Exec(`INSERT into events(title, description, location, start_time, end_time, created_by, created_at, updated_at) values(?, ?, ?, ?, ?, ?, ?, ?)`,
		event.Title, event.Description, event.Location, event.StartTime, event.EndTime, event.CreatedBy, event.CreatedAt, event.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	event.ID = int(id)

	response := Response{
		Message: "Event created",
		Data:    event,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

func listEvent(w http.ResponseWriter) {
	var events []Event
	rows, err := db.Query("SELECT * from events")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.EndTime, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		events = append(events, e)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("header: ", w.Header().Get("Content-Type"))
	json.NewEncoder(w).Encode(events)
}
