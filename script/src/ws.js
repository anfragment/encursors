export function startWS(url, onCursorEnter, onCursorMove, onCursorLeave) {
  const ws = new WebSocket(url);

  const onMove = throttle((e) => {
    const x = e.clientX;
    const y = e.clientY;
    ws.send(JSON.stringify([x, y]));
  }, 3000);

  document.addEventListener('mousemove', onMove);
  document.addEventListener('mouseenter', onMove);

  ws.onmessage = function(event) {
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

  ws.onerror = function(error) {
    console.error('WebSocket error:', error);
  };
}

function throttle(callback, delay) {
  let last = 0;
  return function(...args) {
    const now = new Date().getTime();
    if (now - last < delay) {
      return;
    }
    last = now;
    callback(...args);
  };
}
