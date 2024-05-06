import macCursor from './assets/mac.svg';
import windowsCursor from './assets/win.svg';
import tuxCursor from './assets/tux.svg';

(async () => {
  const scriptTag = document.currentScript;
  
  const apiURL = scriptTag.getAttribute('data-api-url');
  if (!apiURL) {
    throw new Error('No API URL provided');
  }
  
  const res = await fetch(apiURL + '/cursors');
  const json = await res.json();
  const cursors = json.cursors || [];
  renderCursors(cursors);
})();

function renderCursors(cursors) {
  for (const cursor of cursors) {
    let el = document.querySelector(`[data-cursor-id="${cursor.id}"]`);
    if (!el) {
      el = document.createElement('img');
      el.setAttribute('data-cursor-id', cursor.id);
      el.src = macCursor;
    }
    
    el.style.position = 'absolute';
    el.style.left = cursor.x + 'px';
    el.style.top = cursor.y + 'px';
  }
}
