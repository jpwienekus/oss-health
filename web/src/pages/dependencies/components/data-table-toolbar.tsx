import { Input } from "@/components/ui/input"
import type { Table } from "@tanstack/react-table"
import { useState } from "react"
// import { X } from "lucide-react"

// import { Button } from "@/registry/new-york-v4/ui/button"
// import { Input } from "@/registry/new-york-v4/ui/input"
// import { DataTableViewOptions } from "@/app/(app)/examples/tasks/components/data-table-view-options"

// import { priorities, statuses } from "../data/data"
// import { DataTableFacetedFilter } from "./data-table-faceted-filter"

interface DataTableToolbarProps<TData> {
  table: Table<TData>,
  handleSearch: (searchValue: string) => void
}

export function DataTableToolbar<TData>({
  table,
  handleSearch
}: DataTableToolbarProps<TData>) {
  const isFiltered = table.getState().columnFilters.length > 0
  const [search, setSearch] = useState('')
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
          onChange={(event) =>
            handleLocalSearch(event.target.value)
          }
          className="h-8 w-[150px] lg:w-[250px]"
        />
        {/*{table.getColumn("status") && (
          <DataTableFacetedFilter
            column={table.getColumn("status")}
            title="Status"
            options={statuses}
          />
        )} */}
        {/* {table.getColumn("priority") && (
          <DataTableFacetedFilter
            column={table.getColumn("priority")}
            title="Priority"
            options={priorities}
          />
        )} */}
        {/* {isFiltered && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => table.resetColumnFilters()}
          >
            Reset
            <X />
          </Button>
        )} */}
      </div>
      <div className="flex items-center gap-2">
        {/*<DataTableViewOptions table={table} />
        <Button size="sm">Add Task</Button> */}
      </div>
    </div>
  )
}
