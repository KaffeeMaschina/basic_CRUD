package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type Storage struct {
	DB     *Database
	Logger *slog.Logger
}

const (
	noRows = "no rows in result set"
)

var (
	errorNoTitle     = errors.New("title is missing")
	errorZeroId      = errors.New("id shouldn't be 0")
	errorIdNotNumber = errors.New("id should be number")
	errorNoTask      = errors.New("there is no task with id: ")
	errorWrongStatus = errors.New("status should be: 'new', 'in_progress' or 'done'")
	errorNoDelete    = errors.New("there is no delete")
	errorNoUpdate    = errors.New("there is no update")
)

// CreateTask creates a new task and upload it to database and returns the created instance
func (s *Storage) CreateTask(c *fiber.Ctx) error {
	const errorLocation = "internal.app.storage.storage.CreateTask"
	s.Logger.Debug("creating task")

	// Creating task with a timestamp
	timeStamp := time.Now()
	task := &Task{
		CreatedAt: timeStamp,
		UpdatedAt: timeStamp,
	}
	// Parsing body from request
	if err := c.BodyParser(task); err != nil {
		s.Logger.Error("Parsing error: ", err, "at ", errorLocation)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Check if task has a title, if not return BadRequest
	if task.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   errorNoTitle.Error(),
		})
	}

	// Check if status is empty, then status is new
	if task.Status == "" {
		task.Status = "new"
	}
	s.Logger.Debug("body is successfully parsed")

	// Sending data to database
	err := s.DB.Conn.QueryRow(context.Background(), `INSERT INTO tasks (title, description, status, created_at, updated_at) 
		 VALUES ($1,$2,$3,$4,$5) RETURNING ID`, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt).
		Scan(&task.ID)

	if err != nil {
		// Check if error is "SQLSTATE 23514", then status is wrong, return BadRequest
		if strings.Contains(err.Error(), "SQLSTATE 23514") {
			s.Logger.Error("status error: ", errorWrongStatus, "at ", errorLocation)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   errorWrongStatus.Error(),
			})
		}
		// Other errors
		s.Logger.Error("Error inserting task: ", err, "at ", errorLocation)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	// Logging about success, and returning task
	s.Logger.Info("task is successfully created and added to database ")
	return c.JSON(fiber.Map{
		"error": false,
		"task":  task,
	})
}

