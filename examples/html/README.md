# Plain HTML quickstart

The CSS works with static HTML and any server or template engine. No Go, templ,
Node.js, npm or JavaScript runtime is required.

From a checkout of `v2026.2.0-preview.1`, copy these two files into a new folder:

```sh
mkdir my-tendui-page
cp examples/html/index.html static/tendui.css my-tendui-page/
```

Open `my-tendui-page/index.html` in your browser. No server or build step is needed.
Try the text field and disclosure. Change the root `data-theme="light"` to
`data-theme="dark"`; omit the attribute to follow the system preference.

Deploy those two files together using any static web server. Keep the relative
stylesheet link, or change it to your application's asset path. Forms, storage,
authentication and any optional HTMX runtime belong to your application.

The stylesheet is from the same immutable preview tag as the HTML. MIT license:
retain the copyright and permission notice when redistributing the library.
