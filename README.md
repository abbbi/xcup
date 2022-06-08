# DHV XC Uploader

Publish your paraglider flights on DHV-XC via command line. Client
implementation for the new DHV-XC (https://www.dhv-xc.de) API as described
here:

https://www.dhv.de/fileadmin/user_upload/files/2022/05/DHV_XC_Flight_Upload_Interface_Specification.pdf

Before using, you must configure an upload password in your DHV-XC User
Settings.

# Example:

```
./xc -u username -p upload_password -f 16Habla1.igc -g "Niviuk Ikuma 2" -P
INFO[0000] Login OK
INFO[0000] Glider: [Niviuk Ikuma 2]
INFO[0000] Publishing flight during upload.
[..]
```

# Usage:

```
Usage:
  xc [OPTIONS]

Application Options:
  -u, --user=    DHV-XC User name
  -p, --pass=    DHV-XC Upload Password
  -f, --file=    IGC file
  -P, --publish  Publish flight after upload
  -g, --glider=  Glider name

Help Options:
  -h, --help     Show this help message
```
