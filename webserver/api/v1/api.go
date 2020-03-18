package v1

import (
	"github.com/plumbie/plumbie/webserver/api/v1/auth"

	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

func RegisterRoutes(m *macaron.Macaron) {
	m.Group("/v1", func() {
		m.Group("/auth", func() {
			m.Post("/login", binding.Bind(auth.LoginForm{}), auth.Login)
			m.Post("/signup", binding.Bind(auth.SignUpForm{}), auth.SignUp)
			m.Post("/logout", auth.Logout)
		})

		m.Group("/workflow", func() {
			m.Combo("/").
				Get(GetAllWorkflows).
				Post(binding.Bind(NewWorkflowForm{}), NewWorkflow)

			m.Group("/:workflowID", func() {
				m.Combo("/").
					Get(GetWorkflow).
					Patch(binding.Bind(UpdateWorkflowForm{}), UpdateWorkflow).
					Delete(DeleteWorkflow)

				m.Group("/step", func() {
					m.Combo("/").
						Get(GetAllSteps).
						Post(binding.Bind(NewStepForm{}), NewStep)

					m.Combo("/:stepID").
						Get(GetStep).
						Patch(binding.Bind(UpdateStepForm{}), UpdateStep).
						Delete(DeleteStep)
				})
			})
		})
	})
}
