import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { ImportReposDialog } from './ImportRepositoriesDialog'
import { toast } from 'sonner'

import { useGithubRepositoriesLazyQuery } from '@/generated/graphql'

vi.mock('@/generated/graphql.tsx', () => ({
  useGithubRepositoriesLazyQuery: vi.fn(() => [
    vi.fn(),
    {
      data: {
        githubRepositories: [
          {
            githubId: 1,
            name: 'Repo One',
            description: 'First repo',
            updatedAt: new Date().toISOString(),
            stars: 5,
            forks: 2,
            watchers: 3,
            private: false,
          },
          {
            githubId: 2,
            name: 'Repo Two',
            description: 'Second repo',
            updatedAt: new Date().toISOString(),
            stars: 8,
            forks: 1,
            watchers: 4,
            private: true,
          },
        ],
      },
      loading: false,
      error: null,
    },
  ]),
}))

vi.mock('sonner', () => ({
  toast: {
    error: vi.fn(),
  },
}))

describe('ImportReposDialog', () => {
  let onConfirm: ReturnType<typeof vi.fn>

  beforeEach(() => {
    onConfirm = vi.fn()
  })

  it('renders import button if under max limit', () => {
    render(<ImportReposDialog alreadyTracked={[]} onConfirm={onConfirm} />)
    expect(screen.getByText(/import from github/i)).toBeInTheDocument()
  })

  it('shows max limit message when alreadyTracked is at limit', () => {
    render(<ImportReposDialog alreadyTracked={[1, 2]} onConfirm={onConfirm} />)
    expect(
      screen.getByText(/maximum number of repositories/i),
    ).toBeInTheDocument()
  })

  it('opens dialog and shows repositories', async () => {
    render(<ImportReposDialog alreadyTracked={[]} onConfirm={onConfirm} />)

    fireEvent.click(screen.getByRole('button', { name: /import from github/i }))

    await waitFor(() => {
      expect(
        screen.getByPlaceholderText(/search repositories/i),
      ).toBeInTheDocument()
      expect(screen.getByText('Repo One')).toBeInTheDocument()
      expect(screen.getByText('Repo Two')).toBeInTheDocument()
    })
  })

  it('filters repositories by search', async () => {
    render(<ImportReposDialog alreadyTracked={[]} onConfirm={onConfirm} />)
    fireEvent.click(screen.getByRole('button', { name: /import from github/i }))

    await waitFor(() => screen.getByPlaceholderText(/search repositories/i))

    fireEvent.change(screen.getByPlaceholderText(/search repositories/i), {
      target: { value: 'Two' },
    })

    expect(screen.getByText('Repo Two')).toBeInTheDocument()
    expect(screen.queryByText('Repo One')).not.toBeInTheDocument()
  })

  it('selects and confirms repositories', async () => {
    render(<ImportReposDialog alreadyTracked={[]} onConfirm={onConfirm} />)
    fireEvent.click(screen.getByRole('button', { name: /import from github/i }))

    await waitFor(() => screen.getByText('Repo One'))

    fireEvent.click(screen.getByText('Repo One'))
    fireEvent.click(screen.getByRole('button', { name: /confirm/i }))

    expect(onConfirm).toHaveBeenCalledWith([1])
  })

  it('disables confirm button when nothing is selected', async () => {
    render(<ImportReposDialog alreadyTracked={[]} onConfirm={onConfirm} />)
    fireEvent.click(screen.getByRole('button', { name: /import from github/i }))

    await waitFor(() =>
      expect(screen.getByRole('button', { name: /confirm/i })).toBeDisabled(),
    )
  })

  it('shows toast error when GraphQL error is returned', async () => {
    // vi.mocked(require('@/generated/graphql.tsx').useGithubRepositoriesLazyQuery).mockReturnValueOnce([
    //   vi.fn(),
    //   { data: null, loading: false, error: new Error('GraphQL error') },
    // ])
    ;(
      useGithubRepositoriesLazyQuery as ReturnType<typeof vi.fn>
    ).mockReturnValueOnce([
      vi.fn(),
      { data: null, loading: false, error: new Error('GraphQL error') },
    ])

    render(<ImportReposDialog alreadyTracked={[]} onConfirm={onConfirm} />)
    fireEvent.click(screen.getByRole('button', { name: /import from github/i }))

    await waitFor(() => {
      expect(toast.error).toHaveBeenCalledWith('Could not fetch repositories', {
        description: 'GraphQL error',
      })
    })
  })
})
