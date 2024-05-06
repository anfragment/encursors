(function () {
  'use strict';

  var img = "data:image/svg+xml,%3csvg height='32' viewBox='0 0 32 32' width='32' xmlns='http://www.w3.org/2000/svg'%3e%3cg fill='none' fill-rule='evenodd' transform='translate(10 7)'%3e%3cpath d='m6.148 18.473 1.863-1.003 1.615-.839-2.568-4.816h4.332l-11.379-11.408v16.015l3.316-3.221z' fill='white'/%3e%3cpath d='m6.431 17 1.765-.941-2.775-5.202h3.604l-8.025-8.043v11.188l2.53-2.442z' fill='black'/%3e%3c/g%3e%3c/svg%3e";

  const scriptTag = document.currentScript;

  const apiURL = scriptTag.getAttribute('data-api-url');
  if (!apiURL) {
    throw new Error('No API URL provided');
  }

  fetch(apiURL + '/cursors')
    .then(response => response.json())
    .then(renderCursors)
    .catch(err => console.error('Error fetching cursors:', err));

  function renderCursors(cursors) {
    for (const cursor of cursors) {
      let el = document.querySelector(`[data-cursor-id="${cursor.id}"]`);
      if (!el) {
        el = document.createElement('img');
        el.setAttribute('data-cursor-id', cursor.id);
        el.src = img;
      }
      
      el.style.position = 'absolute';
      el.style.left = cursor.x + 'px';
      el.style.top = cursor.y + 'px';
    }
  }

})();
