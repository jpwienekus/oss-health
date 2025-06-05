export interface Repository {
  name: string
  description?: string
  updatedAt: string

  // TODO: implement these on BE
  healthScore: number
  maintenanceScore: number
  communityScore: number
  releaseCadenceScore: number
  stars: number
  forks: number
  watchers: number
  private: boolean
}

export interface GitHubRepo {
  githubId: number
  name: string
  description: string
  stars: number
  // TODO: fetch
  private: boolean
  forks: number
  watchers: number
  updatedAt: string
}
