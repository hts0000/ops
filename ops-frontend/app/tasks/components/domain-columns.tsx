"use client";

import { ColumnDef } from "@tanstack/react-table";

import { Checkbox } from "@/registry/new-york/ui/checkbox";
import { Badge } from "@/registry/new-york/ui/badge";

import { DataTableColumnHeader } from "./data-table-column-header";
import { Domain } from "../data/schema";

export const domainColumns: ColumnDef<Domain>[] = [
  {
    id: "select",
    header: ({ table }) => (
      <Checkbox
        checked={
          table.getIsAllPageRowsSelected() ||
          (table.getIsSomePageRowsSelected() && "indeterminate")
        }
        onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
        aria-label="Select all"
        className="translate-y-[2px]"
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value) => row.toggleSelected(!!value)}
        aria-label="Select row"
        className="translate-y-[2px]"
      />
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "id",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="ID" />
    ),
    cell: ({ row }) => <div className="w-[80px]">{row.getValue("id")}</div>,
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "type",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Type" />
    ),
    cell: ({ row }) => (
      <div className="flex space-x-2">
        <Badge className="max-w-[500px] truncate font-medium" variant="outline">
          {row.getValue("type")}
        </Badge>
      </div>
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "domain",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Domain" />
    ),
    cell: ({ row }) => (
      <div className="flex space-x-2">
        <span className="max-w-[500px] truncate font-medium">
          {row.getValue("domain")}
        </span>
      </div>
    ),
  },
  {
    accessorKey: "record",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Record" />
    ),
    cell: ({ row }) => (
      <div className="flex space-x-2">
        <span className="max-w-[500px] truncate font-medium">
          {row.getValue("record")}
        </span>
      </div>
    ),
  },
  {
    accessorKey: "status",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Status" />
    ),
    cell: ({ row }) => (
      <div className="flex space-x-2">
        <span className="max-w-[500px] truncate font-medium">
          {row.getValue("status") === 1 ? "Active" : "Inactive"}
        </span>
      </div>
    ),
  },
  {
    accessorKey: "isp",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="ISP" />
    ),
    cell: ({ row }) => (
      <div className="flex space-x-2">
        <Badge className="max-w-[500px] truncate font-medium" variant="outline">
          {row.getValue("isp")}
        </Badge>
      </div>
    ),
  },
  {
    accessorKey: "ip",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="IP" />
    ),
    cell: ({ row }) => (
      <div className="flex space-x-2">
        <span className="max-w-[500px] truncate font-medium">
          {row.getValue("ip")}
        </span>
      </div>
    ),
  },
  {
    accessorKey: "city",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="City" />
    ),
    cell: ({ row }) => (
      <div className="flex space-x-2">
        <span className="max-w-[500px] truncate font-medium">
          {row.getValue("city")}
        </span>
      </div>
    ),
  },
  {
    accessorKey: "expire_at",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Expires" />
    ),
    cell: ({ row }) => {
      const now = new Date();
      const date = new Date(row.getValue("expire_at"));
      const diff = date.getTime() - now.getTime();
      const seconds = Math.floor(diff / 1000);
      const days = Math.floor(seconds / (3600 * 24));

      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate font-medium">
            {days <= 30 ? "30天内过期" : days + " 天后过期"}
          </span>
        </div>
      );
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id));
    },
  },
];
