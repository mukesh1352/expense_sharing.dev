import { useEffect, useState } from "react";
import { get, post } from "../api";
import type { ExpenseInput, UserView, GroupView } from "../types";

type Props = {
  onSuccess: () => void;
};

export default function CreateExpense({ onSuccess }: Props) {
  const [groups, setGroups] = useState<GroupView[]>([]);
  const [members, setMembers] = useState<UserView[]>([]);

  const [expense, setExpense] = useState<ExpenseInput>({
    expense_id: crypto.randomUUID(),
    group_id: "",
    paid_by: "",
    total_amount: 0,
    split_type: "EQUAL",
    participants: [],
    splits: [],
  });

  // Load groups ONCE
  useEffect(() => {
    get<GroupView[]>("/groups").then(setGroups);
  }, []);

  // Load members ONLY when group_id exists
  useEffect(() => {
    if (!expense.group_id) return;

    get<UserView[]>(
      `/groups/members?group_id=${expense.group_id}`
    ).then(setMembers);
  }, [expense.group_id]);

  const onGroupChange = (groupId: string) => {
    setMembers([]); // safe here
    setExpense(prev => ({
      ...prev,
      group_id: groupId,
      paid_by: "",
      participants: [],
      splits: [],
    }));
  };

  const submit = async () => {
    if (!expense.group_id) {
      alert("Select a group");
      return;
    }

    if (!expense.paid_by) {
      alert("Select who paid");
      return;
    }

    if (expense.participants.length === 0) {
      alert("Select at least one participant");
      return;
    }

    await post<void>("/expenses", expense);
    alert("Expense created");

    onSuccess();

    // Reset form (keep group)
    setExpense(prev => ({
      ...prev,
      expense_id: crypto.randomUUID(),
      total_amount: 0,
      paid_by: "",
      participants: [],
      splits: [],
    }));
  };

  return (
    <div className="section">
      <h2>Create Expense</h2>

      <div className="form-row">
        {/* Group */}
        <select
          value={expense.group_id}
          onChange={e => onGroupChange(e.target.value)}
        >
          <option value="">Select Group</option>
          {groups.map(g => (
            <option key={g.id} value={g.id}>
              {g.name}
            </option>
          ))}
        </select>

        {/* Paid By */}
        <select
          value={expense.paid_by}
          disabled={!expense.group_id}
          onChange={e =>
            setExpense(prev => ({ ...prev, paid_by: e.target.value }))
          }
        >
          <option value="">Paid By</option>
          {members.map(u => (
            <option key={u.id} value={u.id}>
              {u.name}
            </option>
          ))}
        </select>

        {/* Amount */}
        <input
          type="number"
          placeholder="Total Amount"
          value={expense.total_amount || ""}
          onChange={e =>
            setExpense(prev => ({
              ...prev,
              total_amount: +e.target.value,
            }))
          }
        />

        {/* Split Type */}
        <select
          value={expense.split_type}
          onChange={e =>
            setExpense(prev => ({
              ...prev,
              split_type: e.target.value as ExpenseInput["split_type"],
            }))
          }
        >
          <option value="EQUAL">Equal</option>
          <option value="EXACT">Exact</option>
          <option value="PERCENT">Percentage</option>
        </select>
      </div>

      {/* Participants */}
      <h4>Participants</h4>

      {members.length === 0 && expense.group_id && (
        <p>No members in this group</p>
      )}

      {members.map(u => (
        <label key={u.id} style={{ display: "block" }}>
          <input
            type="checkbox"
            checked={expense.participants.includes(u.id)}
            onChange={e =>
              setExpense(prev => ({
                ...prev,
                participants: e.target.checked
                  ? [...prev.participants, u.id]
                  : prev.participants.filter(p => p !== u.id),
              }))
            }
          />
          {u.name}
        </label>
      ))}

      <button onClick={submit}>Create Expense</button>
    </div>
  );
}
