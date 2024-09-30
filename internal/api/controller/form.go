package controller

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type FormData struct {
	Data   map[string]string
	Errors map[string]string
}

func NewFormData() FormData {
	return FormData{
		Data:   make(map[string]string),
		Errors: make(map[string]string),
	}
}

func GetForm(ctx *gin.Context) FormData {
	ctx.Request.ParseForm()
	form := NewFormData()

	for key, values := range ctx.Request.PostForm {
		if len(values) > 0 {
			form.Data[key] = strings.TrimSpace(values[0])
		}
	}

	return form
}
