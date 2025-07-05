import type { ColumnDef } from "@tanstack/react-table"
import type { DependencyType } from "@/generated/graphql";
import { formatDate } from '@/utils'
import { CheckCircle, LoaderCircle, XCircle } from "lucide-react";

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
  },
  {
    accessorKey: "ecosystem",
    header: "Ecosystem",
    enableHiding: false,
  },
  {
    accessorKey: "status",
    header: "Status",
    enableHiding: false,
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
    header: 'Checked At',
    enableHiding: false,
    cell: ({ row }) => {
      return formatDate(row.getValue('repositoryUrlCheckedAt'))
    }
  },
  {
    accessorKey: "repositoryUrlResolveFailedReason",
    header: "Failed Reason",
    cell: ({ row }) => {
      const value = row.getValue('repositoryUrlResolveFailedReason') 

      return value ? value : '-'
    },
  },
]
