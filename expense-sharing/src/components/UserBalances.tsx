import { useEffect, useState } from "react";
import { get } from "../api";
import type { BalanceView, UserView } from "../types";

type Props = {
  refreshKey: number;
};

export default function UserBalances({ refreshKey }: Props) {
  const [users, setUsers] = useState<UserView[]>([]);
  const [userId, setUserId] = useState("");
  const [balances, setBalances] = useState<BalanceView[]>([]);

  useEffect(() => {
    get<UserView[]>("/users").then(setUsers);
  }, []);

  useEffect(() => {
    if (!userId) return;

    get<BalanceView[]>(`/balances/user?user_id=${userId}`)
      .then(setBalances);
  }, [userId, refreshKey]);

  const nameById = (id: string) =>
    users.find(u => u.id === id)?.name ?? id;

  const owes = balances.filter(
    b => b.from_user_id === userId && b.amount > 0
  );

  const owed = balances.filter(
    b => b.to_user_id === userId && b.amount > 0
  );

  return (
    <div className="section">
      <h2>User Balances</h2>

      <select value={userId} onChange={e => setUserId(e.target.value)}>
        <option value="">Select User</option>
        {users.map(u => (
          <option key={u.id} value={u.id}>{u.name}</option>
        ))}
      </select>

      {userId && owes.length === 0 && owed.length === 0 && (
        <p>No outstanding balances 🎉</p>
      )}

      {owes.length > 0 && (
        <>
          <h4>You Owe</h4>
          <table>
            <thead>
              <tr>
                <th>To</th>
                <th>Amount</th>
              </tr>
            </thead>
            <tbody>
              {owes.map((b, i) => (
                <tr key={i}>
                  <td>{nameById(b.to_user_id)}</td>
                  <td>₹ {b.amount}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}

      {owed.length > 0 && (
        <>
          <h4>You Are Owed</h4>
          <table>
            <thead>
              <tr>
                <th>From</th>
                <th>Amount</th>
              </tr>
            </thead>
            <tbody>
              {owed.map((b, i) => (
                <tr key={i}>
                  <td>{nameById(b.from_user_id)}</td>
                  <td>₹ {b.amount}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}
    </div>
  );
}
