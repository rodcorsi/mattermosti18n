# Mattermosti18n

This tool was developed to convert Json translated files from the project [Mattermost](https://github.com/mattermost/platform) to PO used in [Pootle](http://186.202.167.109/projects/) server.

# To build
```
git clone https://github.com/rodrigocorsi2/mattermosti18n.git
cd mattermosti18n
make
```

# Usage

## Convert PO -> Json

After you translate some phrases in Pootle server you can convert PO files to Json to test in your Mattermost fork.

1 - Download the last version of platform/i18n/en.json and web/static/i18n/en.json
```
wget https://raw.githubusercontent.com/mattermost/platform/master/web/static/i18n/en.json -O web_static.json
wget https://raw.githubusercontent.com/mattermost/platform/master/i18n/en.json -O platform.json
````

2 - Download the PO's, change the **\<CODE\>** for the languange code (eg. es, pt_BR, de, zh_CN, etc)
```
wget "http://186.202.167.109/export/?path=/<CODE>/mattermost/web_static.po"
wget "http://186.202.167.109/export/?path=/<CODE>/mattermost/platform.po"
```

3 - After build Mattermosti18n you can use _po2i18n_ to convert the files
```
mattermosti18n/bin/po2i18n -t web_static.json -o new_web_static.json web_static.po
mattermosti18n/bin/po2i18n -t platform.json -o new_platform.json platform.po
```

4 - Now you can move the new json to your fork, again change the **\<CODE\>** for the languange code
```
mv new_web_static.json <path_to_your_mattermost>platform/web/static/i18n/<CODE>.json
mv new_platform.json <path_to_your_mattermost>platform/i18n/<CODE>.json
```


## Convert Json -> PO

If you have a translated json file and you can convert to PO and then upload in Pootle server.

1 - Download the last version of platform/i18n/en.json and web/static/i18n/en.json
```
wget https://raw.githubusercontent.com/mattermost/platform/master/web/static/i18n/en.json -O web_static.json
wget https://raw.githubusercontent.com/mattermost/platform/master/i18n/en.json -O platform.json
````

2 - After build Mattermosti18n you can use _i18n2po_ to convert the files
```
mattermosti18n/bin/i18n2po -o new_web_static.po web_static.json <your-web_static.json>
mattermosti18n/bin/i18n2po -o new_platform.po platform.json <your-platform.json>
```

3 - Go to the [Pootle](http://186.202.167.109/) server and sign-in

* Click in Mattermost -> (your language) -> webstatic.po
* Upload translations
* Choose your file new_web_static.po
* Repeat this process to _platform.po_
