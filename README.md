kin
===

A friendly CLI for working with AWS Kinesis.

Goals
-----

The primary motivation for Kin is the need to pull records from a Kinesis stream, parse them as JSON, and output them in a way that integrates cleanly with other tools like [jq](https://github.com/stedolan/jq).

Related Work
------------

- [dinsaw/kines](https://github.com/dinsaw/kines) - A pre-existing Python implementation of the same concept. Writes its output with fancy ASCII art that makes it hard to use with other tools.