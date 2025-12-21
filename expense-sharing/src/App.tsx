import { useState } from "react";

import Users from "./components/Users";
import Groups from "./components/Groups";
import CreateExpense from "./components/CreateExpense";
import UserBalances from "./components/UserBalances";
import GroupBalances from "./components/GroupBalances";
import SettleBalance from "./components/SettleBalance";

export default function ExpenseLedgerApp() {
  const [refreshKey, setRefreshKey] = useState(0);

  const refreshLedger = () => {
    setRefreshKey(prev => prev + 1);
  };

  return (
    <main style={{ padding: "24px", maxWidth: "1000px", margin: "0 auto" }}>
      <header style={{ marginBottom: "32px" }}>
        <h1>Expense-Ledger</h1>
        <p style={{ color: "#666" }}>
          Track shared expenses, balances, and settlements
        </p>
      </header>

      {/* Reference Data */}
      <section style={{ marginBottom: "40px" }}>
        <Users />
        <Groups />
      </section>

      {/* Expense Creation */}
      <section style={{ marginBottom: "40px" }}>
        <CreateExpense onSuccess={refreshLedger} />
      </section>

      {/* Ledger Views */}
      <section style={{ marginBottom: "40px" }}>
        <UserBalances refreshKey={refreshKey} />
        <GroupBalances refreshKey={refreshKey} />
      </section>

      {/* Settlements */}
      <section style={{ marginBottom: "40px" }}>
        <SettleBalance onSuccess={refreshLedger} />
      </section>
    </main>
  );
}
