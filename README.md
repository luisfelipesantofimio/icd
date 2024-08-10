# ICD data fetch tool

Clone any chapter from [ICD-11 API](https://icd.who.int/icdapi).

Created for providing required ICD data for [Aronnax](https://github.com/luisfelipesantofimio/aronnax).

## Usage

Go to [ICD-11 API](https://icd.who.int/icdapi), register and get your `clientId` and `clientSecret`.

Create a file called `config.json` with the following content:

```json
{
  "clientId": "YOUR_CLIENT_ID",
  "clientSecret": "YOUR_CLIENT_SECRET"
}
```

Then run `go run main.go` and find the result json file for `es` and `en` in the `output` folder.
