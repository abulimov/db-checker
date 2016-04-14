[![Build Status](https://travis-ci.org/abulimov/db-checker.svg?branch=master)](https://travis-ci.org/abulimov/db-checker)

# db-checker

Utility to run queries on PostgreSQL/MySQL database, and alert on some assertions against
query result.

For example, one can check if query returns data, or query does not return data.

## Installation

Grab latest release on [releases page](https://github.com/abulimov/db-checker/releases),
or build from source.

### Building from source

You need working Go compiler.
Tested against Go 1.5+

On Linux/OSX:

```
# set GOPATH to some valid path
export GOPATH=~/go && mkdir -p ~/go
go get github.com/abulimov/db-checker
```

Now you have *db-checker* binary.

## Checks

Checks are specified as [YAML](http://www.yaml.org) files,
with only 3 mandatory fields:

* query: any SQL query you can imagine
* description: human-readable description of performed check
* assert: type of check assertion, *present*, *absent*, *true* or *false*

### Check example

Check if we have any locks in our database.
We set the assertion type as *absent* because any found lock will result in
non-zero exit status.

```yaml
query: |
    SELECT
      COALESCE(blockingl.relation::regclass::text,blockingl.locktype) as locked_item,
      (now() - blockeda.query_start)::time AS waiting_duration,
      blockeda.pid AS blocked_pid,
      blockeda.query as blocked_query, blockedl.mode as blocked_mode,
      blockinga.pid AS blocking_pid, blockinga.query as blocking_query,
      blockingl.mode as blocking_mode
    FROM pg_catalog.pg_locks blockedl
    JOIN pg_stat_activity blockeda ON blockedl.pid = blockeda.pid
    JOIN pg_catalog.pg_locks blockingl ON(
      ( (blockingl.transactionid=blockedl.transactionid) OR
      (blockingl.relation=blockedl.relation AND blockingl.locktype=blockedl.locktype)
      ) AND blockedl.pid != blockingl.pid)
    JOIN pg_stat_activity blockinga ON blockingl.pid = blockinga.pid
      AND blockinga.datid = blockeda.datid
    WHERE NOT blockedl.granted;
description: Locks in database
assert: absent
```

More examples in **examples** directory.

## Usage

This utility is a Nagios-compatible plugin.

You must at least specify credentials to access the database and a directory
to get checks from.

```console
nagios@example.com:~$ ./db-checker --dbname stupid --dbuser=checker --dbhost=localhost --dbpassword=SomePassword --checks /opt/checks/stupid
WARNING:
* Stupid check
No results found
 | problems=1;0;0;0;0

nagios@example.com:~$ ./db-checker --dbname movies --dbuser=checker --dbhost=localhost --dbpassword=SomePassword --checks /opt/checks/movies --critical
CRITICAL:
* Found movies with zero duration
N. ¦ column1 ¦ orig_title                ¦ rus_title
1. ¦ 1346    ¦ Midnight Express          ¦ Полуночный экспресс
2. ¦ 2165    ¦ In the Loop               ¦ В петле
3. ¦ 2254    ¦ Sex & Drugs & Rock & Roll ¦ Секс, наркотики и рок-н-ролл
4. ¦ 2534    ¦ Resident Evil: Damnation  ¦ Обитель Зла: Проклятие

* Found movies with zero rating
N. ¦ column1 ¦ orig_title  ¦ rus_title
1. ¦ 2165    ¦ In the Loop ¦ В петле
 | problems=5;0;0;0;0
```

## License

Licensed under the [MIT License](http://opensource.org/licenses/MIT),
see **LICENSE**.
