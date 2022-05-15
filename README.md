# Mock Slackbot

A slack bot that creates a mock message on demand.


## Notes

Reason for using channel for processing messages

> De-couple your ingestion of events from processing and reacting to them.
> Especially when working with large workspaces, many workspaces, or subscribing to a large number of events.
> Quickly respond to events with HTTP 200 and add them to a queue before doing amazing things with them.
