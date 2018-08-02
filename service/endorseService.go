package service

import (
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
	"fmt"
)

// 票据背书请求
// args: 0 - Bill_No;  1 - endorserCmId(待背书人ID);	2 - endorserAcct(待背书人名称)
func (setup *FabricSetupService) Endorse(bill_no string, endorseCmId string, endorseAcct string) (string, error) {
	var args []string
	args = append(args, "endorse")

	// 设置查询的请求参数
	req := chclient.Request{ChaincodeID: setup.Fabric.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(bill_no), []byte(endorseCmId), []byte(endorseAcct)}}
	// 执行查询
	response, err := setup.Fabric.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf("票据背书请求失败: %v\n", err)
	}

	return string(response.Payload), nil
}

// 根据待背书人证件号码, 查询当前用户的待背书票据
func (setup *FabricSetupService) FindWaitBills(endorserCmId string) ([]byte, error){
	var args []string
	args = append(args, "queryMyWaitBills")

	// 设置查询的请求参数
	req := chclient.Request{ChaincodeID: setup.Fabric.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(endorserCmId)}}
	// 执行查询
	response, err := setup.Fabric.Client.Query(req)
	if err != nil {
		return []byte{0x00}, fmt.Errorf("查询待背书票据失败: %v\n", err)
	}

	b := response.Payload

	return b[:], nil
}

// 票据背书签收
// args: 0 - Bill_No;  1 - endorserCmId(待背书人ID);	2 - endorserAcct(待背书人名称)
func (setup *FabricSetupService) EndorseAccept(bill_no string, endorseCmId string, endorseAcct string) (string, error){
	var args []string
	args = append(args, "accept")

	req := chclient.Request{ChaincodeID: setup.Fabric.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(bill_no), []byte(endorseCmId), []byte(endorseAcct)}}
	response, err := setup.Fabric.Client.Execute(req)

	if err != nil {
		return "", fmt.Errorf("票据背书签收失败: %v\n", err)
	}
	return string(response.Payload), nil
}

// 票据背书拒签(拒绝背书)
// args: 0 - bill_NO;	1 - endorserCmId(待背书人ID);	2 - endorserAcct(待背书人名称)
func (setup *FabricSetupService) EndorseReject(bill_no string, endorseCmId string, endorseAcct string) (string, error) {
	var args []string
	args = append(args, "reject")

	req := chclient.Request{ChaincodeID: setup.Fabric.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(bill_no), []byte(endorseCmId), []byte(endorseAcct)}}
	response, err := setup.Fabric.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf("票据背书拒签失败: %v\n", err)
	}
	return string(response.Payload), nil
}