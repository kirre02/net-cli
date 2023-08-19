# net-cli
net-cli is a really simple cli tool for making http calls which makes it easy to test your api.
## Reason
Reason for why I made this cli-tool is because I wanted make a cli tool that works in a similar fashion as thunderclient while still being really simple as it does not bear as much functionality as thunderclient.

## Examples
For example we could run the get command that runs a GET request.
```bash
$ net-cli get https://jsonplaceholder.typicode.com/posts/10
```
and you should get a response that looks similar to this.
```bash
$ Response:
{
  "body": "quo et expedita modi cum officia vel magni\ndoloribus qui repudiandae\nvero nisi sit\nquos veniam quod sed accusamus veritatis error",
  "id": 10,
  "title": "optio molestias id quia eum",
  "userId": 1
}
```
or you could just add the url into a config file that is called net-cli.json that can contain both a payload and a url.
we can show an example of the previous example where the url is in the config file

so inside the config file we have

```json
{
  "url":"https://jsonplaceholder.typicode.com/posts/10"
}
```

and then we run
```bash
$ net-cli get --config net-cli.json
```
and we get a the same response as before.

## Installation
### github releases
to install you could download the binary from the github releases like this.
```bash
$ wget https://github.com/kirre02/net-cli/releases/download/0.1.0/net-cli_0.1.0_linux_amd64.tar.gz
```

Then you need to extract the binary from the net-cli_0.1.0_linux_amd64.tar.gz file.
```bash
$ tar -xzf net-cli_0.1.0_linux_amd64.tar.gz
```
After you have unpacked the binary you can then move it to your /bin folder.
```bash
$ sudo mv "./net-cli" "$HOME/bin"
```
### install script
You could also just git clone the repository into a local machine and run the ```install.sh``` script like in this example below.

```bash
$ git clone https://github.com/kirre02/net-cli.git --depth=1
```
and the you can run the install script
```bash
$ ./install.sh
```
if successfull you should get a message that says
```bash
"Installation complete! The binary has been updated."
```
**_NOTE:_** this script only works for Linux and macOS