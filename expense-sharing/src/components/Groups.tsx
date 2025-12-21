import { useEffect, useState } from "react";
import { get } from "../api";
import type { GroupView } from "../types";

export default function Groups() {
  const [groups, setGroups] = useState<GroupView[]>([]);

  useEffect(() => {
    get<GroupView[]>("/groups").then(setGroups);
  }, []);

  return (
    <div className="section">
      <h2>Groups</h2>

      <table>
        <thead>
          <tr>
            <th>Group Name</th>
            <th>Group ID</th>
          </tr>
        </thead>
        <tbody>
          {groups.map(g => (
            <tr key={g.id}>
              <td>{g.name}</td>
              <td>{g.id}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
