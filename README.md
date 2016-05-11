# Mattermosti18n

This tool was developed to convert Json translated files from the project [Mattermost](https://github.com/mattermost/platform) to PO used in [Pootle](http://186.202.167.109/projects/) server.

# To build
```
$ go get github.com/rodrigocorsi2/mattermosti18n
```

# Usage

## Convert PO -> Json

After you translate some phrases in Pootle server you can convert PO files to Json to test in your Mattermost fork.

1 - Download the last version of platform/i18n/en.json and web/static/i18n/en.json
```
$ wget https://raw.githubusercontent.com/mattermost/platform/master/webapp/i18n/en.json -O web_static.json
$ wget https://raw.githubusercontent.com/mattermost/platform/master/i18n/en.json -O platform.json
```

2 - Download the PO's, change the **\<LOCALE\>** for the languange code (eg. es, pt_BR, de, zh_CN, etc)
```
$ wget "http://186.202.167.109/export/?path=/<LOCALE>/mattermost/web_static.po"
$ wget "http://186.202.167.109/export/?path=/<LOCALE>/mattermost/platform.po"
```

3 - After build Mattermosti18n you can use **po2i18n** to convert the files
```
$ po2i18n -t web_static.json -o new_web_static.json web_static.po
$ po2i18n -t platform.json -o new_platform.json platform.po
```

4 - Now you can move the new json to your fork, again change the **\<LOCALE\>** for the languange code
```
$ mv new_web_static.json <path_to_your_mattermost>platform/webapp/i18n/<LOCALE>.json
$ mv new_platform.json <path_to_your_mattermost>platform/i18n/<LOCALE>.json
```


## Convert Json -> PO

If you have a translated json file and you can convert to PO and then upload in Pootle server.

1 - Download the PO's, change the **\<CODE\>** for the languange code (eg. es, pt_BR, de, zh_CN, etc)
```
$ wget "http://186.202.167.109/export/?path=/<CODE>/mattermost/web_static.po"
$ wget "http://186.202.167.109/export/?path=/<CODE>/mattermost/platform.po"
```

2 - After build Mattermosti18n you can use **i18n2po** to convert the files
```
$ i18n2po -o new_web_static.po -t web_static.po <your-web_static.json>
$ i18n2po -o new_platform.po -t platform.po <your-platform.json>
```

3 - Go to the [Pootle](http://186.202.167.109/) server and sign-in

* Click in Mattermost -> (your language) -> webstatic.po
* Upload translations
* Choose your file new_web_static.po
* Repeat this process to _platform.po_
