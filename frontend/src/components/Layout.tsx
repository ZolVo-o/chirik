import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { Home, User, LogOut, Sparkles } from 'lucide-react';

export const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-blue-50 to-white">
      <header className="bg-white/80 backdrop-blur-md border-b border-blue-100/50 sticky top-0 z-50 shadow-sm">
        <div className="max-w-4xl mx-auto px-4 h-16 flex items-center justify-between">
          <Link to="/" className="flex items-center gap-2 group">
            <Sparkles className="w-8 h-8 text-blue-500 group-hover:rotate-12 transition-transform duration-300" />
            <span className="text-xl font-extrabold gradient-text">
              Чирик
            </span>
          </Link>

          <div className="flex items-center gap-1">
            {user && (
              <>
                <Link
                  to="/"
                  className="p-2 rounded-xl hover:bg-blue-50 text-gray-600 hover:text-blue-600 transition-all duration-200"
                  title="Главная"
                >
                  <Home className="w-5 h-5" />
                </Link>
                <Link
                  to="/profile"
                  className="p-2 rounded-xl hover:bg-blue-50 text-gray-600 hover:text-blue-600 transition-all duration-200"
                  title="Профиль"
                >
                  <User className="w-5 h-5" />
                </Link>
                <button
                  onClick={handleLogout}
                  className="p-2 rounded-xl hover:bg-red-50 text-gray-600 hover:text-red-500 transition-all duration-200"
                  title="Выйти"
                >
                  <LogOut className="w-5 h-5" />
                </button>
                <div className="ml-2 flex items-center gap-2">
                  <div className="w-9 h-9 rounded-full bg-gradient-to-r from-blue-500 to-cyan-500 flex items-center justify-center text-white text-sm font-bold shadow-md shadow-blue-500/25">
                    {user.name?.charAt(0).toUpperCase() || 'U'}
                  </div>
                  <span className="text-sm font-semibold text-gray-700 hidden sm:block">
                    {user.name}
                  </span>
                </div>
              </>
            )}
          </div>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-4 py-6">
        {children}
      </main>
    </div>
  );
};
