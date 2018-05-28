package main

import (
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"fmt"
	"encoding/json"
)

const (
	ADMIN_USERID = "Admin"
	USER1_USERID = "User1"
)

type ChainCodeSpec struct {
	channelClient apitxn.ChannelClient
	userId        string
	chaincodeId   string
}

func Initialize(channelId, chaincodeId, userId string) (*ChainCodeSpec, error) {
	configFile := "/development/go_develop/testsystem/src/conf/fastapp-sdk-config.yaml"
	chnlClient, err := getChannelClient(channelId, configFile, userId)
	if err != nil {
		return nil, err
	}
	return &ChainCodeSpec{channelClient: chnlClient, chaincodeId: chaincodeId, userId: userId}, nil
}

func getChannelClient(channelId, configFile, userId string) (apitxn.ChannelClient, error) {
	var chClient apitxn.ChannelClient
	sdkOptions := fabapi.Options{
		ConfigFile: configFile,
	}

	sdk, err := fabapi.NewSDK(sdkOptions)
	if err != nil {
		fmt.Println("ERROR: Failed to create new SDK instance.", err)
		return nil, err
	}

	chClient, err = sdk.NewChannelClient(channelId, userId)
	if err != nil {
		fmt.Println("ERROR: Failed to create new channel client.", err)
		return nil, err
	}

	return chClient, err
}

func (ccs *ChainCodeSpec) ChaincodeInvoke(fcn string, chaincodeArgs [][]byte) (responsePayload []byte, err error) {
	_, err = ccs.channelClient.ExecuteTx(apitxn.ExecuteTxRequest{ChaincodeID: ccs.chaincodeId, Fcn: fcn, Args: chaincodeArgs})
	if err != nil {
		fmt.Println("Error in executing trasaction: %s", err.Error())
		return nil, err
	}
	return nil, nil
}

func (ccs *ChainCodeSpec) ChaincodeQuery(fcn string, chaincodeArgs [][]byte) (responsePayload []byte, err error) {
	value, err := ccs.channelClient.Query(apitxn.QueryRequest{ChaincodeID: ccs.chaincodeId, Fcn: fcn, Args: chaincodeArgs})
	if err != nil {
		fmt.Println("ERROR: Failed to invoke query function of the chaincode: ", err)
		return nil, err
	}
	return value, err
}

func (ccs *ChainCodeSpec) Close() {
	ccs.channelClient.Close()
}

func main() {
	channelId := "testchannel"
	chaincodeId := "chaicodetest"

	ccs, err := Initialize(channelId, chaincodeId, USER1_USERID)
	if err != nil {
		fmt.Println(err)
	}

	defer ccs.Close()

	var chaincodeArgs [][]byte
	chaincodeArgs = append(chaincodeArgs, []byte("testproject"))
	chaincodeArgs = append(chaincodeArgs, []byte("testmodule"))
	chaincodeArgs = append(chaincodeArgs, []byte("test2"))
	chaincodeArgs = append(chaincodeArgs, []byte("test2_desp"))

	result, err := ccs.ChaincodeInvoke("create_testcases", chaincodeArgs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Invoke create test case result: ", result)

	type TestCaseNames struct {
		TestCaseNameList []string `json:"testcasenames"`
	}

	var testcase_list TestCaseNames
	var chaincodeArg1s [][]byte
	chaincodeArg1s = append(chaincodeArg1s, []byte("testproject"))
	chaincodeArg1s = append(chaincodeArg1s, []byte("testmodule"))
	resultQ, errQ := ccs.ChaincodeInvoke("query_testcases", chaincodeArg1s)
	if errQ != nil {
		fmt.Println(err.Error())
	}

	errU := json.Unmarshal(resultQ, testcase_list)
	if errU != nil {
		fmt.Println(errU.Error())
	}

	fmt.Println(testcase_list.TestCaseNameList)
}
