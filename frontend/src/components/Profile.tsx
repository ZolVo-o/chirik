import React, { useEffect, useState, useCallback } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { api } from '../services/api';
import { Tweet } from '../types';
import { useNavigate } from 'react-router-dom';
import { Camera, Edit2, Check, X, Calendar, Heart, Eye } from 'lucide-react';

export const Profile: React.FC = () => {
  const { user, updateUser, logout } = useAuth();
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
      updateUser(updatedUser);
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
    <div className="max-w-3xl mx-auto">
      <div className="surface rounded-[30px] overflow-hidden">
        <div className="h-36 gradient-bg relative"><div className="absolute inset-0 opacity-20 bg-[radial-gradient(circle_at_20%_20%,white,transparent_30%)]" /></div>
        
        <div className="px-6 pb-6 relative">
          <div className="relative -mt-12 mb-4">
            <div className="w-24 h-24 rounded-[28px] gradient-bg flex items-center justify-center text-white text-3xl font-bold border-4 border-[#f4f1eb] shadow-lg shadow-[#ef765e]/25">
              {user.name?.charAt(0) || 'U'}
            </div>
              <button className="absolute bottom-0 right-0 p-1.5 bg-[#17202a] text-white rounded-full hover:bg-[#ef765e] transition shadow-md">
              <Camera className="w-4 h-4" />
            </button>
          </div>

          <div className="flex justify-end gap-2 -mt-12">
            {isEditing ? (
              <>
                <button
                  onClick={handleCancel}
                  className="btn-outline px-4 py-2 border rounded-full transition flex items-center gap-1 text-sm"
                >
                  <X className="w-4 h-4" /> Отмена
                </button>
                <button
                  onClick={handleSave}
                  className="btn-primary px-4 py-2 rounded-full transition flex items-center gap-1 text-sm"
                >
                  <Check className="w-4 h-4" /> Сохранить
                </button>
              </>
            ) : (
              <>
                <button
                  onClick={handleEdit}
                  className="btn-outline px-4 py-2 rounded-full transition flex items-center gap-1 text-sm"
                >
                  <Edit2 className="w-4 h-4" /> Редактировать
                </button>
                <button
                  onClick={handleLogout}
                  className="px-4 py-2 bg-[#17202a] text-white rounded-full hover:bg-[#ef765e] transition flex items-center gap-1 text-sm"
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
                className="w-full px-3 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Имя"
              />
              <textarea
                value={editBio}
                onChange={(e) => setEditBio(e.target.value)}
                className="w-full px-3 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
                rows={2}
                placeholder="О себе"
              />
            </div>
          ) : (
            <>
               <p className="section-label mb-1">Профиль</p>
               <h2 className="text-3xl font-black tracking-[-0.04em] text-[#17202a]">{user.name}</h2>
               <p className="text-[#8d857c]">@{user.username}</p>
               <p className="mt-3 text-[#394550] leading-6">{user.bio || 'Пока ничего не рассказал о себе'}</p>
            </>
          )}

          <div className="flex gap-6 mt-4 text-sm">
            <div>
               <span className="font-black text-[#17202a]">{tweets.length}</span>
               <span className="text-[#8d857c] ml-1">твитов</span>
            </div>
            <div>
               <span className="font-black text-[#17202a]">{followers.length}</span>
               <span className="text-[#8d857c] ml-1">читателей</span>
            </div>
            <div>
               <span className="font-black text-[#17202a]">{following.length}</span>
               <span className="text-[#8d857c] ml-1">читает</span>
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
           <h3 className="section-label mb-4">Мои твиты</h3>
        
        {loading ? (
          <div className="text-center text-gray-400 py-8">Загрузка...</div>
        ) : tweets.length === 0 ? (
           <div className="surface rounded-[26px] p-12 text-center">
            <p className="text-4xl mb-2">🐦</p>
            <p className="text-gray-400">У вас пока нет твитов</p>
          </div>
        ) : (
          <div className="space-y-4">
            {tweets.map((tweet) => (
             <div key={tweet.id} className="surface rounded-[26px] p-5">
                 <p className="text-[#394550] leading-7">{tweet.content}</p>
                 <div className="mt-3 flex items-center gap-4 text-sm text-[#9d958b]">
                  <span className="flex items-center gap-1">
                    <Heart className="w-4 h-4" /> {tweet.likes}
                  </span>
                  <span className="flex items-center gap-1">
                    <Eye className="w-4 h-4" /> {tweet.views || 0}
                  </span>
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
