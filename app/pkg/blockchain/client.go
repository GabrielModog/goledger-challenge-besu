package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	ethClient  *ethclient.Client
	contract   *bind.BoundContract
	address    common.Address
	privateKey *ecdsa.PrivateKey
}

func NewBlockchainClient(rpcURL, contractAddress, privateKey, abiPath string) (*Client, error) {
	ethClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ethereum client: %w", err)
	}

	address := common.HexToAddress(contractAddress)

	abiBytes, err := os.ReadFile(abiPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ABI: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	contract := bind.NewBoundContract(address, parsedABI, ethClient, ethClient, ethClient)

	return &Client{
		ethClient:  ethClient,
		contract:   contract,
		address:    address,
		privateKey: privKey,
	}, nil
}

func (c *Client) GetValue(ctx context.Context) (*big.Int, error) {
	var out []any
	err := c.contract.Call(&bind.CallOpts{Context: ctx}, &out, "get")
	if err != nil {
		return nil, fmt.Errorf("failed to call get: %w", err)
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("unexpected empty response from get")
	}

	val, ok := out[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected type in get return: %T", out[0])
	}

	return val, nil
}

func (c *Client) SetValue(ctx context.Context, value *big.Int) (string, error) {
	chainID, err := c.ethClient.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
	if err != nil {
		return "", fmt.Errorf("failed to create transactor: %w", err)
	}

	tx, err := c.contract.Transact(auth, "set", value)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return tx.Hash().Hex(), nil
}

func (c *Client) Close() {
	c.ethClient.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.ethClient.BlockNumber(ctx)
	return err
}
