package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const WorkflowSchema = `
CREATE TABLE workflow (
  id serial,
  user_id integer NOT NULL REFERENCES user_acc,
  name varchar NOT NULL,
  trigger_app varchar NOT NULL,
  trigger_action varchar NOT NULL,
  PRIMARY KEY (id)
);
`

const WorkflowStructViewSchema = `
CREATE VIEW workflow_step as
SELECT id, user_id, name, trigger_app, trigger_action, (
  SELECT array_to_json(array_agg(to_json(row)))
  FROM (
    SELECT *
    FROM step
    WHERE step.workflow_id = workflow.id
  ) row
) steps
FROM workflow;
`

type Workflow struct {
	ID            int     `json:"id"`
	UserID        int     `json:"user_id"`
	Name          string  `json:"name"`
	TriggerApp    string  `json:"trigger_app"`
	TriggerAction string  `json:"trigger_action"`
	Steps         []*Step `json:"steps,omitempty"`
}

const StepSchema = `
CREATE TABLE step (
  id serial,
  workflow_id integer NOT NULL REFERENCES workflow,
  app varchar NOT NULL,
  action varchar NOT NULL,
  input_variables jsonb NOT NULL,
  PRIMARY KEY (id)
);
`

type InputVariables map[string]string

func (vars *InputVariables) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &vars)
		return nil
	case string:
		json.Unmarshal([]byte(v), &vars)
		return nil
	default:
		return fmt.Errorf("Unsupported type: %T", v)
	}
	return nil
}

func (vars InputVariables) Value() (driver.Value, error) {
	return json.Marshal(&vars)
}

type Step struct {
	ID             int            `json:"id"`
	WorkflowID     int            `json:"workflow_id"`
	App            string         `json:"app"`
	Action         string         `json:"action"`
	InputVariables InputVariables `json:"input_variables"`
}

func GetAllWorkflowsWithSteps(userID int) (*[]Workflow, error) {
	stmt := `
	SELECT array_to_json(array_agg(row))
	FROM (
		SELECT *
		FROM workflow_step
		WHERE user_id = $1
	) row;
	`
	var workflowStr string
	var workflows []Workflow
	if err := db.Get(&workflowStr, stmt, userID); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(workflowStr), &workflows); err != nil {
		return nil, err
	}
	return &workflows, nil
}

func GetAllWorkflows(userID int) (*[]Workflow, error) {
	var workflows []Workflow
	if err := db.Select(&workflows, "SELECT * FROM workflow WHERE user_id = $1", userID); err != nil {
		return nil, err
	}
	return &workflows, nil
}

func GetWorkflowWithSteps(workflowID int, userID int) (*Workflow, error) {
	stmt := `
	SELECT row_to_json(row)
	FROM (
		SELECT *
		FROM workflow_step
		WHERE id = $1 and user_id = $2
	) row
	`
	var workflowStr string
	workflow := Workflow{}
	if err := db.Get(&workflowStr, stmt, workflowID, userID); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(workflowStr), &workflow); err != nil {
		return nil, err
	}
	return &workflow, nil
}

func GetWorkflow(workflowID int, userID int) (*Workflow, error) {
	workflow := Workflow{}
	if err := db.Get(&workflow, "SELECT * FROM workflow WHERE id = $1 and user_id = $2", workflowID, userID); err != nil {
		return nil, err
	}
	return &workflow, nil
}

func NewWorkflow(workflow *Workflow) (int, error) {
	stmt := `
	INSERT INTO workflow (user_id, name, trigger_app, trigger_action)
	VALUES (:user_id, :name, :trigger_app, :trigger_action)
	RETURNING id;
	`
	namedStmt, err := db.PrepareNamed(stmt)
	if err != nil {
		return 0, err
	}

	var workflowID int
	if err := namedStmt.Get(&workflowID, workflow); err != nil {
		return 0, err
	}
	return workflowID, nil
}

func DeleteWorkflow(workflowID int, userID int) error {
	_, err := db.Exec("DELETE FROM workflow WHERE id = $1 AND user_id = $2", workflowID, userID)
	return err
}

func UpdateWorkflow(workflow *Workflow) error {
	stmt := "UPDATE workflow SET name = :name WHERE id = :id AND user_id = :user_id"
	_, err := db.NamedExec(stmt, workflow)
	return err
}

func GetAllSteps(workflowID int, userID int) (*[]Step, error) {
	stmt := `
	SELECT step.id,
		step.workflow_id,
		step.app,
		step.action,
		step.input_variables
	FROM step INNER JOIN workflow
		ON workflow.id = step.workflow_id
	WHERE step.workflow_id = $1
		AND workflow.user_id = $2;
	`
	var steps []Step
	if err := db.Select(&steps, stmt, workflowID, userID); err != nil {
		return nil, err
	}
	return &steps, nil
}

func GetStep(stepID, workflowID, userID int) (*Step, error) {
	stmt := `
	SELECT step.id,
		step.workflow_id,
		step.app,
		step.action,
		step.input_variables
	FROM step INNER JOIN workflow
		ON workflow.id = step.workflow_id
	WHERE step.id = $1
		AND step.workflow_id = $2
		AND workflow.user_id = $3;
	`
	step := Step{}
	if err := db.Get(&step, stmt, stepID, workflowID, userID); err != nil {
		return nil, err
	}
	return &step, nil
}

func NewStep(step *Step, userID int) (int, error) {
	stmt := `
	INSERT INTO step (workflow_id, app, action, input_variables)
	SELECT $1, $2, $3, $4
	WHERE EXISTS (
		SELECT *
		FROM workflow
		WHERE workflow.user_id = $5 AND workflow.id = $1
	) RETURNING id;
	`
	var stepID int
	err := db.Get(&stepID, stmt, step.WorkflowID, step.App, step.Action, step.InputVariables, userID)
	return stepID, err
}

func UpdateStep(step *Step, userID int) error {
	stmt := `
	UPDATE step SET input_variables = $1 FROM workflow
	WHERE step.id = $2 AND step.workflow_id = $3 AND workflow.user_id = $4
	`
	_, err := db.Exec(stmt, step.InputVariables, step.ID, step.WorkflowID, userID)
	return err
}

func DeleteStep(stepID int, workflowID int, userID int) error {
	stmt := `
	DELETE FROM step
	USING workflow
	WHERE step.id = $1 AND step.workflow_id = $2 AND workflow.user_id = $3
	`
	_, err := db.Exec(stmt, stepID, workflowID, userID)
	return err
}
