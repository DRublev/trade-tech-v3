import React, { useState, useEffect } from 'react';
import s from './styles.css';


export default function OfflineBanner() {
  const [isOnline, setIsOnline] = useState(navigator.onLine);
  const [message, setMessage] = useState(isOnline ? 'Подключение восстановлено' : 'Нет подключения к интернету');
  const [showMessage, setShowMessage] = useState(!isOnline);
  const [timeoutId, setTimeoutId] = useState(null);

  const handleOnline = () => {
    setIsOnline(true);
    setMessage('Подключение восстановлено');
    setShowMessage(true);
    const id = setTimeout(() => {
      setShowMessage(false);
    }, 2000);
    setTimeoutId(id);
  };

  const handleOffline = () => {
    setIsOnline(false);
    setMessage('Нет подключения к интернету');
    setShowMessage(true);
  };

  useEffect(() => {
    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    if (isOnline) {
      setShowMessage(false);
    }

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
      if (timeoutId) {
        clearTimeout(timeoutId);
      }
    };
  }, []);

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
