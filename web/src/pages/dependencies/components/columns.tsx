import type { ColumnDef } from "@tanstack/react-table"
import type { DependencyType } from "@/generated/graphql";
import { formatDate } from '@/utils'
import { CheckCircle, LoaderCircle, XCircle } from "lucide-react";
import { DataTableColumnHeader } from "./data-table-column-header";

export const statuses = [
  {
    value: "pending",
    label: "Pending",
    icon: LoaderCircle,
    className: "text-yellow-600"
  },
  {
    value: "completed",
    label: "Completed",
    icon: CheckCircle,
    className: "text-green-600"
  },
  {
    value: "failed",
    label: "Failed",
    icon: XCircle,
    className: "text-red-600"
  },
]

export const columns: ColumnDef<DependencyType>[] = [
  {
    accessorKey: "name",
    header: "Name",
    enableHiding: false,
    size: 100,
  },
  {
    accessorKey: "ecosystem",
    header: "Ecosystem",
    enableHiding: false,
    size: 100,
  },
  {
    accessorKey: "status",
    header: "Status",
    enableHiding: false,
    size: 100,
    cell: ({ row }) => {
      const status = statuses.find(
        (status) => status.value === row.getValue("status")
      )

      if (!status) {
        return ''
      }

      return <div className="flex w-[100px] items-center gap-2">
        {status.icon && (
          <status.icon className={`size-4 ${status.className}`} />
        )}
        <span>{status.label}</span>
      </div>
    },
  },
  {
    accessorKey: 'repositoryUrlCheckedAt',
    header: "Scanned At",
    enableHiding: false,
    size: 100,
    cell: ({ row }) => {
      return formatDate(row.getValue('repositoryUrlCheckedAt'))
    }
  },
  {
    id: "repositoryUrlResolveFailedReason",
    accessorKey: "repositoryUrlResolveFailedReason",
    header: "Failed Reason",
    size: 200,
    cell: ({ row }) => {
      const value = row.getValue('repositoryUrlResolveFailedReason') 

      return value ? value : '-'
    },
  },
]
