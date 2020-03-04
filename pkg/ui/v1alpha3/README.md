# Katib User Interface

This is the source code for the Katib UI. Current version of Katib UI is v1alpha3. On the official Kubeflow website [here](https://www.kubeflow.org/docs/components/hyperparameter-tuning/experiment/#running-the-experiment-from-the-katib-ui) you can find information how to use Katib UI.
We are using [React](https://reactjs.org/) framework to create frontend and Go as a backend.

## Folder structure

1. `Dockerfile` and file to serve the UI `main.go` you can find under [cmd/ui/v1alpha3](https://github.com/kubeflow/katib/tree/master/cmd/ui/v1alpha3).

2. Go backend you can find under [pkg/ui/v1alpha3](https://github.com/kubeflow/katib/tree/master/pkg/ui/v1alpha3).

3. React frontend you can find under [pkg/ui/v1alpha3/frontend](https://github.com/kubeflow/katib/tree/master/pkg/ui/v1alpha3/frontend).

## Requirements

To make changes to the UI you need to install:

- Tools, defined [here](https://github.com/kubeflow/katib/blob/master/docs/developer-guide.md#requirements).

- `Node` (10.13 or later) and `npm` (6.13 or later). You can find [here](https://nodejs.org/en/download/) how to download it.

## Development

While development you have different ways to run Katib UI.

### First time

1. Clone the repository.

2. Go to `/frontend` folder.

3. Run `npm install` to install all dependencies.

It will create `/frontend/node_modules` folder with all dependencies from `package.json`. If you want to add new package, edit `/frontend/package.json` file with new dependency.

### Start frontend server

If you want to edit only frontend without connection to the backend, you can start frontend server in your local environment. For it, run `npm run start` under `/frontend` folder. You can access the UI using this URL: `http://localhost:3000/`.

### Serve UI frontend and backend

You can serve Katib UI locally. To make it you need to follow these steps:

1. Run `npm run build` under `/frontend` folder. It will create `/frontend/build` directory with optimized production build.

2. Go to `cmd/ui/v1alpha3`

3. Run `main.go` file with appropriate flags. For example, if you clone Katib repository to `/home` folder, run this command:

```
go run main.go --build-dir=/home/katib/pkg/ui/v1alpha3/frontend/build --port=8080
```

After that, you can access the UI using this URL: `http://localhost:8080/katib/

## Production

To run Katib UI in Production, after all changes in frontend and backend, you need to create an image for the UI. Under `katib` repository run this: `docker build . -f cmd/ui/v1alpha3/Dockerfile -t <name of your image>` to build image. You can modify UI [deployment](https://github.com/kubeflow/katib/blob/master/manifests/v1alpha3/ui/deployment.yaml#L24) with your new image. After this, follow [these steps](https://www.kubeflow.org/docs/components/hyperparameter-tuning/hyperparameter/#accessing-the-katib-ui) to access Katib UI.

## Code style

To make frontend code consistent and easy to review we use [Prettier](https://prettier.io/). You can find Prettier config [here](https://github.com/kubeflow/katib/tree/master/pkg/ui/v1alpha3/frontend/.prettierrc.yaml).
Check [here](https://prettier.io/docs/en/install.html), how to install Prettier CLI to check and format your code.

### IDE integration

For VSCode you can install plugin: "Prettier - Code formatter" and it will pick Prettier config automatically.

You can edit [settings.json](https://code.visualstudio.com/docs/getstarted/settings#_settings-file-locations) file for VSCode to autoformat on save.

```json
  "settings": {
    "editor.formatOnSave": true
  }
```

For others IDE see [this](https://prettier.io/docs/en/editors.html).

### Check and format code

Before submitting PR check and format your code. To check your code run `npm run format:check` under `/frontend` folder. To format your code run `npm run format:write` under `/frontend` folder.
If all files formatted you can submit the PR.

If you don't want to format some code, [here](https://prettier.io/docs/en/ignore.html) is an instruction how to disable Prettier.