// GetTasks gets all the tasks from database and returns a json with this tasks
func (s *Storage) GetTasks(c *fiber.Ctx) error {
	const errorLocation = "internal.app.storage.storage.GetTasks"
	s.Logger.Debug("getting tasks")
	// Getting all the rows  from database
	rows, err := s.DB.Conn.Query(context.Background(), `SELECT * FROM tasks`)
	if err != nil {
		s.Logger.Error("Error getting tasks: ", err, "at ", errorLocation)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	// Scan each row to "task" and append to "tasks"
	var tasks []Task
	for rows.Next() {
		var task Task

		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			s.Logger.Error("Error getting tasks: ", err, "at ", errorLocation)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
		tasks = append(tasks, task)
	}
	// Return tasks
	return c.JSON(fiber.Map{
		"error": false,
		"tasks": tasks,
	})
}

func (s *Storage) UpdateTask(c *fiber.Ctx) error {

	const errorLocation = "internal.app.storage.storage.UpdateTask"
	s.Logger.Debug("updating task")

	// Get id from url, and parse it to int
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.Logger.Error("Error converting id to int: ", err, "at ", errorLocation)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   errorIdNotNumber.Error(),
		})
	}

	// Check, id shouldn't be 0
	if id == 0 {
		s.Logger.Error("Error id is zero: ", err, "at ", errorLocation)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   errorZeroId.Error(),
		})
	}

	// Creating task
	task := &Task{ID: id}

	// Parsing body from request
	if err = c.BodyParser(task); err != nil {
		s.Logger.Error("Parsing error: ", err, "at ", errorLocation)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//Check if there is a task with given id in the database
	foundedTask, err := s.GetTask(id)
	if err != nil {
		// Check if error contains "no rows in result set", then there is no task with given id
		if strings.Contains(err.Error(), noRows) {
			s.Logger.Error("No such task: ", errorNoTask, "at ", errorLocation)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   errorNoTask.Error(),
			})
		} else {
			s.Logger.Error("Error getting task: ", err, "at ", errorLocation)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	}
	// Get "createdAt" from database
	task.CreatedAt = foundedTask.CreatedAt
	// Create a new timestamp of update
	task.UpdatedAt = time.Now()

	// if new title isn't mentioned in http request, get it from database
	if task.Title == "" {
		task.Title = foundedTask.Title
	}
	// if status isn't mentioned in http request, get it from database
	if task.Status == "" {
		task.Status = foundedTask.Status
	}
	// if description isn't mentioned in http request, get it from database
	if task.Description == nil {
		task.Description = foundedTask.Description
	}
	// Update task in database
	tag, err := s.DB.Conn.Exec(context.Background(),
		`UPDATE tasks SET 
            title = $2, description = $3, status = $4, 
            created_at = $5, updated_at = $6 WHERE id = $1`,
		id, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		// Check if error is "SQLSTATE 23514", then status is wrong, return BadRequest
		if strings.Contains(err.Error(), "SQLSTATE 23514") {
			s.Logger.Error("status error: ", errorWrongStatus, "at ", errorLocation)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   errorWrongStatus.Error(),
			})
		} else {
			s.Logger.Error("Error updating task: ", err, "at ", errorLocation)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	}
	// Check if deleting is successful
	if tag.RowsAffected() == 0 {
		s.Logger.Error("Error updating task: ", errorNoDelete, "at ", errorLocation)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errorNoUpdate.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"task":  task,
	})

}
func (s *Storage) DeleteTask(c *fiber.Ctx) error {

	const errorLocation = "internal.app.storage.storage.DeleteTask"
	// Get id from url, and parse it to int
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.Logger.Error("Error converting id to int: ", err, "at ", errorLocation)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   errorIdNotNumber.Error(),
		})
	}
	// Check, id shouldn't be 0
	if id == 0 {
		s.Logger.Error("Error id is zero: ", err, "at ", errorLocation)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   errorZeroId.Error(),
		})
	}
	//Check if there is a task with given id in the database
	_, err = s.GetTask(id)
	if err != nil {
		// Check if error contains "no rows in result set", then there is no task with given id
		if strings.Contains(err.Error(), noRows) {
			s.Logger.Error("No such task: ", errorNoTask, "at ", errorLocation)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   fmt.Sprintf("%v%v", errorNoTask.Error(), id),
			})
		} else {
			s.Logger.Error("Error getting task: ", err, "at ", errorLocation)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	}
	tag, err := s.DB.Conn.Exec(context.Background(),
		"DELETE FROM tasks WHERE id = $1", id)

	if err != nil {
		s.Logger.Error("Error deleting task: ", err, "at ", errorLocation)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Check if deleting is successful
	if tag.RowsAffected() == 0 {
		s.Logger.Error("Error deleting task: ", errorNoDelete, "at ", errorLocation)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errorNoDelete.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "successful deletion",
	})
}

func (s *Storage) GetTask(id int) (Task, error) {
	const errorLocation = "internal.app.storage.storage.GetTask"

	if id == 0 {
		s.Logger.Error("Error id is zero at: ", errorLocation)
	}
	s.Logger.Debug("getting task: ", id)

	var task Task
	err := s.DB.Conn.QueryRow(context.Background(), `SELECT * FROM tasks WHERE id = $1`, id).
		Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		s.Logger.Error("Error getting task: ", err, "at ", errorLocation)
		return task, err
	}

	s.Logger.Info("task is successfully retrieved: ", task)
	return task, nil
}
