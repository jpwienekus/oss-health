import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { RepositoryOverview } from './RepositoryOverview'
import type { GitHubRepository } from '@/generated/graphql'

vi.mock('@/utils', () => ({
  formatDate: (date: string) => `Formatted: ${date}`,
}))

const baseRepository: GitHubRepository = {
  name: 'test-repo',
  description: 'A test repo',
  lastScannedAt: '2025-07-01T12:00:00Z',
  score: 85,
  private: false,
  dependencies: 10,
  vulnerabilities: 2,
  forks: 1,
  githubId: 2,
  url: 'test',
  stars: 1,
  watchers: 1
}

describe('RepositoryOverview', () => {
  it('renders repository name and description', () => {
    render(<RepositoryOverview repository={baseRepository} />)
    expect(screen.getByText(baseRepository.name)).toBeInTheDocument()
    expect(screen.getByText(baseRepository.description)).toBeInTheDocument()
  })

  it('shows the green CheckCircle icon and green score when score >= 80', () => {
    render(<RepositoryOverview repository={{ ...baseRepository, score: 90 }} />)

    expect(screen.getByText('90/100')).toHaveClass('text-green-600')
    // Icon is SVG with class 'text-green-600'
    const icon = document.querySelector('svg.text-green-600')
    expect(icon).toBeInTheDocument()
  })

  it('shows yellow AlertTriangle icon and color when score is between 60 and 79', () => {
    render(<RepositoryOverview repository={{ ...baseRepository, score: 65 }} />)

    expect(screen.getByText('65/100')).toHaveClass('text-yellow-600')
    const icon = document.querySelector('svg.text-yellow-600')
    expect(icon).toBeInTheDocument()
  })

  it('shows red XCircle icon and color when score is below 60', () => {
    render(<RepositoryOverview repository={{ ...baseRepository, score: 50 }} />)

    expect(screen.getByText('50/100')).toHaveClass('text-red-600')
    const icon = document.querySelector('svg.text-red-600')
    expect(icon).toBeInTheDocument()
  })

  it('shows "Not scanned yet" if lastScannedAt is null or undefined', () => {
    const repo = { ...baseRepository, lastScannedAt: null }
    render(<RepositoryOverview repository={repo} />)

    expect(screen.getByText('Not scanned yet')).toBeInTheDocument()
    // Should NOT show score or icon
    expect(screen.queryByText(/\/100$/)).not.toBeInTheDocument()
  })

  it('shows "Private" badge if repository.private is true', () => {
    render(<RepositoryOverview repository={{ ...baseRepository, private: true }} />)

    expect(screen.getByText('Private')).toBeInTheDocument()
  })

  it('shows dependencies count or "-" if lastScannedAt missing', () => {
    // When lastScannedAt exists, show dependencies count
    render(<RepositoryOverview repository={baseRepository} />)
    expect(screen.getByText(/10 dependencies/i)).toBeInTheDocument()

    // When lastScannedAt missing, show "-"
    render(<RepositoryOverview repository={{ ...baseRepository, lastScannedAt: null }} />)
    expect(screen.getByText(/- dependencies/i)).toBeInTheDocument()
  })

  it('shows vulnerabilities count with red text if > 0', () => {
    render(<RepositoryOverview repository={baseRepository} />)
    const vulnSpan = screen.getByText(/2 vulnerabilities/i)
    expect(vulnSpan).toHaveClass('text-red-600')

    // If vulnerabilities = 0, no red text
    render(
      <RepositoryOverview
        repository={{ ...baseRepository, vulnerabilities: 0 }}
      />,
    )
    const vulnSpanZero = screen.getByText(/0 vulnerabilities/i)
    expect(vulnSpanZero).not.toHaveClass('text-red-600')
  })

  it('shows formatted lastScannedAt date', () => {
    render(<RepositoryOverview repository={baseRepository} />)
    expect(screen.getByText(`Formatted: ${baseRepository.lastScannedAt}`)).toBeInTheDocument()

    render(<RepositoryOverview repository={{ ...baseRepository, lastScannedAt: null }} />)
    expect(screen.getByText('-')).toBeInTheDocument()
  })
})
