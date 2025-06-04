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
