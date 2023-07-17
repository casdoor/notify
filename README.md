<div align="center">
<img
    width=40%
    src="assets/gopher-letter.svg"
    alt="notify logo"
/>

[![codecov](https://codecov.io/gh/nikoksr/notify/branch/v2/graph/badge.svg?token=QDON0KO2WV)](https://codecov.io/gh/nikoksr/notify)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikoksr/notify)](https://goreportcard.com/report/github.com/nikoksr/notify)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/37fdff3c275c4a72a3a061f2d0ec5553)](https://www.codacy.com/gh/nikoksr/notify/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=nikoksr/notify&amp;utm_campaign=Badge_Grade)
[![Maintainability](https://api.codeclimate.com/v1/badges/b3afd7bf115341995077/maintainability)](https://codeclimate.com/github/nikoksr/notify/maintainability)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/notify)

</div>

> <p align="center">A dead simple Go library for sending notifications to various messaging services.</p>

<h1></h1>

## Notify v2 <a id="v2"></a>

This is the home branch of Notify `v2`. It is currently in active development and not yet ready for production use. The [main branch](https://github.com/nikoksr/notify/tree/main) will stay the default branch until `v2` is ready for production use. At this point, the main branch will be renamed to `v1` and the `v2` branch will be merged into the main branch.

Notify `v2` lets you enjoy the simplicity of Notify `v1` with more power and flexibility at your hands. Providing a simple interface that lets you send attachments, define custom message renderers and dynamic message enrichment.

## About <a id="about"></a>

*Notify* was born out of my own need to have my API servers running in production be able to notify me when critical errors occur. Of course, _Notify_ can be used for any other purpose as well. The library is kept as simple as possible for quick integration and ease of use.

## Disclaimer <a id="disclaimer"></a>

Usage of this library should comply with the stated rules and the terms present in the [license](LICENSE) and [code of conduct](CODE_OF_CONDUCT.md). Failure to comply, including but not limited to misuse or spamming, can result in permanent banning from all supported platforms as governed by the respective platform's rules and regulations.

Notify's functionality is determined by the consistency of the supported external services and their corresponding latest client libraries; it may change without notice. This fact, coupled with the inevitable inconsistencies that can arise, dictates that Notify should not be used in situations where its failure or inconsistency could result in significant damage, loss, or disruption. Always have a backup plan and only use Notify after understanding and accepting these conditions.

Please read the [license](LICENSE) for a complete understanding of the permissions, conditions, and limitations of use.

## Install <a id="install"></a>

```sh
go get -u github.com/nikoksr/notify/v2
```

## Example usage <a id="usage"></a>

You can use Notify just like you're used to from `v1`. A simple example in which we send a notification to a Telegram
channel could look like this:

```go

func main() {
    // Create a new telegram service. We're using the new constructor option WithRecipients() to specify the recipients. We can,
    // however, also rely on the old way of doing things and add the recipients to the service later on using the AddRecipients()
    // method.
    svc, _ := telegram.New(token,
        telegram.WithRecipients(recipient),
    )

	// Create the actual notify instance and pass the telegram service to it. Again, we're making use of the new constructor
	// option WithServices() to specify the services. UseServices() is still available and can be used to add services later
	// on.
    n := notify.New(
        notify.WithServices(svc),
    )

	// Send a notification
    _ = n.Send(context.Background(),
        "Subject/Title",
        "The actual message - Hello, you awesome gophers! :)",
    )
}
```

We touched a little bit on what's new in `v2` in the example above. Let's take a deeper dive into the new, more advanced
features.

In this example, we're going to send a notification to a Discord channel. We're going to make use of the new
`discord.Webhook` service, which allows us to send notifications to Discord webhooks. We're also going to define a custom
message renderer, which allows us to define how the message should look like. Lastly, we're going to send a couple of
attachments and metadata along with the notification.

```go
func main() {
    // Create a new discord webhook service.
    svc, _ := discord.NewWebhook(
        discord.WithRecipients(webhookURL),
        discord.WithMessageRenderer(customRenderer),
    )

	// Open a couple of files to send as attachments.
    img, _ := os.Open("/path/to/image.png")
	defer img.Close()

    txt, _ := os.Open("/path/to/text.txt")
	defer txt.Close()

	// Create some example metadata that we make use of in our custom renderer.
    exampleMetadata := map[string]interface{}{
        "foo":  "bar",
    }

	// Send a notification with the attachments and metadata. In this case, we're using the service directly.
    _ = svc.Send(ctx,
        "[Test] Notify v2",
        "Hello, you awesome gophers! :)",
        notify.SendWithAttachments(img, txt),
        notify.SendWithMetadata(exampleMetadata),
    )
}


// The custom renderer allows us to define how the message should look like. The respective SendConfig is passed to the
// renderer, which contains all the information we need to render the message.
func customRenderer(conf discord.SendConfig) string {
    var builder strings.Builder

	// For demo purposes, we're just going to marshal the metadata to human-readable JSON and add it to the message.
    metadata, _ := json.MarshalIndent(conf.Metadata(), "", "  ")

	// Put together the message.
    builder.WriteString(conf.Subject())
    builder.WriteString("\n\n")
    builder.WriteString(conf.Message())
    builder.WriteString("\n\n")
	builder.WriteString("Metadata:\n")
	builder.WriteString(string(metadata))
	builder.WriteString("\n\n")
    builder.WriteString("<-- A super necessary footer -->\n")

    return builder.String()
}
```

## Contributing <a id="contributing"></a>

Yes, please! Contributions of all kinds are very welcome! Feel free to check our [open issues](https://github.com/nikoksr/notify/issues). Please also take a look at the [contribution guidelines](https://github.com/nikoksr/notify/blob/main/CONTRIBUTING.md).

> Psst, don't forget to check the list of [missing services](https://github.com/nikoksr/notify/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3Aaffects%2Fservices+label%3A%22help+wanted%22+no%3Aassignee) waiting to be added by you or create [a new issue](https://github.com/nikoksr/notify/issues/new?assignees=&labels=affects%2Fservices%2C+good+first+issue%2C+hacktoberfest%2C+help+wanted%2C+type%2Fenhancement%2C+up+for+grabs&template=service-request.md&title=feat%28service%29%3A+Add+%5BSERVICE+NAME%5D+service) if you want a new service to be added.

## Supported services <a id="supported_services"></a>

> Click [here](https://github.com/nikoksr/notify/issues/new?assignees=&labels=affects%2Fservices%2C+good+first+issue%2C+hacktoberfest%2C+help+wanted%2C+type%2Fenhancement%2C+up+for+grabs&template=service-request.md&title=feat%28service%29%3A+Add+%5BSERVICE+NAME%5D+service) to request a missing service.

| Service                                     | Path                                 | Credits                                                                                         |       Tested       |
|---------------------------------------------|--------------------------------------|-------------------------------------------------------------------------------------------------|:------------------:|
| [Discord](https://discord.com)              | [service/discord](service/discord)   | [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)                                     | :heavy_check_mark: |
| [Mail](https://en.wikipedia.org/wiki/Email) | [service/mail](service/mail)         | [xhit/go-simple-mail/v2](https://github.com/xhit/go-simple-mail)                                |        :x:         |
| [Ntfy](https://ntfy.sh)                     | [service/ntfy](service/ntfy)         | -                                                                                               | :heavy_check_mark: |
| [Slack](https://slack.com)                  | [service/slack](service/slack)       | [slack-go/slack](https://github.com/slack-go/slack)                                             |        :x:         |
| [Telegram](https://telegram.org)            | [service/telegram](service/telegram) | [go-telegram-bot-api/telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) | :heavy_check_mark: |

## Special Thanks <a id="special_thanks"></a>

### Maintainers <a id="maintainers"></a>

- [@svaloumas](https://github.com/svaloumas)

### Logo <a id="logo"></a>

The [logo](https://github.com/MariaLetta/free-gophers-pack) was made by the amazing [MariaLetta](https://github.com/MariaLetta).

## Similar projects <a id="similar_projects"></a>

> Just to clarify, Notify was not inspired by any other project. I created it as a tiny subpackage of a larger project and only later decided to make it a standalone project. In this section I just want to mention other great projects.

  - [containrrr/shoutrrr](https://github.com/containrrr/shoutrrr)
  - [caronc/apprise](https://github.com/caronc/apprise)

## Show your support <a id="support"></a>

Please give a ⭐️ if you like the project! It draws more attention to the project, which helps us improve it even faster.

