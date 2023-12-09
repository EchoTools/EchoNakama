import { XPlatformId } from './x-platform-id';
import { PlatformCode } from './platform-code';

describe('XPlatformId', () => {
    describe('constructor', () => {
        it('should set default values if no arguments are provided', () => {
            const xPlatformId = new XPlatformId();
            expect(xPlatformId.PlatformCode).toBe(PlatformCode.STM);
            expect(xPlatformId.AccountId).toBe(0);
        });

        it('should set the provided platform code and account ID', () => {
            const platformCode = PlatformCode.PSN;
            const accountId = 12345;
            const xPlatformId = new XPlatformId(platformCode, accountId);
            expect(xPlatformId.PlatformCode).toBe(platformCode);
            expect(xPlatformId.AccountId).toBe(accountId);
        });
    });

    describe('Valid', () => {
        it('should return true if the platform code is valid', () => {
            const xPlatformId = new XPlatformId(PlatformCode.XBX, 12345);
            expect(xPlatformId.Valid()).toBe(true);
        });

        it('should return false if the platform code is invalid', () => {
            // @ts-ignore
            const xPlatformId = new XPlatformId("AAA", 12345);
            expect(xPlatformId.Valid()).toBe(false);
        });
    });

    describe('Parse', () => {
        it('should return null if the input string is invalid', () => {
            const invalidString = 'invalid-string';
            const xPlatformId = XPlatformId.Parse(invalidString);
            expect(xPlatformId).toBeNull();
        });

        it('should return a valid XPlatformId object if the input string is valid', () => {
            const validString = 'PSN-12345';
            const xPlatformId = XPlatformId.Parse(validString);
            expect(xPlatformId).toBeInstanceOf(XPlatformId);
            expect(xPlatformId!.PlatformCode).toBe(PlatformCode.PSN);
            expect(xPlatformId!.AccountId).toBe(12345);
        });
    });

    describe('Equals', () => {
        it('should return true if the two XPlatformId objects are equal', () => {
            const xPlatformId1 = new XPlatformId(PlatformCode.XBX, 12345);
            const xPlatformId2 = new XPlatformId(PlatformCode.XBX, 12345);
            expect(XPlatformId.Equals(xPlatformId1, xPlatformId2)).toBe(true);
        });

        it('should return false if one of the XPlatformId objects is null', () => {
            const xPlatformId1 = new XPlatformId(PlatformCode.XBX, 12345);
            const xPlatformId2 = null;
            expect(XPlatformId.Equals(xPlatformId1, xPlatformId2)).toBe(false);
        });

        it('should return false if the two XPlatformId objects are not equal', () => {
            const xPlatformId1 = new XPlatformId(PlatformCode.XBX, 12345);
            const xPlatformId2 = new XPlatformId(PlatformCode.PSN, 54321);
            expect(XPlatformId.Equals(xPlatformId1, xPlatformId2)).toBe(false);
        });
    });

    describe('NotEquals', () => {
        it('should return true if the two XPlatformId objects are not equal', () => {
            const xPlatformId1 = new XPlatformId(PlatformCode.XBX, 12345);
            const xPlatformId2 = new XPlatformId(PlatformCode.PSN, 54321);
            expect(XPlatformId.NotEquals(xPlatformId1, xPlatformId2)).toBe(true);
        });

        it('should return false if the two XPlatformId objects are equal', () => {
            const xPlatformId1 = new XPlatformId(PlatformCode.XBX, 12345);
            const xPlatformId2 = new XPlatformId(PlatformCode.XBX, 12345);
            expect(XPlatformId.NotEquals(xPlatformId1, xPlatformId2)).toBe(false);
        });
    });

    describe('Equals (instance method)', () => {
        it('should return true if the two XPlatformId objects are equal', () => {
            const xPlatformId1 = new XPlatformId(PlatformCode.XBX, 12345);
            const xPlatformId2 = new XPlatformId(PlatformCode.XBX, 12345);
            expect(xPlatformId1.Equals(xPlatformId2)).toBe(true);
        });

        it('should return false if the two XPlatformId objects are not equal', () => {
            const xPlatformId1 = new XPlatformId(PlatformCode.XBX, 12345);
            const xPlatformId2 = new XPlatformId(PlatformCode.PSN, 54321);
            expect(xPlatformId1.Equals(xPlatformId2)).toBe(false);
        });

        it('should return false if the input object is not an XPlatformId', () => {
            const xPlatformId = new XPlatformId(PlatformCode.XBX, 12345);
            const otherObject = { PlatformCode: PlatformCode.XBX, AccountId: 12345 };
            expect(xPlatformId.Equals(otherObject)).toBe(false);
        });
    });

    describe('ToString', () => {
        it('should return the string representation of the XPlatformId', () => {
            const xPlatformId = new XPlatformId(PlatformCode.PSN, 12345);
            expect(xPlatformId.ToString()).toBe('PSN-12345');
        });
    });

    describe('PlatformCode (getter and setter)', () => {
        it('should get the platform code', () => {
            const platformCode = PlatformCode.XBX;
            const xPlatformId = new XPlatformId(platformCode, 12345);
            expect(xPlatformId.PlatformCode).toBe(platformCode);
        });

        it('should set the platform code', () => {
            const xPlatformId = new XPlatformId();
            const newPlatformCode = PlatformCode.PSN;
            xPlatformId.PlatformCode = newPlatformCode;
            expect(xPlatformId.PlatformCode).toBe(newPlatformCode);
        });
    });
});