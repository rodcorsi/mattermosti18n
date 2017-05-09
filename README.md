# Mattermosti18n

This tool was developed to convert JSON translated files from the project [Mattermost](https://github.com/mattermost/platform) to PO used in [Pootle](http://translate.mattermost.com/projects/) server.

## Installation

### Automatic

Either install a package made for your GNU/Linux distribution ([example for Arch Linux](https://aur.archlinux.org/packages/mattermosti18n-git/)), or if you have `go` installed, execute the build and install process all-in-one with this oneliner:
```
$ go get github.com/rodcorsi/mattermosti18n/...
```

### Manual

The manual way could be useful if, for some reason, you need to cut the build and install process in several steps (like this is usually [needed for Arch Linux packages](https://wiki.archlinux.org/index.php/Creating_packages#PKGBUILD_functions)) or if you simply want to know what the `get` command does under-the-hood.
```
$ git clone https://github.com/rodcorsi/mattermosti18n src/github.com/rodcorsi/mattermosti18n
```
#### Build all executables (files with main) at once:
```
$ GOPATH=$(pwd) go install github.com/rodcorsi/mattermosti18n/...
```

#### Build executables manually:
```
$ GOPATH=$(pwd) go build -o i18n2po github.com/rodcorsi/mattermosti18n/i18n2po
$ GOPATH=$(pwd) go build -o po2i18n github.com/rodcorsi/mattermosti18n/po2i18n
```

We are closing the repository without the `.git` extension in the URL in order to [keep a working repository and not
a bare one](http://stackoverflow.com/a/11069413/3514658). Also, please note that if you are not familiar with the Go concepts yet, [please read this article, this is a must-read](https://golang.org/doc/code.html).

## Usage

### PO -> JSON

After you translate some phrases in Pootle server you can convert PO files to JSON to test in your Mattermost fork.

1. Download the latest version of `platform/i18n/en.json` and `web/static/i18n/en.json`:
   ```
   $ wget https://raw.githubusercontent.com/mattermost/platform/master/webapp/i18n/en.json -O web_static.json
   $ wget https://raw.githubusercontent.com/mattermost/platform/master/i18n/en.json -O platform.json
   ```

2. Download the PO's, change the **\<LOCALE\>** for the language code (eg. es, pt_BR, de, zh_CN, etc)
    ```
    $ wget "https://translate.mattermost.com/export/?path=/<LOCALE>/mattermost/web_static.po" -O web_static.po
    $ wget "https://translate.mattermost.com/export/?path=/<LOCALE>/mattermost/platform.po" -O platform.po
    ```

3. After build Mattermosti18n you can use **po2i18n** to convert the files
    ```
    $ po2i18n -t web_static.json -o new_web_static.json web_static.po
    $ po2i18n -t platform.json -o new_platform.json platform.po
    ```

4. Now you can move the new json to your fork, again change the **\<LOCALE\>** for the language code
    ```
    $ mv new_web_static.json <path_to_your_mattermost>platform/webapp/i18n/<LOCALE>.json
    $ mv new_platform.json <path_to_your_mattermost>platform/i18n/<LOCALE>.json
    ```

### JSON -> PO

If you have a translated JSON file and you can convert to PO and then upload in Pootle server.

1. Download the PO's, change the **\<CODE\>** for the language code (eg. es, pt_BR, de, zh_CN, etc)
    ```
    $ wget "http://translate.mattermost.com/export/?path=/<CODE>/mattermost/web_static.po"
    $ wget "http://translate.mattermost.com/export/?path=/<CODE>/mattermost/platform.po"
    ```

2. After build Mattermosti18n you can use **i18n2po** to convert the files
    ```
    $ i18n2po -o new_web_static.po -t web_static.po <your-web_static.json>
    $ i18n2po -o new_platform.po -t platform.po <your-platform.json>
    ```

3. Go to the [Pootle interface](https://translate.mattermost.com/) and sign-in:

   * Click on Mattermost -> (your language) -> webstatic.po
   * Upload your PO translations
   * Choose your file new_web_static.po
   * Repeat this process to _platform.po_
