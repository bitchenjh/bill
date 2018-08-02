package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

// 票据背书请求
// args: 0 - Bill_No;  1 - endorseCmId(待背书人ID);	2 - endorseAcct(待背书人名称)
func (t *BillChainCode) endorse(stub shim.ChaincodeStubInterface, args []string) peer.Response  {
	// 1. 检查参数长度是否为3(票据号码, 待背书人ID, 待背书人名称)
	if len(args) < 3 {
		res := GetRetString(1, "票据背书失败, 请求背书的参数不能少于3个")
		return shim.Error(res)
	}

	// 2. 根据票据号码获取票据状态
	bill, bl := t.getBill(stub, args[0])
	if !bl {
		res := GetRetString(1, "票据背书失败, 根据票据号码获取票据状态时发生错误")
		return shim.Error(res)
	}

	// 3. 检查待背书人与当前持票人是否为同一人
	if bill.HoldrCmID == args[1] {
		res := GetRetString(1, "票据背书失败, 待背书人与当前持票人不能为同一人")
		return shim.Error(res)
	}

	// 4. 获取票据历史变更数据
	resultIterator, err := stub.GetHistoryForKey(Bill_Prefix + args[0])
	if err != nil {
		res := GetRetString(1, "票据背书失败, 查询该票据的历史变更状态时发生错误")
		return shim.Error(res)
	}
	defer resultIterator.Close()

	// 5. 检查待背书人是否为票据历史持有人
	var hisBill Bill
	for resultIterator.HasNext(){
		historyData, err := resultIterator.Next()
		if err != nil {
			res := GetRetString(1, "票据背书失败, 处理变更历史数据时发生错误")
			return shim.Error(res)
		}

		json.Unmarshal(historyData.Value, &hisBill)
		if historyData.Value == nil{
			var emptyBill Bill
			hisBill = emptyBill
		}

		if hisBill.HoldrCmID == args[1]{
			res := GetRetString(1, "票据背书失败, 待背书人不应是票据历史持有人")
			return shim.Error(res)
		}

	}

	// 6.  更改票据信息与状态: 添加待背书人信息(证件号码与名称), 票据状态更改为待背书, 重置已拒绝背书人
	bill.WaitEndorseCmID = args[1]
	bill.WaitEndorseAcct = args[2]
	bill.State = BillInfo_State_EndorseWaitSign
	bill.RejectEndorseCmID = ""
	bill.RejectEndorseAcct = ""

	// 7. 保存票据
	_, bl = t.putBill(stub, bill)
	if !bl {
		res := GetRetString(1, "票据背书失败, 保存票据状态时发生错误")
		return shim.Error(res)
	}

	// 8. 保存以待背书人ID与票据号码构造的复合键, 以便待背书人批量查询. value为空即可
	holderNameBillNoIndexKey, err := stub.CreateCompositeKey(IndexName, []string{bill.WaitEndorseCmID, bill.BillInfoID})
	if err != nil {
		res := GetRetString(1, "创建待背书人ID与票据复合键失败")
		return shim.Error(res)
	}
	stub.PutState(holderNameBillNoIndexKey, []byte{0x00})

	// 9. 返回
	res := GetRetByte(0, "背书请求成功")
	return shim.Success(res)
}


// 票据背书签收
// args: 0 - Bill_No;  1 - endorseCmId(待背书人ID);	2 - endorseAcct(待背书人名称)
func (t *BillChainCode) accept(stub shim.ChaincodeStubInterface, args []string) peer.Response  {
	// 1. 检查参数长度是否为3(票据号码, 待背书人ID, 待背书人名称)
	if len(args) < 3 {
		res := GetRetString(1, "票据背书签收失败, 参数不能少于3个")
		return shim.Error(res)
	}

	// 2. 根据票据号码获取票据状态
	bill, bl := t.getBill(stub, args[0])
	if !bl {
		res := GetRetString(1, "票据背书签收失败, 根据票据号码查询对应票据状态时发生错误")
		return shim.Error(res)
	}

	// 3. 以前手持票人ID与票据号码构造复合键, 删除该key, 以便前手持票人无法再查到该票据
	holderNameBillNoIndexKey, err := stub.CreateCompositeKey(IndexName, []string{bill.HoldrCmID, bill.BillInfoID})
	if err != nil{
		res := GetRetString(1, "票据背书签收失败, 创建持票人ID与票据号码复合键时发生错误")
		return shim.Error(res)
	}
	stub.DelState(holderNameBillNoIndexKey)

	// 4. 更改票据信息与状态: 将当前持票人更改为待背书人(证件与名称), 票据状态更改为背书签收, 重置待背书人
	bill.HoldrCmID = args[1]
	bill.HoldrAcct = args[2]
	bill.State = BillInfo_State_EndorseSigned
	bill.WaitEndorseCmID = ""
	bill.WaitEndorseAcct = ""

	// 5. 保存票据
	_, bl = t.putBill(stub, bill)
	if !bl {
		res := GetRetString(1, "票据背书签收失败, 保存票据状态时发生错误")
		return shim.Error(res)
	}

	// 6. 返回
	res := GetRetByte(0, "票据背书签收成功")
	return shim.Success(res)
}

// 票据背书拒签(拒绝背书)
// args: 0 - bill_NO;	1 - endorseCmId(待背书人ID);	2 - endorseAcct(待背书人名称)
func (t *BillChainCode) reject(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 1. 检查参数长度是否为3(票据号码, 待背书人ID, 待背书人名称)
	if len(args) < 3 {
		res := GetRetString(1, "票据背书拒签失败, 参数不能少于3个")
		return shim.Error(res)
	}

	// 2. 根据票据号码查询对应的票据状态
	bill, bl := t.getBill(stub, args[0])
	if !bl {
		res := GetRetString(1, "票据背书拒签失败, 根据票据号码查询对应的票据状态时发生错误")
		return shim.Error(res)
	}

	// 3. 以待背书人ID及票据号码构造复合键, 从search中删除该key, 以便当前被背书人无法再次查询到该票据
	holderNameBillNoIndexKey, err := stub.CreateCompositeKey(IndexName, []string{args[1], bill.BillInfoID})
	if err != nil {
		res := GetRetString(1, "票据背书拒签失败, 以待背书人ID及票据号码构造复合键时发生错误")
		return shim.Error(res)
	}
	stub.DelState(holderNameBillNoIndexKey)

	// 4. 更改票据信息与状态: 将拒绝背书人更改为当前待背书人(证件号码与名称), 票据状态更改为背书拒绝, 重置待背书人
	bill.RejectEndorseCmID = args[1]
	bill.RejectEndorseAcct = args[2]
	bill.State = BillInfo_State_EndorseReject
	bill.WaitEndorseCmID = ""
	bill.WaitEndorseAcct = ""

	// 5. 保存票据状态
	_, bl = t.putBill(stub, bill)
	if !bl {
		res := GetRetString(1, "票据背书拒签失败, 保存票据状态时发生错误")
		return shim.Error(res)
	}

	// 6. 返回
	res := GetRetByte(0, "票据背书拒签成功")
	return shim.Success(res)
}
