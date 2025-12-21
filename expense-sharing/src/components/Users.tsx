import { useEffect, useState } from "react";
import { get } from "../api";
import type { UserView } from "../types";

export default function Users() {
  const [users, setUsers] = useState<UserView[]>([]);

  useEffect(() => {
    get<UserView[]>("/users").then(setUsers);
  }, []);

  return (
    <div className="section">
      <h2>Users</h2>

      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>User ID</th>
          </tr>
        </thead>
        <tbody>
          {users.map(u => (
            <tr key={u.id}>
              <td>{u.name}</td>
              <td>{u.id}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
