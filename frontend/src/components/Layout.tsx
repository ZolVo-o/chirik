import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

export const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b sticky top-0 z-50 shadow-sm">
        <div className="max-w-5xl mx-auto px-4 h-16 flex items-center justify-between">
          <Link to="/" className="text-2xl font-bold text-blue-500 hover:text-blue-600 transition flex items-center gap-2">
            🐦 <span className="hidden sm:inline">Чирик</span>
          </Link>
          
          <div className="flex items-center gap-4">
            {user && (
              <>
                <Link 
                  to="/" 
                  className="text-gray-600 hover:text-blue-500 transition text-lg"
                  title="Главная"
                >
                  🏠
                </Link>
                <Link 
                  to="/profile" 
                  className="text-gray-600 hover:text-blue-500 transition text-lg"
                  title="Профиль"
                >
                  👤
                </Link>
                <button 
                  onClick={handleLogout} 
                  className="text-gray-600 hover:text-red-500 transition text-lg"
                  title="Выйти"
                >
                  🚪
                </button>
                <div className="flex items-center gap-2 ml-2">
                  <div className="w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center text-white text-sm font-bold">
                    {user.name?.charAt(0) || 'U'}
                  </div>
                  <span className="text-sm font-medium text-gray-700 hidden sm:block">
                    {user.name}
                  </span>
                </div>
              </>
            )}
          </div>
        </div>
      </header>
      <main className="max-w-5xl mx-auto px-4 py-6">
        {children}
      </main>
    </div>
  );
};
