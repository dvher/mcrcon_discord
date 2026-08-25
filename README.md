# MC Rcon Discord

This project is basically an implementation of the rcon protocol through a discord bot.  

Right now it isn't in active development since it's a personal project with not a lot of use  

Even though it is thought for my personal use only, you may use it with the following considerations

## Running

Right now it doesn't have a Dockerfile, it may be included in the future  

For safety reasons, run this on the same vps or wherever you may run your server, or run everything through an ssh tunnel or vpn to not expose your rcon server to the internet  

First, you must fill your .env, be sure to use your loopback IP if you're running everything through your VPS, if you're using Docker, then use your server's internal IP.

Executionwise I just use  

```bash
docker run cmd/main.go
```

after filling .env.

## Considerations

- The `start` command is working but is severely unsafe right now, I recommend just not implementing it. I recommend using `ping` to make sure the server is up and running  
- I have not considered any role setup since again, this is just a private project, be careful if you use this in a big server or an untrusted server.
