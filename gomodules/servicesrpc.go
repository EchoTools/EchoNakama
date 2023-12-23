package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"

	"echonakama/server/services"
	"echonakama/server/services/login"
)

const (
	// Websocket Error Codes
	StatusOK                 = 0  // StatusOK indicates a successful operation.
	StatusCanceled           = 1  // StatusCanceled indicates the operation was canceled.
	StatusUnknown            = 2  // StatusUnknown indicates an unknown error occurred.
	StatusInvalidArgument    = 3  // StatusInvalidArgument indicates an invalid argument was provided.
	StatusDeadlineExceeded   = 4  // StatusDeadlineExceeded indicates the operation exceeded the deadline.
	StatusNotFound           = 5  // StatusNotFound indicates the requested resource was not found.
	StatusAlreadyExists      = 6  // StatusAlreadyExists indicates the resource already exists.
	StatusPermissionDenied   = 7  // StatusPermissionDenied indicates the operation was denied due to insufficient permissions.
	StatusResourceExhausted  = 8  // StatusResourceExhausted indicates the resource has been exhausted.
	StatusFailedPrecondition = 9  // StatusFailedPrecondition indicates a precondition for the operation was not met.
	StatusAborted            = 10 // StatusAborted indicates the operation was aborted.
	StatusOutOfRange         = 11 // StatusOutOfRange indicates a value is out of range.
	StatusUnimplemented      = 12 // StatusUnimplemented indicates the operation is not implemented.
	StatusInternalError      = 13 // StatusInternal indicates an internal server error occurred.
	StatusUnavailable        = 14 // StatusUnavailable indicates the service is currently unavailable.
	StatusDataLoss           = 15 // StatusDataLoss indicates a loss of data occurred.
	StatusUnauthenticated    = 16 // StatusUnauthenticated indicates the request lacks valid authentication credentials.
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {

	if err := initializer.RegisterRpc("relay/loginrequest", LoginRequestRpc); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	registerIndexes(initializer)

	logger.Info("Initialized module.")
	return nil
}

// Handles the user login request from Echo Relay
func LoginRequestRpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	// Parse the payload into a LoginRequest object

	var request login.LoginRequest
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		logger.WithField("err", err).Error("Unable to unmarshal payload")
		return "", runtime.NewError("Unable to unmarshal payload", StatusInvalidArgument)
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
		return "", runtime.NewError(fmt.Sprintf("error marshalling LoginSuccess response: %v", err), StatusInternalError)
	}

	return string(loginSuccessJson), nil
}

func registerIndexes(initializer runtime.Initializer) error {
	// Register the LinkTicket Index that prevents multiple LinkTickets with the same device_id_str
	name := login.LinkTicketIndex
	collection := login.LinkTicketCollection
	key := ""                                                        // Set to empty string to match all keys instead
	fields := []string{"game_user_id_token", "nk_device_auth_token"} // index on these fields
	maxEntries := 1000000
	indexOnly := false

	if err := initializer.RegisterStorageIndex(name, collection, key, fields, maxEntries, indexOnly); err != nil {
		return err
	}

	// Register the IP Address index for looking up user's by IP Address
	name = login.IpAddressIndex
	collection = login.XPlatformIdStorageCollection
	key = ""                               // Set to empty string to match all keys instead
	fields = []string{"client_ip_address"} // index on these fields
	maxEntries = 1000000
	indexOnly = false

	err := initializer.RegisterStorageIndex(name, collection, key, fields, maxEntries, indexOnly)
	if err != nil {
		return err
	}
	return nil
}
