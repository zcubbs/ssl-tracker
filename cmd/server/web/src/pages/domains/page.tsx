import {columns, Domain, Domains} from "./columns"
import {DataTable} from "./data-table"
import {Tabs, TabsList, TabsTrigger} from "../../components/ui/tabs.tsx";
import {Button} from "../../components/ui/button.tsx";
import {useQuery} from "@tanstack/react-query";
import {useMemo, useState} from "react";

async function getData(): Promise<Domains> {
    return fetch('http://localhost:8000/api/v1/get_domains').then(
        (res) => res.json(),
    );
}

export default function DomainsPage() {
    const {isLoading, error, data} = useQuery({
        queryKey: ['domains'],
        refetchIntervalInBackground: true,
        refetchInterval: 10000,
        queryFn: getData,
    });

    const [tableData, setTableData] = useState<Domain[]>([]);
    useMemo( () => {
        if (data) {
            setTableData(data?.domains);
        }
    }, [data]);

    if (error) {
        return <div>Error: {error.message}</div>;
    }

    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <div className="flex items-center justify-between space-y-2">
                <h2 className="text-3xl font-bold tracking-tight">Tracked Domains</h2>
                <div className="flex items-center space-x-2">
                    <Button variant="outline">Export</Button>
                </div>
            </div>
            <Tabs defaultValue="all" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="all">All</TabsTrigger>
                    <TabsTrigger value="expired">Expired</TabsTrigger>
                    <TabsTrigger value="valid">Valid</TabsTrigger>
                    <TabsTrigger value="unknown">Unknown</TabsTrigger>
                </TabsList>
            </Tabs>
            {isLoading ? (
                <div>Loading...</div>
            ) : (
                <DataTable columns={columns} data={tableData}/>
            )}
        </div>
    )
}
