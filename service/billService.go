package service

import (
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
	"fmt"
	"encoding/json"
)

// 发布票据
func (setup *FabricSetupService) SaveBill(bill Bill) (string, error) {
	var args []string
	args = append(args, "issue")

	b, _ := json.Marshal(bill)

	// 设置交易请求参数
	req := chclient.Request{ChaincodeID: setup.Fabric.ChaincodeID, Fcn: args[0], Args: [][]byte{b}}

	// 执行交易
	response, err := setup.Fabric.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf("保存票据时发生错误: %v\n", err)
	}

	return response.TransactionID.ID, nil
}

// 根据当前持票人证件号码, 批量查询票据
func (setup *FabricSetupService) FindBills(holderCmId string) ([]byte, error) {

	var args []string
	args = append(args, "queryMyBills")
	args = append(args, holderCmId)

	// 设置查询的请求参数
	req := chclient.Request{ChaincodeID: setup.Fabric.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}}

	// 执行查询
	response, err := setup.Fabric.Client.Query(req)
	if err != nil {
		return []byte{0x00}, fmt.Errorf("查询票据失败: %v\n", err)
	}

	b := response.Payload

	return b[:], nil
}

func (setup *FabricSetupService)Delete(bill_no string) error {
	fmt.Println("开始删除票据:"+bill_no)
	var args []string
	args = append(args, "delete")
	req := chclient.Request{ChaincodeID:setup.Fabric.ChaincodeID,Fcn:args[0],Args:[][]byte{[]byte(bill_no)}}

	// 执行交易
	_, err := setup.Fabric.Client.Execute(req)
	if err != nil {
		return fmt.Errorf("删除票据时发生错误: %v\n", err)
	}
	fmt.Println("开始删除票据完成" )
	return nil
}


// 根据票据号码获取票据状态及该票据的背书历史
func (setup *FabricSetupService) FindBillByNo(bill_no string) ([]byte, error) {
	var args []string
	args = append(args, "queryBillByNo")

	req := chclient.Request{ChaincodeID: setup.Fabric.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(bill_no)}}
	response, err := setup.Fabric.Client.Query(req)
	if err != nil {
		return []byte{0x00}, fmt.Errorf("查询指定的票据信息发生错误: %v\n", err)
	}

	b := response.Payload

	return b[:], nil
}
