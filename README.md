# HTML Mailer
[![Build Status](https://travis-ci.com/inkychris/htmlmailer.svg?branch=master)](https://travis-ci.com/inkychris/htmlmailer)

This project is a simple app
which automates the process
of downloading HTML from a web-page
and emailing it to a list of recipients.

```yaml
login:
  credentials:
    username: bob@example.com
    password: wxyz6789
  form:
    action: https://some.site.com/login
    username_field: user[email]
    password_field: user[password]
target_url: https://some.site.com/user/bob/example
schedule: 0 7 * * 2
email:
  to:
    - alice@example.com
    - charlie@example.com
  from: bob@example.com
  subject: Bob's Example

smtp:
  host: example.com
  port: 465
  username: bob@example.com
  password: abcd1234
```
