import { PlatformCode, PlatformCodeExtensions } from './platform-code';

describe('PlatformCode', () => {
  test('PlatformCode enum should have correct values', () => {
    expect(PlatformCode.UNK).toBe(0);
    expect(PlatformCode.STM).toBe(1);
    expect(PlatformCode.PSN).toBe(2);
    expect(PlatformCode.XBX).toBe(3);
    expect(PlatformCode.OVR_ORG).toBe(4);
    expect(PlatformCode.OVR).toBe(5);
    expect(PlatformCode.BOT).toBe(6);
    expect(PlatformCode.DMO).toBe(7);
    expect(PlatformCode.TEN).toBe(8);
  });
});

describe('PlatformCodeExtensions', () => {
  test('GetPrefix should return the correct prefix', () => {
    expect(PlatformCodeExtensions.GetPrefix(PlatformCode.UNK)).toBe("UNK");
    expect(PlatformCodeExtensions.GetPrefix(PlatformCode.STM)).toBe("STM");
    expect(PlatformCodeExtensions.GetPrefix(PlatformCode.PSN)).toBe("PSN");
    expect(PlatformCodeExtensions.GetPrefix(PlatformCode.XBX)).toBe("XBX");
    expect(PlatformCodeExtensions.GetPrefix(PlatformCode.OVR_ORG)).toBe("OVR-ORG");
    expect(PlatformCodeExtensions.GetPrefix(PlatformCode.OVR)).toBe("OVR");
    expect(PlatformCodeExtensions.GetPrefix(PlatformCode.BOT)).toBe("BOT");
    expect(PlatformCodeExtensions.GetPrefix(PlatformCode.DMO)).toBe("DMO");
    expect(PlatformCodeExtensions.GetPrefix(PlatformCode.TEN)).toBe("TEN");
  });

  test('GetDisplayName should return the correct display name', () => {
    expect(PlatformCodeExtensions.GetDisplayName(PlatformCode.UNK)).toBe("Unknown");
    expect(PlatformCodeExtensions.GetDisplayName(PlatformCode.STM)).toBe("Steam");
    expect(PlatformCodeExtensions.GetDisplayName(PlatformCode.PSN)).toBe("Playstation");
    expect(PlatformCodeExtensions.GetDisplayName(PlatformCode.XBX)).toBe("Xbox");
    expect(PlatformCodeExtensions.GetDisplayName(PlatformCode.OVR_ORG)).toBe("Oculus VR (ORG)");
    expect(PlatformCodeExtensions.GetDisplayName(PlatformCode.OVR)).toBe("Oculus VR");
    expect(PlatformCodeExtensions.GetDisplayName(PlatformCode.BOT)).toBe("Bot");
    expect(PlatformCodeExtensions.GetDisplayName(PlatformCode.DMO)).toBe("Demo");
    expect(PlatformCodeExtensions.GetDisplayName(PlatformCode.TEN)).toBe("Tencent");
  });

  test('Parse should return the correct PlatformCode', () => {
    expect(PlatformCodeExtensions.Parse("UNK")).toBe(PlatformCode.UNK);
    expect(PlatformCodeExtensions.Parse("STM")).toBe(PlatformCode.STM);
    expect(PlatformCodeExtensions.Parse("PSN")).toBe(PlatformCode.PSN);
    expect(PlatformCodeExtensions.Parse("XBX")).toBe(PlatformCode.XBX);
    expect(PlatformCodeExtensions.Parse("OVR-ORG")).toBe(PlatformCode.OVR_ORG);
    expect(PlatformCodeExtensions.Parse("OVR")).toBe(PlatformCode.OVR);
    expect(PlatformCodeExtensions.Parse("BOT")).toBe(PlatformCode.BOT);
    expect(PlatformCodeExtensions.Parse("DMO")).toBe(PlatformCode.DMO);
    expect(PlatformCodeExtensions.Parse("TEN")).toBe(PlatformCode.TEN);
    expect(PlatformCodeExtensions.Parse("invalid")).toBe(PlatformCode.UNK);
  });
});