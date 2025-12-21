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
    <div className="min-h-screen bg-gray-100 py-10 px-4">
      <div className="max-w-5xl mx-auto space-y-8">
        
        {/* Header */}
        <h1 className="text-4xl font-bold text-center text-gray-800">
          Expense Ledger
        </h1>

        {/* Users & Groups */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-white rounded-2xl shadow p-6">
            <Users />
          </div>

          <div className="bg-white rounded-2xl shadow p-6">
            <Groups />
          </div>
        </div>

        {/* Create Expense */}
        <div className="bg-white rounded-2xl shadow p-6">
          <CreateExpense onSuccess={triggerRefresh} />
        </div>

        {/* Balances */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-white rounded-2xl shadow p-6">
            <UserBalances refreshKey={refreshKey} />
          </div>

          <div className="bg-white rounded-2xl shadow p-6">
            <GroupBalances refreshKey={refreshKey} />
          </div>
        </div>

        {/* Settle Balance */}
        <div className="bg-white rounded-2xl shadow p-6">
          <SettleBalance onSuccess={triggerRefresh} />
        </div>

      </div>
    </div>
  );
}
