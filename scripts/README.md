# Otto Script

This is a list of scripts I used to help me with the initial setup

Please update your `.env` before using it

## Initial load in the database

Small Golang script that will push data in the chat into the DB.

Not meant to be used in production just to help the setup

Update the `init-data.json` with your information

```bash
make push_init_chat
```
