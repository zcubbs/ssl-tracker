import React from 'react';
import {
    useReactTable,
    getCoreRowModel,
    getPaginationRowModel,
    getSortedRowModel, createColumnHelper, flexRender
} from '@tanstack/react-table'
import {
    useQuery,
} from '@tanstack/react-query'
import {Badge, Loader, Table} from '@mantine/core';

function DomainsTable() {
    const {isLoading, error, data} = useQuery({
        queryKey: ['repoData'],
        queryFn: () =>
            fetch('http://localhost:8000/api/domains').then(
                (res) => res.json(),
            ),
    });

    const columnHelper = createColumnHelper()

    const columns = [
        columnHelper.accessor('name', {
            header: () => 'Domain Name',
        }),
        columnHelper.accessor('status', {
            header: () => 'Status',
            cell: (props) => {
                if (props.getValue() === 'expired') {
                    return (
                        <Badge color="red">{props.getValue()}</Badge>
                    )
                }

                if (props.getValue() === 'valid') {
                    return (
                        <Badge color="green">{props.getValue()}</Badge>
                    )
                }

                return (
                    <Badge color="gray">{props.getValue()}</Badge>
                )
            }
        }),
        columnHelper.accessor('issuer', {
            header: () => 'Issuer',
            cell: (props) => {
                return (
                    <div>
                        {props.getValue().String}
                    </div>
                )
            }
        }),
        columnHelper.accessor('until', {
            header: () => 'Certificate Expiry',
            cell: (props) => {
                return (
                    <div>
                        {props.getValue()}
                    </div>
                )
            }
        })
    ];

    const table = useReactTable({
        columns,
        data,
        getCoreRowModel: getCoreRowModel(),
        getPaginationRowModel: getPaginationRowModel(),
        getSortedRowModel: getSortedRowModel(), //order doesn't matter anymore!
        // etc.
    });

    if (isLoading) return <Loader/>

    if (error) return 'An error has occurred: ' + error.message

    return (
        <div className="p-2">
            <Table>
                <thead>
                {table.getHeaderGroups().map(headerGroup => (
                    <tr key={headerGroup.id}>
                        {headerGroup.headers.map(header => (
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
                {table.getRowModel().rows.map(row => (
                    <tr key={row.id}>
                        {row.getVisibleCells().map(cell => (
                            <td key={cell.id}>
                                {flexRender(cell.column.columnDef.cell, cell.getContext())}
                            </td>
                        ))}
                    </tr>
                ))}
                </tbody>
            </Table>
        </div>
    );
}

export default DomainsTable;
