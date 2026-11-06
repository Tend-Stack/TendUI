// Package components provides typed templ components for GothUI's embedded CSS.
//
// Layout supplies the document shell and stylesheet link. Children and explicit
// templ.Component slots provide composition; ordinary text is escaped by templ.
// Applications own routes, authorization, validation, persistence and unique IDs.
// Do not pass untrusted content through templ.Raw to bypass escaping.
//
// Components use native HTML behavior. Optional HX fields emit attributes but do
// not include an HTMX runtime. Motion is opt-in and respects reduced-motion CSS.
// Field associates a label, help and error with a raw control through its
// ControlAttributes method; avoid nesting already-labeled convenience controls.
// Native navigation components render links, not client-side tab controllers.
package components
