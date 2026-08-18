import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';

export const Auth: React.FC = () => {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [username, setUsername] = useState('');
  const [name, setName] = useState('');
  const [error, setError] = useState('');
  const { login, register, loading } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!email || !password) {
      setError('Заполните все поля');
      return;
    }

    try {
      if (isLogin) {
        await login(email, password);
      } else {
        if (!username || !name) {
          setError('Заполните все поля');
          return;
        }
        if (password.length < 6) {
          setError('Пароль должен быть минимум 6 символов');
          return;
        }
        await register(username, email, password, name);
      }
    } catch (err: any) {
      setError(err.message || 'Ошибка');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-100 to-cyan-100 p-4">
      <div className="bg-white rounded-3xl shadow-xl p-8 w-full max-w-md">
        <div className="text-center mb-8">
          <div className="text-5xl mb-3">🐦</div>
          <h1 className="text-3xl font-bold gradient-text">Чирик</h1>
          <p className="text-gray-500 mt-1">
            {isLogin ? 'Войдите в свой аккаунт' : 'Создайте новый аккаунт'}
          </p>
        </div>

        {error && (
          <div className="bg-red-50 text-red-600 p-3 rounded-xl mb-4 text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          {!isLogin && (
            <>
              <input
                type="text"
                placeholder="Имя пользователя"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="w-full px-4 py-3 mb-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
              />
              <input
                type="text"
                placeholder="Ваше имя"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full px-4 py-3 mb-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
              />
            </>
          )}

          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full px-4 py-3 mb-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />

          <input
            type="password"
            placeholder="Пароль"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full px-4 py-3 mb-4 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />

          <button
            type="submit"
            disabled={loading}
            className="w-full gradient-bg text-white py-3 rounded-xl font-semibold hover:opacity-90 transition disabled:opacity-50"
          >
            {loading ? 'Загрузка...' : (isLogin ? 'Войти' : 'Зарегистрироваться')}
          </button>
        </form>

        <div className="mt-6 text-center">
          <button
            onClick={() => setIsLogin(!isLogin)}
            className="text-blue-600 hover:text-blue-800 transition font-medium"
          >
            {isLogin ? 'Нет аккаунта? → Регистрация' : 'Уже есть аккаунт? → Вход'}
          </button>
        </div>
      </div>
    </div>
  );
};
