import type { ColumnDef } from "@tanstack/react-table"
import type { DependencyType } from "@/generated/graphql";
import { formatDate } from '@/utils'
import { CheckCircle, LoaderCircle, XCircle } from "lucide-react";

const statuses = {
  pending: {
    label: "Pending",
    icon: LoaderCircle,
    class: "text-yellow-600 size-4"
  },
  completed: {
    label: "Completed",
    icon: CheckCircle,
    class: "text-green-600 size-4"
  },
  failed: {
    label: "Failed",
    icon: XCircle,
    class: "text-red-600 size-4"
  },
} as const

export const columns: ColumnDef<DependencyType>[] = [
  {
    accessorKey: "name",
    header: "Name",
  },
  {
    accessorKey: "ecosystem",
    header: "Registry",
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => {
      const status = statuses[row.getValue('status') as keyof typeof statuses]

      return <div className="flex w-[100px] items-center gap-2">
        {status.icon && (
          <status.icon className={status.class} />
        )}
        <span>{status.label}</span>
      </div>
    },
  },
  {
    accessorKey: 'repositoryUrlCheckedAt',
    header: 'Checked At',
    cell: ({ row }) => {
      return formatDate(row.getValue('repositoryUrlCheckedAt'))
    }
  },
  {
    accessorKey: "repositoryUrlResolveFailedReason",
    header: "Failed Reason",
  },
]
