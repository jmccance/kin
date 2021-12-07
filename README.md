kin
===

A friendly CLI for working with AWS Kinesis.

Installation
------------

Assuming you have Go installed locally, you should be able to install with:

```
go install github.com/jmccance/kin@latest
```

By default this will install into `~/go/bin`, which you will need to add to your path.

Usage
-----

```
$ kin help
A friendly CLI for working with Amazon Kinesis

Usage:
  kin [command]

Available Commands:
  help         Help about any command
  list-shards  List shards
  list-streams List Kinesis streams
  tail         Tail records from a Kinesis Data Stream

Flags:
  -h, --help             help for kin
  -p, --profile string   AWS Profile Name
  -r, --region string    AWS Region Name

Use "kin [command] --help" for more information about a command.
```

`kin` needs AWS credentials in order to access Kinesis. It should pick those up through the default credential chain, including environment variables and profiles. If you get an error like the following, it means that either you aren't signed in or `kin` isn't picking up the right profile.

```
operation error Kinesis: ListStreams, failed to sign request: failed to retrieve credentials: no EC2 IMDS role found, operation error ec2imds: GetMetadata, request canceled, context deadline exceeded
```

If you know you're signed in, the easiest thing may be to manually override the profile it's using.

```
AWS_PROFILE=staging kin ls
```

Related Work
------------

- [dinsaw/kines](https://github.com/dinsaw/kines) - A pre-existing Python implementation of the same concept. Writes its output with fancy ASCII art that makes it hard to use with other tools.
