import React, { useState, useEffect, useCallback } from 'react';
import { CreditCard } from 'lucide-react';
import api from '../lib/axios';
import WalletsGrid from './WalletsGrid';
import CreateWalletModal from './CreateWalletModal';

export default function WalletsSection({ user, onWalletsFetched, selectedWalletId, onSelectWallet, refreshKey, addToast }) {
  const [wallets, setWallets] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedCurrency, setSelectedCurrency] = useState('');

  const fetchWallets = useCallback(async () => {
    try {
      const response = await api.get('/wallets');
      const loadedWallets = response.data.wallets || [];
      setWallets(loadedWallets);
      onWalletsFetched(loadedWallets);
    } catch (error) {
      addToast(error.response?.data?.error || 'Failed to fetch wallets', 'error');
    } finally {
      setIsLoading(false);
    }
  }, [onWalletsFetched, addToast]);

  useEffect(() => {
    fetchWallets();
  }, [fetchWallets, refreshKey]);

  const handleCreateWallet = async (currency) => {
    try {
      await api.post('/wallets', {
        user_id: user.id,
        currency,
      });
      addToast(`${currency} wallet activated successfully`, 'success');
      setIsModalOpen(false);
      setSelectedCurrency('');
      fetchWallets();
    } catch (error) {
      addToast(error.response?.data?.error || 'Failed to activate wallet', 'error');
    }
  };

  const handleCopy = (text) => {
    navigator.clipboard.writeText(text);
    addToast('Wallet number copied to clipboard', 'success');
  };

  const handleDeleteWallet = async (walletId) => {
    if (!window.confirm('Are you sure you want to delete this wallet?')) return;
    try {
      await api.delete(`/wallets/${walletId}`);
      addToast('Wallet deleted successfully', 'success');
      fetchWallets();
    } catch (error) {
      addToast(error.response?.data?.error || 'Failed to delete wallet', 'error');
    }
  };

  const getCurrencySymbol = (currency) => {
    switch (currency?.toUpperCase()) {
      case 'TRY':
        return '₺';
      case 'USD':
        return '$';
      case 'EUR':
        return '€';
      default:
        return '¤';
    }
  };

  const getCurrencyIcon = (currency) => {
    switch (currency?.toUpperCase()) {
      case 'TRY':
        return 'fa-lira-sign';
      case 'USD':
        return 'fa-dollar-sign';
      case 'EUR':
        return 'fa-euro-sign';
      default:
        return 'fa-wallet';
    }
  };

  const activeCurrencies = wallets.map((w) => w?.currency?.toUpperCase() || '');
  const missingCurrencies = ['TRY', 'USD', 'EUR'].filter((c) => !activeCurrencies.includes(c));

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
        <CreditCard className="w-5 h-5 text-primary" /> My Wallets
      </h2>
      {isLoading ? (
        <div className="text-white/60">Loading wallets...</div>
      ) : (
        <WalletsGrid
          wallets={wallets}
          onAddWalletClick={(currency) => {
            setSelectedCurrency(currency);
            setIsModalOpen(true);
          }}
          onCopy={handleCopy}
          onDelete={handleDeleteWallet}
          selectedWalletId={selectedWalletId}
          onSelectWallet={onSelectWallet}
          getCurrencySymbol={getCurrencySymbol}
          getCurrencyIcon={getCurrencyIcon}
        />
      )}

      <CreateWalletModal
        isOpen={isModalOpen}
        onClose={() => {
          setIsModalOpen(false);
          setSelectedCurrency('');
        }}
        selectedCurrency={selectedCurrency}
        setSelectedCurrency={setSelectedCurrency}
        onCreateWallet={handleCreateWallet}
        missingCurrencies={missingCurrencies}
      />
    </div>
  );
}
