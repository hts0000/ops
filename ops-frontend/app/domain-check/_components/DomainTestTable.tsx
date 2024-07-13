"use client";

import {
  Column,
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import React from "react";

export type DomainTable = {
  ID: number;
  Type: "A" | "CNAME" | "MX" | "";
  Domain: string;
  Record: string;
  Status: 0 | 1;
  ISP: string;
  IP: string;
  City: string;
  ExpireAt: string;
  Remark: string;
};

const data: DomainTable[] = [
  {
    ID: 1,
    Type: "A",
    Domain: "www.baidu.com",
    Record: "www",
    Status: 0,
    ISP: "baidu",
    IP: "",
    City: "",
    ExpireAt: "2021-01-01",
    Remark: "test",
  },
  {
    ID: 2,
    Type: "CNAME",
    Domain: "www.baidu.com",
    Record: "www",
    Status: 0,
    ISP: "baidu",
    IP: "",
    City: "",
    ExpireAt: "2021-01-01",
    Remark: "test",
  },
  {
    ID: 3,
    Type: "MX",
    Domain: "www.baidu.com",
    Record: "www",
    Status: 0,
    ISP: "baidu",
    IP: "",
    City: "",
    ExpireAt: "2021-01-01",
    Remark: "test",
  },
  {
    ID: 4,
    Type: "A",
    Domain: "www.baidu.com",
    Record: "www",
    Status: 0,
    ISP: "baidu",
    IP: "",
    City: "",
    ExpireAt: "2021-01-01",
    Remark: "test",
  },
  {
    ID: 5,
    Type: "CNAME",
    Domain: "www.baidu.com",
    Record: "www",
    Status: 0,
    ISP: "baidu",
    IP: "",
    City: "",
    ExpireAt: "2021-01-01",
    Remark: "test",
  },
  {
    ID: 6,
    Type: "MX",
    Domain: "www.baidu.com",
    Record: "www",
    Status: 0,
    ISP: "baidu",
    IP: "",
    City: "",
    ExpireAt: "2021-01-01",
    Remark: "test",
  },
  {
    ID: 7,
    Type: "A",
    Domain: "www.baidu.com",
    Record: "www",
    Status: 0,
    ISP: "baidu",
    IP: "",
    City: "",
    ExpireAt: "2021-01-01",
    Remark: "test",
  },
  {
    ID: 8,
    Type: "CNAME",
    Domain: "www.baidu.com",
    Record: "www",
    Status: 0,
    ISP: "baidu",
    IP: "",
    City: "",
    ExpireAt: "2021-01-01",
    Remark: "test",
  },
  {
    ID: 9,
    Type: "MX",
    Domain: "www.baidu.com",
    Record: "www",
    Status: 0,
    ISP: "baidu",
    IP: "",
    City: "",
    ExpireAt: "2021-01-01",
    Remark: "test",
  },
  {
    ID: 10,
    Type: "A",
    Domain: "www.baidu.com",
    Record: "www",
    Status: 0,
    ISP: "baidu",
    IP: "",
    City: "",
    ExpireAt: "2021-01-01",
    Remark: "test",
  },
];

export const columns: ColumnDef<DomainTable>[] = [
  {
    id: "select",
    header: "Select",
    cell: "Select",
    footer: "select",
  },
  {
    id: "select2",
    header: () => <h3>Select2</h3>,
    cell: "Select2",
    footer: "select2",
  },
];

const DomainTestTable = () => {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="p-2">
      <table>
        <thead>
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <th key={header.id}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody>
          {table.getRowModel().rows.map((row) => (
            <tr key={row.id}>
              {row.getVisibleCells().map((cell) => (
                <td key={cell.id}>
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
        <tfoot>
          {table.getFooterGroups().map((footerGroup) => (
            <tr key={footerGroup.id}>
              {footerGroup.headers.map((header) => (
                <th key={header.id}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.footer,
                        header.getContext()
                      )}
                </th>
              ))}
            </tr>
          ))}
        </tfoot>
      </table>
    </div>
  );
};

export default DomainTestTable;
