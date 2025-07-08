import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { RequestLogin } from './RequestLogin'

describe('RequestLogin', () => {
  it('renders the login message and icon', () => {
   const { container } = render(<RequestLogin />)

    expect(
      screen.getByRole('heading', { name: /you're not logged in/i }),
    ).toBeInTheDocument()

    expect(
      screen.getByText(/log in with your github account to view your repositories/i),
    ).toBeInTheDocument()

    const icon = container.querySelector('svg.lucide-log-in')
    expect(icon).toBeInTheDocument()
    expect(icon).toHaveClass('h-12', 'w-12', 'text-gray-400')
  })
})
