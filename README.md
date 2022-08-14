To run the server, first change directory into it

```bash
$ cd server 
```

Then start it

```bash
$ go run main.go --assets_dir ../assets/ --upload_dir /tmp/screenshotdata
Listening on http://localhost:10987
```

Future work:

1. Google Chrome extension that takes a screenshot on Command+Shift+S and uploads it to the server.
2. Management page.
3. Productionization.
