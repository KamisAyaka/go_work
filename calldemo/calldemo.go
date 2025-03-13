// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// CalldemoMetaData contains all meta data concerning the Calldemo contract.
var CalldemoMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"getCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_count\",\"type\":\"uint256\"}],\"name\":\"setCount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// CalldemoABI is the input ABI used to generate the binding from.
// Deprecated: Use CalldemoMetaData.ABI instead.
var CalldemoABI = CalldemoMetaData.ABI

// Calldemo is an auto generated Go binding around an Ethereum contract.
type Calldemo struct {
	CalldemoCaller     // Read-only binding to the contract
	CalldemoTransactor // Write-only binding to the contract
	CalldemoFilterer   // Log filterer for contract events
}

// CalldemoCaller is an auto generated read-only Go binding around an Ethereum contract.
type CalldemoCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CalldemoTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CalldemoTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CalldemoFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CalldemoFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CalldemoSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CalldemoSession struct {
	Contract     *Calldemo         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CalldemoCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CalldemoCallerSession struct {
	Contract *CalldemoCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// CalldemoTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CalldemoTransactorSession struct {
	Contract     *CalldemoTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// CalldemoRaw is an auto generated low-level Go binding around an Ethereum contract.
type CalldemoRaw struct {
	Contract *Calldemo // Generic contract binding to access the raw methods on
}

// CalldemoCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CalldemoCallerRaw struct {
	Contract *CalldemoCaller // Generic read-only contract binding to access the raw methods on
}

// CalldemoTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CalldemoTransactorRaw struct {
	Contract *CalldemoTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCalldemo creates a new instance of Calldemo, bound to a specific deployed contract.
func NewCalldemo(address common.Address, backend bind.ContractBackend) (*Calldemo, error) {
	contract, err := bindCalldemo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Calldemo{CalldemoCaller: CalldemoCaller{contract: contract}, CalldemoTransactor: CalldemoTransactor{contract: contract}, CalldemoFilterer: CalldemoFilterer{contract: contract}}, nil
}

// NewCalldemoCaller creates a new read-only instance of Calldemo, bound to a specific deployed contract.
func NewCalldemoCaller(address common.Address, caller bind.ContractCaller) (*CalldemoCaller, error) {
	contract, err := bindCalldemo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CalldemoCaller{contract: contract}, nil
}

// NewCalldemoTransactor creates a new write-only instance of Calldemo, bound to a specific deployed contract.
func NewCalldemoTransactor(address common.Address, transactor bind.ContractTransactor) (*CalldemoTransactor, error) {
	contract, err := bindCalldemo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CalldemoTransactor{contract: contract}, nil
}

// NewCalldemoFilterer creates a new log filterer instance of Calldemo, bound to a specific deployed contract.
func NewCalldemoFilterer(address common.Address, filterer bind.ContractFilterer) (*CalldemoFilterer, error) {
	contract, err := bindCalldemo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CalldemoFilterer{contract: contract}, nil
}

// bindCalldemo binds a generic wrapper to an already deployed contract.
func bindCalldemo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CalldemoMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Calldemo *CalldemoRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Calldemo.Contract.CalldemoCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Calldemo *CalldemoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Calldemo.Contract.CalldemoTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Calldemo *CalldemoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Calldemo.Contract.CalldemoTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Calldemo *CalldemoCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Calldemo.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Calldemo *CalldemoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Calldemo.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Calldemo *CalldemoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Calldemo.Contract.contract.Transact(opts, method, params...)
}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(uint256)
func (_Calldemo *CalldemoCaller) GetCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Calldemo.contract.Call(opts, &out, "getCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(uint256)
func (_Calldemo *CalldemoSession) GetCount() (*big.Int, error) {
	return _Calldemo.Contract.GetCount(&_Calldemo.CallOpts)
}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(uint256)
func (_Calldemo *CalldemoCallerSession) GetCount() (*big.Int, error) {
	return _Calldemo.Contract.GetCount(&_Calldemo.CallOpts)
}

// SetCount is a paid mutator transaction binding the contract method 0xd14e62b8.
//
// Solidity: function setCount(uint256 _count) returns()
func (_Calldemo *CalldemoTransactor) SetCount(opts *bind.TransactOpts, _count *big.Int) (*types.Transaction, error) {
	return _Calldemo.contract.Transact(opts, "setCount", _count)
}

// SetCount is a paid mutator transaction binding the contract method 0xd14e62b8.
//
// Solidity: function setCount(uint256 _count) returns()
func (_Calldemo *CalldemoSession) SetCount(_count *big.Int) (*types.Transaction, error) {
	return _Calldemo.Contract.SetCount(&_Calldemo.TransactOpts, _count)
}

// SetCount is a paid mutator transaction binding the contract method 0xd14e62b8.
//
// Solidity: function setCount(uint256 _count) returns()
func (_Calldemo *CalldemoTransactorSession) SetCount(_count *big.Int) (*types.Transaction, error) {
	return _Calldemo.Contract.SetCount(&_Calldemo.TransactOpts, _count)
}
