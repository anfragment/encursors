<p align="center">
  <picture>
    <img src="https://github.com/anfragment/encursors/blob/master/logo.png?raw=true" alt="Encursors logo" width="150" />
  </picture>
</p>

<h2 align="center">
  Encursors
</h2>

Ever feel like a lone wanderer on the digital plains of the web? Do your static websites seem just a tad too... *static*? Fear not, for Encursors is here to help! This nifty little tool transforms your website into a bustling village square. With a simple script, Encursors displays each visitor's cursor movements in real time, letting everyone see where everyone else is looking. It's like a party on your page, and everyone's invited!

Note: Encursors only works for desktop users, as it relies on the position of the mouse cursor.

## Demo
You can see Encursors in action on our [demo page](https://anfragment.github.io/encursors/). Open the page in multiple tabs or devices to see the cursors move in real time!

## Installation
To install Encursors, simply add the following script tag to your website's HTML:
```html
<script data-api-url="prod-encursors-ypdi.encr.app" src="https://cdn.jsdelivr.net/gh/anfragment/encursors@latest/script/dist/cursors.min.js"></script>
```

## Configuration options
You can configure Encursors by passing options to the script tag as data attributes. Here are the available options:
- `data-api-url`: The URL of the API. Change this if you're running your own instance.
- `data-z-index`: The z-index of the cursor elements.
