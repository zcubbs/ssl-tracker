import {useEffect, useState} from "react";
import {Table, TableBody, TableCaption, TableCell, TableHeader, TableRow} from "@/components/ui/table.tsx";
import axios from "@/api/axios.ts";
import useRefreshToken from "@/hooks/use-refresh-token.ts";
import {Button} from "@/components/ui/button.tsx";

type User = {
  id: string;
  username: string;
  full_name: string;
  role: string;
  password_changed_at: string;
  created_at: string;
}

const UsersPage = () => {
  const [users, setUsers] = useState<User[]>([]);
  const refresh = useRefreshToken();

  useEffect(() => {
    let isMounted = true;
    const controller = new AbortController();

    const getUsers = async () => {
      try {
        const response = await axios.get('/api/v1/get_users', {
          signal: controller.signal,
        });
        console.log(response.data);
        isMounted && setUsers(response.data);
      } catch (error) {
        console.error(error);
      }
    }

    getUsers();

    return () => {
      isMounted = false;
      controller.abort();
    }

  }, []);

  return (
    <div>
      <h1>Users</h1>

      <Button onClick={() => refresh()}>Refresh</Button>
      <Table>
        <TableCaption>Users</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHeader>Username</TableHeader>
            <TableHeader>Full Name</TableHeader>
            <TableHeader>Role</TableHeader>
            <TableHeader>Password Changed At</TableHeader>
            <TableHeader>Created At</TableHeader>
          </TableRow>
        </TableHeader>
        <TableBody>
          {users?.length ?
            (users.map((user) => (
            <TableRow key={user.id}>
              <TableCell>{user.username}</TableCell>
              <TableCell>{user.full_name}</TableCell>
              <TableCell>{user.role}</TableCell>
              <TableCell>{user.password_changed_at}</TableCell>
              <TableCell>{user.created_at}</TableCell>
            </TableRow>
            ))) : (<TableRow><TableCell>No users found</TableCell></TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}

export default UsersPage;
