import { vi } from "vitest";

// Mock next/headers before importing auth
vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("next/navigation", () => ({
  redirect: vi.fn(() => {
    throw new Error("NEXT_REDIRECT");
  }),
}));

vi.mock("./jwt", () => ({
  validate: vi.fn(),
}));

import { getUser, isAuthenticated } from "./auth";
import { headers } from "next/headers";
import { redirect } from "next/navigation";
import { validate } from "./jwt";

describe("getUser", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.unstubAllEnvs();
  });

  it("returns mock user in development mode", async () => {
    vi.stubEnv("NODE_ENV", "development");

    const user = await getUser();
    expect(user).not.toBeNull();
    expect(user!.firstName).toBe("Hans Kristian");
    expect(user!.email).toContain("@nav.no");
  });

  it("throws when required env vars are missing", async () => {
    vi.stubEnv("NODE_ENV", "production");
    vi.stubEnv("AZURE_APP_CLIENT_ID", "");

    await expect(getUser(false)).rejects.toThrow("Environment variable");
  });

  it("returns null when no Authorization header and shouldRedirect is false", async () => {
    vi.stubEnv("NODE_ENV", "production");
    vi.stubEnv("AZURE_APP_CLIENT_ID", "client-id");
    vi.stubEnv("AZURE_OPENID_CONFIG_JWKS_URI", "https://jwks.example.com");
    vi.stubEnv("AZURE_OPENID_CONFIG_ISSUER", "https://issuer.example.com");

    vi.mocked(headers).mockResolvedValue({
      get: () => null,
    } as unknown as Awaited<ReturnType<typeof headers>>);

    const user = await getUser(false);
    expect(user).toBeNull();
  });

  it("redirects when no Authorization header and shouldRedirect is true", async () => {
    vi.stubEnv("NODE_ENV", "production");
    vi.stubEnv("AZURE_APP_CLIENT_ID", "client-id");
    vi.stubEnv("AZURE_OPENID_CONFIG_JWKS_URI", "https://jwks.example.com");
    vi.stubEnv("AZURE_OPENID_CONFIG_ISSUER", "https://issuer.example.com");

    vi.mocked(headers).mockResolvedValue({
      get: () => null,
    } as unknown as Awaited<ReturnType<typeof headers>>);

    await expect(getUser(true)).rejects.toThrow("NEXT_REDIRECT");
    expect(redirect).toHaveBeenCalledWith("/oauth2/login");
  });

  it("returns null on invalid token when shouldRedirect is false", async () => {
    vi.stubEnv("NODE_ENV", "production");
    vi.stubEnv("AZURE_APP_CLIENT_ID", "client-id");
    vi.stubEnv("AZURE_OPENID_CONFIG_JWKS_URI", "https://jwks.example.com");
    vi.stubEnv("AZURE_OPENID_CONFIG_ISSUER", "https://issuer.example.com");

    vi.mocked(headers).mockResolvedValue({
      get: () => "Bearer invalid-token",
    } as unknown as Awaited<ReturnType<typeof headers>>);

    vi.mocked(validate).mockResolvedValue({ isValid: false, payload: undefined, error: "bad token" });

    const user = await getUser(false);
    expect(user).toBeNull();
  });

  it("parses user from valid JWT payload", async () => {
    vi.stubEnv("NODE_ENV", "production");
    vi.stubEnv("AZURE_APP_CLIENT_ID", "client-id");
    vi.stubEnv("AZURE_OPENID_CONFIG_JWKS_URI", "https://jwks.example.com");
    vi.stubEnv("AZURE_OPENID_CONFIG_ISSUER", "https://issuer.example.com");

    vi.mocked(headers).mockResolvedValue({
      get: () => "Bearer valid-token",
    } as unknown as Awaited<ReturnType<typeof headers>>);

    vi.mocked(validate).mockResolvedValue({
      isValid: true,
      payload: {
        name: "Flaatten, Hans Kristian",
        preferred_username: "Hans.Kristian.Flaatten@Nav.No",
        groups: ["admin", "copilot-users"],
      },
    });

    const user = await getUser(false);
    expect(user).not.toBeNull();
    expect(user!.firstName).toBe("Hans Kristian");
    expect(user!.lastName).toBe("Flaatten");
    expect(user!.email).toBe("hans.kristian.flaatten@nav.no"); // lowercased
    expect(user!.groups).toEqual(["admin", "copilot-users"]);
  });

  it("handles missing name in payload", async () => {
    vi.stubEnv("NODE_ENV", "production");
    vi.stubEnv("AZURE_APP_CLIENT_ID", "client-id");
    vi.stubEnv("AZURE_OPENID_CONFIG_JWKS_URI", "https://jwks.example.com");
    vi.stubEnv("AZURE_OPENID_CONFIG_ISSUER", "https://issuer.example.com");

    vi.mocked(headers).mockResolvedValue({
      get: () => "Bearer valid-token",
    } as unknown as Awaited<ReturnType<typeof headers>>);

    vi.mocked(validate).mockResolvedValue({
      isValid: true,
      payload: { preferred_username: "user@nav.no", groups: [] },
    });

    const user = await getUser(false);
    expect(user).not.toBeNull();
    expect(user!.firstName).toBe("");
    expect(user!.lastName).toBe("");
  });
});

describe("isAuthenticated", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.unstubAllEnvs();
  });

  it("returns true in development mode", async () => {
    vi.stubEnv("NODE_ENV", "development");
    expect(await isAuthenticated()).toBe(true);
  });

  it("returns false when no auth header", async () => {
    vi.stubEnv("NODE_ENV", "production");
    vi.stubEnv("AZURE_APP_CLIENT_ID", "client-id");
    vi.stubEnv("AZURE_OPENID_CONFIG_JWKS_URI", "https://jwks.example.com");
    vi.stubEnv("AZURE_OPENID_CONFIG_ISSUER", "https://issuer.example.com");

    vi.mocked(headers).mockResolvedValue({
      get: () => null,
    } as unknown as Awaited<ReturnType<typeof headers>>);

    expect(await isAuthenticated()).toBe(false);
  });
});
