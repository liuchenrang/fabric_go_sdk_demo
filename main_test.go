package main_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	mspClient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	//"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/lookup"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/test/integration"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/op/go-logging"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
	OrgName    = "Org1"
	ChannelID  = "mychannel"
)

func Test_main(t *testing.T) {
	var e = logging.MustGetLogger("example")
	sdk, err := fabsdk.New(config.FromFile(basepath + "/ca.yaml"))
	if err != nil {
		e.Errorf("Failed to create new resource management client: %s", err)
	}
	org1AdminUser := "admin"
	org1Name := "Org1"
	//prepare contexts
	org1AdminChannelContext := sdk.ChannelContext(ChannelID, fabsdk.WithUser(org1AdminUser), fabsdk.WithOrg(org1Name))
	// Ledger client
	ledgerClient, errr := ledger.New(org1AdminChannelContext)

	if errr != nil {
		e.Errorf("Failed to create new resource management client: %s", err)
	}
	urls := make([]string, 1)
	urls[0] = "peer0.org1.example.com"
	bci, errrr := ledgerClient.QueryInfo(ledger.WithTargetURLs(urls...))
	if errrr != nil {
		t.Fatalf("QueryInfo return error: %v", errrr)
	}
	println("bci %d", 444)
	println("bci ", bci.BCI.String())
}

