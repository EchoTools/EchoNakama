package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"

	"echo-nakama/server/services"
	"echo-nakama/server/services/login"
)

const (
	OK                  = 0  // OK indicates a successful operation.
	CANCELED            = 1  // CANCELED indicates the operation was canceled.
	UNKNOWN             = 2  // UNKNOWN indicates an unknown error occurred.
	INVALID_ARGUMENT    = 3  // INVALID_ARGUMENT indicates an invalid argument was provided.
	DEADLINE_EXCEEDED   = 4  // DEADLINE_EXCEEDED indicates the operation exceeded the deadline.
	NOT_FOUND           = 5  // NOT_FOUND indicates the requested resource was not found.
	ALREADY_EXISTS      = 6  // ALREADY_EXISTS indicates the resource already exists.
	PERMISSION_DENIED   = 7  // PERMISSION_DENIED indicates the operation was denied due to insufficient permissions.
	RESOURCE_EXHAUSTED  = 8  // RESOURCE_EXHAUSTED indicates the resource has been exhausted.
	FAILED_PRECONDITION = 9  // FAILED_PRECONDITION indicates a precondition for the operation was not met.
	ABORTED             = 10 // ABORTED indicates the operation was aborted.
	OUT_OF_RANGE        = 11 // OUT_OF_RANGE indicates a value is out of range.
	UNIMPLEMENTED       = 12 // UNIMPLEMENTED indicates the operation is not implemented.
	INTERNAL            = 13 // INTERNAL indicates an internal server error occurred.
	UNAVAILABLE         = 14 // UNAVAILABLE indicates the service is currently unavailable.
	DATA_LOSS           = 15 // DATA_LOSS indicates a loss of data occurred.
	UNAUTHENTICATED     = 16 // UNAUTHENTICATED indicates the request lacks valid authentication credentials.
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {

	if err := initializer.RegisterRpc("relay/loginrequest", LoginRequestRpc); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	// Register the LinkTicket Index that prevents multiple LinkTickets with the same device_id_str
	name := login.LINKTICKET_INDEX
	collection := login.LINKTICKET_COLLECTION
	key := ""                                               // Set to empty string to match all keys instead
	fields := []string{"xplatform_id_str", "device_id_str"} // index on these fields
	maxEntries := 1000000
	indexOnly := false

	err := initializer.RegisterStorageIndex(name, collection, key, fields, maxEntries, indexOnly)
	if err != nil {
		logger.Error("Unable to register storage index: %v", err)
		return err
	}

	logger.Info("Initialized module.")
	return nil
}

func LoginRequestRpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	// Parse the payload into a LoginRequest object

	var request login.LoginRequest
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		logger.WithField("err", err).Error("Unable to unmarshal payload")
		return "", runtime.NewError("Unable to unmarshal payload", INVALID_ARGUMENT)
	}

	serviceContext := &services.ServiceContext{
		Ctx:          ctx,
		Logger:       logger,
		DbConnection: db,
		NakamaModule: nk,
	}

	// Process the login request
	success, nkerr := login.ProcessLoginRequest(serviceContext, &request)
	if nkerr != nil {
		logger.WithField("err", nkerr).Error("Login Request Error")
		return nkerr.Message, nkerr
	}

	loginSuccessJson, err := json.Marshal(success)
	if err != nil {
		return "", runtime.NewError(fmt.Sprintf("error marshalling LoginSuccess response: %v", err), INTERNAL)
	}

	return string(loginSuccessJson), nil

}
