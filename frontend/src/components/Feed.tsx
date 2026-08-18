import React, { useEffect, useState } from 'react';
import { api } from '../services/api';
import { Tweet } from '../types';
import { Sparkles, Heart, Eye, Send } from 'lucide-react';

export const Feed: React.FC = () => {
  const [tweets, setTweets] = useState<Tweet[]>([]);
  const [loading, setLoading] = useState(true);
  const [content, setContent] = useState('');

  useEffect(() => {
    loadTweets();
  }, []);

  const loadTweets = async () => {
    try {
      const data = await api.getTweets();
      setTweets(data || []);
    } catch (error) {
      console.error('Error loading tweets:', error);
      setTweets([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim()) return;
    try {
      await api.createTweet(content);
      setContent('');
      loadTweets();
    } catch (error) {
      console.error('Error creating tweet:', error);
    }
  };

  const handleView = async (tweetId: number) => {
    try {
      await api.viewTweet(tweetId);
      loadTweets();
    } catch (error) {
      console.error('Error viewing tweet:', error);
    }
  };

  if (loading) {
    return <div className="flex justify-center p-8 text-gray-400">Загрузка...</div>;
  }

  return (
    <div className="max-w-2xl mx-auto">
      {/* Создание поста */}
      <div className="bg-white rounded-2xl shadow-sm p-4 mb-6 border border-blue-100/50">
        <div className="flex items-center gap-3 mb-3">
          <div className="w-10 h-10 rounded-full bg-gradient-to-r from-blue-500 to-cyan-500 flex items-center justify-center text-white font-bold text-sm shadow-md shadow-blue-500/25">
            ✍️
          </div>
          <h2 className="text-lg font-bold text-gray-800">Что нового?</h2>
        </div>
        <form onSubmit={handleSubmit}>
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="Что у вас нового?"
            className="w-full p-3 bg-gray-50 border border-gray-200 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-4 focus:ring-blue-500/20 transition-all duration-300 resize-none"
            rows={3}
            maxLength={280}
          />
          <div className="flex justify-between items-center mt-2">
            <span className="text-sm text-gray-400">{content.length}/280</span>
            <button
              type="submit"
              className="bg-gradient-to-r from-blue-500 to-cyan-500 text-white px-6 py-2 rounded-full font-semibold hover:opacity-90 transition flex items-center gap-2"
            >
              <Send className="w-4 h-4" />
              Чирикнуть
            </button>
          </div>
        </form>
      </div>

      {/* Лента */}
      {tweets.length === 0 ? (
        <div className="bg-white rounded-2xl shadow-sm p-12 text-center border border-blue-100/50">
          <div className="text-6xl mb-4">🐦</div>
          <h3 className="text-xl font-bold text-gray-700">Нет твитов</h3>
          <p className="text-gray-400 mt-2 font-medium">Будьте первым!</p>
        </div>
      ) : (
        <div className="space-y-4">
          {tweets.map((tweet) => (
            <div
              key={tweet.id}
              className="bg-white rounded-2xl shadow-sm p-5 border border-blue-100/50 hover:shadow-md transition cursor-pointer"
              onClick={() => handleView(tweet.id)}
            >
              <div className="flex items-center gap-2">
                <div className="w-10 h-10 rounded-full bg-gradient-to-r from-blue-500 to-cyan-500 flex items-center justify-center text-white font-bold text-sm shadow-md shadow-blue-500/25 flex-shrink-0">
                  {tweet.username?.charAt(0).toUpperCase() || 'U'}
                </div>
                <div className="flex-1">
                  <span className="font-bold text-gray-800">{tweet.username}</span>
                  <span className="text-gray-400 text-sm ml-2">
                    {new Date(tweet.created_at).toLocaleString('ru')}
                  </span>
                </div>
              </div>
              <p className="mt-3 text-gray-700 whitespace-pre-wrap break-words">
                {tweet.content}
              </p>
              <div className="mt-4 flex items-center gap-6 border-t border-gray-100 pt-3">
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    api.likeTweet(tweet.id).then(() => loadTweets());
                  }}
                  className="flex items-center gap-1 text-gray-400 hover:text-red-500 transition text-sm"
                >
                  <Heart className="w-4 h-4" />
                  <span>{tweet.likes}</span>
                </button>
                <span className="flex items-center gap-1 text-gray-400 text-sm">
                  <Eye className="w-4 h-4" />
                  <span>{tweet.views || 0}</span>
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
