import macCursor from './assets/mac.svg';
import windowsCursor from './assets/win.svg';
import tuxCursor from './assets/tux.svg';
import { startWS } from './ws';
import ct from 'countries-and-timezones';

(async () => {
  const prefersReducedMotion = window.matchMedia(`(prefers-reduced-motion: reduce)`) === true || window.matchMedia(`(prefers-reduced-motion: reduce)`).matches === true;
  if (prefersReducedMotion) {
    console.log('Reduced motion is enabled');
    return;
  }

  const os = getOS();
  if (!os) {
    console.log('Unsupported OS');
    return;
  }
  const country = ct.getCountryForTimezone(Intl.DateTimeFormat().resolvedOptions().timeZone);
  if (!country) {
    console.log('Could not determine country');
    return;
  }
  const countryCode = country.id;

  const scriptTag = document.currentScript;

  const apiURL = scriptTag.getAttribute('data-api-url');
  if (!apiURL) {
    throw new Error('No API URL provided');
  }
  const zIndex = scriptTag.getAttribute('data-z-index');

  const url = window.location.href.split('?')[0];

  const response = await fetch(`https://${apiURL}/cursors?url=${url}`);
  const data = await response.json();
  const cursors = {};
  for (const cursor of (data.cursors || [])) {
    cursors[cursor.id] = cursor;
  }
  renderCursors(cursors);

  startWS(
    `wss://${apiURL}/subscribe?url=${url}&country=${countryCode}&os=${os}`,
    (cursor) => {
      cursors[cursor.id] = cursor;
      renderCursors(cursors);
    },
    (cursor) => {
      cursors[cursor.id] = cursor;
      renderCursors(cursors);
    },
    (cursor) => {
      delete cursors[cursor.id];
      deleteCursor(cursors);
    },
  );

  function deleteCursor(cursor) {
    const el = document.querySelector(`[data-cursor-id="${cursor.id}"]`);
    if (el) {
      el.remove();
    }
  }

  function renderCursors(cursors) {
    for (const cursor of Object.values(cursors)) {
      let el = document.querySelector(`[data-cursor-id="${cursor.id}"]`);
      if (!el) {
        el = document.createElement('div');
        el.setAttribute('data-cursor-id', cursor.id);
        el.style.position = 'absolute';

        if (zIndex) {
          el.style.zIndex = zIndex;
        }
        el.style.transition = 'left 0.2s, top 0.2s';
        el.style.display = 'flex';
        el.style.alignItems = 'center';
        el.style.justifyContent = 'center';
        el.style.pointerEvents = 'none';

        const img = document.createElement('img');
        switch (cursor.os) {
        case 0:
          img.src = macCursor;
          break;
        case 1:
          img.src = windowsCursor;
          break;
        case 2:
          img.src = tuxCursor;
          break;
        default:
          img.src = macCursor;
          break;
        }
        img.style.width = '20px';
        img.style.height = '20px';
        el.appendChild(img);
  
        const countryFlag = document.createElement('div');
        countryFlag.textContent = getFlagEmoji(cursor.country);
        countryFlag.style.position = 'relative';
        countryFlag.style.left = '-2px';
        el.appendChild(countryFlag);

        const elWidth = el.clientWidth;
        const elHeight = el.clientHeight;
        el.style.left = `${boundByPageWidth(cursor.posX, elWidth)}px`;
        el.style.top = `${boundByPageHeight(cursor.posY, elHeight)}px`;

        document.body.appendChild(el);
      }

      const elWidth = el.clientWidth;
      const elHeight = el.clientHeight;
      el.style.left = `${boundByPageWidth(cursor.posX, elWidth)}px`;
      el.style.top = `${boundByPageHeight(cursor.posY, elHeight)}px`;

      if (cursor.posX === 0 && cursor.posY === 0) {
        el.style.display = 'none';
      } else {
        el.style.display = 'flex';
      }
    }
  }
})();

function boundByPageHeight(y, elHeight) {
  return Math.min(Math.max(0, y), window.innerHeight - 1.5*elHeight - 1);
}

function boundByPageWidth(x, elWidth) {
  return Math.min(Math.max(0, x), window.innerWidth - 1.5*elWidth - 1);
}

function getFlagEmoji(countryCode) {
  const codePoints = countryCode.toUpperCase().split('').map(char =>  0x1F1E6 + char.charCodeAt(0) - 'A'.charCodeAt(0));
  return String.fromCodePoint(...codePoints);
}

function getOS() {
  const platform = window.navigator?.userAgentData?.platform || window.navigator.platform,
      macosPlatforms = ['macOS', 'Macintosh', 'MacIntel', 'MacPPC', 'Mac68K'],
      windowsPlatforms = ['Win32', 'Win64', 'Windows', 'WinCE'];

  if (macosPlatforms.indexOf(platform) !== -1) {
    return 'mac';
  } else if (windowsPlatforms.indexOf(platform) !== -1) {
    return 'win';
  } else if (/Linux/.test(platform)) {
    return 'linux';
  }
  return null;
}
