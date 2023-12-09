// Originally translated from EchoRelay.Core.Game.XPlatformId.cs
import { PlatformCode, PlatformCodeExtensions } from './platform-code';


/**
 * Represents a cross-platform ID.
 */
export class XPlatformId {
  private _platformCode: number;
  public AccountId: number;

  public static SIZE: number = 16;

  constructor(platformCode?: PlatformCode, accountId?: number) {
    this._platformCode = platformCode !== undefined ? platformCode : PlatformCode.STM;
    this.AccountId = accountId !== undefined ? accountId : 0;
  }

  public Valid(): boolean {
    return Object.values(PlatformCode).includes(this.PlatformCode);
  }

  public static Parse(s: string): XPlatformId | null {
    const dashIndex = s.lastIndexOf('-');
    if (dashIndex < 0) return null;

    const platformCodeStr = s.substring(0, dashIndex);
    const accountIdStr = s.substring(dashIndex + 1);

    const code = PlatformCodeExtensions.Parse(platformCodeStr);

    if (!code) return null;

    const accountId = Number(accountIdStr);
    if (isNaN(accountId)) return null;

    return new XPlatformId(code, accountId);
  }

  public static Equals(a: XPlatformId | null, b: XPlatformId | null): boolean {
    if (a === b) return true;

    if (!a || !b) return false;

    return a.PlatformCode === b.PlatformCode && a.AccountId === b.AccountId;
  }

  public static NotEquals(a: XPlatformId | null, b: XPlatformId | null): boolean {
    return !XPlatformId.Equals(a, b);
  }

  public Equals(obj: any): boolean {
    if (!(obj instanceof XPlatformId)) return false;

    const objP = obj as XPlatformId;
    return this.AccountId === objP.AccountId && this.PlatformCode === objP.PlatformCode;
  }

  public GetHashCode(): number {
    throw new Error("Not implemented");
  }

  public ToString(): string {
    return `${PlatformCodeExtensions.GetPrefix(this.PlatformCode)}-${this.AccountId}`;
  }

  get PlatformCode(): PlatformCode {
    return this._platformCode;
  }

  set PlatformCode(value: PlatformCode) {
    this._platformCode = value;
  }
}