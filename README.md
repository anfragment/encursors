<p align="center">
  <picture>
    <img src="https://github.com/anfragment/encursors/blob/master/logo.png?raw=true" alt="Encursors logo" width="150" />
  </picture>
</p>

<h2 align="center">
  Encursors
</h2>

Ever feel like a lone wanderer on the digital plains of the web? Do your static websites seem just a tad too... *static*? Fear not, for Encursors is here to help! This nifty little tool transforms your website into a bustling village square. With a simple script, Encursors displays each visitor's cursor movements in real time, letting everyone see where everyone else is looking. It's like a party on your page, and everyone's invited!

The backend is built and hosted with [Encore](https://encore.dev), a development platform for building event-driven and distributed systems. You can run your own instance by cloning the repostory.

> [!NOTE]
> Encursors does not display cursors or track users on mobile devices.

## Demo
You can see Encursors in action on our [demo page](https://anfragment.github.io/encursors/). Open the page in multiple tabs or devices to see the cursors move in real time!

## Features
- Displays the flag of the country the visitor is from alongside their cursor.
- Custom cursors based on the visitor's operating system.
- Respects the [prefers-reduced-motion](https://developer.mozilla.org/en-US/docs/Web/CSS/@media/prefers-reduced-motion) setting by not displaying the cursors when it is enabled.
- No cookies or tracking of any kind. Cursor data is permanently deleted after a visitor leaves the page.
- Fully open source and free to use.

## Installation
To install Encursors, simply add the following script tag to your website's HTML:
```html
<script src="https://cdn.jsdelivr.net/gh/anfragment/encursors@release/script/dist/cursors.min.js"></script>
```

## Configuration options
You can configure Encursors by passing options to the script tag as data attributes. Here are the available options:
- `data-api-url`: The base URL of the API. Set if you're running your own instance. Should not include the protocol or the trailing slash.
- `data-z-index`: The z-index of the cursor elements. Optional.
