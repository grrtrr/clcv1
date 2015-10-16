## CenturyLink Cloud v1 API

This is a Go implementation of the [CLC v1 API](https://www.ctl.io/api-docs/v1).

## Getting started

Try some of the examples in the `examples/` folder, which serve to illustrate individual API calls.

Most have help screens (`-h`). The library supports _debug output_ via `-d`.

_Credentials_ can be passed in one of two forms:

1. Via _commandline flags_:
  + `-k <your-API-key>`,
  + `-p <your-API-pass>`.
2. Using _environment variables_:
  + `CLC_V1_API_KEY=<your-API-key>`,
  + `CLC_V1_API_PASS=<your-API-pass>`.

## Todo

Not every single detail of the API is implemented.
For instance, the billing API calls have been left out.
