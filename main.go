package main

import (
	"fmt"
	"github.com/kongyixueyuan.com/bill/blockchain"
	"os"
	"github.com/kongyixueyuan.com/bill/service"
	"encoding/json"
	"github.com/kongyixueyuan.com/bill/web/controller"
	"github.com/kongyixueyuan.com/bill/web"
)

func main() {
	// 定义SDK属性// 链码相关参数
	fSetup := blockchain.FabricSetup{
		OrgAdmin:   "Admin",
		OrgName:    "Org1",
		ConfigFile: "config.yaml",

		// 通道相关
		ChannelID:     "mychannel",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/kongyixueyuan.com/bill/fixtures/artifacts/channel.tx",

		// 链码相关参数
		ChaincodeID:     "billcc",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/kongyixueyuan.com/bill/chaincode",
		UserName:        "User1",
	}

	// 初始化SDK
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("无法初始化Fabric SDK: %v\n", err)
	}

	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("无法安装及实例化链码: %v\n", err)
	}

	// ==========================测试开始==============================
	// 发布票据
	bill :=  service.Bill{
		BillInfoID:		"BOC10000001",
		BillInfoAmt:	"222",
		BillInfoType:	"111",
		BillInfoIsseDate:	"20180501",
		BillInfoDueDate:	"20180503",
		DrwrCmID:			"111",
		DrwrAcct:			"111",
		AccptrCmID:			"111",
		AccptrAcct:			"111",
		PyeeCmID:			"111",
		PyeeAcct:			"111",
		HoldrCmID:			"BCMID",
		HoldrAcct:			"B公司",
	}

	fsservice :=  new(service.FabricSetupService)
	fsservice.Fabric = &fSetup

	// 发布票据
	resp,err := fsservice.SaveBill(bill)
	if err!=nil {
		fmt.Printf("发布票据失败: %v\n",err)
	}else {
		fmt.Println("发布票据成功 =========> " + resp)
	}

	//根据当前用户的证件号码查询票据列表
	result,err := fsservice.FindBills("BCMID")
	if err!=nil {
		fmt.Printf("执行查询失败: %v\n", err)
	} else {
		fmt.Println("根据当前用户的证件号码查询票据列表成功")
		var bills =[]service.Bill{}
		json.Unmarshal([]byte(result),&bills)
		for _,temp := range bills{
			fmt.Println(temp)
		}
	}

	// 根据票据号码获取票据状态及该票据的背书历史
	result, err = fsservice.FindBillByNo("BOC10000001")
	if err != nil {
		fmt.Printf("执行查询失败: %v\n", err)
	}else{
		fmt.Println("根据票据号码获取票据状态及该票据的背书历史查询成功")
		var billInfo = service.Bill{}
		json.Unmarshal([]byte(result), &billInfo)
		fmt.Println("=========> ", billInfo)
		var hisBills = billInfo.History
		fmt.Println("历史信息如下:")
		for _, hisBill := range hisBills {
			fmt.Println("=========> ", hisBill)
		}

	}

	// 票据背书请求
	resp, err = fsservice.Endorse("BOC10000001", "CCMID", "C公司")
	if err != nil {
		fmt.Printf("票据背书请求失败: %v\n", err)
	}else{
		fmt.Println("票据背书请求成功")
		fmt.Println("=========> " + resp)
	}

	// 根据待背书人证件号码, 查询当前用户的待背书票据
	result, err = fsservice.FindWaitBills("CCMID")
	if err != nil {
		fmt.Printf("查询当前用户的待背书票据失败: %v\n", err)
	}else{
		fmt.Println("查询当前用户的待背书票据成功")
		var bills = []service.Bill{}
		json.Unmarshal([]byte(result), &bills)
		for _, temp := range bills{
			fmt.Println("=========> ", temp)
		}
	}

	// 票据背书签收
	resp, err = fsservice.EndorseAccept("BOC10000001", "CCMID", "C公司")

	if err != nil {
		fmt.Printf("票据背书签收失败: %v\n", err)
	}else{
		fmt.Println("票据背书签收成功")
		fmt.Println("=========> " + resp)
	}


/*	err = fsservice.Delete("BOC10000001")
	if err != nil {
		fmt.Printf("票据BOC10000001删除失败: %v\n", err)
	}else{
		fmt.Println("票据BOC10000001删除成功")
	}
*/
	app := controller.Application{Fabric:fsservice}
	web.WebStart(&app)
}