import React, { useState, useEffect } from 'react';
import { Send, Wallet, CreditCard } from 'lucide-react';

export default function TransferForm({ wallets, onTransfer, getCurrencySymbol }) {
  const [fromWalletNumber, setFromWalletNumber] = useState('');
  const [toWalletNumber, setToWalletNumber] = useState('');
  const [transferAmount, setTransferAmount] = useState('');
  const [idempotencyKey, setIdempotencyKey] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const generateIdempotencyKey = () => {
    const key = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
      const r = (Math.random() * 16) | 0;
      const v = c === 'x' ? r : (r & 0x3) | 0x8;
      return v.toString(16);
    });
    setIdempotencyKey(key);
  };

  useEffect(() => {
    generateIdempotencyKey();
  }, []);

  useEffect(() => {
    if (wallets.length > 0 && !fromWalletNumber) {
      setFromWalletNumber(wallets[0].wallet_number?.toString() || '');
    }
  }, [wallets, fromWalletNumber]);

  const selectedWallet = wallets.find((w) => w.wallet_number?.toString() === fromWalletNumber);
  const walletBalance = selectedWallet ? selectedWallet.balance / 100 : 0;
  const enteredAmount = parseFloat(transferAmount) || 0;
  const isInsufficient = !!(selectedWallet && enteredAmount > walletBalance);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (isSubmitting || isInsufficient) return;
    setIsSubmitting(true);

    const success = await onTransfer(
      parseInt(fromWalletNumber, 10),
      parseInt(toWalletNumber, 10),
      transferAmount,
      idempotencyKey
    );
    if (success) {
      setFromWalletNumber(wallets[0]?.wallet_number?.toString() || '');
      setToWalletNumber('');
      setTransferAmount('');
    }
    generateIdempotencyKey();
    setIsSubmitting(false);
  };

  return (
    <div className="glass-panel rounded-3xl p-6 sticky top-28">
      <h3 className="text-lg font-semibold mb-2 flex items-center gap-2">
        <Send className="w-5 h-5 text-accent" /> Transfer Funds
      </h3>
      <p className="text-xs text-white/50 mb-6">Send multi-currency balances instantly.</p>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="transfer-from" className="block text-sm text-white/60 mb-2">Source Wallet</label>
          <select
            id="transfer-from"
            required
            value={fromWalletNumber}
            onChange={(e) => setFromWalletNumber(e.target.value)}
            className="glass-input w-full bg-[#1e293b] text-white cursor-pointer"
          >
            <option value="" disabled>Select sending account...</option>
            {wallets.map((wallet) => (
              <option key={wallet.id} value={wallet.wallet_number?.toString() || ''}>
                {wallet.currency} Wallet ({getCurrencySymbol(wallet.currency)}
                {(wallet.balance / 100).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}) - No: {wallet.wallet_number}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label htmlFor="transfer-to" className="block text-sm text-white/60 mb-2">Destination Wallet Number</label>
          <div className="relative">
            <Wallet className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-white/40" />
            <input
              type="number"
              id="transfer-to"
              required
              placeholder="e.g. 1000000000"
              value={toWalletNumber}
              onChange={(e) => setToWalletNumber(e.target.value)}
              className="glass-input w-full !pl-11"
            />
          </div>
        </div>

        <div>
          <label htmlFor="transfer-amount" className="block text-sm text-white/60 mb-2">Amount</label>
          <div className="relative">
            <CreditCard className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-white/40" />
            <input
              type="number"
              id="transfer-amount"
              required
              min="0.01"
              step="any"
              placeholder="0.00"
              value={transferAmount}
              onChange={(e) => setTransferAmount(e.target.value)}
              className="glass-input w-full !pl-11 text-lg font-medium"
            />
          </div>
          {isInsufficient && (
            <p className="text-xs text-red-400 mt-1.5 font-medium">
              Insufficient funds in the selected wallet.
            </p>
          )}
        </div>

        <button 
          type="submit" 
          disabled={isSubmitting || isInsufficient}
          className="glass-button w-full flex items-center justify-center gap-2 mt-6 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span>{isSubmitting ? 'Sending...' : 'Send Balance'}</span>
          <Send className="w-4 h-4" />
        </button>
      </form>
    </div>
  );
}
