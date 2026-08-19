import React, { useCallback, useEffect, useState } from 'react';
import { api } from '../services/api';
import { Tweet } from '../types';
import { Heart, Eye, Send } from 'lucide-react';
import { useRealtime } from '../contexts/RealtimeContext';

export const Feed: React.FC = () => {
  const [tweets, setTweets] = useState<Tweet[]>([]);
  const [loading, setLoading] = useState(true);
  const [content, setContent] = useState('');
  const [notice, setNotice] = useState('');
  const { lastEvent } = useRealtime();

  const loadTweets = useCallback(async () => {
    try {
      const data = await api.getTweets();
      setTweets(data || []);
    } catch (error) {
      console.error('Error loading tweets:', error);
      setTweets([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadTweets();
  }, [loadTweets]);

  useEffect(() => {
    if (lastEvent && ['tweet.created', 'tweet.liked', 'tweet.viewed'].includes(lastEvent.type)) {
      loadTweets();
    }
  }, [lastEvent, loadTweets]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim()) return;
    try {
      await api.createTweet(content);
      setContent('');
      loadTweets();
    } catch (error) {
      setNotice(error instanceof Error ? error.message : 'Не удалось опубликовать твит');
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
    return <div className="flex justify-center p-16 text-[#8d857c]">Загрузка ленты...</div>;
  }

  return (
    <div className="max-w-3xl mx-auto">
      <div className="mb-8 flex items-end justify-between gap-4">
        <div>
          <p className="section-label mb-3">Публичная лента</p>
          <h1 className="text-4xl sm:text-5xl font-black tracking-[-0.05em] text-[#17202a]">Мысли в движении<span className="text-[#ef765e]">.</span></h1>
        </div>
        <span className="hidden sm:block text-right text-xs text-[#8d857c] leading-relaxed">Коротко.<br />По делу.</span>
      </div>
      {/* Создание поста */}
      <div className="surface rounded-[28px] p-5 sm:p-7 mb-7 card-hover">
        <div className="flex items-center gap-3 mb-3">
          <div className="w-11 h-11 rounded-2xl gradient-bg flex items-center justify-center text-white font-bold text-sm shadow-lg shadow-[#ef765e]/20">
            ✍️
          </div>
          <div><p className="section-label">Новый сигнал</p><h2 className="text-lg font-bold text-[#17202a]">Что нового?</h2></div>
        </div>
        <form onSubmit={handleSubmit}>
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="Что у вас нового?"
             className="input-field min-h-[112px] resize-none"
            rows={3}
            maxLength={280}
          />
          <div className="flex justify-between items-center mt-2">
             <span className="text-xs font-medium text-[#9d958b]">{content.length}/280</span>
            <button
              type="submit"
               className="btn-primary flex items-center gap-2"
            >
              <Send className="w-4 h-4" />
              Чирикнуть
            </button>
          </div>
        </form>
        {notice && <p className="mt-3 text-sm text-[#c8554d]">{notice}</p>}
      </div>

      {/* Лента */}
      {tweets.length === 0 ? (
        <div className="surface rounded-[28px] p-12 text-center">
          <div className="text-6xl mb-4">🐦</div>
           <h3 className="text-xl font-bold text-[#17202a]">Пока тихо</h3>
           <p className="text-[#9d958b] mt-2 font-medium">Станьте первым голосом в ленте.</p>
        </div>
      ) : (
        <div className="space-y-4">
          {tweets.map((tweet) => (
            <div
              key={tweet.id}
              className="surface rounded-[26px] p-5 sm:p-6 card-hover cursor-pointer"
              onClick={() => handleView(tweet.id)}
            >
              <div className="flex items-center gap-2">
                 <div className="w-11 h-11 rounded-2xl gradient-bg flex items-center justify-center text-white font-bold text-sm shadow-md shadow-[#ef765e]/20 flex-shrink-0">
                  {tweet.username?.charAt(0).toUpperCase() || 'U'}
                </div>
                <div className="flex-1">
                   <span className="font-bold text-[#17202a]">{tweet.username}</span>
                   <span className="text-[#9d958b] text-xs ml-2">
                    {new Date(tweet.created_at).toLocaleString('ru')}
                  </span>
                </div>
              </div>
               <p className="mt-4 text-[#394550] whitespace-pre-wrap break-words leading-7 text-[15px]">
                {tweet.content}
              </p>
               <div className="mt-5 flex items-center gap-6 border-t border-[#eee8df] pt-4">
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    api.likeTweet(tweet.id).then(() => loadTweets());
                  }}
                   className="flex items-center gap-1.5 text-[#9d958b] hover:text-[#ef765e] transition text-sm"
                >
                  <Heart className="w-4 h-4" />
                  <span>{tweet.likes}</span>
                </button>
                 <span className="flex items-center gap-1.5 text-[#9d958b] text-sm">
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
