package v1

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/plumbie/plumbie/apps"
	"github.com/plumbie/plumbie/models"

	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

type NewWorkflowForm struct {
	Name       string `binding:"Required" form:"name"`
	TriggerApp string `binding:"Required" form:"trigger_app"`
	Trigger    string `binding:"Required" form:"trigger"`
}

type UpdateWorkflowForm struct {
	Name string `form:"name"`
}

type NewStepForm struct {
	App            string `binding:"Required" form:"app"`
	Action         string `binding:"Required" form:"action"`
	InputVariables string `binding:"Required" form:"input_variables"`
}

type UpdateStepForm struct {
	InputVariables string `form:"input_variables"`
}

func GetAllWorkflows(ctx *macaron.Context, sess session.Store) {
	userID, ok := sess.Get("userID").(int)
	if !ok {
		ctx.JSON(400, "Invalid IDs")
		return
	}
	steps := ctx.QueryBool("steps")

	var workflows *[]models.Workflow
	var err error
	if steps {
		workflows, err = models.GetAllWorkflowsWithSteps(userID)
	} else {
		workflows, err = models.GetAllWorkflows(userID)
	}
	if err != nil {
		ctx.JSON(400, err)
		return
	}
	ctx.JSON(200, workflows)
}

func GetWorkflow(ctx *macaron.Context, sess session.Store) {
	id, err := strconv.Atoi(ctx.Params("workflowID"))
	userID, ok := sess.Get("userID").(int)
	if err != nil || !ok {
		ctx.JSON(400, "Invalid IDs")
		return
	}
	steps := ctx.QueryBool("steps")

	var workflow *models.Workflow
	if steps {
		workflow, _ = models.GetWorkflowWithSteps(id, userID)
	} else {
		workflow, _ = models.GetWorkflow(id, userID)
	}

	if workflow != nil {
		ctx.JSON(200, workflow)
		return
	}
	ctx.JSON(403, "Workflow not found")
}

func NewWorkflow(ctx *macaron.Context, form NewWorkflowForm, sess session.Store) {
	userID, ok := sess.Get("userID").(int)
	if !ok {
		ctx.JSON(400, "Invalid IDs")
		return
	}

	workflow := models.Workflow{
		UserID:     userID,
		Name:       form.Name,
		TriggerApp: form.TriggerApp,
		Trigger:    form.Trigger,
	}
	if apps.TriggerExists(workflow.TriggerApp, workflow.Trigger) {
		id, err := models.NewWorkflow(&workflow)
		if err != nil {
			ctx.JSON(400, err.Error())
			return
		}
		ctx.JSON(200, id)
		return
	}

	ctx.JSON(400, "No such app or action")
}

func UpdateWorkflow(ctx *macaron.Context, form UpdateWorkflowForm, sess session.Store) {
	id, err := strconv.Atoi(ctx.Params("workflowID"))
	userID, ok := sess.Get("userID").(int)
	if err != nil || !ok {
		ctx.JSON(400, "Invalid IDs")
		return
	}
	if form.Name != "" {
		err := models.UpdateWorkflow(&models.Workflow{
			ID:     id,
			UserID: userID,
			Name:   form.Name,
		})
		if err != nil {
			ctx.JSON(400, err.Error())
			return
		}
	}

	ctx.JSON(200, "ok")
}

func DeleteWorkflow(ctx *macaron.Context, sess session.Store) {
	id, err := strconv.Atoi(ctx.Params("workflowID"))
	userID, ok := sess.Get("userID").(int)
	if err != nil || !ok {
		fmt.Println(id, userID)
		ctx.JSON(400, "Invalid IDs")
		return
	}
	if err := models.DeleteWorkflow(id, userID); err != nil {
		ctx.JSON(400, err.Error())
		return
	}
	ctx.JSON(200, "ok")
}

func GetAllSteps(ctx *macaron.Context, sess session.Store) {
	workflowID, err := strconv.Atoi(ctx.Params("workflowID"))
	userID, ok := sess.Get("userID").(int)
	if err != nil || !ok {
		ctx.JSON(400, "Invalid IDs")
		return
	}
	steps, err := models.GetAllSteps(workflowID, userID)
	if err != nil {
		ctx.JSON(400, err)
		return
	}
	ctx.JSON(200, steps)
}

func GetStep(ctx *macaron.Context, sess session.Store) {
	id, err := strconv.Atoi(ctx.Params("stepID"))
	workflowID, err2 := strconv.Atoi(ctx.Params("workflowID"))
	userID, ok := sess.Get("userID").(int)
	if err != nil || err2 != nil || !ok {
		ctx.JSON(400, "Invalid IDs")
		return
	}

	step, _ := models.GetStep(id, workflowID, userID)

	if step != nil {
		ctx.JSON(200, step)
		return
	}
	ctx.JSON(403, "Step not found")

}

func NewStep(ctx *macaron.Context, form NewStepForm, sess session.Store) {
	var inputVariables models.InputVariables
	workflowID, err := strconv.Atoi(ctx.Params("workflowID"))
	userID, ok := sess.Get("userID").(int)
	if err != nil || !ok {
		ctx.JSON(400, "Invalid IDs")
		return
	}

	if err := json.Unmarshal([]byte(form.InputVariables), &inputVariables); err != nil {
		ctx.JSON(400, err)
		return
	}

	if apps.ActionExists(form.App, form.Action) {
		id, err := models.NewStep(&models.Step{
			WorkflowID:     workflowID,
			App:            form.App,
			Action:         form.Action,
			InputVariables: inputVariables,
		}, userID)
		if err != nil {
			ctx.JSON(400, err.Error())
			return
		}
		ctx.JSON(200, id)
		return
	}

	ctx.JSON(400, "No such app or action")

}

func UpdateStep(ctx *macaron.Context, form UpdateStepForm, sess session.Store) {
	id, err := strconv.Atoi(ctx.Params("stepID"))
	workflowID, err2 := strconv.Atoi(ctx.Params("workflowID"))
	userID, ok := sess.Get("userID").(int)
	if err != nil || err2 != nil || !ok {
		ctx.JSON(400, "Invalid IDs")
		return
	}

	var inputVariables models.InputVariables
	if err := json.Unmarshal([]byte(form.InputVariables), &inputVariables); err != nil {
		ctx.JSON(400, err)
		return
	}

	if form.InputVariables != "" {
		err := models.UpdateStep(&models.Step{
			ID:             id,
			WorkflowID:     workflowID,
			InputVariables: inputVariables,
		}, userID)
		if err != nil {
			ctx.JSON(400, err.Error())
			return
		}
	}

	ctx.JSON(200, "ok")

}

func DeleteStep(ctx *macaron.Context, sess session.Store) {
	id, err := strconv.Atoi(ctx.Params("stepID"))
	workflowID, err2 := strconv.Atoi(ctx.Params("workflowID"))
	userID, ok := sess.Get("userID").(int)
	if err != nil || err2 != nil || !ok {
		ctx.JSON(400, "Invalid IDs")
		return
	}

	if err := models.DeleteStep(id, workflowID, userID); err != nil {
		ctx.JSON(400, err.Error())
		return
	}
	ctx.JSON(200, "ok")
}
