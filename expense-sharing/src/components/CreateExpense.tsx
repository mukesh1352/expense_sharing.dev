import { useEffect, useState } from "react";
import { get, post } from "../api";
import type {
  ExpenseInput,
  UserView,
  GroupView,
  SplitInput,
} from "../types";

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

  /* ---------- Load groups ---------- */
  useEffect(() => {
    get<GroupView[]>("/groups").then(setGroups);
  }, []);

  /* ---------- Load members ---------- */
  useEffect(() => {
    if (!expense.group_id) return;
    get<UserView[]>(`/groups/members?group_id=${expense.group_id}`)
      .then(setMembers);
  }, [expense.group_id]);

  /* ---------- Group change ---------- */
  const onGroupChange = (groupId: string) => {
    setMembers([]);
    setExpense(prev => ({
      ...prev,
      group_id: groupId,
      paid_by: "",
      participants: [],
      splits: [],
    }));
  };

  /* ---------- Split updater ---------- */
  const updateSplit = (
    userId: string,
    field: "amount" | "percentage",
    value: number
  ) => {
    setExpense(prev => {
      const splits = [...prev.splits];
      const idx = splits.findIndex(s => s.user_id === userId);

      const updated: SplitInput = {
        user_id: userId,
        [field]: value,
      };

      if (idx >= 0) splits[idx] = updated;
      else splits.push(updated);

      return { ...prev, splits };
    });
  };

  /* ---------- Submit ---------- */
  const submit = async () => {
    if (!expense.group_id) return alert("Select a group");
    if (!expense.paid_by) return alert("Select who paid");
    if (expense.participants.length === 0)
      return alert("Select at least one participant");

    if (
      expense.split_type !== "EQUAL" &&
      expense.splits.length === 0
    ) {
      return alert("Provide split values");
    }

    await post<void>("/expenses", expense);
    alert("Expense created");
    onSuccess();

    setExpense(prev => ({
      ...prev,
      expense_id: crypto.randomUUID(),
      total_amount: 0,
      paid_by: "",
      participants: [],
      splits: [],
    }));
  };

  /* ---------- UI ---------- */
  return (
    <div className="section">
      <h2>Create Expense</h2>

      {/* Expense Details */}
      <fieldset style={{ marginBottom: "16px" }}>
        <legend>Expense Details</legend>

        <div className="form-row">
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

          <select
            value={expense.paid_by}
            disabled={!expense.group_id}
            onChange={e =>
              setExpense(prev => ({
                ...prev,
                paid_by: e.target.value,
              }))
            }
          >
            <option value="">Paid By</option>
            {members.map(u => (
              <option key={u.id} value={u.id}>
                {u.name}
              </option>
            ))}
          </select>

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

          <select
            value={expense.split_type}
            onChange={e =>
              setExpense(prev => ({
                ...prev,
                split_type: e.target.value as ExpenseInput["split_type"],
                splits: [],
              }))
            }
          >
            <option value="EQUAL">Equal Split</option>
            <option value="EXACT">Exact Split</option>
            <option value="PERCENT">Percentage Split</option>
          </select>
        </div>
      </fieldset>

      {/* Participants */}
      <fieldset style={{ marginBottom: "16px" }}>
        <legend>Participants</legend>

        {members.map(u => (
          <label key={u.id} style={{ display: "block", marginBottom: "4px" }}>
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
            />{" "}
            {u.name}
          </label>
        ))}
      </fieldset>

      {/* Split Inputs */}
      {expense.split_type !== "EQUAL" && (
        <fieldset style={{ marginBottom: "16px" }}>
          <legend>
            {expense.split_type === "EXACT"
              ? "Exact Amounts"
              : "Percentages"}
          </legend>

          {expense.participants.map(userId => {
            const user = members.find(u => u.id === userId);
            if (!user) return null;

            return (
              <div key={userId} style={{ marginBottom: "6px" }}>
                <label>
                  {user.name}{" "}
                  <input
                    type="number"
                    placeholder={
                      expense.split_type === "EXACT"
                        ? "Amount"
                        : "Percentage"
                    }
                    onChange={e =>
                      updateSplit(
                        userId,
                        expense.split_type === "EXACT"
                          ? "amount"
                          : "percentage",
                        +e.target.value
                      )
                    }
                  />
                </label>
              </div>
            );
          })}
        </fieldset>
      )}

      <button type="button" onClick={submit}>
        Create Expense
      </button>
    </div>
  );
}
