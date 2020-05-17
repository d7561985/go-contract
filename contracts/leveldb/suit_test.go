// +build unit

package leveldb

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Init initializes the chaincode
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

type Suite struct {
	suite.Suite
	stub     *shimtest.MockStub
	contract *SimpleQueueContract
	ctx      *contractapi.TransactionContext
}

func (s *Suite) SetupTest() {
	uu := uuid.NewV4().String()

	fmt.Println("start TX:", uu)

	s.stub.MockTransactionStart(uu)
}

func (s *Suite) TearDownSuite() {
	fmt.Println("stop TX:", s.stub.TxID)

	s.stub.MockTransactionEnd(s.stub.TxID)
}

func (s *Suite) SetupSuite() {
	s.contract = new(SimpleQueueContract)
	s.stub = shimtest.NewMockStub("levelDB", new(SimpleChaincode))

	s.ctx = new(contractapi.TransactionContext)
	s.ctx.SetStub(s.stub)
}

func TestSimpleQueueContract(t *testing.T) {
	suite.Run(t, new(Suite))

}
