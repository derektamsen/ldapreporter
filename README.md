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

### Building

```shell
make
```
