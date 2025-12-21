import { useState } from "react";
import Users from "./components/Users";
import Groups from "./components/Groups";
import CreateExpense from "./components/CreateExpense";
import UserBalances from "./components/UserBalances";
import GroupBalances from "./components/GroupBalances";
import SettleBalance from "./components/SettleBalance";

export default function App() {
  const [refreshKey, setRefreshKey] = useState(0);

  const triggerRefresh = () => {
    setRefreshKey(k => k + 1);
  };

  return (
    <div style={{ padding: 24 }}>
      <h1>Splitwise Clone</h1>

      <Users />
      <Groups />

      <CreateExpense onSuccess={triggerRefresh} />

      <UserBalances refreshKey={refreshKey} />
      <GroupBalances refreshKey={refreshKey} />

      <SettleBalance onSuccess={triggerRefresh} />
    </div>
  );
}
