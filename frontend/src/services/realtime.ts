export interface RealtimeEvent {
  type: string;
  data: Record<string, unknown>;
}

export const connectToRealtime = (onEvent: (event: RealtimeEvent) => void): (() => void) => {
  const configuredUrl = process.env.REACT_APP_REALTIME_URL;
  const baseUrl = configuredUrl || (
    typeof window !== 'undefined' ? window.location.origin : ''
  );
  if (!baseUrl || typeof EventSource === 'undefined') {
    return () => undefined;
  }

  const source = new EventSource(`${baseUrl.replace(/\/$/, '')}/events`);
  source.onmessage = (message) => {
    try {
      const parsed = JSON.parse(message.data) as Partial<RealtimeEvent>;
      if (parsed.type) {
        onEvent({
          type: parsed.type,
          data: parsed.data && typeof parsed.data === 'object' ? parsed.data : {},
        });
      }
    } catch {
      // Ignore malformed events and let EventSource continue reconnecting.
    }
  };
  source.onerror = () => {
    // EventSource automatically retries while the realtime service is unavailable.
  };

  return () => source.close();
};
