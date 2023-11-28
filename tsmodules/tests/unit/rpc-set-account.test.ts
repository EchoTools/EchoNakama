import { createMock } from 'ts-auto-mock';
import { On, method } from 'ts-auto-mock/extension';
import { errBadInput, errMissingPayload, setAccountRpc } from "../../src/rpc";
import { describe, expect, beforeEach, test } from '@jest/globals';

describe('echoRelayAddAccountRpc', function() {
    let mockCtx: any, mockLogger: any, mockNk: any, mockLoggerError: any,
		mockNkStorageRead: any, mockNkStorageWrite: any,
		mockStorageWriteAck: any;


	beforeEach(function () {
		// Create mock objects to pass to the RPC.
		mockCtx = createMock<nkruntime.Context>({ userId: 'mock-user' });
		mockLogger = createMock<nkruntime.Logger>();
		mockNk = createMock<nkruntime.Nakama>();
		mockStorageWriteAck = createMock<nkruntime.StorageWriteAck>();

		// Configure specific mock functions to use Jest spies via jest-ts-auto-mock
		mockLoggerError = On(mockLogger).get(method(function (mock) {
			return mock.error;
		}))
		mockNkStorageRead = On(mockNk).get(method(function (mock) {
			return mock.storageRead;
		}));
		mockNkStorageWrite = On(mockNk).get(method(function (mock) {
			return mock.storageWrite;
		}));
	});

	test('returns failure if payload is null', function() {
		const payload = null;
		expect(() => setAccountRpc(mockCtx, mockLogger, mockNk, payload)).toThrow(errMissingPayload);
		expect(mockLoggerError).toBeCalled();
	});

  // Further code goes here
});
