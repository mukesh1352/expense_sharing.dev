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
    setError("");

    if (!data.from_user_id || !data.to_user_id) {
      setError("Please select both users");
      return;
    }

    if (data.from_user_id === data.to_user_id) {
      setError("Cannot settle balance with the same user");
      return;
    }

    if (data.amount <= 0) {
      setError("Amount must be greater than 0");
      return;
    }

    try {
      await post<void>("/settle", data);

      alert("Settlement successful");
      onSuccess();

      // Reset form after success
      setData({
        from_user_id: "",
        to_user_id: "",
        amount: 0,
      });
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

      <div className="form-row">
        <select
          value={data.from_user_id}
          onChange={e =>
            setData(prev => ({ ...prev, from_user_id: e.target.value }))
          }
        >
          <option value="">From</option>
          {users.map(u => (
            <option key={u.id} value={u.id}>
              {u.name}
            </option>
          ))}
        </select>

        <select
          value={data.to_user_id}
          onChange={e =>
            setData(prev => ({ ...prev, to_user_id: e.target.value }))
          }
        >
          <option value="">To</option>
          {users.map(u => (
            <option key={u.id} value={u.id}>
              {u.name}
            </option>
          ))}
        </select>

        <input
          type="number"
          placeholder="Amount"
          value={data.amount || ""}
          onChange={e =>
            setData(prev => ({ ...prev, amount: +e.target.value }))
          }
        />

        {/* 🔥 Explicit button type prevents GET requests */}
        <button type="button" onClick={submit}>
          Settle
        </button>
      </div>

      {error && (
        <p style={{ color: "red", marginTop: "8px" }}>
          ❌ {error}
        </p>
      )}
    </div>
  );
}
