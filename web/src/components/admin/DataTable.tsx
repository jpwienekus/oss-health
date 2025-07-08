import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  flexRender,
  getCoreRowModel,
  useReactTable,
  type ColumnDef,
} from '@tanstack/react-table'
import { DataTableToolbar } from './DataTableToolbar'
import { DataTableColumnHeader } from './DataTableColumnHeader'
import { DataTablePagination } from './DataTablePagination'

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
  handleSearch: (searchValue: string) => void
  handleStatusFilter: (selectedValues: string[]) => void
  handleSort: (sortColumn: string, sortDirection: string) => void
  handlePageSize: (size: number) => void
  handleSetPage: (page: number) => void
  pageSize: number
  totalPages: number
  currentPage: number
  statusTotals?: { [key: string]: number | undefined }
}

export function DataTable<TData, TValue>({
  columns,
  data,
  handleSearch,
  handleStatusFilter,
  handleSort,
  handlePageSize,
  handleSetPage,
  pageSize,
  totalPages,
  currentPage,
  statusTotals,
}: DataTableProps<TData, TValue>) {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  })

  return (
    <div className="flex flex-col gap-4">
      <DataTableToolbar
        table={table}
        handleSearch={handleSearch}
        handleStatusFilter={handleStatusFilter}
        statusTotals={statusTotals}
      />
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead
                      key={header.id}
                      style={{ width: header.getSize() }}
                    >
                      {header.isPlaceholder ? null : (
                        <DataTableColumnHeader
                          column={header.column}
                          title={header.column.columnDef.header as string}
                          handleSort={handleSort}
                        />
                      )}
                    </TableHead>
                  )
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  daa-state={row.getIsSelected() && 'selected'}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell
                      key={cell.id}
                      style={{ width: cell.column.getSize() }}
                    >
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext(),
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  No results.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <DataTablePagination
        pageSize={pageSize}
        totalPages={totalPages}
        currentPage={currentPage}
        handlePageSize={handlePageSize}
        handleSetPage={handleSetPage}
      />
    </div>
  )
}
