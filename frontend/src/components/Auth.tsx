import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { Mail, Lock, User, UserCircle, Sparkles } from 'lucide-react';

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

    if (!isLogin && (!username || !name)) {
      setError('Заполните все поля');
      return;
    }

    if (!isLogin && password.length < 6) {
      setError('Пароль должен быть минимум 6 символов');
      return;
    }

    try {
      if (isLogin) {
        await login(email, password);
      } else {
        await register(username, email, password, name);
      }
    } catch (err: any) {
      setError(err.message || 'Ошибка');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4 relative overflow-hidden">
      {/* Анимированный градиентный фон */}
      <div className="absolute inset-0 bg-gradient-to-br from-blue-100 via-cyan-100 to-blue-50">
        <div className="absolute inset-0 opacity-30">
          <div className="absolute top-[-40%] left-[-20%] w-[600px] h-[600px] rounded-full bg-blue-400 blur-3xl animate-pulse"></div>
          <div className="absolute bottom-[-40%] right-[-20%] w-[600px] h-[600px] rounded-full bg-cyan-400 blur-3xl animate-pulse delay-1000"></div>
        </div>
      </div>

      <div className="relative w-full max-w-md">
        <div className="bg-white/80 backdrop-blur-xl rounded-3xl p-8 shadow-2xl border border-white/20">
          <div className="text-center mb-8">
            <div className="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-r from-blue-500 to-cyan-500 rounded-2xl shadow-lg mb-4">
              <Sparkles className="w-10 h-10 text-white" />
            </div>
            <h1 className="text-3xl font-extrabold gradient-text">
              Чирик
            </h1>
            <p className="text-gray-500 mt-1 font-medium">
              {isLogin ? 'Добро пожаловать!' : 'Создайте аккаунт'}
            </p>
          </div>

          {error && (
            <div className="bg-red-50 border-l-4 border-red-500 text-red-700 p-3 rounded-xl mb-4 text-sm font-medium">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            {!isLogin && (
              <>
                <div className="relative">
                  <User className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-blue-400" />
                  <input
                    type="text"
                    placeholder="Имя пользователя"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    className="input-field pl-12"
                    required
                  />
                </div>
                <div className="relative">
                  <UserCircle className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-blue-400" />
                  <input
                    type="text"
                    placeholder="Ваше имя"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    className="input-field pl-12"
                    required
                  />
                </div>
              </>
            )}

            <div className="relative">
              <Mail className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-blue-400" />
              <input
                type="email"
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="input-field pl-12"
                required
              />
            </div>

            <div className="relative">
              <Lock className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-blue-400" />
              <input
                type="password"
                placeholder="Пароль"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="input-field pl-12"
                required
              />
              {!isLogin && (
                <p className="text-xs text-gray-400 mt-1 pl-4">Минимум 6 символов</p>
              )}
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-gradient-to-r from-blue-500 to-cyan-500 text-white py-3 rounded-xl font-semibold hover:opacity-90 transition disabled:opacity-50"
            >
              {loading ? (
                <span className="flex items-center justify-center gap-2">
                  <span className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin" />
                  Загрузка...
                </span>
              ) : (
                isLogin ? 'Войти' : 'Зарегистрироваться'
              )}
            </button>
          </form>

          <div className="mt-6 text-center">
            <button
              onClick={() => {
                setIsLogin(!isLogin);
                setError('');
              }}
              className="text-blue-600 hover:text-blue-800 transition font-semibold"
            >
              {isLogin ? 'Нет аккаунта? → Регистрация' : 'Уже есть аккаунт? → Вход'}
            </button>
          </div>

          <div className="mt-4 text-center text-xs text-gray-400">
            {isLogin ? 'Войдите, чтобы продолжить' : 'Присоединяйтесь к сообществу'}
          </div>
        </div>
      </div>
    </div>
  );
};
