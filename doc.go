/*
 * Copyright (c) 2013-2014, Jeremy Bingham (<jbingham@gmail.com>)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
Goiardi is an implementation of the Chef server (http://www.opscode.com) written
in Go. It can either run entirely in memory with the option to save and load the
in-memory data and search indexes to and from disk, drawing inspiration from
chef-zero, or it can use MySQL or PostgreSQL as its storage backend.

Like all software, it is a work in progress. Goiardi now, though, should have all
the functionality of the open source Chef Server, plus some extras like reporting
and event logging. It does not support other Enterprise Chef type features like
organizations or pushy at this time. When used, knife works, and chef-client runs
complete successfully. Almost all chef-pendant tests successfully successfully
run, with a few disagreements about error messages that don't impact the clients.
It does pretty well against the official chef-pedant, but because goiardi handles
some authentication matters a little differently than the official chef-server,
there is also a fork of chef-pedant located at
https://github.com/ctdk/chef-pedant that's more custom tailored to goiardi.

Many go tests are present as well in different goiardi subdirectories.

Goiardi currently has nine dependencies: go-flags, go-cache, go-trie, toml, the
mysql driver from go-sql-driver, the postgres driver, logger, go-uuid, and serf.

To install them, run:

   go get github.com/jessevdk/go-flags
   go get github.com/pmylund/go-cache
   go get github.com/ctdk/go-trie/gtrie
   go get github.com/BurntSushi/toml
   go get github.com/go-sql-driver/mysql
   go get github.com/lib/pq
   go get github.com/ctdk/goas/v2/logger
   go get github.com/codeskyblue/go-uuid
   go get github.com/hashicorp/serf/client

from your $GOROOT, or just use the -t flag when you go get goiardi.

If you would like to modify the search grammar, you'll need the 'peg' package.
To install that, run

   go get github.com/pointlander/peg

In the 'search/' directory, run 'peg -switch -inline search-parse.peg' to
generate the new grammar. If you don't plan on editing the search grammar,
though, you won't need that.

Installation

1. Install go. (http://golang.org/doc/install.html) You may need to upgrade to
go 1.2 to compile all the dependencies. Go 1.3 is also confirmed to work.

2. Make sure your $GOROOT and PATH are set up correctly per the Go installation
instructions.

3. Download goairdi and its dependencies.

   go get -t github.com/ctdk/goiardi

4. Run tests, if desired. Several goiardi subdirectories have go tests, and
chef-pedant can and should be used for testing goiardi as well.

5. Install the goiardi binaries.

   go install github.com/ctdk/goiardi

6. Run goiardi.

   goiardi <options>

Or, you can look at the goiardi releases page on github at
https://github.com/ctdk/goiardi/releases and see if there are precompiled
binaries available for your platform.

You can get a list of command-line options with the '-h' flag.

Goiardi can also take a config file, run like goiardi -c
/path/to/conf-file. See etc/goiardi.conf-sample for an example documented
configuration file. Options in the configuration file share the same name
as the long command line arguments (so, for example, --ipaddress=127.0.0.1
on the command line would be ipaddress = "127.0.0.1" in the config file.

Currently available command line and config file options:

   -v, --version          Print version info.
   -V, --verbose          Show verbose debug information. Repeat for more
			  verbosity.
   -c, --config=          Specify a config file to use.
   -I, --ipaddress=       Listen on a specific IP address.
   -H, --hostname=        Hostname to use for this server. Defaults to hostname
                          reported by the kernel.
   -P, --port=            Port to listen on. If port is set to 443, SSL will be
                          activated. (default: 4545)
   -i, --index-file=      File to save search index data to.
   -D, --data-file=       File to save data store data to.
   -F, --freeze-interval= Interval in seconds to freeze in-memory data
                          structures to disk (requires -i/--index-file and
                          -D/--data-file options to be set). (Default 300
                          seconds/5 minutes.)
   -L, --log-file=        Log to file X
   -s, --syslog           Log to syslog rather than a log file. Incompatible
                          with -L/--log-file.
       --time-slew=       Time difference allowed between the server's clock and
                          the time in the X-OPS-TIMESTAMP header. Formatted like
                          5m, 150s, etc. Defaults to 15m.
       --conf-root=       Root directory for configs and certificates. Default:
                          the directory the config file is in, or the current
                          directory if no config file is set.
   -A, --use-auth         Use authentication. Default: false.
       --use-ssl          Use SSL for connections. If --port is set to 433, this
                          will automatically be turned on. If it is set to 80,
                          it will automatically be turned off. Default: off.
                          Requires --ssl-cert and --ssl-key.
       --ssl-cert=        SSL certificate file. If a relative path, will be set
                          relative to --conf-root.
       --ssl-key=         SSL key file. If a relative path, will be set relative
                          to --conf-root.
       --https-urls       Use 'https://' in URLs to server resources if goiardi
                          is not using SSL for its connections. Useful when
                          goiardi is sitting behind a reverse proxy that uses
                          SSL, but is communicating with the proxy over HTTP.
       --disable-webui    If enabled, disables connections and logins to goiardi
                          over the webui interface.
       --use-mysql        Use a MySQL database for data storage. Configure
                          database options in the config file.
       --use-postgresql   Use a PostgreSQL database for data storage.
                          Configure database options in the config file.
       --local-filestore-dir= Directory to save uploaded files in. Optional when
                          running in in-memory mode, *mandatory* for SQL
                          mode.
       --log-events       Log changes to chef objects.
   -K, --log-event-keep=  Number of events to keep in the event log. If set,
                          the event log will be checked periodically and
                          pruned to this number of entries.
   -x, --export=          Export all server data to the given file, exiting
                          afterwards. Should be used with caution. Cannot be
                          used at the same time as -m/--import.
   -m, --import=          Import data from the given file, exiting
                          afterwards. Cannot be used at the same time as
                          -x/--export.
   -Q, --obj-max-size=    Maximum object size in bytes for the file store.
                          Default 10485760 bytes (10MB).
   -j, --json-req-max-size= Maximum size for a JSON request from the client.
                          Per chef-pedant, default is 1000000.
       --use-unsafe-mem-store Use the faster, but less safe, old method of
                          storing data in the in-memory data store with
                          pointers, rather than encoding the data with gob
                          and giving a new copy of the object to each
                          requestor. If this is enabled goiardi will run
                          faster in in-memory mode, but one goroutine could
                          change an object while it's being used by
                          another. Has no effect when using an SQL backend.
       --db-pool-size=    Number of idle db connections to maintain. Only
                          useful when using one of the SQL backends.
                          Default is 0 - no idle connections retained
       --max-connections= Maximum number of connections allowed for the
                          database. Only useful when using one of the SQL
                          backends. Default is 0 - unlimited.
       --use-serf         If set, have goidari use serf to send and receive
                          events and queries from a serf cluster. Required
                          for shovey.
       --serf-event-announce Announce log events over serf and joining the serf
                          cluster, as serf events. Requires --use-serf.
       --serf-addr=       IP address and port to use for RPC communication
                          with a serf agent. Defaults to 127.0.0.1:7373.
       --use-shovey       Enable using shovey for sending jobs to nodes.
                          Requires --use-serf.
       --sign-priv-key=   Path to RSA private key used to sign shovey
                          requests.

   Options specified on the command line override options in the config file.

For more documentation on Chef, see http://docs.opscode.com.

If goiardi is not running in use-auth mode, it does not actually care about .pem
files at all. You still need to have one to keep knife and chef-client happy.
It's like chef-zero in that regard.

If goiardi is running in use-auth mode, then proper keys are needed. When
goiardi is started, if the chef-webui and chef-validator clients, and the admin
user, are not present, it will create new keys in the --conf-root directory. Use
them as you would normally for validating clients, performing tasks with the
admin user, or using chef-webui if webui will run in front of goiardi.

In auth mode, goiardi supports versions 1.0, 1.1, and 1.2 of the Chef
authentication protocol.

*Note:* The admin user, when created on startup, does not have a password. This
prevents logging in to the webui with the admin user, so a password will have to
be set for admin before doing so.

Upgrading

Upgrading goiardi is generally a straightforward process. Usually all you should
need to do is get the new sources and rebuild, or download the appropriate new
binary. However, sometimes a little more work is involved. Check the release
notes for the new release in question for any extra steps that may need to be
done. If you're running one of the SQL backends, you may need to apply database
patches (either with sqitch or by hand), and in-memory mode especially may
require using the data import/export functionality to dump and load your chef
data between upgrades if the binary save file compatibility breaks between
releases. However, while it should not happen often, occasionally more serious
preparation will be needed before upgrading. It won't happen without a good
reason, and the needed steps will be clearly outlined to make the process as
painless as possible.

As a special note, if you are upgrading from any release prior to 0.6.1-pre1 to
0.7.0 and are using one of the SQL backends, the upgrade is one of the special
cases. Between those releases the way the complex data structures associated
with cookbook versions, nodes, etc. changed from using gob encoding to json
encoding. It turns out that while gob encoding is indeed faster than json (and
was in all the tests I had thrown at it) in the usual case, in this case json is
actually significantly faster, at least once there are a few thousand coobkooks
in the database. In-memory datastore (including file-backed in-memory datastore)
users are advised to dump and reload their data between upgrading from <=
0.6.1-pre1 and 0.7.0, but people using either MySQL or Postgres *have* to do
these things:

* Export their goiardi server's data with the `-x` flag.
* Either revert all changes to the db with sqitch, then redeploy, or drop the
database manually and recreate it from either the sqitch patches or the full
table dump of the release (provided starting with 0.7.0)
* Reload the goiardi data with the `-m` flag.

It's a fairly quick process (a goiardi dump with the `-x` flag took 15 minutes
or so to load with over 6200 cookbooks) at least, but if you don't do it very
little of your goiardi environment will work correctly. The above steps will
take care of it.

Logging

By default, goiardi logs to standard output. A log file may be specified with the
`-L/--log-file` flag, or goiardi can log to syslog with the `-s/--syslog` flag on
platforms that support syslog. Attempting to use syslog on one of these platforms
(currently Windows and plan9 (although plan9 doesn't build for other reasons))
will result in an error.

Log levels

Log levels can be set in goiardi with either the `log-level` option in the
configuration file, or with one to four -V flags on the command line. Log level
options are "debug", "info", "warning", "error", and "critical". More -V on the
command line means more spewing into the log.

MySQL mode

Goiardi can now use MySQL to store its data, instead of keeping all its data
in memory (and optionally freezing its data to disk for persistence).

If you want to use MySQL, you (unsurprisingly) need a MySQL installation that
goiardi can access. This document assumes that you are able to install,
configure, and run MySQL.

Once the MySQL server is set up to your satisfaction, you'll need to install
sqitch to deploy the schema, and any changes to the database schema that may come
along later. It can be installed out of CPAN or homebrew; see "Installation" on
http://sqitch.org for details.

The sqitch MySQL tutorial at https://metacpan.org/pod/sqitchtutorial-mysql
explains how to deploy, verify, and revert changes to the database with sqitch,
but the basic steps to deploy the schema are:

* Create goiardi's database: `mysql -u root --execute 'CREATE DATABASE goiardi'`

* Optionally, create a separate mysql user for goiardi and give it permissions on that database.

* In sql-files/mysql-bundle, deploy the bundle: `sqitch deploy db:mysql://root[:<password>]@/goiardi`

To update an existing database deployed by sqitch, run the `sqitch deploy`
command above again.

If you really really don't want to install sqitch, apply each SQL patch in
sql-files/mysql-bundle by hand in the same order they're listed in the
sqitch.plan file.

The above values are for illustration, of course; nothing requires goiardi's
database to be named "goiardi". Just make sure the right database is specified in
the config file.

Set `use-mysql = true` in the configuration file, or specify `--use-mysql` on
the command line. It is an error to specify both the `-D`/`--data-file` flag and
`--use-mysql` at the same time.

At this time, the mysql connection options have to be defined in the config
file. An example configuration is available in `etc/goiardi.conf-sample`, and is
given below:

	[mysql]
		username = "foo" # technically optional, although you probably want it
		password = "s3kr1t" # optional, if you have no password set for MySQL
		protocol = "tcp" # optional, but set to "unix" for connecting to MySQL
				 # through a Unix socket.
		address = "localhost"
		port = "3306" # optional, defaults to 3306. Not used with sockets.
		dbname = "goiardi"
		# See https://github.com/go-sql-driver/mysql#parameters for an
		# explanation of available parameters
		[mysql.extra_params]
			tls = "false"

Postgres mode

Goiardi can also use Postgres as a backend for storing its data, instead of using
MySQL or the in-memory data store. The overall procedure is pretty similar to
setting up goiardi to use MySQL. Specifically for Postgres, you may want to
create a database especially for goiardi, but it's not mandatory. If you do, you
may also want to create a user for it. If you decide to do that:

* Create the user: `$ createuser goiardi <additional options>`

* Create the database, if you decided to: `$ createdb goiardi_db <additional options>`. If you created a user, make it the owner of the goiardi db with `-O goiardi`.

After you've done that, or decided to use an existing database and user, deploy
the sqitch bundle in sql-files/postgres-bundle. If you're using the default
Postgres user on the local machine, `sqitch deploy db:pg:<dbname>` will be
sufficient. Otherwise, the deploy command will be something like `sqitch deploy db:pg://user:password@localhost/goairdi_db`.

The Postgres sqitch tutorial at https://metacpan.org/pod/sqitchtutorial explains more about how to use sqitch and Postgres.

Set `use-postgresql` in the configuration file, or specify `--use-postgresql` on
the command line. It's also an error to specify both `-D`/`--data-file` flag and
`--use-postgresql` at the same time like it is in MySQL mode. MySQL and Postgres
cannot be used at the same time, either.

Like MySQL, the Postgres connection options must be specified in the config file
at this time. There is also an example Postgres configuration in the config file,
and can be seen below:

	# PostgreSQL options. If "use-postgres" is set to true on the command line or in
	# the configuration file, connect to postgres with the options in [postgres].
	# These options are all strings. See
	# http://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters for details
	# on the connection parameters. All of these parameters are technically optional,
	# although chances are pretty good that you'd want to set at least some of them.
	[postgresql]
		username = "foo"
		password = "s3kr1t"
		host = "localhost"
		port = "5432"
		dbname = "mydb"
		sslmode = "disable"

General Database Options

There are two general options that can be set for either database:
`--db-pool-size` and `--max-connections` (and their configuration file
equivalents `db-pool-size` and `max-connections`). `--db-pool-size` sets the
number of idle connections to keep open to the database, and `--max-connections`
sets the maximum number of connections to open on the database. If they are not
set, the default behavior is to keep no idle connections alive and to have
unlimited connections to the database.

It should go without saying that these options don't do much if you aren't using
one of the SQL backends.

Event Logging

Goiardi has optional event logging. When enabled with the `--log-events` command
line option, or with the `"log-events"` option in the config file, changes to
clients, users, cookbooks, data bags, environments, nodes, and roles will be
tracked. The event log can be viewed through the /events API endpoint.

If the `-K`/`--log-event-keep` option is set, then once a minute the event log
will be automatically purged, leaving that many events in the log. This is particularly recommended when using the event log in in-memory mode.

The event API endpoints work as follows:

	`GET /events` - optionally taking `offset`, `limit`, `from`, `until`, `object_type`, `object_name`, and `doer` query parameters.

	List the logged events, starting with the most recent. Use the `offset`
	and `limit` query parameters to view smaller chunks of the event log at
	one time.  The `from`, `until`, `object_type`, `object_name`, and `doer`
	query parameters can be used to narrow the results returned further, by
	time range (for `from` and `until`), the type of object and the name of
	the object (for `object_type` and `object_name`) and the name of the
	performer of the action (for `doer`).  These options may be used in
	singly or in concert.

	`DELETE /events?purge=1234` - purge logged events older than the given id from the event log.

	`GET /events/1234` - get a single logged event with the given id.

	`DELETE /events/1234` - delete a single logged event from the event log.

A user or client must be an administrator account to use the `/events` endpoint.

The data returned by an event should look something like this:

	{
	  "actor_info": "{\"username\":\"admin\",\"name\":\"admin\",\"email\":\"\",\"admin\":true}\n",
	  "actor_type": "user",
	  "time": "2014-05-06T07:40:12Z",
	  "action": "delete",
	  "object_type": "*client.Client",
	  "object_name": "pedant_testclient_1399361999-483981000-42305",
	  "extended_info": "{\"name\":\"pedant_testclient_1399361999-483981000-42305\",\"node_name\":\"pedant_testclient_1399361999-483981000-42305\",\"json_class\":\"Chef::ApiClient\",\"chef_type\":\"client\",\"validator\":false,\"orgname\":\"default\",\"admin\":true,\"certificate\":\"\"}\n",
	  "id": 22
	}

The easiest way to use the event log is with the knife-goiardi-event-log knife
plugin. It's available on rubygems, or at github at
https://github.com/ctdk/knife-goiardi-event-log.

Reporting

Goiardi now supports Chef's reporting facilities. Nothing needs to be enabled
in goiardi to use this, but changes are required with the client. See
http://docs.opscode.com/reporting.html for details on how to enable reporting
and how to use it.

There is a goiardi extension to reporting: a "status" query parameter may be
passed in a GET request that lists reports to limit the reports returned to ones
that match the status, so you can read only reports of chef runs that were
successful, failed, or started but haven't completed yet. Valid values for the
"status" parameter are "started", "success", and "failure".

To use reporting, you'll either need the Chef knife-reporting plugin, or use the
knife-goiardi-reporting plugin that supports querying runs by status. It's
available on rubygems, or on github at
https://github.com/ctdk/knife-goiardi-reporting.

Import and Export of Data

Goiardi can now import and export its data in a JSON file. This can help both
when upgrading, when the on-disk data format changes between releases, and to
convert your goiardi installation from in-memory to MySQL (or vice versa). The
JSON file has a version number set (currently 1.0), so that in the future if
there is some sort of incompatible change to the JSON file format the importer
will be able to handle it.

Before importing data, you should back up any existing data and index files (and
take a snapshot of the SQL db, if applicable) if there's any reason you might
want it around later. After exporting, you may wish to hold on to the old
installation data until you're satisfied that the import went well.

Remember that the JSON export file contains the client and user public keys
(which for the purposes of goiardi and chef are private) and the user hashed
passwords and password salts. The export file should be guarded closely.

The `-x/--export` and `-m/--import` flags control importing and exporting data.
To export data, stop goiardi, then run it again with the same options as before
but adding `-x <filename>` to the command. This will export all the data to the
given filename, and goiardi will exit.

Importing is ever so slightly trickier. You should remove any existing data store
and index files, and if using an SQL database use sqitch to revert and deploy all
of the SQL files to set up a completely clean schema for goiardi. Then run
goiardi with the new options like you normally would, but add `-m <filename>`.
Goiardi will run, import the new data, and exit. Assuming it went well, the data
will be all imported. The export dump does not contain the user and client .pem
files, so those will need to be saved and moved as needed.

Theoretically a properly crafted export file could be used to do bulk loading of
data into goiardi, thus goiardi does not wipe out the existing data on its own
but rather leaves that task to the administrator. This functionality is merely
theoretical and completely untested. If you try it, you should back your data
up first.

Berks Universe Endpoint

Starting with version 0.6.1, goiardi supports the berks-api `/universe`
endpoint. It returns a JSON list of all the cookbooks and their versions that
have been uploaded to the server, along with the URL and dependencies of each
version. The requester will need to be properly authenticated with the server to
use the universe endpoint.

The universe endpoint works with all backends, but with a ridiculous number of
cookbooks (like, loading all 6000+ cookbooks in the Chef Supermarket), the
Postgres implementation is able to take advantage of some Postgres specific
functionality to generate that page significantly faster than the in-mem or
MySQL implementations. It's not too bad, but on my laptop at home goiardi could
generate /universe against the full 6000+ cookbooks of the supermarket in ~350
milliseconds, while MySQL took about 1 second and in-mem took about 1.2 seconds.
Normal functionality is OK, but if you have that many cookbooks and expect to
use the universe endpoint often you may wish to consider using Postgres.

Serf

As of version 0.8.0, goiardi has some serf integration. At the moment it's
mainly used for shovey (see below), but it will also announce that it's started
up and joined a serf cluster.

If the `--serf-event-announce` flag is set, goiardi will announce logged events
from the event log and starting up and joining the serf cluster over serf as
serf user events. Be aware that if this is enabled, something will need to read
these events from serf. Otherwise, the logged events will pile up and eventually
take up all the space in the event queue and prevent any new events from being
added.

Shovey

Shovey is a facility for sending jobs to nodes independently of a chef-client
run, like Chef Push but serf based.

Shovey requirements

To use shovey, you will need:

* Serf installed on the server goiardi is running on.
* Serf installed on the node(s) running jobs.
* `schob`, the shovey client, must be installed on the node(s) running jobs.
* The `knife-shove` plugin must be installed on the workstation used to manage
  shovey jobs.

The client can be found at https://github.com/ctdk/schob, and a cookbook for
installing the shovey client on a node is at
https://github.com/ctdk/shovey-jobs. The `knife-shove` plugin can be found at
https://github.com/ctdk/knife-shove or on rubygems.

Shovey Installation

Setting goiardi up to use shovey is pretty straightforward.

* Once goiardi is installed or updated, install serf and run it with
  `serf agent`. Make sure that the serf agent is using the same name for its
  node name that goiardi is using for its server name.
* Generate an RSA public/private keypair. Goiardi will use this to sign its
  requests to the client, and schob will verify the requests with it.

	$ openssl genrsa -out shovey.pem 2048 # generate 2048 bit private key
	$ openssl rsa -in shovey.pem -pubout -out shovey.key # public key

  Obviously, save these keys.

* Run goiardi like you usually would, but add these options:
  `--use-serf --use-shovey --sign-priv-key=/path/to/shovey.pem`
* Install serf and schob on a chef node. Ensure that the serf agent on the node
  is using the same name as the chef node. The `shovey-jobs` cookbook makes
  installing schob easier, but it's not too hard to do by hand by running
  `go get github.com/ctdk/schob` and `go install github.com/ctdk/schob`.
* If you didn't use the shovey-jobs cookbook, make sure that the public key you
  generated earlier is uploaded to the node somewhere.
* Shovey uses a whitelist to allow jobs to run on nodes. The shovey whitelist is
  a simple JSON hash, with job names as the keys and the commands to run as the
  values. There's a sample whitelist file in the schob repo at
  `test/whitelist.json`, and the shovey-jobs cookbook will create a whitelist
  file from Chef node attributes using the usual precedence rules. The whitelist
  is drawn from `node["schob"]["whitelist"]`.
* If you used the shovey-jobs cookbook schob should be running already. If not,
  start it with something like `schob -VVVV -e http://chef-server.local:4545 -n
  node-name.local -k /path/to/node.key -w /path/to/schob/test/whitelist.json -p
  /path/to/public.key --serf-addr=127.0.0.1:7373`. Within a minute, goiardi
  should be aware that the node is up and ready to accept jobs.

At this point you should be able to submit jobs and have them run. The
knife-shove documentation goes into detail on what actions you can take with
shovey, but to start try `knife goiardi job start ls <node name>`. To list jobs,
run `knife goiardi job list`. You can also get information on a shovey job,
detailed information of a shovey job's run on one node, cancel jobs, query node
status, and stream job output from a node with the knife-shove plugin. See the
plugin's documentation for more information.

See the serf docs at http://www.serfdom.io/docs/index.html for more information
on setting up serf. One serf option you may want to use, once you're satisfied
that shovey is working properly, is to use encryption with your serf cluster.

The shovey API is documented in the file `shovey_api.md` in the goiardi
repository.

Shovey In More Detail

Every thirty seconds, schob sends a heartbeat back to goiardi over serf to let
goiardi know that the node is up. Once a minute, goiardi pulls up a list of
nodes that it hasn't seen in the last 10 minutes and marks them as being down.
If a node that is down comes back up and sends a heartbeat back to goiardi, it
is marked as being up again. The node statuses are tracked over time as well, so
a motivated user could track node availability over time.

When a shovey run is submitted, goiardi determines which nodes are to be
included in the run, either via the search function or from being listed on the
command line. It then sees how many of the nodes are believed to be up, and
compares that number with the job's quorum. If there aren't enough nodes up to
satisfy the quorum, the job fails.

If the quorum is satisfied, goiardi sends out a serf query with the job's
parameters to the nodes that will run the shovey job, signed with the shovey
private key. The nodes verify the job's signature and compare the job's command
to the whitelist, and if it checks out begin running the job.

As the job runs, schob will stream the command's output back to goiardi. This
output can in turn be streamed to the workstation managing the shovey jobs, or
viewed at a later time. Meanwhile, schob also watches for the job to complete,
receiving a cancellation command from goiardi, or to timeout because it was
running too long. Once the job finishes or is cancelled or killed, schob sends
a report back to goiardi detailing the job's run on that node.

Tested Platforms

Goiardi has been built and run with the native 6g compiler on Mac OS X (10.7,
10.8, and 10.9), Debian squeeze and wheezy, a fairly recent Arch Linux, FreeBSD
9.2, and Solaris. Using Go's cross compiling capabilities, goiardi builds for
all of Go's supported platforms except Dragonfly BSD and plan9 (because of
issues with the postgres client library). Windows support has not been tested
extensively, but a cross compiled binary has been tested successfully on
Windows.

Goiardi has also been built and run with gccgo (using the "-compiler gccgo"
option with the "go" command) on Arch Linux. Building it with gccgo without
the go command probably works, but it hasn't happened yet. This is a priority,
though, so goiardi can be built on platforms the native compiler doesn't support
yet.

Note regarding goiardi persistence and freezing data

As mentioned above, goiardi can now freeze its in-memory data store and index to
disk if specified. It will save before quitting if the program receives a
SIGTERM or SIGINT signal, along with saving every "freeze-interval" seconds
automatically.

Saving automatically helps guard against the case where the server receives a
signal that it can't handle and forces it to quit. In addition, goiardi will not
replace the old save files until the new one is all finished writing. However,
it's still not anywhere near a real database with transaction protection, etc.,
so while it should work fine in the general case, possibilities for data loss
and corruption do exist. The appropriate caution is warranted.

Documentation

In addition to the aforementioned Chef documentation at http://docs.opscode.com,
more documentation specific to goiardi can be viewed with godoc. See
http://godoc.org/code.google.com/p/go.tools/cmd/godoc for an explanation of how
godoc works.

To do

See the TODO file for an up-to-date list of what needs to be done. There's a
lot.

Bugs

There's going to be a lot of these for a while, so we'll just keep those in a
BUGS file, won't we?

Why

This started as a project to learn Go, and because I thought that an in memory
chef server would be handy. Then I found out about chef-zero, but I still wanted
a project to learn Go, so I kept it up. Chef 11 Server also only runs under
Linux at this time, while Goiardi is developed under Mac OS X and ought to run
under any platform Go supports (only partially at this time though).

If you feel like contributing, great! Just fork the repo, make your
improvements, and submit a pull request. Tests would, of course, be appreciated.
Adding tests where there are no tests currently would be even more appreciated.
At least, though, try and not break anything worse than it is. Test coverage has
improved, but is still an ongoing concern.

Goiardi is authored and copyright (c) Jeremy Bingham, 2013.  Like many Chef
ecosystem programs, goairdi is licensed under the Apache 2.0 License. See the
LICENSE file for details.

Chef is copyright (c) 2008-2013 Opscode, Inc. and its various contributors.

Thanks go out to the fine folks of Opscode and the Chef community for all their
hard work.

Also, if you were wondering, Ettore Boiardi was the man behind Chef Boyardee. Wakka wakka.

*/
package main
