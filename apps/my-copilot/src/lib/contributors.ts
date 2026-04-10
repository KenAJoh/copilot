import type { Contributor } from "@/lib/customization-types";

interface CommitAuthor {
  login: string;
  avatar_url: string;
  type?: string;
}

/**
 * Extract unique contributors from GitHub commit data.
 * Filters out bot accounts and deduplicates by login.
 */
export function extractContributors(commits: { author: CommitAuthor | null }[]): Contributor[] {
  const seen = new Set<string>();
  const contributors: Contributor[] = [];

  for (const commit of commits) {
    const author = commit.author;
    if (!author?.login) continue;
    if (author.type === "Bot") continue;
    if (author.login.endsWith("[bot]")) continue;
    if (seen.has(author.login)) continue;

    seen.add(author.login);
    contributors.push({ login: author.login, avatarUrl: author.avatar_url });
  }

  return contributors;
}

/**
 * Fetch contributors for one or more file paths from the public GitHub API.
 * Uses unauthenticated requests since the repo is public.
 */
export async function getFileContributors(
  owner: string,
  repo: string,
  paths: string[]
): Promise<{ contributors: Contributor[]; error: string | null }> {
  try {
    const allCommits: { author: CommitAuthor | null }[] = [];

    for (const path of paths) {
      const url = `https://api.github.com/repos/${owner}/${repo}/commits?path=${encodeURIComponent(path)}&per_page=100`;
      const res = await fetch(url, {
        headers: {
          Accept: "application/vnd.github+json",
          "X-GitHub-Api-Version": "2022-11-28",
        },
      });

      if (!res.ok) {
        return { contributors: [], error: `GitHub API returned ${res.status}` };
      }

      const data: { author: CommitAuthor | null }[] = await res.json();
      allCommits.push(...data);
    }

    return { contributors: extractContributors(allCommits), error: null };
  } catch (error) {
    return { contributors: [], error: error instanceof Error ? error.message : String(error) };
  }
}
