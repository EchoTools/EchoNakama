import { createMock } from 'ts-auto-mock';
import { On, method } from 'ts-auto-mock/extension';
import AddItemRpc from "./rpc-add-item";
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
		const result = AddItemRpc(mockCtx, mockLogger, mockNk, payload);
		const resultPayload = JSON.parse(result);
		const expectedError = 'no payload provided';

		expect(resultPayload.success).toBe(false);
		expect(resultPayload.error).toBe(expectedError);
		expect(mockLoggerError).toBeCalledWith(expectedError);
	});

  // Further code goes here
});
