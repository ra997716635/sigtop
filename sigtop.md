# SIGTOP(1) General Commands Manual SIGTOP(1)

NAME sigtop -- export messages from Signal Desktop

SYNOPSIS sigtop command [argument ...]

DESCRIPTION sigtop is a utility to export messages, attachments and
other data from Signal Desktop.

sigtop needs access to the directory where Signal Desktop stores its
data. By default, sigtop tries the following locations:

Windows %AppData%ignal

macOS $HOME/Library/Application Support/Signal

Unix $XDG_CONFIG_HOME/Signal, or $HOME/.config/Signal if
XDG_CONFIG_HOME is unset or empty. If neither directory exists, sigtop
tries the locations used by the unofficial Snap and Flatpak versions of
Signal Desktop. These locations are
$HOME/snap/signal-desktop/current/.config/Signal and
$HOME/.var/app/org.signal.Signal/config/Signal, respectively.

A different location for the Signal Desktop directory can be specified
using the -d option (see below).

The commands are as follows:

check-database [-d signal-directory] (Alias: check)

Check the integrity of the encrypted Signal Desktop database. The check
is performed using the SQLite/SQLCipher cipher_integrity_check,
integrity_check and foreign_key_check pragmas.

export-attachments [-iLlMm] [-d signal-directory] [-s interval]
[directory] (Alias: att)

Export attachments. The attachment files are created in separate
directories, one for each conversation. These directories are created in
directory, or in the current directory if directory is not specified.

Signal Desktop stores attachments in unencrypted files. sigtop can copy
these files, or create hard links or symbolic links to them. By default,
attachments are copied. If the -L option is specified, hard links are
created. If the -l option is specified, symbolic links are created.

The -M option may be used to set the file modification time of each
attachment to the time it was sent. The -m option is similar, but uses
the time the attachment was received. These options are ignored if -L is
also specified.

By default, all attachments are exported. If the -i option is specified,
an incremental export is performed. This means that only new attachments
are exported; attachments that were exported in a previous run are
skipped. The .incremental file in directory is used to keep track of
exported attachments.

The -s option may be used to export only the attachments that were sent
in the specified time interval. See the TIME INTERVALS section below for
details.

export-database [-d signal-directory] file (Alias: db)

Decrypt and export the Signal Desktop database to file. The exported
database is a regular SQLite database.

export-messages [-i] [-d signal-directory] [-f format] [-s
interval] [directory] (Alias: msg)

Export messages. The messages are written to separate files, one for
each conversation. These files are created in directory, or in the
current directory if directory is not specified.

The -f option may be used to specify the output format. The following
output formats are supported:

json Messages are written in JSON format. The JSON data is copied
directly from the Signal Desktop database, so its structure may differ
between Signal Desktop versions.

text Messages are written as plain text. This is the default.

text-short Messages are written as plain text, in short form. Every
message is written on a single line.

By default, existing files in directory are not overwritten. If the -i
option is specified, an incremental export is performed. This means that
existing conversation files are updated.

By default, all messages are exported. The -s option may be used to
export only the messages that were sent in the specified time interval.
See the TIME INTERVALS section below details.

query-database [-d signal-directory] query (Alias: query)

Query the Signal Desktop database. The results of the query, if any, are
written to standard output. Only the first statement in query is
evaluated. Furthermore, because the Signal Desktop database is opened in
read-only mode, statements attempting to modify the database will fail.

TIME INTERVALS A time is specified as
'yyyy[-mm[-dd[Thh[:mm[:ss]]]]]'. For example:

2023-01-23T12:34:56 2023-01-23T12:34 2023-01 2023

A time interval is specified either as '[min-time],[max-time]' or as
'time'. In the first form, min-time and max-time are the endpoints of
the time interval. The endpoints are inclusive.

Each omitted time field in min-time defaults to the smallest possible
value for that time field. Analogously, each omitted time field in
max-time defaults to the largest possible value for that time field. For
example, the interval '2023-02,2023' is equivalent to:

2023-02-01T00:00:00,2023-12-31T23:59:59

Furthermore, either endpoint of the time interval may be omitted. For
example, the interval from the start of February 2023 to now may be
specified as '2023-02,'.

Time intervals may also be specified in a second form, consisting of a
single time specification. In this form, the same time specification is
used for both endpoints. For example, the time interval '2023' is
equivalent to '2023,2023', which is equivalent to:

2023-01-01T00:00:00,2023-12-31T23:59:59

EXIT STATUS The sigtop utility exits 0 on success, and >0 if an error
occurs.

EXAMPLES Export all messages to the directory messages:

$ sigtop export-messages messages

Use the shorter command alias:

$ sigtop msg messages

Export all messages in JSON format to the directory json:

$ sigtop msg -f json json

Export all messages in CSV format to the directory csv_output:

$ sigtop msg -f csv csv_output

Export all attachments sent at or after 12:34:56 on 23 January 2021 to
the directory attachments:

$ sigtop att -s 2021-01-23T12:34:56, attachments

Export the database from a Signal Desktop directory on a Windows disk
mounted at /mnt:

$ sigtop db -d /mnt/Users/Alice/AppData/Roaming/Signal signal.db

SEE ALSO <https://github.com/tbvdm/sigtop>

AUTHORS The sigtop utility was written by Tim van der Molen
<tim@kariliq.nl>.

macOS 13.5 July 22, 2023 macOS 13.5
