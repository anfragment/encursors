export function startWS(url, onCursorEnter, onCursorMove, onCursorLeave) {
  const ws = new WebSocket(url, "json");

  const onMove = (e) => {
    const x = e.clientX;
    const y = e.clientY;
    ws.send(JSON.stringify([x, y]));
  }

  document.addEventListener('mousemove', onMove);
  document.addEventListener('mouseenter', onMove);

  ws.onopen = function() {
    console.log('Connected to the server');
  };

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