export interface GitHubRepository {
  id: number
  githubId: number
  name: string
  description: string
  stars: number
  private: boolean
  forks: number
  watchers: number
  score: number
  vulnerabilities: number
  dependencies: number
  lastScannedAt: string
  updatedAt: string
}
