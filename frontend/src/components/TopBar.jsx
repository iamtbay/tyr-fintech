import React from 'react';
import { Wallet, LogOut } from 'lucide-react';

export default function TopBar({ user, onLogout }) {
  return (
    <header className="glass-panel border-b-0 border-x-0 rounded-none px-8 py-4 flex justify-between items-center sticky top-0 z-50">
      <div className="flex items-center gap-3">
        <div className="w-10 h-10 rounded-xl bg-gradient-to-tr from-primary to-accent flex items-center justify-center shadow-lg shadow-primary/20">
          <Wallet className="w-6 h-6 text-white" />
        </div>
        <span className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-white to-white/70 tracking-wider">
          TYR FINTECH
        </span>
      </div>
      <div className="flex items-center gap-6">
        <div className="hidden md:flex flex-col text-right">
          <span className="text-sm font-semibold text-white">{user?.name || user?.username}</span>
          <span className="text-xs text-white/50">{user?.email}</span>
        </div>
        <button 
          onClick={onLogout} 
          className="glass-panel bg-red-500/10 hover:bg-red-500/20 text-red-200 border-red-500/30 font-semibold rounded-xl px-4 py-2 transition-colors flex items-center gap-2 text-sm cursor-pointer"
        >
          <LogOut className="w-4 h-4" />
          <span>Logout</span>
        </button>
      </div>
    </header>
  );
}
