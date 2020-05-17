package leveldb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ====CHAINCODE EXECUTION SAMPLES (CLI) ==================

// peer chaincode instantiate -n queue -v 0 -c '{"Args":[]}' -C myc
//
// peer chaincode invoke -n mycc -c '{"Args":["InitLedger"]}' -C myc
//
// peer chaincode invoke -n mycc -c '{"Args":["Get", "1589702933-757936000"]}' -C myc
//
// peer chaincode invoke -n mycc -c '{"Args":["Update", "1589702933-757936000", "{\"country\":\"RU\"}"]}' -C myc
//
// peer chaincode invoke -n mycc -c '{"Args":["Delete", "1589702933-757936000"]}' -C myc
//
// peer chaincode invoke -n mycc -c '{"Args":["GetAll"]}' -C myc
//
// peer chaincode invoke -n mycc -c '{"Args":["GetRange", "0", "1558080533-00000000"]}' -C myc
// peer chaincode invoke -n mycc -c '{"Args":["GetRange", "0", "1305619733-758090001"]}' -C myc
//
//
//  ==== START QUERY ====
// peer chaincode invoke -n mycc -c '{"Args":["Query", "from=0&to=1558080533-00000000&sort=country&filter=country=BY"]}' -C myc
// peer chaincode invoke -n mycc -c '{"Args":["Query", "sort=country"]}' -C myc
// peer chaincode invoke -n mycc -c '{"Args":["Query", "filter=country=RU2"]}' -C myc
//  ==== END QUERY ====
//
// push empty string
// peer chaincode invoke -n mycc -c '{"Args":["PushBack", ""]}' -C myc
//
// push empty json string
// peer chaincode invoke -n mycc -c '{"Args":["PushBack", "{}"]}' -C myc
//
// push extra json context. Take look on character escaping
// peer chaincode invoke -n mycc -c '{"Args":["PushBack", "{\"country\":\"BY\"}"]}' -C myc
//
// access first element
// peer chaincode invoke -n mycc -c '{"Args":["Front"]}' -C myc
//
// access last element
// peer chaincode invoke -n mycc -c '{"Args":["Back"]}' -C myc
//
// get last element and remove them last element
// peer chaincode invoke -n mycc -c '{"Args":["Pop"]}' -C myc
//
// get last element and remove them last element
// peer chaincode invoke -n mycc -c '{"Args":["Swap","1305619733-758090000", "1337242133-758089000"]}' -C myc
//


// WithUniqueProperty smart contract /codechain/ which internally don't handle unique entity property
type SimpleQueueContract struct {
	contractapi.Contract
}

// just example using composite key, for us this is not suitable as we not use search via prefix.
func (s *SimpleQueueContract) compositeKey(stub shim.ChaincodeStubInterface) (string, error) {
	t := time.Now()
	ut := fmt.Sprintf("%d", t.Unix())
	ns := fmt.Sprintf("%d", t.Nanosecond())

	key, err := stub.CreateCompositeKey(ut, []string{ns})
	if err != nil {
		return "", fmt.Errorf("create composite key [%s-%s] error: %w", ut, ns, err)
	}

	return key, nil
}

// InitLedger adds a base set of assets to the ledger
func (s *SimpleQueueContract) InitLedger(ctx contractapi.TransactionContextInterface) ([]Query, error) {
	list := []SimpleQueue{
		// 1589702933-757936000
		{mustParse("2020-05-17T11:08:53.757936+03:00"), map[string]interface{}{"country": "BY"}},
		// 1558080533-758077000
		{mustParse("2019-05-17T11:08:53.758077+03:00"), map[string]interface{}{"country": "RU"}},
		// 1526544533-758079000
		{mustParse("2018-05-17T11:08:53.758079+03:00"), map[string]interface{}{"country": "UA"}},
		// 1495008533-758081000
		{mustParse("2017-05-17T11:08:53.758081+03:00"), map[string]interface{}{"country": "BY"}},
		// 1463472533-758082000
		{mustParse("2016-05-17T11:08:53.758082+03:00"), map[string]interface{}{"country": "BY", "num": 10_000_000}},
		// 1431850133-758084000
		{mustParse("2015-05-17T11:08:53.758084+03:00"), map[string]interface{}{"country": "UA"}},
		// 1400314133-758086000
		{mustParse("2014-05-17T11:08:53.758086+03:00"), map[string]interface{}{"country": "BY"}},
		// 1368778133-758087000
		{mustParse("2013-05-17T11:08:53.758087+03:00"), map[string]interface{}{"country": "RU2"}},
		// 1337242133-758089000
		{mustParse("2012-05-17T11:08:53.758089+03:00"), map[string]interface{}{"country": "BY"}},
		// 1305619733-758090000
		{mustParse("2011-05-17T11:08:53.75809+03:00"), map[string]interface{}{"country": "UA"}},
	}

	res := make([]Query, len(list))

	for i := range list {
		data, err := json.Marshal(list[i])
		if err != nil {
			return nil, fmt.Errorf("marshal data error: %w", err)
		}

		res[i].Object = list[i]
		res[i].Key = TimedKey(list[i].Time)

		if err := ctx.GetStub().PutState(res[i].Key, data); err != nil {
			return nil, fmt.Errorf("write state error: %w", err)
		}
	}

	return res, nil
}

