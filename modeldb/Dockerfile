FROM mitdbg/modeldb-frontend AS build-env
ADD modeldb/frontend/app.js /modeldb/frontend/app.js
ADD modeldb/frontend/public /modeldb/frontend/public
ADD modeldb/frontend/routes /modeldb/frontend/routes
ADD modeldb/frontend/util /modeldb/frontend/util
ADD modeldb/frontend/views /modeldb/frontend/views
ADD modeldb/frontend/package.json /modeldb/frontend/package.json

FROM node:alpine
COPY --from=build-env /modeldb /modeldb
WORKDIR /modeldb/frontend
RUN npm install
EXPOSE 3000
ENTRYPOINT ["sh","/modeldb/dockerbuild/wait_for_backend.sh"]
