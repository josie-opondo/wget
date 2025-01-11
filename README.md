WGET:

REQUIREMENTS:

- start time
- response status
- content length
- content size
- saving path
- show the amount of data being downloaded
- show the percentage of data being downloaded
- show the remaining time
- read from a text file

FLAGS

- `-O` : output file
- `--rate-limit`
- `-B`: background downloads
  `start at <date that the download started>
sending request, awaiting response... status 200 OK
content size: <56370 [~0.06MB]>
saving file to: ./<name-of-the-file-downloaded>
Downloaded [<link-downloaded>]
finished at <date that the download finished>
`

CONS

- does not show estimated time
- rate limit flag not effeciently working. adjust to be dynamic
- not able to read from a file
- check background download, takes longer

MIRROR
Try to run the following command `./wget --mirror --convert-links http://corndog.io/`, then try to open the index.html with a browser
Is the site working?
Try to run the following command `./wget --mirror https://oct82.com/`, then try to open the index.html with a browser
Is the site working?
Try to run the following command `./wget --mirror --reject=gif https://oct82.com/`, then try to open the index.html with a browser
Did the program download the site without the GIFs?
Try to run the following command `./wget --mirror https://trypap.com/`, then use the command ls to see the file system of the created folder.

css img index.html

Try to run the following command `./wget --mirror https://trypap.com/`, then use the command ls to see the file system of the created folder.

css img index.html

Does the created folder has the same fs as above?
Try to run the following command `./wget --mirror -X=/img https://trypap.com/`, then use the command ls to see the file system of the created folder.

css index.html

Does the created folder has the files above?
Try to run the following command `./wget --mirror https://theuselessweb.com/`
Is the site working?
Try to run the following command to mirror a website at your choice `./wget --mirror <https://link_of_your_choice.com>`
Did the program mirror the website?
