import React, { useCallback, useEffect, useState } from 'react';
import { ArrowLeft, MessageCircle, Pencil, Search, Send, Trash2, UserPlus } from 'lucide-react';
import { api } from '../services/api';
import { Conversation, Message, User } from '../types';
import { useAuth } from '../contexts/AuthContext';
import { useRealtime } from '../contexts/RealtimeContext';

export const Messenger: React.FC = () => {
  const { user } = useAuth();
  const { lastEvent } = useRealtime();
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [active, setActive] = useState<Conversation | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [search, setSearch] = useState('');
  const [results, setResults] = useState<User[]>([]);
  const [content, setContent] = useState('');
  const [editing, setEditing] = useState<number | null>(null);
  const [editContent, setEditContent] = useState('');
  const [loading, setLoading] = useState(true);

  const loadConversations = useCallback(async () => {
    try {
      const data = await api.getConversations();
      setConversations(data);
      if (active) {
        const updated = data.find((conversation) => conversation.id === active.id);
        if (updated) setActive(updated);
      }
    } catch (error) {
      console.error('Error loading conversations:', error);
    } finally {
      setLoading(false);
    }
  }, [active]);

  const loadMessages = useCallback(async () => {
    if (!active) return;
    try {
      setMessages(await api.getMessages(active.id));
    } catch (error) {
      console.error('Error loading messages:', error);
    }
  }, [active]);

  useEffect(() => {
    loadConversations();
  }, [loadConversations]);

  useEffect(() => {
    if (lastEvent?.type.startsWith('message.')) {
      loadConversations();
      loadMessages();
    }
  }, [lastEvent, loadConversations, loadMessages]);

  useEffect(() => {
    const timeout = window.setTimeout(async () => {
      if (search.trim().length < 2) {
        setResults([]);
        return;
      }
      try {
        setResults(await api.searchUsers(search.trim()));
      } catch (error) {
        console.error('Error searching users:', error);
      }
    }, 250);
    return () => window.clearTimeout(timeout);
  }, [search]);

  const openConversation = async (conversation: Conversation) => {
    setActive(conversation);
    setSearch('');
    setResults([]);
    setMessages(await api.getMessages(conversation.id));
  };

  const startConversation = async (otherUser: User) => {
    const conversation = await api.createConversation(otherUser.id);
    setConversations((current) => [conversation, ...current.filter((item) => item.id !== conversation.id)]);
    await openConversation(conversation);
  };

  const sendMessage = async (event: React.FormEvent) => {
    event.preventDefault();
    if (!active || !content.trim()) return;
    const message = await api.sendMessage(active.id, content.trim());
    setMessages((current) => [...current, message]);
    setContent('');
    loadConversations();
  };

  const saveEdit = async (message: Message) => {
    if (!editContent.trim()) return;
    const updated = await api.updateMessage(message.id, editContent.trim());
    setMessages((current) => current.map((item) => item.id === updated.id ? updated : item));
    setEditing(null);
  };

  const removeMessage = async (message: Message) => {
    if (!window.confirm('Удалить это сообщение?')) return;
    const deleted = await api.deleteMessage(message.id);
    setMessages((current) => current.map((item) => item.id === deleted.id ? deleted : item));
  };

  return (
    <div className="max-w-5xl mx-auto">
      <div className="mb-7">
        <p className="section-label mb-3">Личные сообщения</p>
        <h1 className="text-4xl sm:text-5xl font-black tracking-[-0.05em] text-[#17202a]">Свои люди<span className="text-[#ef765e]">.</span></h1>
      </div>

      <div className="surface rounded-[30px] overflow-hidden min-h-[640px] grid md:grid-cols-[300px_1fr]">
        <aside className={`${active ? 'hidden md:flex' : 'flex'} flex-col border-r border-[#e7dfd4] bg-[#fbfaf7]/80`}>
          <div className="p-5 border-b border-[#e7dfd4]">
            <div className="relative">
              <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-[#9d958b]" />
              <input value={search} onChange={(event) => setSearch(event.target.value)} className="input-field pl-10 py-2.5 text-sm" placeholder="Найти человека" />
            </div>
            {results.length > 0 && (
              <div className="mt-2 space-y-1 rounded-2xl bg-white p-2 shadow-xl border border-[#e7dfd4]">
                {results.map((result) => <button key={result.id} onClick={() => startConversation(result)} className="w-full flex items-center gap-3 p-2 rounded-xl hover:bg-[#f4f1eb] text-left">
                  <Avatar user={result} /><span className="min-w-0"><strong className="block text-sm text-[#17202a] truncate">{result.name}</strong><small className="text-[#9d958b]">@{result.username}</small></span><UserPlus className="ml-auto h-4 w-4 text-[#ef765e]" />
                </button>)}
              </div>
            )}
          </div>
          <div className="px-5 pt-5 pb-2 flex items-center justify-between"><span className="section-label">Диалоги</span><MessageCircle className="w-4 h-4 text-[#ef765e]" /></div>
          <div className="overflow-y-auto px-3 pb-4">
            {loading ? <p className="p-4 text-sm text-[#9d958b]">Загрузка...</p> : conversations.length === 0 ? <div className="p-4 text-sm text-[#9d958b]">Найдите пользователя, чтобы начать разговор.</div> : conversations.map((conversation) => <button key={conversation.id} onClick={() => openConversation(conversation)} className={`w-full flex gap-3 p-3 rounded-2xl text-left transition ${active?.id === conversation.id ? 'bg-[#17202a] text-white' : 'hover:bg-[#f4f1eb]'}`}>
              <Avatar user={conversation.other_user} /><span className="min-w-0 flex-1"><strong className={`block truncate text-sm ${active?.id === conversation.id ? 'text-white' : 'text-[#17202a]'}`}>{conversation.other_user.name}</strong><small className={`block truncate mt-1 ${active?.id === conversation.id ? 'text-white/50' : 'text-[#9d958b]'}`}>{conversation.last_message?.deleted ? 'Сообщение удалено' : conversation.last_message?.content || 'Новый диалог'}</small></span>
            </button>)}
          </div>
        </aside>

        <section className={`${active ? 'flex' : 'hidden md:flex'} flex-col min-w-0 bg-white/45`}>
          {!active ? <div className="m-auto text-center p-8"><div className="mx-auto grid place-items-center h-16 w-16 rounded-3xl bg-[#f4f1eb] text-[#ef765e] mb-5"><MessageCircle /></div><h2 className="text-xl font-black text-[#17202a]">Выберите диалог</h2><p className="text-[#9d958b] mt-2 text-sm">Или найдите нового собеседника слева.</p></div> : <>
            <header className="flex items-center gap-3 p-5 border-b border-[#e7dfd4]"><button onClick={() => setActive(null)} className="md:hidden p-2 rounded-xl hover:bg-[#f4f1eb]"><ArrowLeft className="w-5 h-5" /></button><Avatar user={active.other_user} /><div><h2 className="font-black text-[#17202a]">{active.other_user.name}</h2><p className="text-xs text-[#9d958b]">@{active.other_user.username}</p></div></header>
            <div className="flex-1 overflow-y-auto p-5 space-y-3">
              {messages.length === 0 && <p className="text-center text-sm text-[#9d958b] py-12">Начните разговор первым.</p>}
              {messages.map((message) => { const mine = message.sender_id === user?.id; return <div key={message.id} className={`group flex ${mine ? 'justify-end' : 'justify-start'}`}>
                <div className={`max-w-[78%] sm:max-w-[65%] ${mine ? 'items-end' : 'items-start'} flex flex-col`}>
                  <div className={`px-4 py-3 rounded-2xl text-sm leading-6 ${mine ? 'bg-[#17202a] text-white rounded-br-md' : 'bg-[#f4f1eb] text-[#394550] rounded-bl-md'}`}>{message.deleted ? <i className="text-white/45">Сообщение удалено</i> : editing === message.id ? <div className="space-y-2"><textarea value={editContent} onChange={(event) => setEditContent(event.target.value)} className="w-full min-w-[180px] rounded-xl p-2 text-[#17202a]" /><button onClick={() => saveEdit(message)} className="text-xs text-[#ef765e] font-bold">Сохранить</button></div> : message.content}</div>
                  <div className="flex items-center gap-2 mt-1 px-1"><span className="text-[10px] text-[#9d958b]">{new Date(message.created_at).toLocaleTimeString('ru', { hour: '2-digit', minute: '2-digit' })}</span>{mine && !message.deleted && editing !== message.id && <><button onClick={() => { setEditing(message.id); setEditContent(message.content); }} className="opacity-0 group-hover:opacity-100 text-[#9d958b] hover:text-[#ef765e]"><Pencil className="w-3 h-3" /></button><button onClick={() => removeMessage(message)} className="opacity-0 group-hover:opacity-100 text-[#9d958b] hover:text-[#ef765e]"><Trash2 className="w-3 h-3" /></button></>}</div>
                </div>
              </div>; })}
            </div>
            <form onSubmit={sendMessage} className="p-4 border-t border-[#e7dfd4] flex gap-2"><input value={content} onChange={(event) => setContent(event.target.value)} className="input-field py-3" placeholder="Напишите сообщение..." maxLength={2000} /><button className="btn-primary px-4 shrink-0" aria-label="Отправить"><Send className="w-4 h-4" /></button></form>
          </>}
        </section>
      </div>
    </div>
  );
};

const Avatar: React.FC<{ user: User }> = ({ user }) => <div className="shrink-0 h-10 w-10 rounded-2xl gradient-bg grid place-items-center text-white font-black text-sm">{user.name?.charAt(0).toUpperCase() || user.username?.charAt(0).toUpperCase() || 'U'}</div>;
