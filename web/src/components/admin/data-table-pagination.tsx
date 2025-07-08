import {
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
} from 'lucide-react'

import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useState } from 'react'

interface DataTablePaginationProps {
  pageSize: number
  totalPages: number
  currentPage: number
  handlePageSize: (size: number) => void
  handleSetPage: (page: number) => void
}

export function DataTablePagination({
  pageSize,
  totalPages,
  currentPage,
  handlePageSize,
  handleSetPage,
}: DataTablePaginationProps) {
  const [pageSizeLocal, setPageSizeLocal] = useState(pageSize.toString())
  return (
    <div className="flex items-center justify-between px-2">
      <div className="text-muted-foreground flex-1 text-sm">
        {pageSize} of {totalPages * pageSize} row(s) selected.
      </div>
      <div className="flex items-center space-x-6 lg:space-x-8">
        <div className="flex items-center space-x-2">
          <p className="text-sm font-medium">Rows per page</p>
          <Select
            value={pageSizeLocal}
            onValueChange={(value) => {
              setPageSizeLocal(value)
              handlePageSize(+value)
            }}
          >
            <SelectTrigger className="h-8 w-[70px]">
              <SelectValue placeholder={pageSizeLocal} />
            </SelectTrigger>
            <SelectContent side="top">
              {[10, 20, 25, 30, 40, 50].map((p) => (
                <SelectItem key={p} value={`${p}`}>
                  {p}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex w-[100px] items-center justify-center text-sm font-medium">
          Page {currentPage} of {totalPages}
        </div>
        <div className="flex items-center space-x-2">
          <Button
            variant="outline"
            size="icon"
            className="hidden size-8 lg:flex"
            onClick={() => handleSetPage(1)}
            disabled={currentPage <= 1}
          >
            <span className="sr-only">Go to first page</span>
            <ChevronsLeft />
          </Button>
          <Button
            variant="outline"
            size="icon"
            className="size-8"
            onClick={() => handleSetPage(currentPage - 1)}
            disabled={currentPage <= 1}
          >
            <span className="sr-only">Go to previous page</span>
            <ChevronLeft />
          </Button>
          <Button
            variant="outline"
            size="icon"
            className="size-8"
            onClick={() => handleSetPage(currentPage + 1)}
            disabled={currentPage >= totalPages}
          >
            <span className="sr-only">Go to next page</span>
            <ChevronRight />
          </Button>
          <Button
            variant="outline"
            size="icon"
            className="hidden size-8 lg:flex"
            onClick={() => handleSetPage(totalPages)}
            disabled={currentPage >= totalPages}
          >
            <span className="sr-only">Go to last page</span>
            <ChevronsRight />
          </Button>
        </div>
      </div>
    </div>
  )
}
