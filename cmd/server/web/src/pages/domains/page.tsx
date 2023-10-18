import {columns, Domain} from "./columns"
import {DataTable} from "./data-table"
import {useMemo, useState} from "react";
import {addTwoHours} from "../../lib/utils.ts";

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
    <div className="hidden h-full flex-1 flex-col space-y-8 p-8 md:flex">
      <div className="flex items-center justify-between space-y-2">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Tracked Domains</h2>
          <p className="text-muted-foreground">
            Tracked domains and their expiration dates.
          </p>
        </div>
        <div className="flex items-center space-x-2">
          {/*<UserNav />*/}
        </div>
      </div>
      <DataTable columns={columns} data={data} />
    </div>
  )
}
