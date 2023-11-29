# thesis generator

Generates a thesis from the template.

## Usage 

```sh
thesis-generator [OPTIONS] <SETTING_JSON>
OPTIONS
    -t, --template <TEMPLATE>    specify the template. available: latex, html, markdown, word.
                                 default: latex.
    -o, --output <OUTPUT>          specify the output directory or archive file (zip or tar.gz).
    -h, --help                     print this help message.
SETTING_JSON
    specify the setting file for the thesis template.
```

### Setting file

```json
{
    "degree":     "bachelor",        // or "master"
    "title":      "Thesis Title",
    "supervisor": "Supervisor Name",
    "author": {
        "name": "Author Name in Japanese",
        "email": "Author Email",
        "affiliation": {
            "university": "University Name",
            "department": "Department Name"
        }
    },
    "repository": {
        "owner": "tamada", // GitHub or GitLab username
        "type": "github"
    }
}
```
