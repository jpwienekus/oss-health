import { describe, it, expect, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { flexRender, getCoreRowModel, useReactTable } from '@tanstack/react-table'
import { Columns } from './Columns'
import type { DependencyType } from '@/generated/graphql'
import { TooltipProvider } from '@/components/ui/tooltip'

function TestTable({ data }: { data: DependencyType[] }) {
  const table = useReactTable({
    data,
    columns: Columns,
    getCoreRowModel: getCoreRowModel(),
  })

  return (
    <TooltipProvider>
      <table>
        <thead>
          {table.getHeaderGroups().map(headerGroup => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map(header => (
                <th key={header.id}>
                  {flexRender(header.column.columnDef.header, header.getContext())}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody>
          {table.getRowModel().rows.map(row => (
            <tr key={row.id}>
              {row.getVisibleCells().map(cell => (
                <td key={cell.id}>
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </TooltipProvider>
  )
}

describe('columns', () => {
  const mockData: DependencyType[] = [
    {
      id: 1,
      name: 'lodash',
      ecosystem: 'npm',
      scanStatus: 'pending',
      scannedAt: '2025-07-08T10:00:00Z',
      repositoryUrl: 'https://github.com/lodash/lodash',
      errorMessage: null,
    },
    {
      id: 2,
      name: 'express',
      ecosystem: 'npm',
      scanStatus: 'completed',
      scannedAt: '2025-07-08T12:00:00Z',
      repositoryUrl: null,
      errorMessage: null,
    },
    {
      id: 3,
      name: 'broken-lib',
      ecosystem: 'pypi',
      scanStatus: 'failed',
      scannedAt: '2025-07-08T08:00:00Z',
      repositoryUrl: 'https://github.com/example/broken-lib',
      errorMessage: 'Failed to fetch metadata',
    },
  ]


  beforeEach(() => {
    render(<TestTable data={mockData} />)
  })

  it('renders dependency names', () => {
    expect(screen.getByText('lodash')).toBeInTheDocument()
    expect(screen.getByText('express')).toBeInTheDocument()
    expect(screen.getByText('broken-lib')).toBeInTheDocument()
  })

  it('renders ecosystem names', () => {
    expect(screen.getAllByText('npm')).toHaveLength(2)
    expect(screen.getByText('pypi')).toBeInTheDocument()
  })

  it('renders scan statuses with correct labels', () => {
    expect(screen.getByText('Pending')).toBeInTheDocument()
    expect(screen.getByText('Completed')).toBeInTheDocument()
    expect(screen.getByText('Failed')).toBeInTheDocument()
  })

  it('renders formatted scan date', () => {
    const dates = screen.getAllByText(/Jul \d{1,2}, 2025/i)
    expect(dates.length).toBe(3)
  })

  it('renders fallback dash for missing repository url', () => {
    expect(screen.getAllByText('-')).toHaveLength(1)
  })

  it('renders actual repository url', () => {
    expect(screen.getByText('https://github.com/lodash/lodash')).toBeInTheDocument()
  })
})
