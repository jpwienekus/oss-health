import { Input } from '@/components/ui/input'
import type { Table } from '@tanstack/react-table'
import { useState } from 'react'
import { DataTableFacetedFilter } from './DataTableFacetedFilter'
import { Button } from '@/components/ui/button'
import { X } from 'lucide-react'
import { DataTableViewOptions } from './DataTableViewOptions'
import { statuses } from './Columns'

interface DataTableToolbarProps<TData> {
  table: Table<TData>
  handleSearch: (searchValue: string) => void
  handleStatusFilter: (selectedValues: string[]) => void
  statusTotals?: { [key: string]: number | undefined }
}

export function DataTableToolbar<TData>({
  table,
  handleSearch,
  handleStatusFilter,
  statusTotals,
}: DataTableToolbarProps<TData>) {
  const [search, setSearch] = useState('')
  const isFiltered =
    search.length > 0 || table.getState().columnFilters.length > 0
  const handleLocalSearch = (value: string) => {
    setSearch(value)
    handleSearch(value)
  }

  return (
    <div className="flex items-center justify-between">
      <div className="flex flex-1 items-center gap-2">
        <Input
          placeholder="Filter dependencies..."
          value={search}
          onChange={(event) => handleLocalSearch(event.target.value)}
          className="h-8 w-[150px] lg:w-[250px]"
        />
        {table.getColumn('scanStatus') && (
          <DataTableFacetedFilter
            column={table.getColumn('scanStatus')}
            title="Status"
            options={statuses}
            handleStatusFilter={handleStatusFilter}
            facetTotals={statusTotals}
          />
        )}
        {isFiltered && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => {
              handleLocalSearch('')
              handleStatusFilter([])
              table.resetColumnFilters()
            }}
          >
            Reset
            <X />
          </Button>
        )}
      </div>
      <div className="flex items-center gap-2">
        <DataTableViewOptions table={table} />
      </div>
    </div>
  )
}
