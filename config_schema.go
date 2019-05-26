package main

const configSchema = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "id": "html_mailer.json",
  "title": "HTML Mailer configuration schema",
  "description": "Configuration for HTML Mailer",
  "type": "object",
  "required": [
    "target_url",
    "email",
    "smtp"
  ],
  "additionalProperties": false,
  "properties": {
    "login": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "credentials": {
          "type": "object",
          "additionalProperties": false,
          "properties": {
            "username": {"type": "string"},
            "password": {"type": "string"}
          }
        },
        "form": {
          "type": "object",
          "additionalProperties": false,
          "properties": {
            "action": {"type": "string"},
            "username_field": {"type": "string"},
            "password_field": {"type": "string"}
          }
        }
      }
    },
    "target_url": {"type": "string"},
    "schedule": {"type": "string"},
    "email": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "to": {
          "type": "array",
          "items": {"$ref": "#/definitions/email_address"},
          "minItems": 1
        },
        "from": {"$ref": "#/definitions/email_address"},
        "subject": {"type": "string"}
      }
    },
    "smtp": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "host": {"type": "string"},
        "port": {"type": "integer"},
        "username": {"type": "string"},
        "password": {"type": "string"}
      }
    }
  },

  "definitions": {
    "email_address": {
      "type": "string",
      "pattern": "^(?i)[A-Z0-9._%+-]+@[A-Z0-9.-]+\\.[A-Z]{2,}$"
    }
  }
}
`
