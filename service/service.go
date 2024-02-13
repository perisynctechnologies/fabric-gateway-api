package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc/status"
)

type IService interface {
	InitLedger() error
	GetAllAssets() (any, error)
	CreateAsset(data string) (any, error)
	DeleteAsset(id string) (any, error)
	VoidAsset(id string) (any, error)
	ReleaseAsset(id string) (any, error)
	ExpireAsset(id string) (any, error)
	ReadAssetByID(id string) (any, error)
}

type Service struct {
	contract *client.Contract
}

func New(contract *client.Contract) IService {
	return &Service{
		contract: contract,
	}
}

// This type of transaction would typically only be run once by an application the first time it was started after its initial deployment.
// A new version of the chaincode deployed later would likely not need to run an "init" function.
func (s *Service) InitLedger() error {
	_, err := s.contract.SubmitTransaction("Init")
	if err != nil {
		return err
	}
	return nil
}

// Evaluate a transaction to query ledger state.
func (s *Service) GetAllAssets() (any, error) {
	res, err := s.contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		s.handleError("GetAllAssets:EvaluateTransaction err", err)
		return nil, err
	}
	return res, nil
}

// Evaluate a transaction by assetID to query ledger state.
func (s *Service) ReadAssetByID(id string) (any, error) {
	res, err := s.contract.EvaluateTransaction("ReadAsset", id)
	if err != nil {
		s.handleError("ReadAssetByID:EvaluateTransaction err", err)
		return nil, err
	}

	info := new(map[string]any)
	if err := json.Unmarshal(res, info); err != nil {
		return nil, err
	}

	return info, nil
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func (s *Service) CreateAsset(data string) (any, error) {
	res, err := s.contract.SubmitTransaction("CreateAsset", data)
	if err != nil {
		s.handleError("SubmitTransaction:CreateAsset", err)
		return nil, err
	}

	result := map[string]any{}
	if err := json.Unmarshal(res, &result); err != nil {
		log.Printf("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("response: %q", res)
		return nil, err
	}

	return result, nil
}

func (s *Service) DeleteAsset(id string) (any, error) {
	res, err := s.contract.SubmitTransaction("DeleteAsset", id)
	if err != nil {
		s.handleError("SubmitTransaction:DeleteAsset err:", err)
		return nil, err
	}

	return res, nil
}

func (s *Service) VoidAsset(data string) (any, error) {
	res, err := s.contract.SubmitTransaction("VodataAsset", data)
	if err != nil {
		s.handleError("SubmitTransaction:VoidAsset err:", err)
		return nil, err
	}

	result := map[string]any{}
	if err := json.Unmarshal(res, &result); err != nil {
		log.Printf("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("response: %q", res)
		return nil, err
	}

	return result, nil
}

func (s *Service) ReleaseAsset(data string) (any, error) {
	res, err := s.contract.SubmitTransaction("ReleaseAsset", data)
	if err != nil {
		s.handleError("SubmitTransaction:ReleaseAsset err:", err)
		return nil, err
	}

	result := map[string]any{}
	if err := json.Unmarshal(res, &result); err != nil {
		log.Printf("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("response: %q", res)
		return nil, err
	}

	return result, nil
}

func (s *Service) ExpireAsset(data string) (any, error) {
	res, err := s.contract.SubmitTransaction("ExpireAsset", data)
	if err != nil {
		s.handleError("SubmitTransaction:ExpireAsset err:", err)
		return nil, err
	}

	result := map[string]any{}
	if err := json.Unmarshal(res, &result); err != nil {
		log.Printf("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("response: %q", res)
		return nil, err
	}

	return result, nil
}

func (s *Service) handleError(tag string, err error) {
	fmt.Println(tag, "---")
	defer fmt.Println("---", tag)

	switch err := err.(type) {
	case *client.EndorseError:
		fmt.Printf("Endorse error for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)
	case *client.SubmitError:
		fmt.Printf("Submit error for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)
	case *client.CommitStatusError:
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Timeout waiting for transaction %s commit status: %s", err.TransactionID, err)
		} else {
			fmt.Printf("Error obtaining commit status for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)
		}
	case *client.CommitError:
		fmt.Printf("Transaction %s failed to commit with status %d: %s\n", err.TransactionID, int32(err.Code), err)
	default:
		panic(fmt.Errorf("unexpected error type %T: %w", err, err))
	}

	// Any error that originates from a peer or orderer node external to the gateway will have its details
	// embedded within the gRPC status error. The following code shows how to extract that.
	statusErr := status.Convert(err)

	details := statusErr.Details()
	if len(details) > 0 {
		fmt.Println("Error Details:")

		for _, detail := range details {
			switch detail := detail.(type) {
			case *gateway.ErrorDetail:
				fmt.Printf("- address: %s, mspId: %s, message: %s\n", detail.Address, detail.MspId, detail.Message)
			}
		}
	}
}
