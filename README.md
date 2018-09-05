# Can I Pair Bot

A GitHub bot ([difference between bot and app][diff-bot-app]) to find available pairs in a
Slack channel.

[diff-bot-app]: https://tutorials.botsfloor.com/slack-app-or-bot-user-integration-842c3843eea8

## Prepare Slack

1. [Add Configuration][bots-page] in the Bots application page. There you'll be able to create
   a value that can used for the `SLACK_API_TOKEN` configuration variable.

[bots-page]: https://pivotal.slack.com/apps/A0F7YS25R-bots?next_id=0

## Run

1. export SLACK_API_TOKEN=[value found here](https://pivotal.slack.com/services/B9BKT74G2)
1. go run main.go

### Run on CF

1. cf api https://api.run.pcfbeta.io
1. cf login --sso
1. cf target -o pivot-smuc -s development
1. cf push --no-start
1. cf set-env can-i-pair SLACK_API_TOKEN [value found here](https://pivotal.slack.com/services/B9BKT74G2)
1. cf restage can-i-pair

# Usage

1. Invite `@can-i-pair` to your channel
1. `@can-i-pair help`

# Bugs

- Bot caches peoples status and does not update them. (Problem is probably in the library we use)
  (https://brandur.org/sdk)
- Does not work on private channels
- Does not consider time zones in scope of availability of the person

## TODO

1. Choose a place to store code
1. Be able to run in a separate org / space
1. Make the `@can-i-pair this week?` work
1. Need a concourse pipeline
1. Find mechanism to address status caching in underlying slack go library
1. Establish token rotation schedule
1. Broadcast interests in projects/pairing activities
1. integrate some timezone magic

## Maintainers

* @ScottMuc

