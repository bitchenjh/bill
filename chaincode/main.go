package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"fmt"
)

type BillChainCode struct {

}

func (t *BillChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *BillChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "issue" {
		return t.issue(stub, args)  //發布票據
	}else if function == "queryMyBills" {
		return t.queryMyBills(stub, args)
	}else if function == "queryBillByNo" {
		return t.queryBillByNo(stub, args)
	}else  if function == "queryMyWaitBills" {
		return t.queryMyWaitBills(stub, args)
	}else if function == "endorse" {
		return t.endorse(stub, args)
	}else if function == "accept" {
		return t.accept(stub, args)
	}else if function == "reject" {
		return t.reject(stub, args)
	}else if function == "delete"{
		return t.delete(stub,args)
	}

	return shim.Error("指定的函数名称错误")
}

func main()  {
	err := shim.Start(new(BillChainCode))
	if err != nil {
		fmt.Println("启动链码错误: ", err)
	}
}