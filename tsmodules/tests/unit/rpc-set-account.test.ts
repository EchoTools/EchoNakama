import fs from 'fs';
import path from 'path';
import { createMock } from 'ts-auto-mock';
import { On, method } from 'ts-auto-mock/extension';
import { errBadInput, errMissingPayload, setAccountRpc } from "../../src/rpc";
import { describe, expect, beforeEach, test } from '@jest/globals';


// Load and parse fixture files
const loadFixture = (fileName) => {
	const filePath = path.join(__dirname, '__fixtures__', fileName);
	const fileContent = fs.readFileSync(filePath, 'utf-8');
	return JSON.parse(fileContent);
  };

// Example: Load a user fixture
const badAccountFixture = loadFixture('invalid-acct-schema.json');

// Set up global variables or functions for fixtures
global.userFixture = badAccountFixture;

  
describe('echoRelaySetAccountRpc', function() {
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

	test('returns failure if nk.storageWrite returns null', function () {
		(mockNkStorageWrite as jest.Mock).mockReturnValueOnce([mockStorageWriteAck]);
			  
		const payload = JSON.stringify(badAccountFixture);
		const result = setAccountRpc(mockCtx, mockLogger, mockNk, payload);
		expect( () => setAccountRpc(mockCtx, mockLogger, mockNk, payload)).toThrow(errBadInput);
	  });

	test('returns failure if account data has invalid schema', function() {
		const payload = JSON.stringify(badAccountFixture);
		expect(() => setAccountRpc(mockCtx, mockLogger, mockNk, payload)).toThrow(errBadInput);
		expect(mockLoggerError).toBeCalled();
	});

  // Further code goes here
});
