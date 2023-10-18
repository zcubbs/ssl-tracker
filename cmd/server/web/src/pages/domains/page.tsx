import {columns, Domain} from "./columns"
import {DataTable} from "./data-table"
import {useMemo, useState} from "react";
import {addTwoHours} from "../../lib/utils.ts";
import {Tabs, TabsList, TabsTrigger} from "../../components/ui/tabs.tsx";
import {Button} from "../../components/ui/button.tsx";

async function getData(): Promise<Domain[]> {
  // Fetch data from your API here.
  return [
    {
      id: "1",
      name: "example.com",
      status: "valid",
      issuer: "Let's Encrypt",
      expirationDate: addTwoHours(),
    },
    {
      id: "2",
      name: "example2.com",
      status: "valid",
      issuer: "Let's Encrypt",
      expirationDate: addTwoHours(),
    },
    {
      id: "3",
      name: "example3.com",
      status: "expired",
      issuer: "GoDaddy",
      expirationDate: addTwoHours(),
    },
    // ...
  ]
}

export default function DomainsPage() {
  const [data, setData] = useState<Domain[]>([]);

  useMemo( () => {
    const data = getData();
    data.then(data => setData(data));
  }, []);

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
      <DataTable columns={columns} data={data} />
    </div>
  )
}
