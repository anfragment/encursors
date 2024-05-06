import macCursor from './assets/mac.svg';
import windowsCursor from './assets/win.svg';
import tuxCursor from './assets/tux.svg';
import { startWS } from './ws';

(async () => {
  const os = getOS();
  if (os === null) {
    console.log('Unsupported OS')
    return;
  }

  const scriptTag = document.currentScript;

  const apiURL = scriptTag.getAttribute('data-api-url');
  if (!apiURL) {
    throw new Error('No API URL provided');
  }

  const path = window.location.pathname;

  const response = await fetch(`https://${apiURL}/cursors?path=${path}`);
  const data = await response.json();
  const cursors = {};
  for (const cursor of (data.cursors || [])) {
    cursors[cursor.id] = cursor;
  }
  renderCursors(Object.values(cursors));
  startWS(
    `ws://${apiURL}/subscribe?path=${path}&country=KZ&os=mac`, 
    (cursor) => {
      cursors[cursor.id] = cursor;
      renderCursors(Object.values(cursors));
    },
    (cursor) => {
      cursors[cursor.id] = cursor;
      renderCursors(Object.values(cursors));
    },
    (cursor) => {
      delete cursors[cursor.id];
      renderCursors(Object.values(cursors));
    },
  );
})();

function renderCursors(cursors) {
  for (const cursor of cursors) {
    let el = document.querySelector(`[data-cursor-id="${cursor.id}"]`);
    if (!el) {
      el = document.createElement('img');
      el.setAttribute('data-cursor-id', cursor.id);
      el.style.position = 'absolute';
      el.style.left = `${cursor.posX}px`;
      el.style.top = `${cursor.posY}px`;
      el.style.width = '20px';
      el.style.height = '20px';
      el.src = tuxCursor;
      el.style.zIndex = '9999';
      document.body.appendChild(el);
    }

    el.style.left = `${cursor.posX}px`;
    el.style.top = `${cursor.posY}px`;
  }
}

function getOS() {
  const platform = window.navigator?.userAgentData?.platform || window.navigator.platform,
      macosPlatforms = ['macOS', 'Macintosh', 'MacIntel', 'MacPPC', 'Mac68K'],
      windowsPlatforms = ['Win32', 'Win64', 'Windows', 'WinCE'];

  if (macosPlatforms.indexOf(platform) !== -1) {
    return 'Mac';
  } else if (windowsPlatforms.indexOf(platform) !== -1) {
    return 'Windows';
  } else if (/Linux/.test(platform)) {
    return 'Linux';
  }
  return null;
}
