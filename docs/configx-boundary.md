# configx Boundary

postgresx does not depend on a config-loading L2 module.

Allowed boundary:

- callers load environment, files, and secret stores outside postgresx;
- callers construct `postgresx.Config` explicitly;
- postgresx validates and redacts the supplied config;
- examples may demonstrate env loading as caller-side code, not core package
  behavior.

Disallowed boundary:

- core package reading env or secret files;
- core package depending on external config loaders;
- core package owning service-specific schema, repository, or migration config.
