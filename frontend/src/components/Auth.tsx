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
    <div className="min-h-screen flex items-center justify-center p-4 relative overflow-hidden bg-[#17202a]">
      {/* Анимированный градиентный фон */}
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_right,_#ef765e_0,_transparent_38%),radial-gradient(circle_at_bottom_left,_#687b8b_0,_transparent_35%)] opacity-60">
        <div className="absolute inset-0 opacity-30">
          <div className="absolute top-[-40%] left-[-20%] w-[600px] h-[600px] rounded-full bg-[#ef765e] blur-3xl animate-pulse"></div>
          <div className="absolute bottom-[-40%] right-[-20%] w-[600px] h-[600px] rounded-full bg-[#a84267] blur-3xl animate-pulse delay-1000"></div>
        </div>
      </div>

      <div className="relative w-full max-w-md">
        <div className="bg-[#f8f6f1]/95 backdrop-blur-xl rounded-[32px] p-7 sm:p-9 shadow-2xl border border-white/20">
          <div className="text-center mb-8">
            <div className="inline-flex items-center justify-center w-20 h-20 gradient-bg rounded-[24px] shadow-lg shadow-[#ef765e]/25 mb-5 text-3xl">
              <Sparkles className="w-10 h-10 text-white" />
            </div>
            <h1 className="text-4xl font-black tracking-[-0.05em] text-[#17202a]">
              Чирик
            </h1>
            <p className="text-[#8d857c] mt-2 font-medium">
              {isLogin ? 'Добро пожаловать!' : 'Создайте аккаунт'}
            </p>
          </div>

          {error && (
            <div className="bg-[#fff0ed] border-l-4 border-[#ef765e] text-[#b7464d] p-3 rounded-xl mb-4 text-sm font-medium">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            {!isLogin && (
              <>
                <div className="relative">
                  <User className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-[#ef765e]" />
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
                  <UserCircle className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-[#ef765e]" />
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
                <Mail className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-[#ef765e]" />
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
                <Lock className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-[#ef765e]" />
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
              className="w-full btn-primary py-3 rounded-2xl"
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
              className="text-[#b74668] hover:text-[#ef765e] transition font-semibold"
            >
              {isLogin ? 'Нет аккаунта? → Регистрация' : 'Уже есть аккаунт? → Вход'}
            </button>
          </div>

          <div className="mt-4 text-center text-xs text-[#9d958b]">
            {isLogin ? 'Войдите, чтобы продолжить' : 'Присоединяйтесь к сообществу'}
          </div>
        </div>
      </div>
    </div>
  );
};
