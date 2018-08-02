package controller

import (
	"net/http"
	"github.com/kongyixueyuan.com/bill/service"
	"encoding/json"
	"fmt"
)

type Application struct {
	Fabric *service.FabricSetupService
}

var cuser User

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request)  {
	response(w, r, "login.html", nil)
}

func (app *Application) Login(w http.ResponseWriter, r *http.Request)  {
	userName := r.FormValue("userName")
	password := r.FormValue("password")

	var flag = false
	for _, user := range Users {
		if userName == user.UserName && password == user.Password {
			cuser = user
			flag = true
			break
		}
	}

	if flag{
		app.FindBills(w, r)
	}else {
		data := &struct {
			Name string
			Flag bool
		}{
			Name:userName,
			Flag:true,
		}
		response(w, r, "login.html", data)
	}

}

// 显示发布票据页面
func (app *Application) IssueView(w http.ResponseWriter, r *http.Request)  {
	data := &struct {
		Flag bool
		Msg string
		Cuser User
	}{
		Flag:false,
		Msg:"",
		Cuser:cuser,
	}
	response(w, r, "issue.html", data)
}

// 发布新票据
func (app *Application) SaveBill(w http.ResponseWriter, r *http.Request)  {
	bill := service.Bill{
		BillInfoID: r.FormValue("BillInfoID"),
		BillInfoType:r.FormValue("BillInfoType"),
		BillInfoAmt:r.FormValue("BillInfoAmt"),

		BillInfoDueDate:r.FormValue("BillInfoDueDate"),
		BillInfoIsseDate:r.FormValue("BillInfoIsseDate"),

		DrwrAcct: r.FormValue("DrwrAcct"),
		DrwrCmID: r.FormValue("DrwrCmID"),

		AccptrAcct: r.FormValue("AccptrAcct"),
		AccptrCmID: r.FormValue("AccptrCmID"),

		PyeeAcct: r.FormValue("PyeeAcct"),
		PyeeCmID: r.FormValue("PyeeCmID"),

		HoldrAcct: r.FormValue("HoldrAcct"),
		HoldrCmID: r.FormValue("HoldrCmID"),
	}

	transactionId, err := app.Fabric.SaveBill(bill)
	var msg string
	if err != nil {
		msg = "票据发布失败: " + err.Error()
	}else{
		msg = "票据发布成功: " + transactionId
	}

	data := &struct {
		Msg string
		Flag bool
		Cuser User
	}{
		Msg: msg,
		Flag:true,
		Cuser:cuser,
	}
	response(w, r, "issue.html", data)

}

// 查询我的票据列表
func (app *Application) FindBills(w http.ResponseWriter, r *http.Request)  {
	result, err := app.Fabric.FindBills(cuser.CmId)
	if err != nil{
		fmt.Println("查询票据列表错误: %v", err)
	}
	var bills []service.Bill
	json.Unmarshal(result, &bills)

	data := &struct {
		Bills []service.Bill
		Cuser User
	}{
		Bills:bills,
		Cuser:cuser,
	}

	response(w, r, "bills.html", data)
}

func (app *Application) FindBillInfoByNo(w http.ResponseWriter, r *http.Request)  {
	billInfoNo := r.FormValue("billInfoNo")
	result, err := app.Fabric.FindBillByNo(billInfoNo)
	if err != nil {
		fmt.Println(err.Error())
	}

	var bill service.Bill
	json.Unmarshal(result, &bill)

	data := &struct {
		Bill service.Bill
		Cuser User
		Flag bool
		Msg string
	}{
		Bill:bill,
		Cuser:cuser,
		Flag:false,
		Msg:"",
	}

	flag := r.FormValue("flag")
	if flag == "t"{
		data.Flag = true
		data.Msg = r.FormValue("Msg")
	}

	response(w, r, "billInfo.html", data)

}

// 发起背书请求
func (app *Application) Endorse(w http.ResponseWriter, r *http.Request)  {
	waitEndorseAcct := r.FormValue("waitEndorseAcct")
	waitEndorseCmId := r.FormValue("waitEndorseCmId")
	billNo := r.FormValue("billNo")

	result, err := app.Fabric.Endorse(billNo, waitEndorseCmId, waitEndorseAcct)
	if err != nil {
		fmt.Println(err.Error())
	}

	r.Form.Set("billInfoNo", billNo)

	r.Form.Set("flag", "t")
	r.Form.Set("Msg", result)


	app.FindBillInfoByNo(w, r)

	//response(w, r, "billInfo.html", data)
}

// 待背书票据列表
func (app *Application) WaitEndorseBills(w http.ResponseWriter, r *http.Request){
	waitEndorseCmId := cuser.CmId
	result, err := app.Fabric.FindWaitBills(waitEndorseCmId)
	if err != nil {
		fmt.Println(err.Error())
	}
	var bills []service.Bill
	json.Unmarshal(result, &bills)

	data := &struct {
		Bills []service.Bill
		Cuser User
	}{
		Bills:bills,
		Cuser:cuser,
	}
	response(w, r, "waitBills.html", data)
}

// 待背书票据详情
func (app *Application) WaitEndorseBillInfo(w http.ResponseWriter, r *http.Request)  {
	billNo := r.FormValue("billNo")
	result, err := app.Fabric.FindBillByNo(billNo)
	if err != nil {
		fmt.Println(err.Error())
	}

	var bill service.Bill
	json.Unmarshal(result, &bill)

	data := &struct {
		Bill service.Bill
		Cuser User
		Flag bool
		Msg string
	}{
		Bill:bill,
		Cuser:cuser,
		Flag:false,
		Msg:"",
	}

	flag := r.FormValue("flag")
	if flag == "t" {
		data.Flag = true
		data.Msg = r.FormValue("Msg")
	}

	response(w, r, "waitBillInfo.html", data)
}

func (app *Application) Loginout(w http.ResponseWriter, r *http.Request)  {
	cuser = User{}
	app.LoginView(w, r)
}

func (app *Application)Delete(w http.ResponseWriter,r *http.Request)  {
	billNo := r.FormValue("billNo")
	err := app.Fabric.Delete(billNo)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("删除票据结束")
	response(w, r, "bills.html", nil)
	fmt.Println("删除票据结束跳转")
}

// 签收
func (app *Application) Accept(w http.ResponseWriter, r *http.Request){
	billNo := r.FormValue("billNo")
	cmid := cuser.CmId
	acct := cuser.Acct

	result, err := app.Fabric.EndorseAccept(billNo,cmid, acct)
	if err != nil {
		fmt.Println(err.Error())
	}

	r.Form.Set("billNo", billNo)
	r.Form.Set("flag", "t")
	r.Form.Set("Msg", result)

	app.WaitEndorseBillInfo(w, r)

}

// 拒签
func (app *Application) Reject(w http.ResponseWriter, r *http.Request)  {
	billNo := r.FormValue("billNo")
	cmid := cuser.CmId
	acct := cuser.Acct
	result, err := app.Fabric.EndorseReject(billNo, cmid, acct)
	if err != nil {
		fmt.Println(err.Error())
	}

	r.Form.Set("billNo", billNo)
	r.Form.Set("flag", "t")
	r.Form.Set("Msg", result)

	app.WaitEndorseBillInfo(w, r)
}

