import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { Home, User, LogOut, Radio, MessageCircle } from 'lucide-react';
import { useRealtime } from '../contexts/RealtimeContext';

export const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, logout } = useAuth();
  const { notifications, dismissNotification } = useRealtime();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  return (
    <div className="min-h-screen bg-[#f4f1eb] relative overflow-hidden">
      <div className="pointer-events-none fixed -top-40 -right-32 h-96 w-96 rounded-full bg-[#ef765e]/10 blur-3xl" />
      <div className="pointer-events-none fixed bottom-0 -left-40 h-96 w-96 rounded-full bg-[#687b8b]/10 blur-3xl" />
      <header className="bg-[#17202a]/95 text-white border-b border-white/10 sticky top-0 z-50 backdrop-blur-xl">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 h-[72px] flex items-center justify-between">
          {/* Логотип слева */}
          <Link to="/" className="flex items-center gap-3 group shrink-0">
            <span className="grid place-items-center h-9 w-9 rounded-xl bg-[#ef765e] text-xl shadow-lg shadow-[#ef765e]/20 group-hover:rotate-6 transition-transform">✦</span>
            <span className="text-xl font-black tracking-tight">Чирик<span className="text-[#ef765e]">.</span></span>
          </Link>

          {/* Кнопки — только справа, компактно */}
          <div className="flex items-center gap-2">
            {user && (
              <>
                <span className="hidden sm:flex items-center gap-2 text-xs text-white/50 mr-2"><Radio className="w-3.5 h-3.5 text-[#ef765e]" /> live</span>
                <Link
                  to="/messages"
                  className="p-2.5 rounded-xl hover:bg-white/10 text-white/60 hover:text-white transition-all duration-200"
                  title="Сообщения"
                >
                  <MessageCircle className="w-5 h-5" />
                </Link>
                <Link
                  to="/"
                   className="p-2.5 rounded-xl hover:bg-white/10 text-white/60 hover:text-white transition-all duration-200"
                  title="Главная"
                >
                  <Home className="w-5 h-5" />
                </Link>
                <Link
                  to="/profile"
                   className="p-2.5 rounded-xl hover:bg-white/10 text-white/60 hover:text-white transition-all duration-200"
                  title="Профиль"
                >
                  <User className="w-5 h-5" />
                </Link>
                <button
                  onClick={handleLogout}
                   className="p-2.5 rounded-xl hover:bg-[#ef765e]/20 text-white/60 hover:text-[#ef765e] transition-all duration-200"
                  title="Выйти"
                >
                  <LogOut className="w-5 h-5" />
                </button>
              </>
            )}
          </div>
        </div>
      </header>

      <div className="fixed right-4 top-20 z-[60] flex w-[min(22rem,calc(100vw-2rem))] flex-col gap-2">
        {notifications.map((notification) => (
          <button
            key={notification.id}
            onClick={() => dismissNotification(notification.id)}
            className="surface rounded-2xl border-l-4 border-[#ef765e] px-4 py-3 text-left text-sm font-semibold text-[#17202a] shadow-xl"
          >
            {notification.text}
          </button>
        ))}
      </div>

      <main className="relative max-w-6xl mx-auto px-4 sm:px-6 py-8 sm:py-12">
        {children}
      </main>
    </div>
  );
};
