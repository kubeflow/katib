# Katib User Interface

This is the source code for the Katib UI. Current version of Katib UI is v1beta1. On the official Kubeflow website [here](https://www.kubeflow.org/docs/components/katib/experiment/#running-the-experiment-from-the-katib-ui) you can find information how to use Katib UI.
We are using [Angular](https://angular.io/) framework to create frontend and Go as a backend.

We are using [Material UI](https://material.angular.io/) to design frontend. Try to use Material UI components to implement new Katib UI features.

## Folder structure

1. You can find `Dockerfile` and `main.go` - file to serve the UI under [`cmd/ui/v1beta1`](../../../cmd/ui/v1beta1/)

1. You can find Go backend under [`pkg/ui/v1beta1`](../../../pkg/ui/v1beta1)

1. You can find Angular frontend under [`pkg/ui/v1beta1/frontend`](../../../pkg/ui/v1beta1/frontend/)

## Requirements

To make changes to the UI you need to install:

- Tools, defined [here](https://github.com/kubeflow/katib/blob/master/docs/developer-guide.md#requirements).

- `node` (v12.18.1) and `npm` (v6.13). Recommended to install `node` and `npm` using [`nvm`](https://github.com/nvm-sh/nvm). After installing `nvm`, you can run `nvm install 12.18.1` to install `node` version 12.18.1 and run `nvm use 12.18.1` to use that version.

## Development

While development you have different ways to run Katib UI.

1. Build and serve only the frontend. The dev server will also be proxying requests to the backend
2. Build the frontend and serve it via the backend locally

### Serve only the frontend

You can run a webpack dev server that only exposes the frontend files, which can be useful for testing only the UI of the app. There's also a `proxy.conf.json` file which configures the dev server to send the backend requests to port `8000`.

In order to build the UI locally, and expose it with a webpack dev server you will need to:

1. Create a module from the [common library](https://github.com/kubeflow/kubeflow/tree/master/components/crud-web-apps/common/frontend/kubeflow-common-lib)
2. Install the node modules of the app and also link the common-library module

You can build the common library with:

```bash
COMMIT=$(cat ./frontend/COMMIT) \
  && cd /tmp && git clone https://github.com/kubeflow/kubeflow.git \
  && cd kubeflow \
  && git checkout $COMMIT \
  && cd components/crud-web-apps/common/frontend/kubeflow-common-lib

# build the common library module
npm i
npm run build

# link the module to your npm packages
# depending on where you npm stores the global packages you
# might need to use sudo
npm link dist/kubeflow
```

And then build and run the UI locally, on `localhost:4200`, with:

```bash
# If you've already cloned the repo then skip this step and just
# navigate to the pkg/ui/v1beta1/frontend dir
cd /tmp && git clone https://github.com/kubeflow/katib.git \
  && cd katib/pkg/ui/v1beta1/frontend

npm i
npm link kubeflow
npm run start
```

### Serve the UI from the backend

This is the recommended way to test the web app e2e. In order to build the UI and serve it via the backend, locally, you will need to:

1. Build the UI locally. You have to follow the steps from the previous section, but instead of running `npm run start` you need to run `npm run build:prod`. It builds the frontend artifacts under `frontend/dist/static` folder.

   Moreover, you are able to run `npm run build:watch` instead of `npm run build:prod`. In that case, it starts a process which is watching the source code changes and building the frontend artifacts under `frontend/dist/static` folder.

   Learn more about Angular scripts in the [official guide](https://angular.io/cli/build).

1. Run `kubectl port-forward svc/katib-db-manager 6789 -n kubeflow` to expose `katib-db-manager` service for external access. You can use [different ways](https://kubernetes.io/docs/tasks/access-application-cluster/) to get external address for Kubernetes service. After exposing service, you should be able to receive information by running `wget <external-ip>:<service-port>`. In case of port-forwarding above, you have to run `wget localhost:6789`.

1. Go to [`cmd/ui/v1beta1`](../../../cmd/ui/v1beta1/)

1. Run `main.go` file with appropriate flags, where:

   - `--build-dir` - directory with the frontend artifacts.
   - `--port` - port to access Katib UI.
   - `--db-manager-address` - Katib DB manager external IP and port address.

   For example, if you use port-forwarding to expose `katib-db-manager`, run this command:

   ```
   export APP_DISABLE_AUTH=true
   go run main.go --build-dir=../../../pkg/ui/v1beta1/frontend/dist --port=8080 --db-manager-address=localhost:6789
   ```

After that, you can access the UI using this URL: `http://localhost:8080/katib/`.

## Production

To run Katib UI in Production, after all changes in frontend and backend, you need to create an image for the UI. Under `/katib` directory run this: `docker build . -f cmd/ui/v1beta1/Dockerfile -t <name of your image>` to build the image.

After that, you can modify the [UI Deployment](../../../manifests/v1beta1/components/ui/ui.yaml) with
your new image. Then, follow
[these steps](https://www.kubeflow.org/docs/components/katib/hyperparameter/#accessing-the-katib-ui) to access Katib UI.

## Code style

To make frontend code consistent and easy to review we use [Prettier](https://prettier.io/).
You can find Prettier config [here](../../../pkg/ui/v1beta1/frontend/.prettierrc.yaml).
Check [here](https://prettier.io/docs/en/install.html), how to install Prettier CLI to check and format your code.

### IDE integration

For VSCode you can install plugin: "Prettier - Code formatter" and it picks Prettier config automatically.

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
