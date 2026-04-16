import { VALID_DOMAINS, type Domain } from "./manifest-types";

describe("manifest-types", () => {
  it("VALID_DOMAINS contains all Domain union members", () => {
    const expected: Domain[] = [
      "platform",
      "frontend",
      "backend",
      "auth",
      "observability",
      "general",
      "testing",
      "design",
    ];
    expect(VALID_DOMAINS).toEqual(expected);
  });

  it("VALID_DOMAINS is frozen at 8 entries", () => {
    expect(VALID_DOMAINS).toHaveLength(8);
  });
});

describe("metadata validation", () => {
  it("accepts valid domain values", () => {
    for (const domain of VALID_DOMAINS) {
      expect(VALID_DOMAINS.includes(domain)).toBe(true);
    }
  });

  it("rejects invalid domain values", () => {
    const invalid = "security" as Domain;
    expect(VALID_DOMAINS.includes(invalid)).toBe(false);
  });
});
