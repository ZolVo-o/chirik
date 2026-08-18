import React, { useEffect, useState, useCallback } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { api } from '../services/api';
import { Tweet, User } from '../types';
import { useNavigate } from 'react-router-dom';
import { Camera, Edit2, Check, X, Calendar } from 'lucide-react';

export const Profile: React.FC = () => {
  const { user, logout, login } = useAuth();
  const navigate = useNavigate();
  const [tweets, setTweets] = useState<Tweet[]>([]);
  const [loading, setLoading] = useState(true);
  const [isEditing, setIsEditing] = useState(false);
  const [editName, setEditName] = useState('');
  const [editBio, setEditBio] = useState('');
  const [followers, setFollowers] = useState<any[]>([]);
  const [following, setFollowing] = useState<any[]>([]);

  const loadTweets = useCallback(async () => {
    if (!user) return;
    try {
      const data = await api.getTweetsByUser(user.id);
      setTweets(data || []);
    } catch (error) {
      console.error('Error loading tweets:', error);
      setTweets([]);
    } finally {
      setLoading(false);
    }
  }, [user]);

  const loadFollowStats = useCallback(async () => {
    if (!user) return;
    try {
      const [followersData, followingData] = await Promise.all([
        api.getFollowers(user.id),
        api.getFollowing(user.id)
      ]);
      setFollowers(followersData || []);
      setFollowing(followingData || []);
    } catch (error) {
      console.error('Error loading follow stats:', error);
    }
  }, [user]);

  useEffect(() => {
    loadTweets();
    loadFollowStats();
  }, [loadTweets, loadFollowStats]);

  const handleEdit = () => {
    setEditName(user?.name || '');
    setEditBio(user?.bio || '');
    setIsEditing(true);
  };

  const handleSave = async () => {
    try {
      const updatedUser = await api.updateProfile(editName, editBio);
      // Обновляем пользователя в контексте
      if (user) {
        user.name = updatedUser.name;
        user.bio = updatedUser.bio;
      }
      setIsEditing(false);
    } catch (error) {
      console.error('Error updating profile:', error);
      alert('Ошибка при обновлении профиля');
    }
  };

  const handleCancel = () => {
    setIsEditing(false);
  };

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  if (!user) return null;

  return (
    <div className="max-w-3xl mx-auto p-4">
      <div className="bg-white rounded-xl shadow overflow-hidden">
        <div className="h-32 bg-gradient-to-r from-blue-500 to-cyan-500"></div>
        
        <div className="px-6 pb-6 relative">
          <div className="relative -mt-12 mb-4">
            <div className="w-24 h-24 rounded-full bg-gradient-to-r from-blue-500 to-cyan-500 flex items-center justify-center text-white text-3xl font-bold border-4 border-white shadow-lg">
              {user.name?.charAt(0) || 'U'}
            </div>
            <button className="absolute bottom-0 right-0 p-1.5 bg-blue-500 text-white rounded-full hover:bg-blue-600 transition">
              <Camera className="w-4 h-4" />
            </button>
          </div>

          <div className="flex justify-end gap-2 -mt-12">
            {isEditing ? (
              <>
                <button
                  onClick={handleCancel}
                  className="px-4 py-2 border border-gray-300 text-gray-700 rounded-full hover:bg-gray-50 transition flex items-center gap-1"
                >
                  <X className="w-4 h-4" /> Отмена
                </button>
                <button
                  onClick={handleSave}
                  className="px-4 py-2 bg-blue-500 text-white rounded-full hover:bg-blue-600 transition flex items-center gap-1"
                >
                  <Check className="w-4 h-4" /> Сохранить
                </button>
              </>
            ) : (
              <>
                <button
                  onClick={handleEdit}
                  className="px-4 py-2 border border-gray-300 text-gray-700 rounded-full hover:bg-gray-50 transition flex items-center gap-1"
                >
                  <Edit2 className="w-4 h-4" /> Редактировать
                </button>
                <button
                  onClick={handleLogout}
                  className="px-4 py-2 bg-red-500 text-white rounded-full hover:bg-red-600 transition flex items-center gap-1"
                >
                  <X className="w-4 h-4" /> Выйти
                </button>
              </>
            )}
          </div>

          {isEditing ? (
            <div className="space-y-2">
              <input
                type="text"
                value={editName}
                onChange={(e) => setEditName(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Имя"
              />
              <textarea
                value={editBio}
                onChange={(e) => setEditBio(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
                rows={2}
                placeholder="О себе"
              />
            </div>
          ) : (
            <>
              <h2 className="text-2xl font-bold">{user.name}</h2>
              <p className="text-gray-500">@{user.username}</p>
              <p className="mt-2 text-gray-700">{user.bio || '✍️ Пока ничего не рассказал о себе'}</p>
            </>
          )}

          <div className="flex gap-6 mt-4 text-sm">
            <div>
              <span className="font-bold">{tweets.length}</span>
              <span className="text-gray-500 ml-1">твитов</span>
            </div>
            <div>
              <span className="font-bold">{followers.length}</span>
              <span className="text-gray-500 ml-1">читателей</span>
            </div>
            <div>
              <span className="font-bold">{following.length}</span>
              <span className="text-gray-500 ml-1">читает</span>
            </div>
          </div>

          <div className="mt-3 text-sm text-gray-400 flex items-center gap-2">
            <Calendar className="w-4 h-4" />
            <span>Присоединился: {new Date(user.created_at).toLocaleDateString('ru', {
              day: 'numeric',
              month: 'long',
              year: 'numeric'
            })}</span>
          </div>
        </div>
      </div>

      <div className="mt-6">
        <h3 className="text-xl font-bold mb-4">Твиты</h3>
        
        {loading ? (
          <div className="text-center text-gray-400 py-8">Загрузка...</div>
        ) : tweets.length === 0 ? (
          <div className="text-center text-gray-400 py-12 bg-white rounded-xl shadow">
            <p className="text-4xl mb-2">🐦</p>
            <p>У вас пока нет твитов</p>
            <p className="text-sm">Напишите что-нибудь!</p>
          </div>
        ) : (
          <div className="space-y-3">
            {tweets.map((tweet) => (
              <div key={tweet.id} className="bg-white rounded-xl shadow p-4 hover:shadow-md transition">
                <p className="text-gray-800">{tweet.content}</p>
                <div className="mt-2 flex items-center gap-4 text-sm text-gray-400">
                  <span>❤️ {tweet.likes}</span>
                  <span>{new Date(tweet.created_at).toLocaleString('ru')}</span>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
