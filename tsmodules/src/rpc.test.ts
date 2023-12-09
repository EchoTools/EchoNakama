import fs from 'fs';
import path from 'path';
import { createMock } from 'ts-auto-mock';
import { On, method } from 'ts-auto-mock/extension';
import { describe, expect, beforeEach, test } from '@jest/globals';
import { getDeviceLinkCodeRpc } from '../src/rpc';

describe('getDeviceLinkCodeRpc', () => {
  let ctx: nkruntime.Context;
  let logger: nkruntime.Logger;
  let nk: nkruntime.Nakama;
  let payload: any;


  beforeEach(() => {
    ctx = {} as nkruntime.Context;
    logger = {} as nkruntime.Logger;
    nk = {} as nkruntime.Nakama;
    payload = { id: 'OVR-ORG-12341234123433'};
  });

  it('should generate a new link code and store it in storage', () => {
    // Mock the necessary functions and objects
    const parsePayloadMock = jest.fn().mockReturnValue({ id: 'testId' });
    const generateLinkCodeMock = jest.fn().mockReturnValue('testCode');
    const getStorageObjectMock = jest.fn().mockReturnValue(null);
    const loggerInfoMock = jest.fn().mockReturnValue(null);
    const storageWriteMock = jest.fn();
    nk.storageWrite = storageWriteMock;

    // Set up the expected values
    const expectedLinkData = { code: 'testCode' };
    const expectedStorageWriteParams = [
      {
        collection: 'linkCodes',
        key: 'testCode',
        value: { deviceId: 'testId', code: 'testCode' },
        userId: 'systemUserId',
        version: '*',
        permissionRead: 1,
        permissionWrite: 1,
      },
    ];

    // Set up the function under test
    const result = getDeviceLinkCodeRpc(ctx, logger, nk, payload);

    // Verify the function calls and assertions
    expect(parsePayloadMock).toHaveBeenCalledWith(payload);
    expect(generateLinkCodeMock).toHaveBeenCalled();
    expect(getStorageObjectMock).toHaveBeenCalledWith(nk, logger, 'linkCodes', 'testCode', 'systemUserId');
    expect(storageWriteMock).toHaveBeenCalledWith(expectedStorageWriteParams);
    expect(result).toBe(JSON.stringify(expectedLinkData));
  });

  it('should throw an error if failed to create link code in storage', () => {
    // Mock the necessary functions and objects
    const parsePayloadMock = jest.fn().mockReturnValue({ id: 'testId' });
    const generateLinkCodeMock = jest.fn().mockReturnValue('testCode');
    const loggerInfoMock = jest.fn().mockReturnValue(null);
    const getStorageObjectMock = jest.fn().mockReturnValue(null);
    const storageWriteMock = jest.fn().mockImplementation(() => {
      throw new Error('Failed to create link code');
    });
    nk.storageWrite = storageWriteMock;

    // Set up the function under test
    expect(() => getDeviceLinkCodeRpc(ctx, logger, nk, payload)).toThrowError('Failed to create link code');
  });
});