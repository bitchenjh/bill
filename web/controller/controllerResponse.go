package controller

import (
	"net/http"
	"path/filepath"
	"html/template"
	"fmt"
)

func response(w http.ResponseWriter, r *http.Request, templateName string, data interface{})  {
	path := filepath.Join("web", "tpl", templateName)

	// 创建模板实例
	result, err := template.ParseFiles(path)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	// 融合数据
	err = result.Execute(w, data)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
}
