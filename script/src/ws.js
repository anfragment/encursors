import { EVENT_THROTTLE_TIMEOUT_MS } from './config';

export function startWS(url, onCursorEnter, onCursorMove, onCursorLeave) {
  let ws = new WebSocket(url);
  let open = false;

  let clientX = 0,
    clientY = 0;
  const onMove = throttle((e) => {
    if (typeof e.clientX === 'number') {
      clientX = e.clientX;
      clientY = e.clientY;
    }
    if (!open) {
      return;
    }
    const x = Math.floor(clientX + window.scrollX);
    const y = Math.floor(clientY + window.scrollY);
    ws.send(JSON.stringify([x, y]));
  }, EVENT_THROTTLE_TIMEOUT_MS);
  document.addEventListener('mousemove', onMove);
  document.addEventListener('mouseenter', onMove);
  document.addEventListener('scroll', onMove);

  ws.onopen = function () {
    open = true;
  };

  ws.onclose = function () {
    open = false;
    setTimeout(() => {
      ws = new WebSocket(url);
    }, 5000);
  };

  ws.onmessage = function (event) {
    const data = JSON.parse(event.data);
    switch (data.type) {
      case 'enter':
        onCursorEnter(data.payload);
        break;
      case 'move':
        onCursorMove(data.payload);
        break;
      case 'leave':
        onCursorLeave(data.payload);
        break;
    }
  };

  ws.onerror = function (error) {
    console.error('WebSocket error:', error);
  };
}

function throttle(cb, delay) {
  let timeoutId,
    lastArgs,
    last = 0;
  return function (...args) {
    lastArgs = args;
    if (timeoutId) {
      return;
    }

    const now = Date.now();
    if (now - last > delay) {
      last = now;
      cb(...lastArgs);
      return;
    }

    timeoutId = setTimeout(
      () => {
        last = Date.now();
        cb(...lastArgs);
        timeoutId = null;
      },
      delay - (now - last)
    );
  };
}
