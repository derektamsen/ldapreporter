# LDAPReporter

The `ldapreporter` is a tool to query ldap and fetch all groups and their members. The data is output as json.

## Usage

```shell
./ldapreporter -h
Usage of ./ldapreporter:
  -basedn string
        LDAP base DN. Ex: dc=planetexpress,dc=com (default "dc=planetexpress,dc=com")
  -loglevel string
        Logging level (default "WARN")
  -password string
        LDAP Bind Password. Ex: GoodNewsEveryone
  -searchfilter string
        LDAP search filter (default "(&(objectclass=Group))")
  -server string
        LDAP Server. Ex: ldap://localhost:389
  -user string
        LDAP Bind User. Ex: cn=admin,dc=planetexpress,dc=com
  -version
        version of ldapreporter
```

## Export Format

```shell
{"admin_staff":["Hubert J. Farnsworth","Hermes Conrad"],"ship_crew":["Philip J. Fry","Turanga Leela","Bender Bending Rodríguez"]}
```

## Development

Using <https://github.com/rroemhild/docker-test-openldap> as a local ldap test host. This can be run locally with `docker-compose up` and is available at `ldap://localhost:8389` and `ldaps://localhost:8636`. You can use `make dev` to start a development environment that you can run `ldapreporter` against.

You can use `make run` to connect to the local ldap instance. Make run uses the following connection details:

```shell
$ ./ldapreporter \
  -server "ldap://localhost:8389" \
  -user "cn=admin,dc=planetexpress,dc=com" \
  -password "GoodNewsEveryone" \
  -basedn "dc=planetexpress,dc=com" \
  -searchfilter "(&(objectclass=Group))"
```

### Building

Builds are setup in circleci and will run automatically with a PR

- If you want to build locally:

```shell
make
```

### Testing

Testing is configured to run automatically with a circleci PR.

- If you want to run them locally:

```shell
make test
```

### Releasing

Merge the release please pr. Merging the pr will trigger the [release-please action](https://github.com/googleapis/release-please-action) to create a new release.
