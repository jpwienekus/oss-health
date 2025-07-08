import type { ColumnDef } from '@tanstack/react-table'
import type { DependencyType } from '@/generated/graphql'
import { formatDate } from '@/utils'
import { CheckCircle, LoaderCircle, XCircle, Info } from 'lucide-react'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'

export const statuses = [
  {
    value: 'pending',
    label: 'Pending',
    icon: LoaderCircle,
    className: 'text-yellow-600',
  },
  {
    value: 'completed',
    label: 'Completed',
    icon: CheckCircle,
    className: 'text-green-600',
  },
  {
    value: 'failed',
    label: 'Failed',
    icon: XCircle,
    className: 'text-red-600',
  },
]

export const columns: ColumnDef<DependencyType>[] = [
  {
    accessorKey: 'name',
    header: 'Name',
    enableHiding: false,
    size: 100,
  },
  {
    accessorKey: 'ecosystem',
    header: 'Ecosystem',
    enableHiding: false,
    size: 100,
  },
  {
    accessorKey: 'scanStatus',
    header: 'Scan Status',
    enableHiding: false,
    size: 100,
    cell: ({ row }) => {
      const foundStatus = statuses.find(
        (status) => status.value === row.getValue('scanStatus'),
      )

      if (!foundStatus) {
        return ''
      }

      return (
        <div className="flex w-[100px] items-center gap-2">
          {foundStatus.icon && (
            <foundStatus.icon className={`size-4 ${foundStatus.className}`} />
          )}
          <span>{foundStatus.label}</span>
          {foundStatus.value === 'failed' && (
            <Tooltip>
              <TooltipTrigger asChild>
                <Info className="size-4 text-muted-foreground" />
              </TooltipTrigger>
              <TooltipContent>
                <p>{row.original.errorMessage}</p>
              </TooltipContent>
            </Tooltip>
          )}
        </div>
      )
    },
  },
  {
    accessorKey: 'scannedAt',
    header: 'Scanned At',
    enableHiding: false,
    size: 100,
    cell: ({ row }) => {
      return formatDate(row.getValue('scannedAt'))
    },
  },
  {
    accessorKey: 'repositoryUrl',
    header: 'Repository URL',
    enableHiding: true,
    size: 100,
    cell: ({ row }) => {
      return row.getValue('repositoryUrl') ?? '-'
    },
  },
]
