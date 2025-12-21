import { useEffect, useState } from "react";
import { get } from "../api";
import type { BalanceView, UserView, GroupView } from "../types";

type Props = {
  refreshKey: number;
};

export default function GroupBalances({ refreshKey }: Props) {
  const [groups, setGroups] = useState<GroupView[]>([]);
  const [users, setUsers] = useState<UserView[]>([]);
  const [groupId, setGroupId] = useState("");
  const [balances, setBalances] = useState<BalanceView[]>([]);

  useEffect(() => {
    get<GroupView[]>("/groups").then(setGroups);
    get<UserView[]>("/users").then(setUsers);
  }, []);

  useEffect(() => {
    if (!groupId) return;

    get<BalanceView[]>(`/balances/groups?group_id=${groupId}`)
      .then(setBalances);
  }, [groupId, refreshKey]);

  const nameById = (id: string) =>
    users.find(u => u.id === id)?.name ?? id;

  const validBalances = balances.filter(b => b.amount > 0);

  return (
    <div className="section">
      <h2>Group Balances</h2>

      <select value={groupId} onChange={e => setGroupId(e.target.value)}>
        <option value="">Select Group</option>
        {groups.map(g => (
          <option key={g.id} value={g.id}>{g.name}</option>
        ))}
      </select>

      {groupId && (
        validBalances.length === 0 ? (
          <p>No outstanding balances 🎉</p>
        ) : (
          <table>
            <thead>
              <tr>
                <th>Owes From</th>
                <th>Owes To</th>
                <th>Amount</th>
              </tr>
            </thead>
            <tbody>
              {validBalances.map((b, i) => (
                <tr key={i}>
                  <td>{nameById(b.from_user_id)}</td>
                  <td>{nameById(b.to_user_id)}</td>
                  <td>₹ {b.amount}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )
      )}
    </div>
  );
}
