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
    <div className="min-h-screen bg-gradient-to-b from-blue-50/50 to-white">
      {/* Прозрачная панель */}
      <header className="bg-white/40 backdrop-blur-md border-b border-blue-100/30 sticky top-0 z-50">
        <div className="max-w-4xl mx-auto px-4 h-14 flex items-center justify-between">
          {/* Логотип слева */}
          <Link to="/" className="flex items-center gap-2 group shrink-0">
            <Sparkles className="w-6 h-6 text-blue-500 group-hover:rotate-12 transition-transform duration-300" />
            <span className="text-lg font-extrabold gradient-text">
              Чирик
            </span>
          </Link>

          {/* Кнопки — только справа, компактно */}
          <div className="flex items-center gap-0.5">
            {user && (
              <>
                <Link
                  to="/"
                  className="p-2 rounded-xl hover:bg-blue-100/60 text-gray-600 hover:text-blue-600 transition-all duration-200"
                  title="Главная"
                >
                  <Home className="w-5 h-5" />
                </Link>
                <Link
                  to="/profile"
                  className="p-2 rounded-xl hover:bg-blue-100/60 text-gray-600 hover:text-blue-600 transition-all duration-200"
                  title="Профиль"
                >
                  <User className="w-5 h-5" />
                </Link>
                <button
                  onClick={handleLogout}
                  className="p-2 rounded-xl hover:bg-red-100/60 text-gray-600 hover:text-red-500 transition-all duration-200"
                  title="Выйти"
                >
                  <LogOut className="w-5 h-5" />
                </button>
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