// Get extract existing queue element by  it's key
func (s *SimpleQueueContract) Get(ctx contractapi.TransactionContextInterface, key string) (*Query, error) {
	v, err := ctx.GetStub().GetState(key)
	switch {
	case err != nil:
		return nil, fmt.Errorf("error extracting object with provided key: %w", err)
	case v == nil:
		return nil, fmt.Errorf("asset with key %s not exists", key)
	}

	q := &Query{Key: key}

	if err = json.Unmarshal(v, &q.Object); err != nil {
		return nil, fmt.Errorf("unmarshal old asset: %w", err)
	}

	return q, nil
}

// Update existing queue element by  it's key
// @js - expect correct JSON valid extra context data. Can be empty
func (s *SimpleQueueContract) Update(ctx contractapi.TransactionContextInterface, key string, js string) (*Query, error) {
	old, err := s.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(js) > 0 {
		if err = json.Unmarshal([]byte(js), &old.Object.Context); err != nil {
			return nil, fmt.Errorf("unmarshal extra context data: %w", err)
		}
	}

	blob, err := json.Marshal(old.Object)
	if err != nil {
		return nil, fmt.Errorf("marshal updated object: %w", err)
	}

	if err = ctx.GetStub().PutState(old.Key, blob); err != nil {
		return nil, fmt.Errorf("save state: %w", err)
	}

	return old, nil
}

// Delete asset by key
func (s *SimpleQueueContract) Delete(ctx contractapi.TransactionContextInterface, key string) error {
	_, err := s.Get(ctx, key)
	if err != nil {
		return err
	}

	if err = ctx.GetStub().DelState(key); err != nil {
		return fmt.Errorf("delete object: %w", err)
	}

	return nil
}

// GetAll list of queue
// very expensive operation which read all queue till last element
func (s *SimpleQueueContract) GetAll(ctx contractapi.TransactionContextInterface) (res []Query, err error) {
	return s.GetRange(ctx, "", TimedKey(time.Now()))
}

// GetRange get range [from, to)
func (s *SimpleQueueContract) GetRange(ctx contractapi.TransactionContextInterface, from, to string) (res SimpleQuery, err error) {
	if to == "" {
		to = TimedKey(time.Now())
	}

	// support backport extraction
	if to < from{
		from, to = to, from
	}

	itr, err := ctx.GetStub().GetStateByRange(from, to)
	if err != nil {
		return nil, fmt.Errorf("can't get range state")
	}

	for itr.HasNext() {
		i, err := itr.Next()
		if err != nil {
			return nil, fmt.Errorf("next result error: %w", err)
		}

		obj := SimpleQueue{}
		if err = json.Unmarshal(i.Value, &obj); err != nil {
			return nil, fmt.Errorf("unmarshal error: %w", err)
		}

		res = append(res, Query{i.Key, obj})
	}

	return res, nil
}

// Query extract list of element using operation query
// Supported operations uses url query syntax and support followed arguments:
// @from - select from which key should performed result extraction. Empty uses as from beggining
// @to - select to which key should be performed range extraction. (provided value excluded). Empty till NOW
// @Filter - extra context filtering result. Uses = separator and support only equation.
//  example: Filter=country=RU
// @Sort - order result with some provided context field, if field not exists result will be in the end of slice
//
// ascending example: Sort=country
// descending example: Sort=-country
//
// Sort require all context data provided with type consistency
func (s *SimpleQueueContract) Query(ctx contractapi.TransactionContextInterface, operation string) (res []Query, err error) {
	op, err := ParseOperation(operation)
	if err != nil {
		return nil, fmt.Errorf("read operation parameter error: %w", err)
	}

	v, err := s.GetRange(ctx, op.Selector.From, op.Selector.To)
	if err != nil {
		return nil, fmt.Errorf("extract range error: %w", err)
	}

	if op.Filter.Key != "" {
		v, err = v.Filter(op.Filter)
		if err != nil {
			return nil, fmt.Errorf("filtering error: %w", err)
		}
	}

	if op.Sort.Field != "" {
		v, err = v.Sort(op.Sort)
	}

	return v, err
}

