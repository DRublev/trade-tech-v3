import React, { useState, useEffect } from 'react';
import s from './styles.css';

export default function OfflineBanner() {
  const [isOnline, setIsOnline] = useState(navigator.onLine);

  const connectionRestoredMessage = 'Подключение восстановлено';
  const noConnectionMessage = 'Нет подключения к интернету';
  const message = isOnline ? connectionRestoredMessage : noConnectionMessage;

  const [showMessage, setShowMessage] = useState(!isOnline);
  const [timeoutId, setTimeoutId] = useState(null);

  const handleOnline = () => {
    setIsOnline(true);
    setShowMessage(true);

    // Если уже есть активный таймер, очищаем его
    if (timeoutId) {
      clearTimeout(timeoutId);
    }

    // Запускаем новый таймер
    const id = setTimeout(() => {
      setShowMessage(false);
    }, 2000);
    setTimeoutId(id);
  };

  const handleOffline = () => {
    setIsOnline(false);
    setShowMessage(true);
  };

  useEffect(() => {
    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    if (isOnline) {
      // Если онлайн при загрузке компонента, устанавливаем таймер
      const id = setTimeout(() => {
        setShowMessage(false);
      }, 2000);
      setTimeoutId(id);
    }

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
      if (timeoutId) {
        clearTimeout(timeoutId);
      }
    };
  }, [isOnline]); // Добавлен isOnline как зависимость

  return (
    <>
      {showMessage && (
        <div className={s.internetStatus}>
          <p className={isOnline ? s.online : s.offline}>{message}</p>
        </div>
      )}
    </>
  );
}
