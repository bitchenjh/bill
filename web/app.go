package web

import (
	"net/http"
	"fmt"
	"github.com/kongyixueyuan.com/bill/web/controller"
)

func WebStart(app *controller.Application) error {
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))


	// http://localhost:9000/
	http.HandleFunc("/", app.LoginView)
	http.HandleFunc("/login.html", app.LoginView)

	http.HandleFunc("/login", app.Login)

	http.HandleFunc("/issue.html", app.IssueView)
	http.HandleFunc("/issue", app.SaveBill)

	http.HandleFunc("/findBills", app.FindBills)
	http.HandleFunc("/billInfoByNo", app.FindBillInfoByNo)
	http.HandleFunc("/endorse", app.Endorse)
	http.HandleFunc("/waitEndorseBills", app.WaitEndorseBills)
	http.HandleFunc("/waitEndorseBillInfo", app.WaitEndorseBillInfo)

	http.HandleFunc("/accept", app.Accept)
	http.HandleFunc("/reject", app.Reject)

	http.HandleFunc("/loginout", app.Loginout)

	http.HandleFunc("/delete",app.Delete)

	fmt.Println("启动应用程序, 监听端口号为: 9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		return fmt.Errorf("启动Web服务失败: %v", err)
	}

	return nil


}

