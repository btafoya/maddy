### [\#775] Custom (configurable) SMTP Error Messages
**URL**: https://github.com/foxcpp/maddy/issues/775
**State**: OPEN
**Labels**: new feature, 

**Body:**
# Use case

I work with a small team, and we are often contacting other external small teams. This means that email deliverability across the board can be quite low, compared to larger organisations. I think I have worked out that more often than not, the issue is that an externel sender/recipient has not configured their mailserver properly - so issues with DMARC, SPF, key issues, etc. on the other side are causing Maddy to reject/quarantine incoming messages, or refuse to send outgoing messages. Secure by default is definitely good in this regard, however my local users look to me to fix all mail delivery problems. Without being an expert in the protocols involved, it can be quite difficult to pin down exactly what is happening and why a message was rejected; and users also tend to assume something is wrong with the setup, which is not always the case. 

Consider for example:
> SMTP error 550: Failed to establish the MX record authenticity (REQUIRETLS)

vs

> SMTP error 550: The recipient you're trying to contact has specified email must always be encrypted, however, their server is not accepting encrypted connections. This  risks leaking the contents of your message or suggests their server may have been compromised. Please try sending again later, or contact Maurice Moss form the IT Department if it is a persistent issue.

(I'm not sure if this is an exact description of the problem, but hopefully you get the idea - it explains to the user what to do and that it's not our (Maddy's) fault, and the risks of ignoring it).



The error messages Maddy uses are currently string literals, meaning that I can't customise them to provide more/custom information for things like bounceback messages without a recompile. They are also in English, and it is possible someone might want to run a Maddy server for a userbase where account owners are not proficient in English, and in this case a bounceback message would be difficult or impossible for them to read.

For communications with an external server it does of course make sense to stick to standarised verbiage in English for the benefit of the other side of the communication. However I think for internal messages, and possibly in select cases if the admin wants to provide more explicit information to the external party, it would be useful to be able to override default messages with custom ones, be it for language or customisation reasons. I have the same problem externally, if someone from outside the organisation can't contact us, it's rare they consider their server or DNS records might be at fault and just expect everything to work all the time so start pushing things onto my plate.


# Your idea for a solution

I expect an easy way to achieve this would be with i18n or some similar translation framework. I've very little experience with them, but assume there would be an idiomatic and easy-enough way to achieve it in Go, as in many languages, with little change to the source code at the use sites (am think about like Qt's `tr("Some hardcoded string")` function. Then, if I can override the translations in my config file and point at runtime to a new file in JSON or whatever, it would let me customise the error messages for specific situations so that I can give more guided help to my users, and reduce some of the forwarded bounces that are not caused by something I can do anything about.

Thoughts?

- [x] I'm willing to help with the implementation
- [x] I quite possibly would not have time to help with the implementation, depending on my todo list and if this remains a priority on my todo list.