// PushBack create new queue element and put it to the end of queue
// @js - expect correct JSON valid extra context data. Can be empty
func (s *SimpleQueueContract) PushBack(ctx contractapi.TransactionContextInterface, js string) (*Query, error) {
	item := NewSimpleQueue()

	if len(js) > 0 {
		if err := json.Unmarshal([]byte(js), &item.Context); err != nil {
			return nil, fmt.Errorf("parameter has bad JSON format")
		}
	}

	// test composition key how it uses
	fmt.Println(s.compositeKey(ctx.GetStub()))

	out := &Query{Key: TimedKey(item.Time), Object: item}

	blob, err := json.Marshal(&item)
	if err != nil {
		return nil, fmt.Errorf("marhaling error :%w", err)
	}

	if err = ctx.GetStub().PutState(out.Key, blob); err != nil {
		return nil, fmt.Errorf("")
	}

	return out, nil
}

// Front extract first element of queue
// very expensive operation which read all queue till last element
func (s *SimpleQueueContract) Front(ctx contractapi.TransactionContextInterface) (*Query, error) {
	itr, err := ctx.GetStub().GetStateByRange("", TimedKey(time.Now()))
	if err != nil {
		return nil, fmt.Errorf("can't get range state")
	}

	if !itr.HasNext() {
		return nil, nil
	}

	i, err := itr.Next()
	if err != nil {
		return nil, fmt.Errorf("next result error: %w", err)
	}

	res := SimpleQueue{}
	if err = json.Unmarshal(i.Value, &res); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &Query{i.Key, res}, nil
}

// Back extract last element of queue
// very expensive operation which read all queue till last element
func (s *SimpleQueueContract) Back(ctx contractapi.TransactionContextInterface) (*Query, error) {
	itr, err := ctx.GetStub().GetStateByRange("", TimedKey(time.Now()))
	if err != nil {
		return nil, fmt.Errorf("can't get range state")
	}

	if !itr.HasNext() {
		return nil, nil
	}

	for {
		i, err := itr.Next()
		if err != nil {
			return nil, fmt.Errorf("next result error: %w", err)
		}

		// we interested only in last element
		if itr.HasNext() {
			continue
		}

		res := SimpleQueue{}
		if err = json.Unmarshal(i.Value, &res); err != nil {
			return nil, fmt.Errorf("unmarshal error: %w", err)
		}

		return &Query{i.Key, res}, nil
	}
}

// Pop extract and remove last element of queue
// very expensive operation which read all queue till last element
func (s *SimpleQueueContract) Pop(ctx contractapi.TransactionContextInterface) (*Query, error) {
	q, err := s.Back(ctx)
	if err != nil {
		return nil, err
	}

	if err = ctx.GetStub().DelState(q.Key); err != nil {
		return nil, fmt.Errorf("delete state key %q error: %w", q.Key, err)
	}

	return q, nil
}

// Swap replace between 2 elements their context
// Swap performed only with context data
func (s *SimpleQueueContract) Swap(ctx contractapi.TransactionContextInterface, a, b string) (bool, error) {
	if a == b {
		return true, nil
	}

	first, err := s.Get(ctx, a)
	if err != nil {
		return false, fmt.Errorf("first element %q exstraction error: %w", a, err)
	}

	second, err := s.Get(ctx, b)
	if err != nil {
		return false, fmt.Errorf("second element %q extraction error: %w", b, err)
	}

	first.Object.Time, second.Object.Time = second.Object.Time, first.Object.Time

	firstBlob, err := json.Marshal(&first.Object)
	if err != nil {
		return false, fmt.Errorf("marshal element %q error: %w", a, err)
	}

	secondBlob, err := json.Marshal(&second.Object)
	if err != nil {
		return false, fmt.Errorf("marshal element %q error: %w", b, err)
	}

	fmt.Println(second.Key, string(firstBlob))

	if err = ctx.GetStub().PutState(second.Key, firstBlob); err != nil {
		return false, fmt.Errorf("key %q put context error: %w", second.Key, err)
	}

	// is that ROLLBACK previous operation
	if err = ctx.GetStub().PutState(first.Key, secondBlob); err != nil {
		return false, fmt.Errorf("key %q put first context error: %w", first.Key, err)
	}

	return true, nil
}
