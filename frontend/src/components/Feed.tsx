import React, { useEffect, useState } from 'react';
import { api } from '../services/api';
import { Tweet } from '../types';

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
    <div className="max-w-2xl mx-auto p-4">
      <div className="bg-white rounded-xl shadow p-4 mb-4">
        <h2 className="text-xl font-bold mb-2">🐦 Что нового?</h2>
        <form onSubmit={handleSubmit}>
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="Что у вас нового?"
            className="w-full p-3 border rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
            rows={3}
            maxLength={280}
          />
          <div className="flex justify-between items-center mt-2">
            <span className="text-sm text-gray-400">{content.length}/280</span>
            <button
              type="submit"
              className="bg-blue-500 text-white px-6 py-2 rounded-xl font-semibold hover:bg-blue-600 transition"
            >
              Чирикнуть 🐦
            </button>
          </div>
        </form>
      </div>

      {tweets.length === 0 ? (
        <div className="text-center text-gray-400 py-12">
          <p className="text-4xl mb-2">🐦</p>
          <p>Нет твитов</p>
          <p className="text-sm">Будьте первым!</p>
        </div>
      ) : (
        tweets.map((tweet) => (
          <div 
            key={tweet.id} 
            className="bg-white rounded-xl shadow p-4 mb-3 hover:shadow-md transition cursor-pointer"
            onClick={() => handleView(tweet.id)}
          >
            <div className="flex items-center gap-2">
              <span className="font-bold">{tweet.username}</span>
              <span className="text-gray-400 text-sm">
                {new Date(tweet.created_at).toLocaleString('ru')}
              </span>
            </div>
            <p className="mt-2">{tweet.content}</p>
            <div className="mt-3 flex items-center gap-4 text-sm text-gray-400">
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  api.likeTweet(tweet.id).then(() => loadTweets());
                }}
                className="hover:text-red-500 transition"
              >
                ❤️ {tweet.likes}
              </button>
              <span>👁️ {tweet.views || 0}</span>
            </div>
          </div>
        ))
      )}
    </div>
  );
};
