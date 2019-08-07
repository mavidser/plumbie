package auth

import (
	"fmt"

	"github.com/plumbie/plumbie/models"

	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

type LoginForm struct {
	Username string `binding:"Required"`
	Password string `binding:"Required"`
}

func Login(ctx *macaron.Context, form LoginForm, sess session.Store) {
	fmt.Println("XXX")
	if user, err := models.UserLogin(form.Username, form.Password); err != nil {
		ctx.JSON(200, map[string]interface{}{
			"ok":    false,
			"error": err.Error(),
		})
	} else {
		sess.Set("userid", user.ID)
		ctx.JSON(200, map[string]interface{}{
			"ok":   true,
			"user": user,
		})
	}
}

type SignUpForm struct {
	Username string `binding:"Required"`
	Password string `binding:"Required"`
}

func SignUp(ctx *macaron.Context, form SignUpForm) {
	user := models.User{
		Username: form.Username,
		Password: form.Password,
	}
	if err := user.Create(); err != nil {
		ctx.JSON(200, map[string]interface{}{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, map[string]interface{}{
		"ok": true,
	})
}

func Logout(ctx *macaron.Context, sess session.Store) {
	// err := sess.Destory(ctx)
	fmt.Println(sess.Get("userid").(int64))
	err := sess.Delete("userid")
	if err != nil {
		ctx.JSON(200, map[string]interface{}{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, map[string]interface{}{
		"ok": true,
	})
}
