# LDAPReporter

The `ldapreporter` is a tool to query ldap and fetch all groups and their
members. The data is exported as json.

## Export Format

```shell
{
  group1: [
    user1,
    user2,
  ],
  group2: [
    user1,
    user3,
  ]
}
```

## Development

Using <https://github.com/rroemhild/docker-test-openldap> as a local ldap test
host. This can be run locally with `docker-compose up` and is available at
`ldap://localhost:8389` and `ldaps://localhost:8636`. You can use `make dev`
to start a development environment that you can run `ldapreporter` against.

### Building

```shell
make
```

### Testing

```shell
make test
```
