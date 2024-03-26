import React, { useState, useEffect } from 'react';
import { Theme } from "@radix-ui/themes";

export default function OfflineBanner() {
  const [isOnline, setIsOnline] = useState(navigator.onLine);
  const [message, setMessage] = useState(isOnline ? 'Подключение восстановлено' : 'Нет подключения к интернету');
  const [showMessage, setShowMessage] = useState(false);

  useEffect(() => {
    const updateOnlineStatus = () => {
      setIsOnline(navigator.onLine);
      setMessage(navigator.onLine ? 'Подключение восстановлено' : 'Нет подключения к интернету');

      if (!navigator.onLine) {
        setShowMessage(true);
        setTimeout(() => {
          setShowMessage(false);
        }, 5000);
      }
    };

    window.addEventListener('online', updateOnlineStatus);
    window.addEventListener('offline', updateOnlineStatus);

    return () => {
      window.removeEventListener('online', updateOnlineStatus);
      window.removeEventListener('offline', updateOnlineStatus);
    };
  }, []);

  return (
    <div className="App">
      {showMessage && (
        <Theme>
          <div className="internet-status">
            <p>{message}</p>
          </div>
        </Theme>
      )}
    </div>
  );
}

