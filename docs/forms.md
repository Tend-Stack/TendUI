# Composable fields and native multi-step forms

GothUI keeps form ownership in the Go application. Components supply semantic
labels, controls and feedback; your handler owns validation, storage and request
protection. The `/wizard` example exercises a complete native form flow without
client JavaScript, accounts, cookies or a database.

## Field and InputGroup

```templ
templ projectSlug(value, errorMessage string) {
    {{ field := components.FieldProps{
        ID: "project-slug",
        Label: "Project URL slug",
        Description: "The /projects/ prefix is supplied for you.",
        Error: errorMessage,
        Required: true,
    } }}
    @components.Field(field) {
        @components.InputGroup(components.InputGroupProps{Prefix: "/projects/"}) {
            <input type="text" name="slug" value={ value }
                { field.ControlAttributes()... }/>
        }
    }
}
```

Field renders a label, optional help and optional error around its children.
Spread `ControlAttributes()` onto the single raw input, select or textarea. It
supplies the ID, required boolean, aria-invalid on error, and an aria-describedby
list containing the help/error IDs actually rendered by Field. Each call returns
a fresh attribute map. The ID must be unique within the document and contain no
whitespace; an empty ID or label fails rendering. Applications own uniqueness.

Use a raw control inside Field, not the already-labeled Input, Select or Textarea
components. Those existing convenience components remain available unchanged.
Do not duplicate or override the ID/description attributes in the spread. Put
required-field instructions in the label or description as appropriate to the
form. Use ErrorSummary once above the form; individual Field errors are static
text to avoid repeating live announcements.

InputGroup is a layout wrapper around one raw input/select, with escaped text
prefix/suffix decorations. Those decorations are hidden from assistive technology.
Include any meaningful unit, currency or URL prefix in the label or description,
as the example does. Controls keep their native semantics and visible focus ring.
The wrapper owns no values, validation or submission logic.

## Submit actions that can go back

```go
components.ButtonProps{
    Text: "Back", Type: components.ButtonTypeSubmit,
    Name: "action", Value: "back", FormNoValidate: true,
}
```

FormNoValidate emits the native `formnovalidate` attribute only for a submit
button. It does not apply to reset/button types or native overlay openers, which
are always type=button. This lets Back send current edits without native required
fields blocking navigation. It does not disable server validation. Keep the
forward submitter first in DOM order so implicit Enter submission does not
unexpectedly navigate backward.

## Project brief wizard

The example has details, preferences and review steps. Its ordinary form POSTs
carry both visible controls and hidden values from other steps. Back preserves
incomplete edits. Continue validates details and then preferences. Finish checks
all fields again, including values that came from hidden controls. Invalid
content returns HTTP 422 with linked field errors on the relevant step.

The handler bounds request bodies to 8 KB, rejects duplicate draft fields and
invalid UTF-8, validates allowed step/action combinations and enum choices, and
shares the example server's cross-origin protection. Names and summaries use
Unicode character counts in server validation. Slugs use a bounded lowercase
ASCII pattern that is also supplied as a browser hint. Native validation is a
convenience; all final checks are performed in Go.

The result is a validated preview only. No project is created, published, saved
or sent. The URL carries no draft values, and there is no library-owned state.
A new GET starts fresh. Browser history may retain submitted values, and refresh
may ask to resubmit the form; use fictional data in this demonstration. Use the
in-form Back action to preserve the current edits.

This pattern proves field validation, not workflow authorization. A caller with
all valid values can submit directly to the review/finish stage. Production
workflows that require ordered approvals must verify server-owned state and
permissions, use durable storage, and implement idempotency plus an appropriate
POST/redirect/GET flow for mutations. Never trust hidden fields as previously
validated or authoritative. Validate and authorize everything at the mutation
boundary, regardless of the visible step.

## Review checklist

Automated component/HTTP tests exercise label/help/error associations, escaping,
fresh attribute maps, submitter semantics, actual hidden-input round trips,
Back edit preservation, final revalidation, malformed/oversized input, method
handling, cross-origin rejection and independence between fresh requests.

Browser review is still required: keyboard Enter defaults to Continue/Finish;
Back bypasses native constraints; error links focus the relevant field; assistive
technology announces its label and descriptions; long prefixes and review text
fit mobile widths; light/dark/reduced-motion settings remain usable. No browser
verification is claimed inside the restricted environment.
