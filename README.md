This runs an SMTP server, and forwards all received emails to a hardcoded channel.
Plain-text emails are supported. HTML emails (or anything using multipart messages)
are not yet supported.

The userid/channelid used for posting are hardcoded. You can get sample values from the slash command "/email_address", which is intended to be used later on for generating the appropriate target email addresses.
The SMTP server listens specifically on 127.0.0.1:10025 right now (I'm using an SSH reverse proxy to forward remote port 25 from an online server to port 10025 on my dev machine). This address needs to be made configurable

Aside from parsing HTML emails and adding configuration options, the main thing missing is proper routing of emails to the suitable channels.
