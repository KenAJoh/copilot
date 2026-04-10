import { describe, it, expect } from "vitest";
import { extractContributors } from "./contributors";

function commit(login: string, avatar_url = `https://avatars.githubusercontent.com/u/${login}`, type = "User") {
  return { author: { login, avatar_url, type } };
}

describe("extractContributors", () => {
  it("should extract unique contributors from commits", () => {
    const commits = [commit("alice"), commit("bob"), commit("alice")];
    const result = extractContributors(commits);

    expect(result).toEqual([
      { login: "alice", avatarUrl: "https://avatars.githubusercontent.com/u/alice" },
      { login: "bob", avatarUrl: "https://avatars.githubusercontent.com/u/bob" },
    ]);
  });

  it("should filter out bot accounts by type", () => {
    const commits = [commit("alice"), commit("github-actions[bot]", "https://avatar.example.com", "Bot")];
    const result = extractContributors(commits);

    expect(result).toEqual([{ login: "alice", avatarUrl: "https://avatars.githubusercontent.com/u/alice" }]);
  });

  it("should filter out bot accounts by login suffix", () => {
    const commits = [commit("alice"), commit("dependabot[bot]"), commit("copilot[bot]")];
    const result = extractContributors(commits);

    expect(result).toEqual([{ login: "alice", avatarUrl: "https://avatars.githubusercontent.com/u/alice" }]);
  });

  it("should skip commits with null author", () => {
    const commits = [{ author: null }, commit("alice"), { author: null }];
    const result = extractContributors(commits);

    expect(result).toEqual([{ login: "alice", avatarUrl: "https://avatars.githubusercontent.com/u/alice" }]);
  });

  it("should return empty array for empty input", () => {
    expect(extractContributors([])).toEqual([]);
  });

  it("should return empty array when all commits are from bots", () => {
    const commits = [commit("dependabot[bot]"), commit("github-actions[bot]", "https://example.com", "Bot")];
    expect(extractContributors(commits)).toEqual([]);
  });

  it("should preserve order of first appearance", () => {
    const commits = [commit("charlie"), commit("alice"), commit("bob"), commit("charlie"), commit("alice")];
    const result = extractContributors(commits);

    expect(result.map((c) => c.login)).toEqual(["charlie", "alice", "bob"]);
  });
});
