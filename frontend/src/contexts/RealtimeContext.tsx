import React, { createContext, useContext, useEffect, useState } from 'react';
import { useAuth } from './AuthContext';
import { connectToRealtime, RealtimeEvent } from '../services/realtime';

interface Notification {
  id: number;
  text: string;
}

interface RealtimeContextValue {
  lastEvent: RealtimeEvent | null;
  notifications: Notification[];
  dismissNotification: (id: number) => void;
}

const RealtimeContext = createContext<RealtimeContextValue | undefined>(undefined);

export const RealtimeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user } = useAuth();
  const [lastEvent, setLastEvent] = useState<RealtimeEvent | null>(null);
  const [notifications, setNotifications] = useState<Notification[]>([]);

  useEffect(() => {
    if (!user) {
      setLastEvent(null);
      setNotifications([]);
      return undefined;
    }

    return connectToRealtime((event) => {
      setLastEvent(event);
      const targetUserID = Number(event.data.target_user_id);
      if (targetUserID !== user.id) return;

      let text = '';
      if (event.type === 'tweet.liked') {
        text = `${String(event.data.username || 'Кто-то')} отметил ваш твит`;
      } else if (event.type === 'user.followed') {
        text = `${String(event.data.username || 'Кто-то')} подписался на вас`;
      } else if (event.type === 'message.created') {
        text = 'Новое сообщение';
      }
      if (text) {
        setNotifications((current) => [...current, { id: Date.now(), text }].slice(-3));
      }
    });
  }, [user]);

  const dismissNotification = (id: number) => {
    setNotifications((current) => current.filter((item) => item.id !== id));
  };

  return (
    <RealtimeContext.Provider value={{ lastEvent, notifications, dismissNotification }}>
      {children}
    </RealtimeContext.Provider>
  );
};

export const useRealtime = () => {
  const context = useContext(RealtimeContext);
  if (!context) throw new Error('useRealtime must be used within RealtimeProvider');
  return context;
};
