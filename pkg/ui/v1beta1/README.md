# Katib User Interface

This is the source code for the Katib UI. Current version of Katib UI is v1beta1. On the official Kubeflow website [here](https://www.kubeflow.org/docs/components/katib/experiment/#running-the-experiment-from-the-katib-ui) you can find information how to use Katib UI.
We are using [React](https://reactjs.org/) framework to create frontend and Go as a backend.

We are using [Material UI](https://material-ui.com/) to design frontend. Try to use Material UI components to implement new Katib UI features.

## Folder structure

1. You can find `Dockerfile` and `main.go` - file to serve the UI under [cmd/ui/v1beta1](https://github.com/kubeflow/katib/tree/master/cmd/ui/v1beta1).

1. You can find Go backend under [pkg/ui/v1beta1](https://github.com/kubeflow/katib/tree/master/pkg/ui/v1beta1).

1. You can find React frontend under [pkg/ui/v1beta1/frontend](https://github.com/kubeflow/katib/tree/master/pkg/ui/v1beta1/frontend).

## Requirements

To make changes to the UI you need to install:

- Tools, defined [here](https://github.com/kubeflow/katib/blob/master/docs/developer-guide.md#requirements).

- `node` (v12 or later) and `npm` (v6.13 or later). Recommended to install `node` and `npm` using [`nvm`](https://github.com/nvm-sh/nvm). After installing `nvm`, you can run `nvm install 12.18.1` to install `node` version 12.18.1 and run `nvm use 12.18.1` to use that version.

## Development

While development you have different ways to run Katib UI.

### First time

1. Clone the repository.

1. Go to `/frontend` folder.

1. Run `npm install` to install all dependencies.

It creates `/frontend/node_modules` folder with all dependencies from [`package.json`](https://github.com/kubeflow/katib/blob/master/pkg/ui/v1beta1/frontend/package.json). If you want to add new package, run `npm install <package>@<version>`. That should update `/frontend/package.json` and `/frontend/package-lock.json` with the new dependency.

### Start frontend server

If you want to edit only frontend without connection to the backend, you can start frontend server in your local environment. For it, run `npm run start` under `/frontend` folder. You can access the UI using this URL: `http://localhost:3000/`.

### Serve UI frontend and backend

You can serve Katib UI locally. To make it you need to follow these steps:

1. Run `npm run build` under `/frontend` folder. It creates `/frontend/build` directory with optimized production build.

   If your `node` memory limit is not enough to build the frontend, you may see this error while building: `FATAL ERROR: Ineffective mark-compacts near heap limit Allocation failed - JavaScript heap out of memory`. To fix it, you can try to increase `node` memory limit. For that, change [`build`](https://github.com/kubeflow/katib/blob/master/pkg/ui/v1beta1/frontend/package.json#L28) script to `react-scripts --max_old_space_size=4096 build` to increase `node` memory up to 4 Gb.

1. Run `kubectl port-forward svc/katib-db-manager 6789 -n kubeflow` to expose `katib-db-manager` service for external access. You can use [different ways](https://kubernetes.io/docs/tasks/access-application-cluster/) to get external address for Kubernetes service. After exposing service, you should be able to receive information by running `wget <external-ip>:<service-port>`. In case of port-forwarding above, you have to run `wget localhost:6789`.

1. Go to `cmd/ui/v1beta1`.

1. Run `main.go` file with appropriate flags, where:

   - `--build-dir` - builded frontend directory.
   - `--port` - port to access Katib UI.
   - `--db-manager-address` - Katib DB manager external IP and port address.

   For example, if you clone Katib repository to `/home` folder and use port-forwarding to expose `katib-db-manager`, run this command:

   ```
   go run main.go --build-dir=/home/katib/pkg/ui/v1beta1/frontend/build --port=8080 --db-manager-address=localhost:6789
   ```

After that, you can access the UI using this URL: `http://localhost:8080/katib/`.

## Production

To run Katib UI in Production, after all changes in frontend and backend, you need to create an image for the UI. Under `/katib` directory run this: `docker build . -f cmd/ui/v1beta1/Dockerfile -t <name of your image>` to build the image. If Docker resources are not enough to build the frontend, you get `node` out of memory error. You can try to increase Docker resources or modify `package.json` as detailed [above](https://github.com/kubeflow/katib/tree/master/pkg/ui/v1beta1#serve-ui-frontend-and-backend) at step 1.

After that, you can modify UI [deployment](https://github.com/kubeflow/katib/blob/master/manifests/v1beta1/components/ui/ui.yaml#L21) with your new image. Then, follow [these steps](https://www.kubeflow.org/docs/components/katib/hyperparameter/#accessing-the-katib-ui) to access Katib UI.

## Code style

To make frontend code consistent and easy to review we use [Prettier](https://prettier.io/). You can find Prettier config [here](https://github.com/kubeflow/katib/tree/master/pkg/ui/v1beta1/frontend/.prettierrc.yaml).
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