func Test_Chaincode(t *testing.T) {
	var e = logging.MustGetLogger("example")
	sdk, err := fabsdk.New(config.FromFile(basepath + "/ca.yaml"))
	if err != nil {
		e.Errorf("Failed to create new resource management client: %s", err)
	}
	org1AdminUser := "Admin"
	//prepare contexts
	clientContext := sdk.Context(fabsdk.WithUser(org1AdminUser), fabsdk.WithOrg(OrgName))
	urls := make([]string, 1)
	urls[0] = "peer0.org1.example.com"
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		t.Fatalf("Failed to create channel management client: %s", err)
	}
	ops1 := resmgmt.WithTargetURLs(urls...)
	resp, error := resMgmtClient.QueryChannels(ops1)
	if error != nil {
		t.Fatalf("Failed to create channel management client: %s", err)
	}
	println(resp.String())

	ccPkg, cc := gopackager.NewCCPackage("chaincode/example_cc", os.Getenv("GOPATH"))
	if cc != nil {
		t.Fatal(cc)
	}
	ccID := "e2eExampleCC"
	// Install example cc to org peers

	if false {
		installCCReq := resmgmt.InstallCCRequest{Name: ccID, Path: os.Getenv("GOPATH") + "/src/chaincode/example_cc", Version: "0", Package: ccPkg}
		println(installCCReq.Name, installCCReq.Path, installCCReq.Version)
		_, err = resMgmtClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))

		if err != nil {
			t.Fatal(err)
			os.Exit(1)
		}
	}
	if false {
		ccPolicy := &common.SignaturePolicyEnvelope{}
		_, err = resMgmtClient.InstantiateCC(
			ChannelID,
			resmgmt.InstantiateCCRequest{Name: ccID, Path: os.Getenv("GOPATH") + "/src/chaincode/example_cc", Version: "0", Args: integration.ExampleCCInitArgs(), Policy: ccPolicy},
			resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	clientChannelContext := sdk.ChannelContext(ChannelID, fabsdk.WithUser("Admin"), fabsdk.WithOrg(OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	channel.WithTargetURLs(urls...)

	client, err := channel.New(clientChannelContext)
	if err != nil {
		t.Fatalf("Failed to create new channel client: %s", err)
	}
	var queryArgs = [][]byte{[]byte("Michel")}

	response, err := client.Query(channel.Request{ChaincodeID: "fabcar", Fcn: "queryCar", Args: queryArgs},
		channel.WithRetry(retry.DefaultChClientOpts), channel.WithTargetURLs(urls...))
	if err != nil {
		t.Fatalf("Failed to query funds: %s", err)
	}
	value := response.Payload
	println("value")
	println(value)
	println(response.ChaincodeStatus)

}

func Test_InstantiateCC(t *testing.T) {
	var e = logging.MustGetLogger("example")
	sdk, err := fabsdk.New(config.FromFile(basepath + "/ca.yaml"))
	if err != nil {
		e.Errorf("Failed to create new resource management client: %s", err)
	}
	org1AdminUser := "Admin"
	//prepare contexts
	clientContext := sdk.Context(fabsdk.WithUser(org1AdminUser), fabsdk.WithOrg(OrgName))
	urls := make([]string, 1)
	urls[0] = "peer0.org1.example.com"
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		t.Fatalf("Failed to create channel management client: %s", err)
	}
	ops1 := resmgmt.WithTargetURLs(urls...)
	resp, error := resMgmtClient.QueryChannels(ops1)
	if error != nil {
		t.Fatalf("Failed to create channel management client: %s", err)
	}
	println(resp.String())
	ccID := "chaincode_example02"
	ccPkg, cc := gopackager.NewCCPackage("chaincode/chaincode_example02", os.Getenv("GOPATH") )
	if cc != nil {
		t.Fatal(cc)
	}
	// Install example cc to org peers

	if true {
		installCCReq := resmgmt.InstallCCRequest{Name: ccID, Path: "chaincode/chaincode_example02", Version: "0", Package: ccPkg}
		println(installCCReq.Name, installCCReq.Path, installCCReq.Version)
		_, err = resMgmtClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))

		if err != nil {
			t.Fatal(err)
			os.Exit(1)
		}
	}
	if true {
		ccPolicy := &common.SignaturePolicyEnvelope{}
		_, err = resMgmtClient.InstantiateCC(
			ChannelID,
			resmgmt.InstantiateCCRequest{Name: ccID, Path: "chaincode_example02", Version: "0", Args: integration.ExampleCCInitArgs(), Policy: ccPolicy},
			resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		)
		if err != nil {
			t.Fatal(err)
		}
	}

}
func Test_Invoke(t *testing.T) {
	var e = logging.MustGetLogger("example")
	sdk, err := fabsdk.New(config.FromFile(basepath + "/ca.yaml"))
	if err != nil {
		e.Errorf("Failed to create new resource management client: %s", err)
	}
	org1AdminUser := "Admin"
	//prepare contexts
	clientContext := sdk.Context(fabsdk.WithUser(org1AdminUser), fabsdk.WithOrg(OrgName))
	urls := make([]string, 1)
	urls[0] = "peer0.org1.example.com"
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		t.Fatalf("Failed to create channel management client: %s", err)
	}
	ops1 := resmgmt.WithTargetURLs(urls...)
	resp, error := resMgmtClient.QueryChannels(ops1)
	if error != nil {
		t.Fatalf("Failed to create channel management client: %s", err)
	}
	println(resp.String())

	clientChannelContext := sdk.ChannelContext(ChannelID, fabsdk.WithUser("Admin"), fabsdk.WithOrg(OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	channel.WithTargetURLs(urls...)

	client, err := channel.New(clientChannelContext)
	if err != nil {
		t.Fatalf("Failed to create new channel client: %s", err)
	}
	var queryArgs = [][]byte{[]byte("a")}

	response, err := client.Query(channel.Request{ChaincodeID: "chaincode_example02", Fcn: "query", Args: queryArgs},
		channel.WithRetry(retry.DefaultChClientOpts), channel.WithTargetURLs(urls...))
	if err != nil {
		t.Fatalf("Failed to query funds: %s", err)
	}
	value := response.Payload
	str := string(value[:])
	println("value")
	println(str)
}
func Test_CaEnroll(t *testing.T) {
	backend, err := config.FromFile(basepath + "/ca.yaml")()
	if err != nil {
		panic(err)
	}
	//Override ca matchers for this test
	customBackend := getCustomBackend(backend)

	configProvider := func() (core.ConfigBackend, error) {
		return customBackend, nil
	}

	// Instantiate the SDK
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		panic(fmt.Sprintf("SDK init failed: %v", err))
	}

	// configBackend, err := sdk.Config()
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to get config: %v", err))
	// }

	if err != nil {
		t.Fatalf("failed to create CA client: %v", err)
	}
	ctxProvider := sdk.Context()
	msp, errr := mspClient.New(ctxProvider)
	if errr != nil {
		t.Fatalf("failed to create CA client: %v", errr)
	}
	mspClient.WithOrg("Org1")(msp)

	err = msp.Enroll("admin", mspClient.WithSecret("adminpw"))
	if err != nil {
		t.Fatalf("Enroll should return error for empty enrollment ID" + err.Error())
	}

}
//
func Test_Ca_Regist(t *testing.T) {
	backend, err := config.FromFile(basepath + "/ca.yaml")()
	if err != nil {
		panic(err)
	}
	//Override ca matchers for this test
	customBackend := getCustomBackend(backend)

	configProvider := func() (core.ConfigBackend, error) {
		return customBackend, nil
	}

	// Instantiate the SDK
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		panic(fmt.Sprintf("SDK init failed: %v", err))
	}


	if err != nil {
		t.Fatalf("failed to create CA client: %v", err)
	}
	ctxProvider := sdk.Context()
	msp, errr := mspClient.New(ctxProvider)
	if errr != nil {
		t.Fatalf("failed to create CA client: %v", errr)
	}
	mspClient.WithOrg("Org1")(msp)
	register := mspClient.RegistrationRequest{}
	register.Name = "aliniuniu";
	register.Affiliation = "org1.department1";
	register.Type = "user";
	register.CAName = "ca.org1.example.com"


	secret, errrr  := msp.Register(&register)
	if errrr != nil {
		t.Fatalf("Register should return error for empty enrollment ID" + errrr.Error())
	}else{
		println(secret)
	}
	//xinghuo LRtJKLdsFaqM
}
func Test_EnrolleUser(t *testing.T){
	backend, err := config.FromFile(basepath + "/ca.yaml")()
	if err != nil {
		panic(err)
	}
	//Override ca matchers for this test
	customBackend := getCustomBackend(backend)

	configProvider := func() (core.ConfigBackend, error) {
		return customBackend, nil
	}

	// Instantiate the SDK
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		panic(fmt.Sprintf("SDK init failed: %v", err))
	}

	// configBackend, err := sdk.Config()
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to get config: %v", err))
	// }

	if err != nil {
		t.Fatalf("failed to create CA client: %v", err)
	}
	ctxProvider := sdk.Context()
	msp, errr := mspClient.New(ctxProvider)
	if errr != nil {
		t.Fatalf("failed to create CA client: %v", errr)
	}
	mspClient.WithOrg("Org1")(msp)
	err = msp.Enroll("alibaba", mspClient.WithSecret("sITCesVszqWo"))
	if err != nil {
		t.Errorf("eee %s",err.Error())
	}
}
func Test_InvokeWithUser(t *testing.T) {
	var e = logging.MustGetLogger("example")
	sdk, err := fabsdk.New(config.FromFile(basepath + "/ca.yaml"))
	if err != nil {
		e.Errorf("Failed to create new resource management client: %s", err)
	}
	org1Name := "Org1"
	org1UserName  := "alibaba";
	//prepare contexts
	clientContext := sdk.Context(fabsdk.WithUser(org1UserName), fabsdk.WithOrg(org1Name))
	urls := make([]string, 1)
	urls[0] = "peer0.org1.example.com"
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		t.Fatalf("Failed to create channel management client: %s", err)
	}
	ops1 := resmgmt.WithTargetURLs(urls...)
	resp, error := resMgmtClient.QueryChannels(ops1)
	if error != nil {
		t.Fatalf("Failed to create channel management client: %s", err)
	}
	println(resp.String())

	clientChannelContext := sdk.ChannelContext(ChannelID, fabsdk.WithUser(org1UserName), fabsdk.WithOrg(OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	channel.WithTargetURLs(urls...)

	client, err := channel.New(clientChannelContext)
	if err != nil {
		t.Fatalf("Failed to create new channel client: %s", err)
	}
	var queryArgs = [][]byte{[]byte("a")}

	response, err := client.Query(channel.Request{ChaincodeID: "chaincode_example02", Fcn: "query", Args: queryArgs},
		channel.WithRetry(retry.DefaultChClientOpts), channel.WithTargetURLs(urls...))
	if err != nil {
		t.Fatalf("Failed to query funds: %s", err)
	}
	value := response.Payload
	str := string(value[:])
	println("value")
	println(str)
}
func getCustomBackend(backend core.ConfigBackend) *lookup.ConfigLookup {
	networkConfig := fab.NetworkConfig{}
	configLookup := lookup.New(backend)
	configLookup.UnmarshalKey("certificateAuthorities", &networkConfig.CertificateAuthorities)
	return configLookup;
}
