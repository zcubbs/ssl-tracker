"use client"

import {ColumnDef} from "@tanstack/react-table"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../../components/ui/dropdown-menu.tsx";
import {Button} from "../../components/ui/button.tsx";
import {ArrowUpDown, MoreHorizontal} from "lucide-react"
import {parseDate, timeLeftUntil} from "../../lib/utils";

// This type is used to define the shape of our data.
// You can use a Zod schema here if you want.
export type Domain = {
  id: string
  name: string
  status: "pending" | "valid" | "expired" | "unknown"
  issuer: string
  expirationDate: Date
}

export const columns: ColumnDef<Domain>[] = [
  {
    accessorKey: "name",
    header: "Name",
  },
  {
    accessorKey: "status",

    header: ({column}) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Status
          <ArrowUpDown className="ml-2 h-4 w-4"/>
        </Button>
      )
    },
  },
  {
    accessorKey: "issuer",
    header: "Issuer",
  },
  {
    accessorKey: "expirationDate",
    header: () => <div className="text-right">Expires</div>,
    cell: ({row}) => {
      const expirationDate = parseDate(row.getValue("expirationDate"))
      if (!expirationDate) return <div className="text-right">-</div>

      // format time until expiration
      const formatted = timeLeftUntil(expirationDate)

      return <div className="text-right font-medium">{formatted}</div>
    },
  },
  {
    id: "actions",
    cell: ({row}) => {
      const domain = row.original

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="h-4 w-4"/>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            <DropdownMenuItem
              onClick={() => navigator.clipboard.writeText(domain.id)}
            >
              Copy domain ID
            </DropdownMenuItem>
            <DropdownMenuSeparator/>
            <DropdownMenuItem>View domain details</DropdownMenuItem>
            <DropdownMenuItem>Delete</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      )
    },
  },
]
