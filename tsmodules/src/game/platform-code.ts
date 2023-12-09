export enum PlatformCode {
    UNK = 0,
    STM = 1,
    PSN = 2,
    XBX = 3,
    OVR_ORG = 4,
    OVR = 5,
    BOT = 6,
    DMO = 7,
    TEN = 8,
  }
  
  export class PlatformCodeExtensions {
    static GetPrefix(code: PlatformCode): string {
      const name = PlatformCode[code];
      if (name !== undefined) {
        return name.replace("_", "-");
      }
      return "???";
    }
  
    static GetDisplayName(code: PlatformCode): string {
      switch (code) {
        case PlatformCode.UNK:
          return "Unknown";
        case PlatformCode.STM:
          return "Steam";
        case PlatformCode.PSN:
          return "Playstation";
        case PlatformCode.XBX:
          return "Xbox";
        case PlatformCode.OVR_ORG:
          return "Oculus VR (ORG)";
        case PlatformCode.OVR:
          return "Oculus VR";
        case PlatformCode.BOT:
          return "Bot";
        case PlatformCode.DMO:
          return "Demo";
        case PlatformCode.TEN:
          return "Tencent"; // TODO: Verify, this is only suspected to be the target of "TEN".
        default:
          return "Unknown";
      }
    }
  
    static Parse(s: string): PlatformCode {
      s = s.replace("-", "_");
      try {
        return PlatformCode[s as keyof typeof PlatformCode] || PlatformCode.UNK;
      } catch {
        return PlatformCode.UNK;
      }
    }
  }