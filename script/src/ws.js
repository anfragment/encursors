export function startWS(url, onCursorEnter, onCursorMove, onCursorLeave) {
  let ws = new WebSocket(url);
  let open = false;

  const onMove = throttle((e) => {
    if (!open) {
      return;
    }
    const x = Math.floor(e.clientX + window.scrollX);
    const y = Math.floor(e.clientY + window.scrollY);
    ws.send(JSON.stringify([x, y]));
  }, 1500);

  document.addEventListener('mousemove', onMove);
  document.addEventListener('mouseenter', onMove)

  ws.onopen = function() {
    open = true;
  };

  ws.onclose = function() {
    open = false;
    setTimeout(() => {
      ws = new WebSocket(url);
    }, 5000);
  }

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

function throttle(cb, delay) {
  let last = 0;
  return function(...args) {
    const now = Date.now();
    if (now - last < delay) {
      return;
    }
    last = now;
    cb(...args);
  }
}