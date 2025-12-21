import { useEffect, useState } from "react";
import { get, post } from "../api";
import type { SettlementInput, UserView } from "../types";

type Props = {
  onSuccess: () => void;
};

export default function SettleBalance({ onSuccess }: Props) {
  const [users, setUsers] = useState<UserView[]>([]);
  const [error, setError] = useState("");
  const [data, setData] = useState<SettlementInput>({
    from_user_id: "",
    to_user_id: "",
    amount: 0,
  });

  useEffect(() => {
    get<UserView[]>("/users").then(setUsers);
  }, []);

  const submit = async () => {
    try {
      if (!data.from_user_id || !data.to_user_id || data.amount <= 0) {
        setError("Invalid settlement input");
        return;
      }

      await post<void>("/settle", data);

      setError("");
      alert("Settlement successful");
      onSuccess();
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Something went wrong");
      }
    }
  };

  return (
    <div className="section">
      <h2>Settle Balance</h2>

      <select
        value={data.from_user_id}
        onChange={e =>
          setData({ ...data, from_user_id: e.target.value })
        }
      >
        <option value="">From</option>
        {users.map(u => (
          <option key={u.id} value={u.id}>{u.name}</option>
        ))}
      </select>

      <select
        value={data.to_user_id}
        onChange={e =>
          setData({ ...data, to_user_id: e.target.value })
        }
      >
        <option value="">To</option>
        {users.map(u => (
          <option key={u.id} value={u.id}>{u.name}</option>
        ))}
      </select>

      <input
        type="number"
        placeholder="Amount"
        onChange={e =>
          setData({ ...data, amount: +e.target.value })
        }
      />

      <button onClick={submit}>Settle</button>

      {error && (
        <p style={{ color: "red", marginTop: "8px" }}>
           {error}
        </p>
      )}
    </div>
  );
}
