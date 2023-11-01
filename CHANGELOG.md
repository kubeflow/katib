# Changelog

# [v0.16.0](https://github.com/kubeflow/katib/tree/v0.16.0) (2023-10-31)

## Breaking Changes

- Implement KatibConfig API ([#2176](https://github.com/kubeflow/katib/pull/2176) by [@tenzen-y](https://github.com/tenzen-y))
- Drop Kubernetes v1.24 and support Kubernetes v1.27 ([#2182](https://github.com/kubeflow/katib/pull/2182) by [@tenzen-y](https://github.com/tenzen-y))
- Drop Kubernetes v1.23 and support Kubernetes v1.26 ([#2177](https://github.com/kubeflow/katib/pull/2177) by [@tenzen-y](https://github.com/tenzen-y))
- Change failurePolicy to Fail for Katib Webhooks ([#2018](https://github.com/kubeflow/katib/pull/2018) by [@andreyvelich](https://github.com/andreyvelich))

## New Features

### Core Features

- Consolidate the Katib Cert Generator to the Katib Controller ([#2185](https://github.com/kubeflow/katib/pull/2185) by [@tenzen-y](https://github.com/tenzen-y))
- Containerize tests for Katib Conformance ([#2146](https://github.com/kubeflow/katib/pull/2146) by [@nagar-ajay](https://github.com/nagar-ajay))

### UI Improvements

- [UI] Default Resume Policy to never from UI ([#2195](https://github.com/kubeflow/katib/pull/2195) by [@mChowdhury-91](https://github.com/mChowdhury-91))
- [UI] Remove Deprecated Katib UI ([#2179](https://github.com/kubeflow/katib/pull/2179) by [@andreyvelich](https://github.com/andreyvelich))
- [UI] Fix Trial Logs when Kubernetes Job Fails ([#2164](https://github.com/kubeflow/katib/pull/2164) by [@andreyvelich](https://github.com/andreyvelich))
- kwa(front): Support all namespaces ([#2119](https://github.com/kubeflow/katib/pull/2119) by [@elenzio9](https://github.com/elenzio9))
- kwa(front): Update the use of SnackBarService ([#2113](https://github.com/kubeflow/katib/pull/2113) by [@orfeas-k](https://github.com/orfeas-k))
- UI: Remove an unsed import, EventV1beta1Api ([#2116](https://github.com/kubeflow/katib/pull/2116) by [@tenzen-y](https://github.com/tenzen-y))

### SDK Improvements

- [SDK] Enable resource specification for trial containers ([#2192](https://github.com/kubeflow/katib/pull/2192) by [@droctothorpe](https://github.com/droctothorpe))
- [SDK] Add namespace parameter to KatibClient ([#2183](https://github.com/kubeflow/katib/pull/2183) by [@droctothorpe](https://github.com/droctothorpe))
- [SDK] Import all Kubernetes Models ([#2148](https://github.com/kubeflow/katib/pull/2148) by [@andreyvelich](https://github.com/andreyvelich))

## Bug fixes

- Bug: Wait for the certs to be mounted inside the container ([#2213](https://github.com/kubeflow/katib/pull/2213) by [@tenzen-y](https://github.com/tenzen-y))
- Start waiting for certs to be ready before sending data to the channel ([#2215](https://github.com/kubeflow/katib/pull/2215) by [@tenzen-y](https://github.com/tenzen-y))
- E2E: Add additional checks to verify if the components are ready ([#2212](https://github.com/kubeflow/katib/pull/2212) by [@tenzen-y](https://github.com/tenzen-y))
- Remove a katib-webhook-cert Secret from components ([#2214](https://github.com/kubeflow/katib/pull/2214) by [@tenzen-y](https://github.com/tenzen-y))
- Skip to inject the metrics-collector pods to the Katib controller ([#2211](https://github.com/kubeflow/katib/pull/2211) by [@tenzen-y](https://github.com/tenzen-y))
- Sending an empty data to the certsReady channel ([#2196](https://github.com/kubeflow/katib/pull/2196) by [@tenzen-y](https://github.com/tenzen-y))
- Fix conformance docker image ([#2147](https://github.com/kubeflow/katib/pull/2147) by [@nagar-ajay](https://github.com/nagar-ajay))

## Documentation

- Add PITS Global Data Recovery Services to the adopters list ([#2160](https://github.com/kubeflow/katib/pull/2160) by [@ghost](https://github.com/ghost))
- Add SDK Breaking Change to Changelog ([#2133](https://github.com/kubeflow/katib/pull/2133) by [@andreyvelich](https://github.com/andreyvelich))
- Add Changelog for Katib v0.15.0 ([#2129](https://github.com/kubeflow/katib/pull/2129) by [@andreyvelich](https://github.com/andreyvelich))
- Add Changelog for Katib v0.15.0-rc.1 ([#2123](https://github.com/kubeflow/katib/pull/2123) by [@andreyvelich](https://github.com/andreyvelich))
- Add Changelog for Katib v0.15.0-rc.0 ([#2106](https://github.com/kubeflow/katib/pull/2106) by [@andreyvelich](https://github.com/andreyvelich))

## Misc

- Upgrade Tensorflow version to v2.13.0 ([#2216](https://github.com/kubeflow/katib/pull/2216) by [@tenzen-y](https://github.com/tenzen-y))
- Upgrade Go version to v1.20 ([#2190](https://github.com/kubeflow/katib/pull/2190) by [@tenzen-y](https://github.com/tenzen-y))
- Replace grpc_health_probe with the built-in gRPC container probe feature ([#2189](https://github.com/kubeflow/katib/pull/2189) by [@tenzen-y](https://github.com/tenzen-y))
- Allow install binaries for the arm64 in the envtest ([#2188](https://github.com/kubeflow/katib/pull/2188) by [@tenzen-y](https://github.com/tenzen-y))
- Replace action to setup minikube with medyagh/setup-minikube ([#2178](https://github.com/kubeflow/katib/pull/2178) by [@tenzen-y](https://github.com/tenzen-y))
- Remove Charmed Operators for Katib ([#2161](https://github.com/kubeflow/katib/pull/2161) by [@ca-scribner](https://github.com/ca-scribner))
- Namespace and trial pod annotations as CLI argument ([#2138](https://github.com/kubeflow/katib/pull/2138) by [@nagar-ajay](https://github.com/nagar-ajay))
- Relax dependencies restriction for the gRPC libraries ([#2140](https://github.com/kubeflow/katib/pull/2140) by [@tenzen-y](https://github.com/tenzen-y))
- Add SDK Breaking Change to Changelog ([#2133](https://github.com/kubeflow/katib/pull/2133) by [@andreyvelich](https://github.com/andreyvelich))
- Increase the free spaces in CI ([#2131](https://github.com/kubeflow/katib/pull/2131) by [@tenzen-y](https://github.com/tenzen-y))
- Reformat katib-operators ([#2114](https://github.com/kubeflow/katib/pull/2114) by [@tenzen-y](https://github.com/tenzen-y))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.15.0...v0.16.0)

# [v0.16.0-rc.1](https://github.com/kubeflow/katib/tree/v0.16.0-rc.1) (2023-08-16)

## New Features

- Upgrade Tensorflow version to v2.13.0 ([#2216](https://github.com/kubeflow/katib/pull/2216) by [@tenzen-y](https://github.com/tenzen-y))

## Bug Fixes

- Bug: Wait for the certs to be mounted inside the container ([#2213](https://github.com/kubeflow/katib/pull/2213) by [@tenzen-y](https://github.com/tenzen-y))
- Start waiting for certs to be ready before sending data to the channel ([#2215](https://github.com/kubeflow/katib/pull/2215) by [@tenzen-y](https://github.com/tenzen-y))
- E2E: Add additional checks to verify if the components are ready ([#2212](https://github.com/kubeflow/katib/pull/2212) by [@tenzen-y](https://github.com/tenzen-y))
- Remove a katib-webhook-cert Secret from components ([#2214](https://github.com/kubeflow/katib/pull/2214) by [@tenzen-y](https://github.com/tenzen-y))
- Skip to inject the metrics-collector pods to the Katib controller ([#2211](https://github.com/kubeflow/katib/pull/2211) by [@tenzen-y](https://github.com/tenzen-y))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.16.0-rc.0...v0.16.0-rc.1)

# [v0.16.0-rc.0](https://github.com/kubeflow/katib/tree/v0.16.0-rc.0) (2023-08-05)

## Breaking Changes

- Implement KatibConfig API ([#2176](https://github.com/kubeflow/katib/pull/2176) by [@tenzen-y](https://github.com/tenzen-y))
- Drop Kubernetes v1.24 and support Kubernetes v1.27 ([#2182](https://github.com/kubeflow/katib/pull/2182) by [@tenzen-y](https://github.com/tenzen-y))
- Drop Kubernetes v1.23 and support Kubernetes v1.26 ([#2177](https://github.com/kubeflow/katib/pull/2177) by [@tenzen-y](https://github.com/tenzen-y))
- Change failurePolicy to Fail for Katib Webhooks ([#2018](https://github.com/kubeflow/katib/pull/2018) by [@andreyvelich](https://github.com/andreyvelich))

## New Features

### Core Features

- Consolidate the Katib Cert Generator to the Katib Controller ([#2185](https://github.com/kubeflow/katib/pull/2185) by [@tenzen-y](https://github.com/tenzen-y))
- Containerize tests for Katib Conformance ([#2146](https://github.com/kubeflow/katib/pull/2146) by [@nagar-ajay](https://github.com/nagar-ajay))

### UI Improvements

- [UI] Default Resume Policy to never from UI ([#2195](https://github.com/kubeflow/katib/pull/2195) by [@mChowdhury-91](https://github.com/mChowdhury-91))
- [UI] Remove Deprecated Katib UI ([#2179](https://github.com/kubeflow/katib/pull/2179) by [@andreyvelich](https://github.com/andreyvelich))
- [UI] Fix Trial Logs when Kubernetes Job Fails ([#2164](https://github.com/kubeflow/katib/pull/2164) by [@andreyvelich](https://github.com/andreyvelich))
- kwa(front): Support all namespaces ([#2119](https://github.com/kubeflow/katib/pull/2119) by [@elenzio9](https://github.com/elenzio9))
- kwa(front): Update the use of SnackBarService ([#2113](https://github.com/kubeflow/katib/pull/2113) by [@orfeas-k](https://github.com/orfeas-k))
- UI: Remove an unsed import, EventV1beta1Api ([#2116](https://github.com/kubeflow/katib/pull/2116) by [@tenzen-y](https://github.com/tenzen-y))

### SDK Improvements

- [SDK] Enable resource specification for trial containers ([#2192](https://github.com/kubeflow/katib/pull/2192) by [@droctothorpe](https://github.com/droctothorpe))
- [SDK] Add namespace parameter to KatibClient ([#2183](https://github.com/kubeflow/katib/pull/2183) by [@droctothorpe](https://github.com/droctothorpe))
- [SDK] Import all Kubernetes Models ([#2148](https://github.com/kubeflow/katib/pull/2148) by [@andreyvelich](https://github.com/andreyvelich))

## Bug fixes

- Sending an empty data to the certsReady channel ([#2196](https://github.com/kubeflow/katib/pull/2196) by [@tenzen-y](https://github.com/tenzen-y))
- Fix conformance docker image ([#2147](https://github.com/kubeflow/katib/pull/2147) by [@nagar-ajay](https://github.com/nagar-ajay))

## Documentation

- Add PITS Global Data Recovery Services to the adopters list ([#2160](https://github.com/kubeflow/katib/pull/2160) by [@ghost](https://github.com/ghost))
- Add SDK Breaking Change to Changelog ([#2133](https://github.com/kubeflow/katib/pull/2133) by [@andreyvelich](https://github.com/andreyvelich))
- Add Changelog for Katib v0.15.0 ([#2129](https://github.com/kubeflow/katib/pull/2129) by [@andreyvelich](https://github.com/andreyvelich))
- Add Changelog for Katib v0.15.0-rc.1 ([#2123](https://github.com/kubeflow/katib/pull/2123) by [@andreyvelich](https://github.com/andreyvelich))
- Add Changelog for Katib v0.15.0-rc.0 ([#2106](https://github.com/kubeflow/katib/pull/2106) by [@andreyvelich](https://github.com/andreyvelich))

## Misc

- Upgrade Go version to v1.20 ([#2190](https://github.com/kubeflow/katib/pull/2190) by [@tenzen-y](https://github.com/tenzen-y))
- Replace grpc_health_probe with the built-in gRPC container probe feature ([#2189](https://github.com/kubeflow/katib/pull/2189) by [@tenzen-y](https://github.com/tenzen-y))
- Allow install binaries for the arm64 in the envtest ([#2188](https://github.com/kubeflow/katib/pull/2188) by [@tenzen-y](https://github.com/tenzen-y))
- Replace action to setup minikube with medyagh/setup-minikube ([#2178](https://github.com/kubeflow/katib/pull/2178) by [@tenzen-y](https://github.com/tenzen-y))
- Remove Charmed Operators for Katib ([#2161](https://github.com/kubeflow/katib/pull/2161) by [@ca-scribner](https://github.com/ca-scribner))
- Namespace and trial pod annotations as CLI argument ([#2138](https://github.com/kubeflow/katib/pull/2138) by [@nagar-ajay](https://github.com/nagar-ajay))
- Relax dependencies restriction for the gRPC libraries ([#2140](https://github.com/kubeflow/katib/pull/2140) by [@tenzen-y](https://github.com/tenzen-y))
- Add SDK Breaking Change to Changelog ([#2133](https://github.com/kubeflow/katib/pull/2133) by [@andreyvelich](https://github.com/andreyvelich))
- Increase the free spaces in CI ([#2131](https://github.com/kubeflow/katib/pull/2131) by [@tenzen-y](https://github.com/tenzen-y))
- Reformat katib-operators ([#2114](https://github.com/kubeflow/katib/pull/2114) by [@tenzen-y](https://github.com/tenzen-y))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.15.0...v0.16.0-rc.0)

# [v0.15.0](https://github.com/kubeflow/katib/tree/v0.15.0) (2023-03-22)

## Breaking Changes

- Use **Never** Resume Policy as Default ([#2102](https://github.com/kubeflow/katib/pull/2102) by [@andreyvelich](https://github.com/andreyvelich))
- Chocolate Suggestion Service is removed ([#2071](https://github.com/kubeflow/katib/pull/2071) by [@tenzen-y](https://github.com/tenzen-y))
- `request_number` is removed from the GRPC APIs ([#1994](https://github.com/kubeflow/katib/pull/1994) by [@johnugeorge](https://github.com/johnugeorge))
- Enabling Authorization in Katib UI ([#1983](https://github.com/kubeflow/katib/pull/1983) and [#2041](https://github.com/kubeflow/katib/pull/2041) by [@apo-ger](https://github.com/apo-ger))
- The new improved and refactored Katib SDK is not backward compatible ([#2075](https://github.com/kubeflow/katib/pull/2075) by [@andreyvelich](https://github.com/andreyvelich))

## New Features

### Major Features

- Narrow down Katib RBAC rules ([#2091](https://github.com/kubeflow/katib/pull/2091) by [@johnugeorge](https://github.com/johnugeorge))
- Support Postgres as a Katib DB ([#1921](https://github.com/kubeflow/katib/pull/1921) by [@anencore94](https://github.com/anencore94))
- More Suggestion container fields in Katib Config ([#2000](https://github.com/kubeflow/katib/pull/2000) by [@fischor](https://github.com/fischor))
- Katib UI: Create the LOGS tab of Trial's details page ([#2117](https://github.com/kubeflow/katib/pull/2117) by [@elenzio9](https://github.com/elenzio9))
- Katib UI: Enable pagination/sorting/filtering ([#2017](https://github.com/kubeflow/katib/pull/2017) and [#2040](https://github.com/kubeflow/katib/pull/2040) by [@elenzio9](https://github.com/elenzio9))
- [SDK] Create Tune API in the Katib SDK ([#1951](https://github.com/kubeflow/katib/pull/1951) by [@andreyvelich](https://github.com/andreyvelich))
- [SDK] Get Trial Metrics from Katib DB ([#2050](https://github.com/kubeflow/katib/pull/2050) by [@andreyvelich](https://github.com/andreyvelich))

### Core Features

- Add Conformance Program Doc for AutoML and Training WG ([#2048](https://github.com/kubeflow/katib/pull/2048) by [@andreyvelich](https://github.com/andreyvelich))
- Support for grid search algorithm in Optuna Suggestion Service ([#2060](https://github.com/kubeflow/katib/pull/2060) by [@tenzen-y](https://github.com/tenzen-y))
- Add Trial Labels During Pod Mutation ([#2047](https://github.com/kubeflow/katib/pull/2047) by [@andreyvelich](https://github.com/andreyvelich))
- Support for k8s v1.25 in CI ([#1997](https://github.com/kubeflow/katib/pull/1997) by [@johnugeorge](https://github.com/johnugeorge))
- Add the CI to build multi-platform container images ([#1956](https://github.com/kubeflow/katib/pull/1956) by [@tenzen-y](https://github.com/tenzen-y))
- Drop Kubernetes v1.21 and introduce Kubernetes v1.24 ([#1953](https://github.com/kubeflow/katib/pull/1953) by [@tenzen-y](https://github.com/tenzen-y))
- Add --connect-timeout flag to katib-db-manager ([#1937](https://github.com/kubeflow/katib/pull/1937) by [@tenzen-y](https://github.com/tenzen-y))
- Implement validations for DARTS suggestion service ([#1926](https://github.com/kubeflow/katib/pull/1926) by [@tenzen-y](https://github.com/tenzen-y))
- Implement validation for Optuna suggestion service ([#1924](https://github.com/kubeflow/katib/pull/1924) by [@tenzen-y](https://github.com/tenzen-y))

### UI Improvements

- Make links in KWA's tables actual links ([#2090](https://github.com/kubeflow/katib/pull/2090) by [@elenzio9](https://github.com/elenzio9))
- frontend: Rework the trial graph using ECharts in KWA ([#2089](https://github.com/kubeflow/katib/pull/2089) by [@elenzio9](https://github.com/elenzio9))
- kwa(front): Add UI tests with Cypress ([#2088](https://github.com/kubeflow/katib/pull/2088) by [@orfeas-k](https://github.com/orfeas-k))
- frontend: Enable actions in experiment graph ([#2065](https://github.com/kubeflow/katib/pull/2065) by [@elenzio9](https://github.com/elenzio9))
- frontend: Show message in case of uncompleted trial instead of the graph ([#2063](https://github.com/kubeflow/katib/pull/2063) by [@elenzio9](https://github.com/elenzio9))
- frontend: Add source maps in the browser ([#2043](https://github.com/kubeflow/katib/pull/2043) by [@elenzio9](https://github.com/elenzio9))
- Backend for getting logs of a trial ([#2039](https://github.com/kubeflow/katib/pull/2039) by [@d-gol](https://github.com/d-gol))
- frontend: Show the successful trials in the experiment graph (#2013) ([#2033](https://github.com/kubeflow/katib/pull/2033) by [@elenzio9](https://github.com/elenzio9))
- frontend: Migrate from tslint to eslint in KWA ([#2042](https://github.com/kubeflow/katib/pull/2042) by [@elenzio9](https://github.com/elenzio9))
- Dedicated yaml tab for Trials ([#2034](https://github.com/kubeflow/katib/pull/2034) by [@elenzio9](https://github.com/elenzio9))
- KWA: Use new Editor component (Monaco) ([#2023](https://github.com/kubeflow/katib/pull/2023) by [@orfeas-k](https://github.com/orfeas-k))
- kwa(build): Introduce COMMIT file for building KWA ([#2014](https://github.com/kubeflow/katib/pull/2014) by [@orfeas-k](https://github.com/orfeas-k))
- frontend: Fix 500 error after detail page refresh (#1967) ([#2001](https://github.com/kubeflow/katib/pull/2001) by [@elenzio9](https://github.com/elenzio9))
- Introduce KWA's frontend component for kfp links ([#1991](https://github.com/kubeflow/katib/pull/1991) by [@elenzio9](https://github.com/elenzio9))
- UI: Rename and right align the age column ([#1989](https://github.com/kubeflow/katib/pull/1989) by [@elenzio9](https://github.com/elenzio9))
- Show the trials table's status column first ([#1990](https://github.com/kubeflow/katib/pull/1990) by [@elenzio9](https://github.com/elenzio9))
- UI: Make KWA's main table responsive and add toolbar ([#1982](https://github.com/kubeflow/katib/pull/1982) by [@elenzio9](https://github.com/elenzio9))
- UI: Fix unit tests ([#1977](https://github.com/kubeflow/katib/pull/1977) by [@elenzio9](https://github.com/elenzio9))
- UI: Format code ([#1979](https://github.com/kubeflow/katib/pull/1979) by [@orfeas-k](https://github.com/orfeas-k))
- Recreate the Experiments Parallel Coordinates Graph ([#1974](https://github.com/kubeflow/katib/pull/1974) by [@elenzio9](https://github.com/elenzio9))
- Improve UI API/controller logging to ease troubleshooting ([#1966](https://github.com/kubeflow/katib/pull/1966) by [@lukeogg](https://github.com/lukeogg))

### SDK Improvements

- [SDK] Use Katib SDK for E2E Tests ([#2075](https://github.com/kubeflow/katib/pull/2075) by [@andreyvelich](https://github.com/andreyvelich))
- [SDK] Use Katib Client without Kube Config ([#2098](https://github.com/kubeflow/katib/pull/2098) by [@andreyvelich](https://github.com/andreyvelich))
- [SDK] Fix namespace parameter in tune API ([#1981](https://github.com/kubeflow/katib/pull/1981) by [@andreyvelich](https://github.com/andreyvelich))
- [SDK] Remove Final Keyword from constants ([#1980](https://github.com/kubeflow/katib/pull/1980) by [@andreyvelich](https://github.com/andreyvelich))

## Bug fixes

- Fix Release Script for Updating SDK Version ([#2104](https://github.com/kubeflow/katib/pull/2104) by [@andreyvelich](https://github.com/andreyvelich))
- [Fix] add early stopped trials in converter ([#2004](https://github.com/kubeflow/katib/pull/2004) by [@shaowei-su](https://github.com/shaowei-su))
- [bugfix] Fix value passing bug in New Experiment form ([#2027](https://github.com/kubeflow/katib/pull/2027) by [@orfeas-k](https://github.com/orfeas-k))
- Fix main process retrieve logic for early stopping ([#1988](https://github.com/kubeflow/katib/pull/1988) by [@shaowei-su](https://github.com/shaowei-su))
- [hotfix]: filter by name of experiment ([#1920](https://github.com/kubeflow/katib/pull/1920) by [@anencore94](https://github.com/anencore94))
- Fix push script to include new images ([#1911](https://github.com/kubeflow/katib/pull/1911) by [@johnugeorge](https://github.com/johnugeorge))
- fix: only validate Kubernetes Job ([#2025](https://github.com/kubeflow/katib/pull/2025) by [@zhixian82](https://github.com/zhixian82))
- Upgrade grpc-health-probe version to fix some security issues ([#2093](https://github.com/kubeflow/katib/pull/2093) by [@tenzen-y](https://github.com/tenzen-y))
- Format Katib Charm Operator ([#2115](https://github.com/kubeflow/katib/pull/2115) by [@tenzen-y](https://github.com/tenzen-y))

## Documentation

- Add CERN to adopters ([#2010](https://github.com/kubeflow/katib/pull/2010) by [@d-gol](https://github.com/d-gol))
- Add More Katib Presentations 2022 ([#2009](https://github.com/kubeflow/katib/pull/2009) by [@andreyvelich](https://github.com/andreyvelich))
- Add the documentation for simple-pbt ([#1978](https://github.com/kubeflow/katib/pull/1978) by [@tenzen-y](https://github.com/tenzen-y))
- Add the license to pbt ([#1958](https://github.com/kubeflow/katib/pull/1958) by [@tenzen-y](https://github.com/tenzen-y))
- Update the Katib version in docs ([#1950](https://github.com/kubeflow/katib/pull/1950) by [@tenzen-y](https://github.com/tenzen-y))
- Update CHANGELOG for v0.14.0 release ([#1932](https://github.com/kubeflow/katib/pull/1932) by [@johnugeorge](https://github.com/johnugeorge))

## Misc

- Update Training operator Image in CI ([#2103](https://github.com/kubeflow/katib/pull/2103) by [@johnugeorge](https://github.com/johnugeorge))
- Upgrade Go libraries to resolve security issues ([#2094](https://github.com/kubeflow/katib/pull/2094) by [@tenzen-y](https://github.com/tenzen-y))
- Run e2e with various Python versions to verify Python SDK ([#2092](https://github.com/kubeflow/katib/pull/2092) by [@tenzen-y](https://github.com/tenzen-y))
- Add a --prefer-binary flag to 'pip install' command ([#2096](https://github.com/kubeflow/katib/pull/2096) by [@tenzen-y](https://github.com/tenzen-y))
- Upgrade PyTorch version to v1.13.0 ([#2082](https://github.com/kubeflow/katib/pull/2082) by [@tenzen-y](https://github.com/tenzen-y))
- Upgrade Tensorflow version ([#2079](https://github.com/kubeflow/katib/pull/2079) by [@tenzen-y](https://github.com/tenzen-y))
- Upgrade Python version to 3.10 ([#2057](https://github.com/kubeflow/katib/pull/2057) by [@tenzen-y](https://github.com/tenzen-y))
- Pin the NumPy version with v1.23.5 in some images ([#2070](https://github.com/kubeflow/katib/pull/2070) by [@tenzen-y](https://github.com/tenzen-y))
- Upgrade the actions-setup-minikube version to v2.7.2 ([#2064](https://github.com/kubeflow/katib/pull/2064) by [@tenzen-y](https://github.com/tenzen-y))
- Remove Certificate Chain from Cert Generator ([#2045](https://github.com/kubeflow/katib/pull/2045) by [@andreyvelich](https://github.com/andreyvelich))
- Add resources to earlystopping container ([#2038](https://github.com/kubeflow/katib/pull/2038) by [@zhixian82](https://github.com/zhixian82))
- Add scripts to verify generated codes and Go Modules ([#1999](https://github.com/kubeflow/katib/pull/1999) by [@tenzen-y](https://github.com/tenzen-y))
- [Test] Reduce Katib GitHub Action Runs ([#2036](https://github.com/kubeflow/katib/pull/2036) by [@andreyvelich](https://github.com/andreyvelich))
- gh-actions: Extend action to run Frontend Unit tests ([#1998](https://github.com/kubeflow/katib/pull/1998) by [@orfeas-k](https://github.com/orfeas-k))
- [chore] Upgrade docker/metadata-action, actions/checkout, and actions/setup-python version ([#1996](https://github.com/kubeflow/katib/pull/1996) by [@tenzen-y](https://github.com/tenzen-y))
- [chore] Upgrade Go version to v1.19 ([#1995](https://github.com/kubeflow/katib/pull/1995) by [@tenzen-y](https://github.com/tenzen-y))
- Support for arm64 in simple-pbt image ([#1948](https://github.com/kubeflow/katib/pull/1948) by [@tenzen-y](https://github.com/tenzen-y))
- Support arm64 in darts-cnn-cifar10 image ([#1947](https://github.com/kubeflow/katib/pull/1947) by [@tenzen-y](https://github.com/tenzen-y))
- Support for arm64 in enas-cnn-cifar10 image ([#1944](https://github.com/kubeflow/katib/pull/1944) by [@tenzen-y](https://github.com/tenzen-y))
- Support for arm64 in pytorch-mnist image ([#1943](https://github.com/kubeflow/katib/pull/1943) by [@tenzen-y](https://github.com/tenzen-y))
- Support for arm64 in mxnet-mnist image ([#1940](https://github.com/kubeflow/katib/pull/1940) by [@tenzen-y](https://github.com/tenzen-y))
- Use the katib-new-ui for Charmed gh-actions ([#1987](https://github.com/kubeflow/katib/pull/1987) by [@tenzen-y](https://github.com/tenzen-y))
- [feat] health check for katib-controller ([#1934](https://github.com/kubeflow/katib/pull/1934) by [@anencore94](https://github.com/anencore94))
- Upgrade Optuna from v2.x.x to v3.0.0 ([#1942](https://github.com/kubeflow/katib/pull/1942) by [@keisuke-umezawa](https://github.com/keisuke-umezawa))
- Add validation webhooks for maxFailedTrialCount and parallelTrialCount ([#1936](https://github.com/kubeflow/katib/pull/1936) by [@tenzen-y](https://github.com/tenzen-y))
- Introduce Automatic platform ARGs ([#1935](https://github.com/kubeflow/katib/pull/1935) by [@tenzen-y](https://github.com/tenzen-y))
- Update training operator image in CI ([#1933](https://github.com/kubeflow/katib/pull/1933) by [@johnugeorge](https://github.com/johnugeorge))
- Update Katib SDK version ([#1931](https://github.com/kubeflow/katib/pull/1931) by [@johnugeorge](https://github.com/johnugeorge))
- [chore] Upgrade Go version to v1.18 ([#1925](https://github.com/kubeflow/katib/pull/1925) by [@tenzen-y](https://github.com/tenzen-y))
- Add the pytorch-mnist with GPU support container image ([#1916](https://github.com/kubeflow/katib/pull/1916) by [@tenzen-y](https://github.com/tenzen-y))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.14.0...v0.15.0)

# [v0.15.0-rc.1](https://github.com/kubeflow/katib/tree/v0.15.0-rc.1) (2023-02-15)

## New Features

- UI: Create the LOGS tab of Trial's details page ([#2117](https://github.com/kubeflow/katib/pull/2117) by [@elenzio9](https://github.com/elenzio9))

## Bug Fixes

- Format Katib Charm Operator ([#2115](https://github.com/kubeflow/katib/pull/2115) by [@tenzen-y](https://github.com/tenzen-y))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.15.0-rc.0...v0.15.0-rc.1)

# [v0.15.0-rc.0](https://github.com/kubeflow/katib/tree/v0.15.0-rc.0) (2023-01-27)

## Breaking Changes

- Use **Never** Resume Policy as Default ([#2102](https://github.com/kubeflow/katib/pull/2102) by [@andreyvelich](https://github.com/andreyvelich))
- Chocolate Suggestion Service is removed ([#2071](https://github.com/kubeflow/katib/pull/2071) by [@tenzen-y](https://github.com/tenzen-y))
- `request_number` is removed from the GRPC APIs ([#1994](https://github.com/kubeflow/katib/pull/1994) by [@johnugeorge](https://github.com/johnugeorge))
- The new improved and refactored Katib SDK is not backward compatible ([#2075](https://github.com/kubeflow/katib/pull/2075) by [@andreyvelich](https://github.com/andreyvelich))

## New Features

### Major Features

- Narrow down Katib RBAC rules ([#2091](https://github.com/kubeflow/katib/pull/2091) by [@johnugeorge](https://github.com/johnugeorge))
- Support Postgres as a Katib DB ([#1921](https://github.com/kubeflow/katib/pull/1921) by [@anencore94](https://github.com/anencore94))
- More Suggestion container fields in Katib Config ([#2000](https://github.com/kubeflow/katib/pull/2000) by [@fischor](https://github.com/fischor))
- Katib UI: Enable pagination/sorting/filtering ([#2017](https://github.com/kubeflow/katib/pull/2017) and [#2040](https://github.com/kubeflow/katib/pull/2040) by [@elenzio9](https://github.com/elenzio9))
- Katib UI: Add authorization mechanisms ([#1983](https://github.com/kubeflow/katib/pull/1983) by [@apo-ger](https://github.com/apo-ger))
- [SDK] Create Tune API in the Katib SDK ([#1951](https://github.com/kubeflow/katib/pull/1951) by [@andreyvelich](https://github.com/andreyvelich))
- [SDK] Get Trial Metrics from Katib DB ([#2050](https://github.com/kubeflow/katib/pull/2050) by [@andreyvelich](https://github.com/andreyvelich))

### Core Features

- Add Conformance Program Doc for AutoML and Training WG ([#2048](https://github.com/kubeflow/katib/pull/2048) by [@andreyvelich](https://github.com/andreyvelich))
- Support for grid search algorithm in Optuna Suggestion Service ([#2060](https://github.com/kubeflow/katib/pull/2060) by [@tenzen-y](https://github.com/tenzen-y))
- Add Trial Labels During Pod Mutation ([#2047](https://github.com/kubeflow/katib/pull/2047) by [@andreyvelich](https://github.com/andreyvelich))
- Support for k8s v1.25 in CI ([#1997](https://github.com/kubeflow/katib/pull/1997) by [@johnugeorge](https://github.com/johnugeorge))
- Add the CI to build multi-platform container images ([#1956](https://github.com/kubeflow/katib/pull/1956) by [@tenzen-y](https://github.com/tenzen-y))
- Drop Kubernetes v1.21 and introduce Kubernetes v1.24 ([#1953](https://github.com/kubeflow/katib/pull/1953) by [@tenzen-y](https://github.com/tenzen-y))
- Add --connect-timeout flag to katib-db-manager ([#1937](https://github.com/kubeflow/katib/pull/1937) by [@tenzen-y](https://github.com/tenzen-y))
- Implement validations for DARTS suggestion service ([#1926](https://github.com/kubeflow/katib/pull/1926) by [@tenzen-y](https://github.com/tenzen-y))
- Implement validation for Optuna suggestion service ([#1924](https://github.com/kubeflow/katib/pull/1924) by [@tenzen-y](https://github.com/tenzen-y))

### UI Improvements

- Make links in KWA's tables actual links ([#2090](https://github.com/kubeflow/katib/pull/2090) by [@elenzio9](https://github.com/elenzio9))
- frontend: Rework the trial graph using ECharts in KWA ([#2089](https://github.com/kubeflow/katib/pull/2089) by [@elenzio9](https://github.com/elenzio9))
- kwa(front): Add UI tests with Cypress ([#2088](https://github.com/kubeflow/katib/pull/2088) by [@orfeas-k](https://github.com/orfeas-k))
- Update manifests to enable authorization check mechanisms for Katib UI in Kubeflow mode ([#2041](https://github.com/kubeflow/katib/pull/2041) by [@apo-ger](https://github.com/apo-ger))
- frontend: Enable actions in experiment graph ([#2065](https://github.com/kubeflow/katib/pull/2065) by [@elenzio9](https://github.com/elenzio9))
- frontend: Show message in case of uncompleted trial instead of the graph ([#2063](https://github.com/kubeflow/katib/pull/2063) by [@elenzio9](https://github.com/elenzio9))
- frontend: Add source maps in the browser ([#2043](https://github.com/kubeflow/katib/pull/2043) by [@elenzio9](https://github.com/elenzio9))
- Backend for getting logs of a trial ([#2039](https://github.com/kubeflow/katib/pull/2039) by [@d-gol](https://github.com/d-gol))
- frontend: Show the successful trials in the experiment graph (#2013) ([#2033](https://github.com/kubeflow/katib/pull/2033) by [@elenzio9](https://github.com/elenzio9))
- frontend: Migrate from tslint to eslint in KWA ([#2042](https://github.com/kubeflow/katib/pull/2042) by [@elenzio9](https://github.com/elenzio9))
- Dedicated yaml tab for Trials ([#2034](https://github.com/kubeflow/katib/pull/2034) by [@elenzio9](https://github.com/elenzio9))
- KWA: Use new Editor component (Monaco) ([#2023](https://github.com/kubeflow/katib/pull/2023) by [@orfeas-k](https://github.com/orfeas-k))
- kwa(build): Introduce COMMIT file for building KWA ([#2014](https://github.com/kubeflow/katib/pull/2014) by [@orfeas-k](https://github.com/orfeas-k))
- frontend: Fix 500 error after detail page refresh (#1967) ([#2001](https://github.com/kubeflow/katib/pull/2001) by [@elenzio9](https://github.com/elenzio9))
- Introduce KWA's frontend component for kfp links ([#1991](https://github.com/kubeflow/katib/pull/1991) by [@elenzio9](https://github.com/elenzio9))
- UI: Rename and right align the age column ([#1989](https://github.com/kubeflow/katib/pull/1989) by [@elenzio9](https://github.com/elenzio9))
- Show the trials table's status column first ([#1990](https://github.com/kubeflow/katib/pull/1990) by [@elenzio9](https://github.com/elenzio9))
- UI: Make KWA's main table responsive and add toolbar ([#1982](https://github.com/kubeflow/katib/pull/1982) by [@elenzio9](https://github.com/elenzio9))
- UI: Fix unit tests ([#1977](https://github.com/kubeflow/katib/pull/1977) by [@elenzio9](https://github.com/elenzio9))
- UI: Format code ([#1979](https://github.com/kubeflow/katib/pull/1979) by [@orfeas-k](https://github.com/orfeas-k))
- Recreate the Experiments Parallel Coordinates Graph ([#1974](https://github.com/kubeflow/katib/pull/1974) by [@elenzio9](https://github.com/elenzio9))
- Improve UI API/controller logging to ease troubleshooting ([#1966](https://github.com/kubeflow/katib/pull/1966) by [@lukeogg](https://github.com/lukeogg))

### SDK Improvements

- [SDK] Use Katib SDK for E2E Tests ([#2075](https://github.com/kubeflow/katib/pull/2075) by [@andreyvelich](https://github.com/andreyvelich))
- [SDK] Use Katib Client without Kube Config ([#2098](https://github.com/kubeflow/katib/pull/2098) by [@andreyvelich](https://github.com/andreyvelich))
- [SDK] Fix namespace parameter in tune API ([#1981](https://github.com/kubeflow/katib/pull/1981) by [@andreyvelich](https://github.com/andreyvelich))
- [SDK] Remove Final Keyword from constants ([#1980](https://github.com/kubeflow/katib/pull/1980) by [@andreyvelich](https://github.com/andreyvelich))

## Bug fixes

- Fix Release Script for Updating SDK Version ([#2104](https://github.com/kubeflow/katib/pull/2104) by [@andreyvelich](https://github.com/andreyvelich))
- [Fix] add early stopped trials in converter ([#2004](https://github.com/kubeflow/katib/pull/2004) by [@shaowei-su](https://github.com/shaowei-su))
- [bugfix] Fix value passing bug in New Experiment form ([#2027](https://github.com/kubeflow/katib/pull/2027) by [@orfeas-k](https://github.com/orfeas-k))
- Fix main process retrieve logic for early stopping ([#1988](https://github.com/kubeflow/katib/pull/1988) by [@shaowei-su](https://github.com/shaowei-su))
- [hotfix]: filter by name of experiment ([#1920](https://github.com/kubeflow/katib/pull/1920) by [@anencore94](https://github.com/anencore94))
- Fix push script to include new images ([#1911](https://github.com/kubeflow/katib/pull/1911) by [@johnugeorge](https://github.com/johnugeorge))
- fix: only validate Kubernetes Job ([#2025](https://github.com/kubeflow/katib/pull/2025) by [@zhixian82](https://github.com/zhixian82))
- Upgrade grpc-health-probe version to fix some security issues ([#2093](https://github.com/kubeflow/katib/pull/2093) by [@tenzen-y](https://github.com/tenzen-y))

## Documentation

- Add CERN to adopters ([#2010](https://github.com/kubeflow/katib/pull/2010) by [@d-gol](https://github.com/d-gol))
- Add More Katib Presentations 2022 ([#2009](https://github.com/kubeflow/katib/pull/2009) by [@andreyvelich](https://github.com/andreyvelich))
- Add the documentation for simple-pbt ([#1978](https://github.com/kubeflow/katib/pull/1978) by [@tenzen-y](https://github.com/tenzen-y))
- Add the license to pbt ([#1958](https://github.com/kubeflow/katib/pull/1958) by [@tenzen-y](https://github.com/tenzen-y))
- Update the Katib version in docs ([#1950](https://github.com/kubeflow/katib/pull/1950) by [@tenzen-y](https://github.com/tenzen-y))
- Update CHANGELOG for v0.14.0 release ([#1932](https://github.com/kubeflow/katib/pull/1932) by [@johnugeorge](https://github.com/johnugeorge))

## Misc

- Update Training operator Image in CI ([#2103](https://github.com/kubeflow/katib/pull/2103) by [@johnugeorge](https://github.com/johnugeorge))
- Upgrade Go libraries to resolve security issues ([#2094](https://github.com/kubeflow/katib/pull/2094) by [@tenzen-y](https://github.com/tenzen-y))
- Run e2e with various Python versions to verify Python SDK ([#2092](https://github.com/kubeflow/katib/pull/2092) by [@tenzen-y](https://github.com/tenzen-y))
- Add a --prefer-binary flag to 'pip install' command ([#2096](https://github.com/kubeflow/katib/pull/2096) by [@tenzen-y](https://github.com/tenzen-y))
- Upgrade PyTorch version to v1.13.0 ([#2082](https://github.com/kubeflow/katib/pull/2082) by [@tenzen-y](https://github.com/tenzen-y))
- Upgrade Tensorflow version ([#2079](https://github.com/kubeflow/katib/pull/2079) by [@tenzen-y](https://github.com/tenzen-y))
- Upgrade Python version to 3.10 ([#2057](https://github.com/kubeflow/katib/pull/2057) by [@tenzen-y](https://github.com/tenzen-y))
- Pin the NumPy version with v1.23.5 in some images ([#2070](https://github.com/kubeflow/katib/pull/2070) by [@tenzen-y](https://github.com/tenzen-y))
- Upgrade the actions-setup-minikube version to v2.7.2 ([#2064](https://github.com/kubeflow/katib/pull/2064) by [@tenzen-y](https://github.com/tenzen-y))
- Remove Certificate Chain from Cert Generator ([#2045](https://github.com/kubeflow/katib/pull/2045) by [@andreyvelich](https://github.com/andreyvelich))
- Add resources to earlystopping container ([#2038](https://github.com/kubeflow/katib/pull/2038) by [@zhixian82](https://github.com/zhixian82))
- Add scripts to verify generated codes and Go Modules ([#1999](https://github.com/kubeflow/katib/pull/1999) by [@tenzen-y](https://github.com/tenzen-y))
- [Test] Reduce Katib GitHub Action Runs ([#2036](https://github.com/kubeflow/katib/pull/2036) by [@andreyvelich](https://github.com/andreyvelich))
- gh-actions: Extend action to run Frontend Unit tests ([#1998](https://github.com/kubeflow/katib/pull/1998) by [@orfeas-k](https://github.com/orfeas-k))
- [chore] Upgrade docker/metadata-action, actions/checkout, and actions/setup-python version ([#1996](https://github.com/kubeflow/katib/pull/1996) by [@tenzen-y](https://github.com/tenzen-y))
- [chore] Upgrade Go version to v1.19 ([#1995](https://github.com/kubeflow/katib/pull/1995) by [@tenzen-y](https://github.com/tenzen-y))
- Support for arm64 in simple-pbt image ([#1948](https://github.com/kubeflow/katib/pull/1948) by [@tenzen-y](https://github.com/tenzen-y))
- Support arm64 in darts-cnn-cifar10 image ([#1947](https://github.com/kubeflow/katib/pull/1947) by [@tenzen-y](https://github.com/tenzen-y))
- Support for arm64 in enas-cnn-cifar10 image ([#1944](https://github.com/kubeflow/katib/pull/1944) by [@tenzen-y](https://github.com/tenzen-y))
- Support for arm64 in pytorch-mnist image ([#1943](https://github.com/kubeflow/katib/pull/1943) by [@tenzen-y](https://github.com/tenzen-y))
- Support for arm64 in mxnet-mnist image ([#1940](https://github.com/kubeflow/katib/pull/1940) by [@tenzen-y](https://github.com/tenzen-y))
- Use the katib-new-ui for Charmed gh-actions ([#1987](https://github.com/kubeflow/katib/pull/1987) by [@tenzen-y](https://github.com/tenzen-y))
- [feat] health check for katib-controller ([#1934](https://github.com/kubeflow/katib/pull/1934) by [@anencore94](https://github.com/anencore94))
- Upgrade Optuna from v2.x.x to v3.0.0 ([#1942](https://github.com/kubeflow/katib/pull/1942) by [@keisuke-umezawa](https://github.com/keisuke-umezawa))
- Add validation webhooks for maxFailedTrialCount and parallelTrialCount ([#1936](https://github.com/kubeflow/katib/pull/1936) by [@tenzen-y](https://github.com/tenzen-y))
- Introduce Automatic platform ARGs ([#1935](https://github.com/kubeflow/katib/pull/1935) by [@tenzen-y](https://github.com/tenzen-y))
- Update training operator image in CI ([#1933](https://github.com/kubeflow/katib/pull/1933) by [@johnugeorge](https://github.com/johnugeorge))
- Update Katib SDK version ([#1931](https://github.com/kubeflow/katib/pull/1931) by [@johnugeorge](https://github.com/johnugeorge))
- [chore] Upgrade Go version to v1.18 ([#1925](https://github.com/kubeflow/katib/pull/1925) by [@tenzen-y](https://github.com/tenzen-y))
- Add the pytorch-mnist with GPU support container image ([#1916](https://github.com/kubeflow/katib/pull/1916) by [@tenzen-y](https://github.com/tenzen-y))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.14.0...v0.15.0-rc.0)

# [v0.14.0](https://github.com/kubeflow/katib/tree/v0.14.0) (2022-08-18)

## New Features

### Core Features

- Population based training ([#1833](https://github.com/kubeflow/katib/pull/1833) by [@a9p](https://github.com/a9p))
- Support JSON format logs in `file-metrics-collector` ([#1765](https://github.com/kubeflow/katib/pull/1765) by [@tenzen-y](https://github.com/tenzen-y))
- Include MetricsUnavailable condition to Complete in Trial ([#1877](https://github.com/kubeflow/katib/pull/1877) by [@tenzen-y](https://github.com/tenzen-y))
- Allow running examples on Apple Silicon M1 and fix image build errors for arm64 ([#1898](https://github.com/kubeflow/katib/pull/1898) by [@tenzen-y](https://github.com/tenzen-y))
- Configurable job name and service name for cert generator ([#1889](https://github.com/kubeflow/katib/pull/1889) by [@shaowei-su](https://github.com/shaowei-su))

### UI Features and Enhancements

- Add PBT to experiment creation form ([#1909](https://github.com/kubeflow/katib/pull/1909) by [@a9p](https://github.com/a9p))
- Distinct page for each Trial in the UI ([#1783](https://github.com/kubeflow/katib/pull/1783) by [@d-gol](https://github.com/d-gol))

## Bug fixes

- Add the pytorch-mnist with GPU support container image ([#1917](https://github.com/kubeflow/katib/pull/1917) by [@tenzen-y](https://github.com/tenzen-y))
- Fix push script to include new images ([#1912](https://github.com/kubeflow/katib/pull/1912) by [@johnugeorge](https://github.com/johnugeorge))
- Fixes lint warnings in YAML files ([#1902](https://github.com/kubeflow/katib/pull/1902) by [@Rishit-dagli](https://github.com/Rishit-dagli))
- Fix errors when running the test on Apple Silicon M1 ([#1886](https://github.com/kubeflow/katib/pull/1886) by [@tenzen-y](https://github.com/tenzen-y))
- Reconcile trial assignments by comparing suggestion and trials being executed ([#1831](https://github.com/kubeflow/katib/pull/1831) by [@henrysecond1](https://github.com/henrysecond1))
- Increate the probes seconds in manifests ([#1845](https://github.com/kubeflow/katib/pull/1845) by [@haoxins](https://github.com/haoxins))
- Set upper constraint for Optuna ([#1852](https://github.com/kubeflow/katib/pull/1852) by [@himkt](https://github.com/himkt))
- Don't check if trial's metadata is in spec.parameters ([#1848](https://github.com/kubeflow/katib/pull/1848) by [@alexeygorobets](https://github.com/alexeygorobets))

## Documentation

- Fix the FPGA examples documentation ([#1841](https://github.com/kubeflow/katib/pull/1841) by [@eliaskoromilas](https://github.com/eliaskoromilas))
- Add CyberAgent to adopters ([#1894](https://github.com/kubeflow/katib/pull/1894) by [@tenzen-y](https://github.com/tenzen-y))

## Misc

- Updating the training operator image in CI ([#1910](https://github.com/kubeflow/katib/pull/1910) by [@johnugeorge](https://github.com/johnugeorge))
- Upgrade Python and Pytorch versions for some examples ([#1906](https://github.com/kubeflow/katib/pull/1906) by [@tenzen-y](https://github.com/tenzen-y))
- Linting for K8s YAML files ([#1901](https://github.com/kubeflow/katib/pull/1901) by [@Rishit-dagli](https://github.com/Rishit-dagli))
- Change integration test sysytem from KinD Cluster to Minikube Cluster ([#1899](https://github.com/kubeflow/katib/pull/1899) by [@tenzen-y](https://github.com/tenzen-y))
- Upgrade mysql version to v8.0.29 ([#1897](https://github.com/kubeflow/katib/pull/1897) by [@tenzen-y](https://github.com/tenzen-y))
- Upgrade tensorflow-aarch64 version to v2.9.1 ([#1891](https://github.com/kubeflow/katib/pull/1891) by [@tenzen-y](https://github.com/tenzen-y))
- chore: Upgrade Go libraries to resolve some security issues in the katib-controller ([#1888](https://github.com/kubeflow/katib/pull/1888) by [@tenzen-y](https://github.com/tenzen-y))
- Migrate kubeflow-katib-presubmit to GitHub Actions ([#1882](https://github.com/kubeflow/katib/pull/1882) by [@tenzen-y](https://github.com/tenzen-y))
- Add semicolon when using `command` command in Makefile ([#1885](https://github.com/kubeflow/katib/pull/1885) by [@tenzen-y](https://github.com/tenzen-y))
- Fix `HAS_SHELLCHECK` and `HAS_SETUP_ENVTEST` in Makefile ([#1884](https://github.com/kubeflow/katib/pull/1884) by [@tenzen-y](https://github.com/tenzen-y))
- Remove presubmit tests depending on optional-test-infra ([#1871](https://github.com/kubeflow/katib/pull/1871) by [@aws-kf-ci-bot](https://github.com/aws-kf-ci-bot))
- Upgrade the Tensorflow version to address some security issues ([#1870](https://github.com/kubeflow/katib/pull/1870) by [@tenzen-y](https://github.com/tenzen-y))
- Upgrade the grpc_health_probe version to v0.4.11 to resolve security vulnerability CVE-2022-27191 ([#1875](https://github.com/kubeflow/katib/pull/1875) by [@tenzen-y](https://github.com/tenzen-y))
- additional metric names should not include objective metric name ([#1874](https://github.com/kubeflow/katib/pull/1874) by [@henrysecond1](https://github.com/henrysecond1))
- Upgrade the Kubernetes Python client to 22.6.0 ([#1869](https://github.com/kubeflow/katib/pull/1869) by [@tenzen-y](https://github.com/tenzen-y))
- Upgrade the kubebuilder to v3.2.0 and Kubernetes Go libraries to v1.22.2 ([#1861](https://github.com/kubeflow/katib/pull/1861) by [@tenzen-y](https://github.com/tenzen-y))
- Update FPGA XGBoost example ([#1865](https://github.com/kubeflow/katib/pull/1865) by [@eliaskoromilas](https://github.com/eliaskoromilas))
- Fix kubeflowkatib/mxnet-mnist image ([#1866](https://github.com/kubeflow/katib/pull/1866) by [@tenzen-y](https://github.com/tenzen-y))
- pins pip and setuptools versions operators to avoid installation issues ([#1867](https://github.com/kubeflow/katib/pull/1867) by [@DnPlas](https://github.com/DnPlas))
- Add shellcheck ([#1857](https://github.com/kubeflow/katib/pull/1857) by [@tenzen-y](https://github.com/tenzen-y))
- Bump kubeflow-katib and kfp version in notebook examples ([#1849](https://github.com/kubeflow/katib/pull/1849) by [@tenzen-y](https://github.com/tenzen-y))
- Add prometheus scraping and grafana support to charmed katib-controller operator ([#1839](https://github.com/kubeflow/katib/pull/1839) by [@jardon](https://github.com/jardon))
- Upgrade Black to fix linting ([#1842](https://github.com/kubeflow/katib/pull/1842) by [@jardon](https://github.com/jardon))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.13.0...v0.14.0).

# [v0.13.0](https://github.com/kubeflow/katib/tree/v0.13.0) (2022-03-04)

## New Features

### Algorithms and Components

- Implement validation for Early Stopping ([#1709](https://github.com/kubeflow/katib/pull/1709) by [@tenzen-y](https://github.com/tenzen-y))
- Change namespace label for Metrics Collector injection ([#1740](https://github.com/kubeflow/katib/pull/1740) by [@andreyvelich](https://github.com/andreyvelich))
- Modify gRPC API with Current Request Number ([#1728](https://github.com/kubeflow/katib/pull/1728) by [@andreyvelich](https://github.com/andreyvelich))
- Allow to remove each resource in Katib config ([#1729](https://github.com/kubeflow/katib/pull/1729) by [@andreyvelich](https://github.com/andreyvelich))
- Support leader election for Katib Controller ([#1713](https://github.com/kubeflow/katib/pull/1713) by [@tenzen-y](https://github.com/tenzen-y))
- Change default Metrics Collect format ([#1707](https://github.com/kubeflow/katib/pull/1707) by [@anencore94](https://github.com/anencore94))
- Bump Python version to 3.9 ([#1731](https://github.com/kubeflow/katib/pull/1731) by [@tenzen-y](https://github.com/tenzen-y))
- Update Go version to 1.17 ([#1683](https://github.com/kubeflow/katib/pull/1683) by [@andreyvelich](https://github.com/andreyvelich))
- Create Python script to run e2e Argo Workflow ([#1674](https://github.com/kubeflow/katib/pull/1674) by [@andreyvelich](https://github.com/andreyvelich))
- Reimplement Katib Cert Generator in Go ([#1662](https://github.com/kubeflow/katib/pull/1662) by [@tenzen-y](https://github.com/tenzen-y))
- SDK: change list apis to return objects as default ([#1630](https://github.com/kubeflow/katib/pull/1630) by [@anencore94](https://github.com/anencore94))

### UI Features

- Enhance Katib UI feasible space ([#1721](https://github.com/kubeflow/katib/pull/1721) by [@seong7](https://github.com/seong7))
- Handle missing TrialTemplates in Katib UI ([#1652](https://github.com/kubeflow/katib/pull/1652) by [@kimwnasptd](https://github.com/kimwnasptd))
- Add Prettier devDependency in Katib UI ([#1629](https://github.com/kubeflow/katib/pull/1629) by [@seong7](https://github.com/seong7))

## Documentation

- Fix a link for GRPC API documentation ([#1786](https://github.com/kubeflow/katib/pull/1786) by [@tenzen-y](https://github.com/tenzen-y))
- Add my presentations that include Katib ([#1753](https://github.com/kubeflow/katib/pull/1753) by [@terrytangyuan](https://github.com/terrytangyuan))
- Add Akuity to list of adopters ([#1749](https://github.com/kubeflow/katib/pull/1749) by [@terrytangyuan](https://github.com/terrytangyuan))
- Change Argo -> Argo Workflows ([#1741](https://github.com/kubeflow/katib/pull/1741) by [@terrytangyuan](https://github.com/terrytangyuan))
- Update Algorithm Service Doc for the new CI script ([#1724](https://github.com/kubeflow/katib/pull/1724) by [@andreyvelich](https://github.com/andreyvelich))
- Update link to Training Operator ([#1699](https://github.com/kubeflow/katib/pull/1699) by [@terrytangyuan](https://github.com/terrytangyuan))
- Refactor examples folder structure ([#1691](https://github.com/kubeflow/katib/pull/1691) by [@andreyvelich](https://github.com/andreyvelich))
- Fix README in examples directory ([#1687](https://github.com/kubeflow/katib/pull/1687) by [@tenzen-y](https://github.com/tenzen-y))
- Add Kubeflow MXJob example ([#1688](https://github.com/kubeflow/katib/pull/1688) by [@andreyvelich](https://github.com/andreyvelich))
- Update FPGA examples ([#1685](https://github.com/kubeflow/katib/pull/1685) by [@eliaskoromilas](https://github.com/eliaskoromilas))
- Refactor README ([#1667](https://github.com/kubeflow/katib/pull/1667) by [@andreyvelich](https://github.com/andreyvelich))
- Change the minimal Kustomize version in the developer guide ([#1675](https://github.com/kubeflow/katib/pull/1675) by [@tenzen-y](https://github.com/tenzen-y))
- Add Katib release process guide ([#1641](https://github.com/kubeflow/katib/pull/1641) by [@andreyvelich](https://github.com/andreyvelich))

## Bug Fixes

- Remove unrecognized keys from metadata.yaml in Charmed operators ([#1759](https://github.com/kubeflow/katib/pull/1759) by [@DnPlas](https://github.com/DnPlas))
- Fix the default Metrics Collector regex ([#1755](https://github.com/kubeflow/katib/pull/1755) by [@andreyvelich](https://github.com/andreyvelich))
- Fix Status Handling in Charmed Operators ([#1743](https://github.com/kubeflow/katib/pull/1743) by [@DomFleischmann](https://github.com/DomFleischmann))
- Fix bug on list type HP in Katib UI ([#1704](https://github.com/kubeflow/katib/pull/1704) by [@seong7](https://github.com/seong7))
- Fix Range for Int and Double values in Grid search ([#1732](https://github.com/kubeflow/katib/pull/1732) by [@andreyvelich](https://github.com/andreyvelich))
- Check if parameter references exist in Experiment parameters ([#1726](https://github.com/kubeflow/katib/pull/1726) by [@henrysecond1](https://github.com/henrysecond1))
- Fix same set for HyperParameters in Bayesian Optimization algorithm ([#1701](https://github.com/kubeflow/katib/pull/1701) by [@fabianvdW](https://github.com/fabianvdW))
- Close MySQL statement and rows resources when SQL exec ends ([#1720](https://github.com/kubeflow/katib/pull/1720) by [@chenwenjun-github](https://github.com/chenwenjun-github))
- Fix Cluster Role of Katib Controller to access image pull secrets ([#1725](https://github.com/kubeflow/katib/pull/1725) by [@henrysecond1](https://github.com/henrysecond1))
- Emit events when fails to reconcile all Trials ([#1706](https://github.com/kubeflow/katib/pull/1706) by [@henrysecond1](https://github.com/henrysecond1))
- Missing metrics port annotation ([#1715](https://github.com/kubeflow/katib/pull/1715) by [@alexeykaplin](https://github.com/alexeykaplin))
- Fix absolute value in Katib UI ([#1676](https://github.com/kubeflow/katib/pull/1676) by [@anencore94](https://github.com/anencore94))
- Add missing omitempty parameter to APIs ([#1645](https://github.com/kubeflow/katib/pull/1645) by [@andreyvelich](https://github.com/andreyvelich))
- Reconcile semantics for Suggestion Algorithms ([#1633](https://github.com/kubeflow/katib/pull/1633) by [@johnugeorge](https://github.com/johnugeorge))
- Fix default label for Training Operators ([#1813](https://github.com/kubeflow/katib/pull/1813) by [@andreyvelich](https://github.com/andreyvelich))
- Update supported Python version for Katib SDK ([#1798](https://github.com/kubeflow/katib/pull/1798) by [@tenzen-y](https://github.com/tenzen-y))

## Misc

- Use release tags for Trial images ([#1757](https://github.com/kubeflow/katib/pull/1757) by [@andreyvelich](https://github.com/andreyvelich))
- Upgrade cert-manager API from v1alpha2 to v1 ([#1752](https://github.com/kubeflow/katib/pull/1752) by [@haoxins](https://github.com/haoxins))
- Add Workflow to Publish Katib Images ([#1746](https://github.com/kubeflow/katib/pull/1746) by [@andreyvelich](https://github.com/andreyvelich))
- Update Charmed Katib Operators + CI to 0.12 ([#1717](https://github.com/kubeflow/katib/pull/1717) by [@knkski](https://github.com/knkski))
- Updating Katib CI to use Training Operator ([#1710](https://github.com/kubeflow/katib/pull/1710) by [@midhun1998](https://github.com/midhun1998))
- Update OWNERS for Charmed operators ([#1718](https://github.com/kubeflow/katib/pull/1718) by [@ca-scribner](https://github.com/ca-scribner))
- Implement some unit tests for the Katib Config package ([#1690](https://github.com/kubeflow/katib/pull/1690) by [@tenzen-y](https://github.com/tenzen-y))
- Add GitHub Actions for Python unit tests ([#1677](https://github.com/kubeflow/katib/pull/1677) by [@andreyvelich](https://github.com/andreyvelich))
- Add OWNERS file for the Katib new UI ([#1681](https://github.com/kubeflow/katib/pull/1681) by [@kimwnasptd](https://github.com/kimwnasptd))
- Add envtest to check `reconcileRBAC` ([#1678](https://github.com/kubeflow/katib/pull/1678) by [@tenzen-y](https://github.com/tenzen-y))
- Use golangci-lint as linter for Go ([#1671](https://github.com/kubeflow/katib/pull/1671) by [@tenzen-y](https://github.com/tenzen-y))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.12.0...v0.13.0)

# [v0.13.0-rc.1](https://github.com/kubeflow/katib/tree/v0.13.0-rc.1) (2022-02-15)

## Bug fixes

- Fix default label for Training Operators ([#1813](https://github.com/kubeflow/katib/pull/1813) by [@andreyvelich](https://github.com/andreyvelich))
- Update supported Python version for Katib SDK ([#1798](https://github.com/kubeflow/katib/pull/1798) by [@tenzen-y](https://github.com/tenzen-y))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.13.0-rc.0...v0.13.0-rc.1)

# [v0.13.0-rc.0](https://github.com/kubeflow/katib/tree/v0.13.0-rc.0) (2022-01-25)

## New Features

### Algorithms and Components

- Implement validation for Early Stopping ([#1709](https://github.com/kubeflow/katib/pull/1709) by [@tenzen-y](https://github.com/tenzen-y))
- Change namespace label for Metrics Collector injection ([#1740](https://github.com/kubeflow/katib/pull/1740) by [@andreyvelich](https://github.com/andreyvelich))
- Modify gRPC API with Current Request Number ([#1728](https://github.com/kubeflow/katib/pull/1728) by [@andreyvelich](https://github.com/andreyvelich))
- Allow to remove each resource in Katib config ([#1729](https://github.com/kubeflow/katib/pull/1729) by [@andreyvelich](https://github.com/andreyvelich))
- Support leader election for Katib Controller ([#1713](https://github.com/kubeflow/katib/pull/1713) by [@tenzen-y](https://github.com/tenzen-y))
- Change default Metrics Collect format ([#1707](https://github.com/kubeflow/katib/pull/1707) by [@anencore94](https://github.com/anencore94))
- Bump Python version to 3.9 ([#1731](https://github.com/kubeflow/katib/pull/1731) by [@tenzen-y](https://github.com/tenzen-y))
- Update Go version to 1.17 ([#1683](https://github.com/kubeflow/katib/pull/1683) by [@andreyvelich](https://github.com/andreyvelich))
- Create Python script to run e2e Argo Workflow ([#1674](https://github.com/kubeflow/katib/pull/1674) by [@andreyvelich](https://github.com/andreyvelich))
- Reimplement Katib Cert Generator in Go ([#1662](https://github.com/kubeflow/katib/pull/1662) by [@tenzen-y](https://github.com/tenzen-y))
- SDK: change list apis to return objects as default ([#1630](https://github.com/kubeflow/katib/pull/1630) by [@anencore94](https://github.com/anencore94))

### UI Features

- Enhance Katib UI feasible space ([#1721](https://github.com/kubeflow/katib/pull/1721) by [@seong7](https://github.com/seong7))
- Handle missing TrialTemplates in Katib UI ([#1652](https://github.com/kubeflow/katib/pull/1652) by [@kimwnasptd](https://github.com/kimwnasptd))
- Add Prettier devDependency in Katib UI ([#1629](https://github.com/kubeflow/katib/pull/1629) by [@seong7](https://github.com/seong7))

## Documentation

- Fix a link for GRPC API documentation ([#1786](https://github.com/kubeflow/katib/pull/1786) by [@tenzen-y](https://github.com/tenzen-y))
- Add my presentations that include Katib ([#1753](https://github.com/kubeflow/katib/pull/1753) by [@terrytangyuan](https://github.com/terrytangyuan))
- Add Akuity to list of adopters ([#1749](https://github.com/kubeflow/katib/pull/1749) by [@terrytangyuan](https://github.com/terrytangyuan))
- Change Argo -> Argo Workflows ([#1741](https://github.com/kubeflow/katib/pull/1741) by [@terrytangyuan](https://github.com/terrytangyuan))
- Update Algorithm Service Doc for the new CI script ([#1724](https://github.com/kubeflow/katib/pull/1724) by [@andreyvelich](https://github.com/andreyvelich))
- Update link to Training Operator ([#1699](https://github.com/kubeflow/katib/pull/1699) by [@terrytangyuan](https://github.com/terrytangyuan))
- Refactor examples folder structure ([#1691](https://github.com/kubeflow/katib/pull/1691) by [@andreyvelich](https://github.com/andreyvelich))
- Fix README in examples directory ([#1687](https://github.com/kubeflow/katib/pull/1687) by [@tenzen-y](https://github.com/tenzen-y))
- Add Kubeflow MXJob example ([#1688](https://github.com/kubeflow/katib/pull/1688) by [@andreyvelich](https://github.com/andreyvelich))
- Update FPGA examples ([#1685](https://github.com/kubeflow/katib/pull/1685) by [@eliaskoromilas](https://github.com/eliaskoromilas))
- Refactor README ([#1667](https://github.com/kubeflow/katib/pull/1667) by [@andreyvelich](https://github.com/andreyvelich))
- Change the minimal Kustomize version in the developer guide ([#1675](https://github.com/kubeflow/katib/pull/1675) by [@tenzen-y](https://github.com/tenzen-y))
- Add Katib release process guide ([#1641](https://github.com/kubeflow/katib/pull/1641) by [@andreyvelich](https://github.com/andreyvelich))

## Bug Fixes

- Remove unrecognized keys from metadata.yaml in Charmed operators ([#1759](https://github.com/kubeflow/katib/pull/1759) by [@DnPlas](https://github.com/DnPlas))
- Fix the default Metrics Collector regex ([#1755](https://github.com/kubeflow/katib/pull/1755) by [@andreyvelich](https://github.com/andreyvelich))
- Fix Status Handling in Charmed Operators ([#1743](https://github.com/kubeflow/katib/pull/1743) by [@DomFleischmann](https://github.com/DomFleischmann))
- Fix bug on list type HP in Katib UI ([#1704](https://github.com/kubeflow/katib/pull/1704) by [@seong7](https://github.com/seong7))
- Fix Range for Int and Double values in Grid search ([#1732](https://github.com/kubeflow/katib/pull/1732) by [@andreyvelich](https://github.com/andreyvelich))
- Check if parameter references exist in Experiment parameters ([#1726](https://github.com/kubeflow/katib/pull/1726) by [@henrysecond1](https://github.com/henrysecond1))
- Fix same set for HyperParameters in Bayesian Optimization algorithm ([#1701](https://github.com/kubeflow/katib/pull/1701) by [@fabianvdW](https://github.com/fabianvdW))
- Close MySQL statement and rows resources when SQL exec ends ([#1720](https://github.com/kubeflow/katib/pull/1720) by [@chenwenjun-github](https://github.com/chenwenjun-github))
- Fix Cluster Role of Katib Controller to access image pull secrets ([#1725](https://github.com/kubeflow/katib/pull/1725) by [@henrysecond1](https://github.com/henrysecond1))
- Emit events when fails to reconcile all Trials ([#1706](https://github.com/kubeflow/katib/pull/1706) by [@henrysecond1](https://github.com/henrysecond1))
- Missing metrics port annotation ([#1715](https://github.com/kubeflow/katib/pull/1715) by [@alexeykaplin](https://github.com/alexeykaplin))
- Fix absolute value in Katib UI ([#1676](https://github.com/kubeflow/katib/pull/1676) by [@anencore94](https://github.com/anencore94))
- Add missing omitempty parameter to APIs ([#1645](https://github.com/kubeflow/katib/pull/1645) by [@andreyvelich](https://github.com/andreyvelich))
- Reconcile semantics for Suggestion Algorithms ([#1633](https://github.com/kubeflow/katib/pull/1633) by [@johnugeorge](https://github.com/johnugeorge))

## Misc

- Use release tags for Trial images ([#1757](https://github.com/kubeflow/katib/pull/1757) by [@andreyvelich](https://github.com/andreyvelich))
- Upgrade cert-manager API from v1alpha2 to v1 ([#1752](https://github.com/kubeflow/katib/pull/1752) by [@haoxins](https://github.com/haoxins))
- Add Workflow to Publish Katib Images ([#1746](https://github.com/kubeflow/katib/pull/1746) by [@andreyvelich](https://github.com/andreyvelich))
- Update Charmed Katib Operators + CI to 0.12 ([#1717](https://github.com/kubeflow/katib/pull/1717) by [@knkski](https://github.com/knkski))
- Updating Katib CI to use Training Operator ([#1710](https://github.com/kubeflow/katib/pull/1710) by [@midhun1998](https://github.com/midhun1998))
- Update OWNERS for Charmed operators ([#1718](https://github.com/kubeflow/katib/pull/1718) by [@ca-scribner](https://github.com/ca-scribner))
- Implement some unit tests for the Katib Config package ([#1690](https://github.com/kubeflow/katib/pull/1690) by [@tenzen-y](https://github.com/tenzen-y))
- Add GitHub Actions for Python unit tests ([#1677](https://github.com/kubeflow/katib/pull/1677) by [@andreyvelich](https://github.com/andreyvelich))
- Add OWNERS file for the Katib new UI ([#1681](https://github.com/kubeflow/katib/pull/1681) by [@kimwnasptd](https://github.com/kimwnasptd))
- Add envtest to check `reconcileRBAC` ([#1678](https://github.com/kubeflow/katib/pull/1678) by [@tenzen-y](https://github.com/tenzen-y))
- Use golangci-lint as linter for Go ([#1671](https://github.com/kubeflow/katib/pull/1671) by [@tenzen-y](https://github.com/tenzen-y))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.12.0...v0.13.0-rc.0)

# [v0.12.0](https://github.com/kubeflow/katib/tree/v0.12.0) (2021-10-05)

## New Features

### Algorithms and Components

- Add Optuna based suggestion service ([#1613](https://github.com/kubeflow/katib/pull/1613) by [@g-votte](https://github.com/g-votte))
- Support Sobol's Quasirandom Sequence using Goptuna. ([#1523](https://github.com/kubeflow/katib/pull/1523) by [@c-bata](https://github.com/c-bata))
- Bump the Goptuna version up to v0.8.0 with IPOP-CMA-ES and BIPOP-CMA-ES support. ([#1519](https://github.com/kubeflow/katib/pull/1519) by [@c-bata](https://github.com/c-bata))
- Validate possible operations for Grid suggestion ([#1205](https://github.com/kubeflow/katib/pull/1205) by [@andreyvelich](https://github.com/andreyvelich))
- Validate for Bayesian Optimization algorithm settings ([#1600](https://github.com/kubeflow/katib/pull/1600) by [@anencore94](https://github.com/anencore94))
- Add Support for Argo Workflows ([#1605](https://github.com/kubeflow/katib/pull/1605) by [@andreyvelich](https://github.com/andreyvelich))
- Add Support for XGBoost Operator with LightGBM example ([#1603](https://github.com/kubeflow/katib/pull/1603) by [@andreyvelich](https://github.com/andreyvelich))
- Allow empty resources for CPU and Memory in Katib config ([#1564](https://github.com/kubeflow/katib/pull/1564) by [@andreyvelich](https://github.com/andreyvelich))
- Add kustomization overlay: katib-openshift ([#1513](https://github.com/kubeflow/katib/pull/1513) by [@maanur](https://github.com/maanur))
- Switch to SDI in Katib Charm ([#1555](https://github.com/kubeflow/katib/pull/1555) by [@knkski](https://github.com/knkski))

### UI Features

- Add Multivariate TPE to Katib UI ([#1625](https://github.com/kubeflow/katib/pull/1625) by [@andreyvelich](https://github.com/andreyvelich))
- Update Katib UI with Optuna Algorithm Settings ([#1626](https://github.com/kubeflow/katib/pull/1626) by [@andreyvelich](https://github.com/andreyvelich))
- Change the default image for the new Katib UI ([#1608](https://github.com/kubeflow/katib/pull/1608) by [@andreyvelich](https://github.com/andreyvelich))

## Documentation

- Add Katib 2021 ROADMAP ([#1524](https://github.com/kubeflow/katib/pull/1524) by [@andreyvelich](https://github.com/andreyvelich))
- Add AutoML and Training WG Summit July 2021 ([#1615](https://github.com/kubeflow/katib/pull/1615) by [@andreyvelich](https://github.com/andreyvelich))
- Add the new Katib presentations 2021 ([#1539](https://github.com/kubeflow/katib/pull/1539) by [@andreyvelich](https://github.com/andreyvelich))
- Add Doc checklist to PR template ([#1568](https://github.com/kubeflow/katib/pull/1568) by [@andreyvelich](https://github.com/andreyvelich))
- Fix typo in operators/README ([#1557](https://github.com/kubeflow/katib/pull/1557) by [@evilnick](https://github.com/evilnick))
- Adds docs on how to use Katib Charm within KF ([#1556](https://github.com/kubeflow/katib/pull/1556) by [@RFMVasconcelos](https://github.com/RFMVasconcelos))
- Fix a link to Kustomize manifest for new Katib UI ([#1521](https://github.com/kubeflow/katib/pull/1521) by [@c-bata](https://github.com/c-bata))

## Bug Fixes

- Fix UI for handling missing params ([#1657](https://github.com/kubeflow/katib/pull/1657) by [@kimwnasptd](https://github.com/kimwnasptd))
- Reconcile semantics for Suggestion Algorithms ([#1644](https://github.com/kubeflow/katib/pull/1644) by [@johnugeorge](https://github.com/johnugeorge))
- Fix Metrics Collector error in case of non-existing Process ([#1614](https://github.com/kubeflow/katib/pull/1614) by [@andreyvelich](https://github.com/andreyvelich))
- Fix mysql version in docker image ([#1594](https://github.com/kubeflow/katib/pull/1594) by [@munagekar](https://github.com/munagekar))
- Fix grep in Tekton Experiment Doc ([#1578](https://github.com/kubeflow/katib/pull/1578) by [@andreyvelich](https://github.com/andreyvelich))
- Error messages corrected ([#1522](https://github.com/kubeflow/katib/pull/1522) by [@himanshu007-creator](https://github.com/himanshu007-creator))
- Install charmcraft 1.0.0 ([#1593](https://github.com/kubeflow/katib/pull/1593) by [@DomFleischmann](https://github.com/DomFleischmann))

## Misc

- Modify XGBoostJob example for the new Controller ([#1623](https://github.com/kubeflow/katib/pull/1623) by [@andreyvelich](https://github.com/andreyvelich))
- Modify Labels for controller resources ([#1621](https://github.com/kubeflow/katib/pull/1621) by [@andreyvelich](https://github.com/andreyvelich))
- Modify Labels for Katib Components ([#1611](https://github.com/kubeflow/katib/pull/1611) by [@andreyvelich](https://github.com/andreyvelich))
- Upgrade CRDs to apiextensions.k8s.io/v1 ([#1610](https://github.com/kubeflow/katib/pull/1610) by [@andreyvelich](https://github.com/andreyvelich))
- Update Katib SDK with OpenAPI generator ([#1572](https://github.com/kubeflow/katib/pull/1572) by [@andreyvelich](https://github.com/andreyvelich))
- Disable default PV for Experiment with resume from volume ([#1552](https://github.com/kubeflow/katib/pull/1552) by [@andreyvelich](https://github.com/andreyvelich))
- Remove PV from MySQL component ([#1527](https://github.com/kubeflow/katib/pull/1527) by [@andreyvelich](https://github.com/andreyvelich))
- feat: add naming regex check on validating webhook ([#1541](https://github.com/kubeflow/katib/pull/1541) by [@anencore94](https://github.com/anencore94))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.11.1...v0.12.0)

# [v0.12.0-rc.1](https://github.com/kubeflow/katib/tree/v0.12.0-rc.1) (2021-09-07)

## Bug Fixes

- Fix UI for handling missing params ([#1657](https://github.com/kubeflow/katib/pull/1657) by [@kimwnasptd](https://github.com/kimwnasptd))
- Reconcile semantics for Suggestion Algorithms ([#1644](https://github.com/kubeflow/katib/pull/1644) by [@johnugeorge](https://github.com/johnugeorge))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.12.0-rc.0...v0.12.0-rc.1)

# [v0.12.0-rc.0](https://github.com/kubeflow/katib/tree/v0.12.0-rc.0) (2021-08-19)

## New Features

### Algorithms and Components

- Add Optuna based suggestion service ([#1613](https://github.com/kubeflow/katib/pull/1613) by [@g-votte](https://github.com/g-votte))
- Support Sobol's Quasirandom Sequence using Goptuna. ([#1523](https://github.com/kubeflow/katib/pull/1523) by [@c-bata](https://github.com/c-bata))
- Bump the Goptuna version up to v0.8.0 with IPOP-CMA-ES and BIPOP-CMA-ES support. ([#1519](https://github.com/kubeflow/katib/pull/1519) by [@c-bata](https://github.com/c-bata))
- Validate possible operations for Grid suggestion ([#1205](https://github.com/kubeflow/katib/pull/1205) by [@andreyvelich](https://github.com/andreyvelich))
- Validate for Bayesian Optimization algorithm settings ([#1600](https://github.com/kubeflow/katib/pull/1600) by [@anencore94](https://github.com/anencore94))
- Add Support for Argo Workflows ([#1605](https://github.com/kubeflow/katib/pull/1605) by [@andreyvelich](https://github.com/andreyvelich))
- Add Support for XGBoost Operator with LightGBM example ([#1603](https://github.com/kubeflow/katib/pull/1603) by [@andreyvelich](https://github.com/andreyvelich))
- Allow empty resources for CPU and Memory in Katib config ([#1564](https://github.com/kubeflow/katib/pull/1564) by [@andreyvelich](https://github.com/andreyvelich))
- Add kustomization overlay: katib-openshift ([#1513](https://github.com/kubeflow/katib/pull/1513) by [@maanur](https://github.com/maanur))
- Switch to SDI in Katib Charm ([#1555](https://github.com/kubeflow/katib/pull/1555) by [@knkski](https://github.com/knkski))

### UI Features

- Add Multivariate TPE to Katib UI ([#1625](https://github.com/kubeflow/katib/pull/1625) by [@andreyvelich](https://github.com/andreyvelich))
- Update Katib UI with Optuna Algorithm Settings ([#1626](https://github.com/kubeflow/katib/pull/1626) by [@andreyvelich](https://github.com/andreyvelich))
- Change the default image for the new Katib UI ([#1608](https://github.com/kubeflow/katib/pull/1608) by [@andreyvelich](https://github.com/andreyvelich))

## Documentation

- Add Katib 2021 ROADMAP ([#1524](https://github.com/kubeflow/katib/pull/1524) by [@andreyvelich](https://github.com/andreyvelich))
- Add AutoML and Training WG Summit July 2021 ([#1615](https://github.com/kubeflow/katib/pull/1615) by [@andreyvelich](https://github.com/andreyvelich))
- Add the new Katib presentations 2021 ([#1539](https://github.com/kubeflow/katib/pull/1539) by [@andreyvelich](https://github.com/andreyvelich))
- Add Doc checklist to PR template ([#1568](https://github.com/kubeflow/katib/pull/1568) by [@andreyvelich](https://github.com/andreyvelich))
- Fix typo in operators/README ([#1557](https://github.com/kubeflow/katib/pull/1557) by [@evilnick](https://github.com/evilnick))
- Adds docs on how to use Katib Charm within KF ([#1556](https://github.com/kubeflow/katib/pull/1556) by [@RFMVasconcelos](https://github.com/RFMVasconcelos))
- Fix a link to Kustomize manifest for new Katib UI ([#1521](https://github.com/kubeflow/katib/pull/1521) by [@c-bata](https://github.com/c-bata))

## Bug Fixes

- Fix Metrics Collector error in case of non-existing Process ([#1614](https://github.com/kubeflow/katib/pull/1614) by [@andreyvelich](https://github.com/andreyvelich))
- Fix mysql version in docker image ([#1594](https://github.com/kubeflow/katib/pull/1594) by [@munagekar](https://github.com/munagekar))
- Fix grep in Tekton Experiment Doc ([#1578](https://github.com/kubeflow/katib/pull/1578) by [@andreyvelich](https://github.com/andreyvelich))
- Error messages corrected ([#1522](https://github.com/kubeflow/katib/pull/1522) by [@himanshu007-creator](https://github.com/himanshu007-creator))
- Install charmcraft 1.0.0 ([#1593](https://github.com/kubeflow/katib/pull/1593) by [@DomFleischmann](https://github.com/DomFleischmann))

## Misc

- Modify XGBoostJob example for the new Controller ([#1623](https://github.com/kubeflow/katib/pull/1623) by [@andreyvelich](https://github.com/andreyvelich))
- Modify Labels for controller resources ([#1621](https://github.com/kubeflow/katib/pull/1621) by [@andreyvelich](https://github.com/andreyvelich))
- Modify Labels for Katib Components ([#1611](https://github.com/kubeflow/katib/pull/1611) by [@andreyvelich](https://github.com/andreyvelich))
- Upgrade CRDs to apiextensions.k8s.io/v1 ([#1610](https://github.com/kubeflow/katib/pull/1610) by [@andreyvelich](https://github.com/andreyvelich))
- Update Katib SDK with OpenAPI generator ([#1572](https://github.com/kubeflow/katib/pull/1572) by [@andreyvelich](https://github.com/andreyvelich))
- Disable default PV for Experiment with resume from volume ([#1552](https://github.com/kubeflow/katib/pull/1552) by [@andreyvelich](https://github.com/andreyvelich))
- Remove PV from MySQL component ([#1527](https://github.com/kubeflow/katib/pull/1527) by [@andreyvelich](https://github.com/andreyvelich))
- feat: add naming regex check on validating webhook ([#1541](https://github.com/kubeflow/katib/pull/1541) by [@anencore94](https://github.com/anencore94))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.11.1...v0.12.0-rc.0)

# [v0.11.1](https://github.com/kubeflow/katib/tree/v0.11.1) (2021-06-09)

## Bug fixes

- Fix Katib manifest for Kubeflow 1.3 ([#1503](https://github.com/kubeflow/katib/pull/1503) by [@yanniszark](https://github.com/yanniszark))
- Fix Katib release script ([#1510](https://github.com/kubeflow/katib/pull/1510) by [@andreyvelich](https://github.com/andreyvelich))

## Enhancements

- Remove Application CR ([#1509](https://github.com/kubeflow/katib/pull/1509) by [@yanniszark](https://github.com/yanniszark))
- Modify Katib manifest to support newer Kustomize version ([#1515](https://github.com/kubeflow/katib/pull/1515) by [@DavidSpek](https://github.com/DavidSpek))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.11.0...v0.11.1)

# [v0.11.0](https://github.com/kubeflow/katib/tree/v0.11.0) (2021-03-22)

## New Features

### Core Features

- Disable dynamic Webhook creation ([#1450](https://github.com/kubeflow/katib/pull/1450) by [@andreyvelich](https://github.com/andreyvelich))
- Add the `waitAllProcesses` flag to the Katib config ([#1394](https://github.com/kubeflow/katib/pull/1394) by [@robbertvdg](https://github.com/robbertvdg))
- Migrate Katib to Go modules ([#1438](https://github.com/kubeflow/katib/pull/1438) by [@andreyvelich](https://github.com/andreyvelich))
- Update Katib SDK with the `get_success_trial_details` API ([#1442](https://github.com/kubeflow/katib/pull/1442) by [@Adarsh2910](https://github.com/Adarsh2910))
- Add release process script ([#1473](https://github.com/kubeflow/katib/pull/1473) by [@andreyvelich](https://github.com/andreyvelich))
- Refactor the Katib installation using Kustomize ([#1464](https://github.com/kubeflow/katib/pull/1464) by [@andreyvelich](https://github.com/andreyvelich))

### UI Features and Enhancements

- First step for the Katib new UI implementation ([#1427](https://github.com/kubeflow/katib/pull/1427) by [@kimwnasptd](https://github.com/kimwnasptd))
- Add missing fields to the Katib new UI ([#1463](https://github.com/kubeflow/katib/pull/1463) by [@kimwnasptd](https://github.com/kimwnasptd))
- Add instructions to install the new Katib UI ([#1476](https://github.com/kubeflow/katib/pull/1476) by [@kimwnasptd](https://github.com/kimwnasptd))

### Katib Juju operator

- Add Juju operator support for Katib ([#1403](https://github.com/kubeflow/katib/pull/1403) by [@knkski](https://github.com/knkski))
- Add GitHub Actions for the Juju operator ([#1407](https://github.com/kubeflow/katib/pull/1407) by [@knkski](https://github.com/knkski))
- Add install docs for the Juju operator ([#1411](https://github.com/kubeflow/katib/pull/1411) by [@RFMVasconcelos](https://github.com/RFMVasconcelos))
- Modify ClusterRoles for the Juju operator ([#1426](https://github.com/kubeflow/katib/pull/1426) by [@DomFleischmann](https://github.com/DomFleischmann))
- Update the Juju operator with the new Katib Webhooks ([#1465](https://github.com/kubeflow/katib/pull/1465) by [@knkski](https://github.com/knkski))

## Bug fixes

- Fix compare step for Early Stopping ([#1386](https://github.com/kubeflow/katib/pull/1386) by [@andreyvelich](https://github.com/andreyvelich))
- Fix Early Stopping in the Goptuna Suggestion ([#1404](https://github.com/kubeflow/katib/pull/1404) by [@andreyvelich](https://github.com/andreyvelich))
- Fix SDK examples to work with the Katib 0.10 ([#1402](https://github.com/kubeflow/katib/pull/1402) by [@andreyvelich](https://github.com/andreyvelich))
- Fix links in the TFEvent Metrics Collector ([#1417](https://github.com/kubeflow/katib/pull/1417) by [@zuston](https://github.com/zuston))
- Fix the gRPC build script ([#1492](https://github.com/kubeflow/katib/pull/1492) by [@andreyvelich](https://github.com/andreyvelich))

## Documentation

- Modify docs for the Katib 0.10 ([#1392](https://github.com/kubeflow/katib/pull/1392) by [@andreyvelich](https://github.com/andreyvelich))
- Add Katib presentation list ([#1446](https://github.com/kubeflow/katib/pull/1446) by [@andreyvelich](https://github.com/andreyvelich))
- Add Canonical to the Katib Adopters ([#1401](https://github.com/kubeflow/katib/pull/1401) by [@RFMVasconcelos](https://github.com/RFMVasconcelos))
- Update developer guide with the Katib controller flags ([#1449](https://github.com/kubeflow/katib/pull/1449) by [@annajung](https://github.com/annajung))
- Add Fuzhi to the Katib Adopters ([#1451](https://github.com/kubeflow/katib/pull/1451) by [@Planck0591](https://github.com/Planck0591))
- Fix Katib broken links to the Kubeflow guides ([#1477](https://github.com/kubeflow/katib/pull/1477) by [@theofpa](https://github.com/theofpa))
- Add the Katib Webhook docs ([#1486](https://github.com/kubeflow/katib/pull/1486) by [@andreyvelich](https://github.com/andreyvelich))

## Misc

- Add recreate strategy for the MySQL deployment ([#1393](https://github.com/kubeflow/katib/pull/1393) by [@andreyvelich](https://github.com/andreyvelich))
- Modify worker image for the Katib AWS CI/CD ([#1423](https://github.com/kubeflow/katib/pull/1423) by [@PatrickXYS](https://github.com/PatrickXYS))
- Add the SVG logo for Katib ([#1414](https://github.com/kubeflow/katib/pull/1414) by [@knkski](https://github.com/knkski))
- Verify empty Objective in the Experiment defaults ([#1445](https://github.com/kubeflow/katib/pull/1445) by [@andreyvelich](https://github.com/andreyvelich))
- Move the Katib manifests upstream ([#1432](https://github.com/kubeflow/katib/pull/1432) by [@yanniszark](https://github.com/yanniszark))
- Build the Trial images in the Katib CI ([#1457](https://github.com/kubeflow/katib/pull/1457) by [@andreyvelich](https://github.com/andreyvelich))
- Add script to update the boilerplates ([#1491](https://github.com/kubeflow/katib/pull/1491) by [@andreyvelich](https://github.com/andreyvelich))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.10.1...v0.11.0)

# [v0.10.1](https://github.com/kubeflow/katib/tree/v0.10.1) (2021-03-02)

## Features and Bug Fixes

- add adopter ([#1451](https://github.com/kubeflow/katib/pull/1451) by [@Planck0591](https://github.com/Planck0591))
- Add katib controller flags to developers guide ([#1449](https://github.com/kubeflow/katib/pull/1449) by [@annajung](https://github.com/annajung))
- Enhance katib client by adding get_success_trial_details() ([#1442](https://github.com/kubeflow/katib/pull/1442) by [@Adarsh2910](https://github.com/Adarsh2910))
- Add Katib presentations and community information ([#1446](https://github.com/kubeflow/katib/pull/1446) by [@andreyvelich](https://github.com/andreyvelich))
- Verify nil objective in Experiment defaults ([#1445](https://github.com/kubeflow/katib/pull/1445) by [@andreyvelich](https://github.com/andreyvelich))
- Migrate to Go modules ([#1438](https://github.com/kubeflow/katib/pull/1438) by [@andreyvelich](https://github.com/andreyvelich))
- Change roles to clusterroles for operators ([#1426](https://github.com/kubeflow/katib/pull/1426) by [@DomFleischmann](https://github.com/DomFleischmann))
- Migrate katib to new test-infra ([#1423](https://github.com/kubeflow/katib/pull/1423) by [@PatrickXYS](https://github.com/PatrickXYS))
- Add SVG logo traced from bitmap logo ([#1414](https://github.com/kubeflow/katib/pull/1414) by [@knkski](https://github.com/knkski))
- Invalid example url ([#1417](https://github.com/kubeflow/katib/pull/1417) by [@zuston](https://github.com/zuston))
- Fix SDK examples for 0.10 version ([#1402](https://github.com/kubeflow/katib/pull/1402) by [@andreyvelich](https://github.com/andreyvelich))
- Add Github Actions CI for charm operators ([#1407](https://github.com/kubeflow/katib/pull/1407) by [@knkski](https://github.com/knkski))
- Add Juju install commands to operators README ([#1411](https://github.com/kubeflow/katib/pull/1411) by [@RFMVasconcelos](https://github.com/RFMVasconcelos))
- Fix indentation in the OWNERS file ([#1408](https://github.com/kubeflow/katib/pull/1408) by [@andreyvelich](https://github.com/andreyvelich))
- Bump Prettier to 2.2.0 for the Katib UI ([#1409](https://github.com/kubeflow/katib/pull/1409) by [@andreyvelich](https://github.com/andreyvelich))
- Add Katib Bundle for Juju ([#1403](https://github.com/kubeflow/katib/pull/1403) by [@knkski](https://github.com/knkski))
- Remove duecredit pkg from the Suggestions ([#1406](https://github.com/kubeflow/katib/pull/1406) by [@andreyvelich](https://github.com/andreyvelich))
- Fix Early Stopped Trials in Goptuna Suggestion ([#1404](https://github.com/kubeflow/katib/pull/1404) by [@andreyvelich](https://github.com/andreyvelich))
- Remove v1alpha3 version ([#1396](https://github.com/kubeflow/katib/pull/1396) by [@andreyvelich](https://github.com/andreyvelich))
- Update docs for Katib 0.10 ([#1392](https://github.com/kubeflow/katib/pull/1392) by [@andreyvelich](https://github.com/andreyvelich))
- Adding to ADOPTERS.md ([#1401](https://github.com/kubeflow/katib/pull/1401) by [@RFMVasconcelos](https://github.com/RFMVasconcelos))
- Feature/waitallprocesses config ([#1394](https://github.com/kubeflow/katib/pull/1394) by [@robbertvdg](https://github.com/robbertvdg))
- Add recreate strategy to MySQL deployment ([#1393](https://github.com/kubeflow/katib/pull/1393) by [@andreyvelich](https://github.com/andreyvelich))
- Move Adopters file ([#1391](https://github.com/kubeflow/katib/pull/1391) by [@andreyvelich](https://github.com/andreyvelich))
- Add Stale config to close inactivity issues ([#1390](https://github.com/kubeflow/katib/pull/1390) by [@andreyvelich](https://github.com/andreyvelich))
- Remove new Trial kind doc ([#1388](https://github.com/kubeflow/katib/pull/1388) by [@andreyvelich](https://github.com/andreyvelich))
- Fix compare step for the early stopping ([#1386](https://github.com/kubeflow/katib/pull/1386) by [@andreyvelich](https://github.com/andreyvelich))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.10.0...v0.10.1)

# [v0.10.0](https://github.com/kubeflow/katib/tree/v0.10.0) (2020-11-07)

## New Features

### Core Features

- The new Trial template design ([#1208](https://github.com/kubeflow/katib/issues/1208) by [@andreyvelich](https://github.com/andreyvelich))
- Support custom Kubernetes CRD in the Trial template ([#1214](https://github.com/kubeflow/katib/issues/1214) by [@andreyvelich](https://github.com/andreyvelich))
  - Add example for the [Tekton `Pipeline`](https://github.com/tektoncd/pipeline) ([#1339](https://github.com/kubeflow/katib/pull/1339) by [@andreyvelich](https://github.com/andreyvelich))
  - Add example for the [Kubeflow `MPIJob`](https://github.com/kubeflow/mpi-operator) ([#1342](https://github.com/kubeflow/katib/pull/1342) by [@andreyvelich](https://github.com/andreyvelich))
- Support early stopping with the Median Stopping Rule ([#1344](https://github.com/kubeflow/katib/pull/1344) by [@andreyvelich](https://github.com/andreyvelich))
- Resume Experiment from the volume ([#1275](https://github.com/kubeflow/katib/pull/1275) by [@andreyvelich](https://github.com/andreyvelich))
  - Support volume settings in the Katib config ([#1291](https://github.com/kubeflow/katib/pull/1291) by [@andreyvelich](https://github.com/andreyvelich))
- Extract the Experiment metrics in multiple ways ([#1140](https://github.com/kubeflow/katib/pull/1140) by [@sperlingxx](https://github.com/sperlingxx))
- Update the Python SDK for the v1beta1 version ([#1252](https://github.com/kubeflow/katib/pull/1252) by [@sperlingxx](https://github.com/sperlingxx))

### UI Features and Enhancements

- Show the Trial parameters on the submit Experiment page ([#1224](https://github.com/kubeflow/katib/pull/1224) by [@andreyvelich](https://github.com/andreyvelich))
- Enable to set the Trial template YAML from the submit Experiment page ([#1363](https://github.com/kubeflow/katib/pull/1363) by [@andreyvelich](https://github.com/andreyvelich))
- Optimise the Katib UI image ([#1232](https://github.com/kubeflow/katib/pull/1232) by [@andreyvelich](https://github.com/andreyvelich))
- Enable sorting in the Trial list table ([#1251](https://github.com/kubeflow/katib/pull/1251) by [@andreyvelich](https://github.com/andreyvelich))
- Add pages to the Trial list table ([#1262](https://github.com/kubeflow/katib/pull/1262) by [@andreyvelich](https://github.com/andreyvelich))
- Use the V4 version for the Material UI ([#1254](https://github.com/kubeflow/katib/pull/1254) by [@andreyvelich](https://github.com/andreyvelich))
- Automatically delete an empty ConfigMap with Trial templates ([#1260](https://github.com/kubeflow/katib/pull/1260) by [@andreyvelich](https://github.com/andreyvelich))
- Create a ConfigMap with Trial templates ([#1265](https://github.com/kubeflow/katib/pull/1265) by [@andreyvelich](https://github.com/andreyvelich))
- Support metrics strategies on the submit Experiment page ([#1364](https://github.com/kubeflow/katib/pull/1364) by [@andreyvelich](https://github.com/andreyvelich))
- Add the resume policy to the submit Experiment page ([#1362](https://github.com/kubeflow/katib/pull/1362) by [@andreyvelich](https://github.com/andreyvelich))
- Enable to create an early stopping Experiment from the submit Experiment page ([#1373](https://github.com/kubeflow/katib/pull/1373) by [@andreyvelich](https://github.com/andreyvelich))

## Bug fixes

- Check the Trials count before deleting it ([#1223](https://github.com/kubeflow/katib/pull/1223) by [@gaocegege](https://github.com/gaocegege))
- Check that Trials are deleted ([#1288](https://github.com/kubeflow/katib/pull/1288) by [@andreyvelich](https://github.com/andreyvelich))
- Fix the out of range error in the Hyperopt suggestion ([#1315](https://github.com/kubeflow/katib/pull/1315) by [@andreyvelich](https://github.com/andreyvelich))
- Fix the pod ownership to inject the metrics collector ([#1303](https://github.com/kubeflow/katib/pull/1303) by [@andreyvelich](https://github.com/andreyvelich))

## Misc

- Switch the test infra to the AWS ([#1356](https://github.com/kubeflow/katib/pull/1356) by [@andreyvelich](https://github.com/andreyvelich))
- Use the `docker.io/kubeflowkatib` registry to release images ([#1372](https://github.com/kubeflow/katib/pull/1372) by [@andreyvelich](https://github.com/andreyvelich))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.9.0...v0.10.0)

# [v0.9.0](https://github.com/kubeflow/katib/tree/v0.9.0) (2020-06-10)

## Features and Bug Fixes

- add clientset/lister/informer generation ([#1194](https://github.com/kubeflow/katib/pull/1194) by [@sperlingxx](https://github.com/sperlingxx))
- New Trial Template API controller implementation ([#1202](https://github.com/kubeflow/katib/pull/1202) by [@andreyvelich](https://github.com/andreyvelich))
- Add citation information ([#1210](https://github.com/kubeflow/katib/pull/1210) by [@terrytangyuan](https://github.com/terrytangyuan))
- Python SDK for katib ([#1177](https://github.com/kubeflow/katib/pull/1177) by [@prem0912](https://github.com/prem0912))
- Rename algorithm_setting to algorithm_settings in manager ([#1204](https://github.com/kubeflow/katib/pull/1204) by [@andreyvelich](https://github.com/andreyvelich))
- Update doc for training container images with DARTS ([#1201](https://github.com/kubeflow/katib/pull/1201) by [@andreyvelich](https://github.com/andreyvelich))
- Re: Support string metrics values in Controller ([#1200](https://github.com/kubeflow/katib/pull/1200) by [@andreyvelich](https://github.com/andreyvelich))
- Modify new algorithm service doc ([#1198](https://github.com/kubeflow/katib/pull/1198) by [@andreyvelich](https://github.com/andreyvelich))
- Katib v1beta1 version ([#1197](https://github.com/kubeflow/katib/pull/1197) by [@andreyvelich](https://github.com/andreyvelich))
- Add more algorithm settings to DARTS ([#1195](https://github.com/kubeflow/katib/pull/1195) by [@andreyvelich](https://github.com/andreyvelich))
- Fix additional metrics in TF Event metrics collector ([#1191](https://github.com/kubeflow/katib/pull/1191) by [@andreyvelich](https://github.com/andreyvelich))
- UI: Fix comparison of metric values in Metrics Info Plot ([#1192](https://github.com/kubeflow/katib/pull/1192) by [@andreyvelich](https://github.com/andreyvelich))
- Support one and two NN layers in DARTS ([#1185](https://github.com/kubeflow/katib/pull/1185) by [@andreyvelich](https://github.com/andreyvelich))
- Revert 1176 PR (Support string metric values) ([#1189](https://github.com/kubeflow/katib/pull/1189) by [@andreyvelich](https://github.com/andreyvelich))
- Fix Never Resume Policy for Experiment ([#1184](https://github.com/kubeflow/katib/pull/1184) by [@andreyvelich](https://github.com/andreyvelich))
- Change scikit-learn version to 0.22.0 for BO Suggestion ([#1187](https://github.com/kubeflow/katib/pull/1187) by [@andreyvelich](https://github.com/andreyvelich))
- DARTS documentation ([#1180](https://github.com/kubeflow/katib/pull/1180) by [@andreyvelich](https://github.com/andreyvelich))
- Unittest for DARTS Suggestion ([#1179](https://github.com/kubeflow/katib/pull/1179) by [@andreyvelich](https://github.com/andreyvelich))
- Build image for DARTS Suggestion ([#1178](https://github.com/kubeflow/katib/pull/1178) by [@andreyvelich](https://github.com/andreyvelich))
- DARTS Suggestion ([#1175](https://github.com/kubeflow/katib/pull/1175) by [@andreyvelich](https://github.com/andreyvelich))
- Support string metrics values in Controller ([#1176](https://github.com/kubeflow/katib/pull/1176) by [@andreyvelich](https://github.com/andreyvelich))
- Delete Suggestion deployment after Experiment is finished ([#1150](https://github.com/kubeflow/katib/pull/1150) by [@sperlingxx](https://github.com/sperlingxx))
- Fix Cuda version in training container for ENAS ([#1172](https://github.com/kubeflow/katib/pull/1172) by [@andreyvelich](https://github.com/andreyvelich))
- Rename chocolate algorithm names for consistency ([#1164](https://github.com/kubeflow/katib/pull/1164) by [@c-bata](https://github.com/c-bata))
- restructure algorithm configuration for hyperopt_service ([#1161](https://github.com/kubeflow/katib/pull/1161) by [@sperlingxx](https://github.com/sperlingxx))
- Refactor suggestion services folder structure ([#1166](https://github.com/kubeflow/katib/pull/1166) by [@andreyvelich](https://github.com/andreyvelich))
- UI: Suggestion view from Experiment page ([#1162](https://github.com/kubeflow/katib/pull/1162) by [@andreyvelich](https://github.com/andreyvelich))
- Add support Kubeflow operators to ROADMAP ([#1145](https://github.com/kubeflow/katib/pull/1145) by [@andreyvelich](https://github.com/andreyvelich))
- Remove Suggestion Request from Update Suggestion ([#1158](https://github.com/kubeflow/katib/pull/1158) by [@andreyvelich](https://github.com/andreyvelich))
- E2E test for CMA-ES Suggestion ([#1157](https://github.com/kubeflow/katib/pull/1157) by [@andreyvelich](https://github.com/andreyvelich))
- Build Goptuna Suggestion image in CI ([#1154](https://github.com/kubeflow/katib/pull/1154) by [@andreyvelich](https://github.com/andreyvelich))
- Add an example of Goptuna suggestion service ([#1155](https://github.com/kubeflow/katib/pull/1155) by [@c-bata](https://github.com/c-bata))
- ENAS enable to add None values in algorithm settings ([#1153](https://github.com/kubeflow/katib/pull/1153) by [@andreyvelich](https://github.com/andreyvelich))
- Support Categorical string values in Chocolate Suggestion ([#1149](https://github.com/kubeflow/katib/pull/1149) by [@andreyvelich](https://github.com/andreyvelich))
- katib ui: adapt environment in which cluster role is unavailable ([#1141](https://github.com/kubeflow/katib/pull/1141) by [@sperlingxx](https://github.com/sperlingxx))
- Add Goptuna based suggestion service for CMA-ES. ([#1131](https://github.com/kubeflow/katib/pull/1131) by [@c-bata](https://github.com/c-bata))
- ENAS Check Algorithm Settings in Validate Function ([#1146](https://github.com/kubeflow/katib/pull/1146) by [@andreyvelich](https://github.com/andreyvelich))
- Change folder structure for NAS algorithms, rename NASRL to ENAS ([#1143](https://github.com/kubeflow/katib/pull/1143) by [@andreyvelich](https://github.com/andreyvelich))
- Update ENAS Algorithm Settings in Katib UI ([#1142](https://github.com/kubeflow/katib/pull/1142) by [@andreyvelich](https://github.com/andreyvelich))
- Refactor NAS RL Suggestion ([#1134](https://github.com/kubeflow/katib/pull/1134) by [@andreyvelich](https://github.com/andreyvelich))
- Fix duplicated imports ([#1133](https://github.com/kubeflow/katib/pull/1133) by [@c-bata](https://github.com/c-bata))
- Remove Anneal from supported algorithms ([#1139](https://github.com/kubeflow/katib/pull/1139) by [@c-bata](https://github.com/c-bata))
- Refactor file-metricscollector ([#1137](https://github.com/kubeflow/katib/pull/1137) by [@c-bata](https://github.com/c-bata))
- Fix typo in suggestion packages ([#1138](https://github.com/kubeflow/katib/pull/1138) by [@c-bata](https://github.com/c-bata))
- Bump up the Go version to 1.14.2 at Travis CI ([#1132](https://github.com/kubeflow/katib/pull/1132) by [@c-bata](https://github.com/c-bata))
- Fix NotImplementedError for TPE and Random suggestion. ([#1130](https://github.com/kubeflow/katib/pull/1130) by [@c-bata](https://github.com/c-bata))
- Add ENAS enhancements to ROADMAP ([#1129](https://github.com/kubeflow/katib/pull/1129) by [@andreyvelich](https://github.com/andreyvelich))
- feat: Add 2020 roadmap ([#1121](https://github.com/kubeflow/katib/pull/1121) by [@gaocegege](https://github.com/gaocegege))
- Optimize Chocolate Suggestion ([#1116](https://github.com/kubeflow/katib/pull/1116) by [@andreyvelich](https://github.com/andreyvelich))
- Support step for int parameter in Chocolate and Hyperopt Suggestion ([#1123](https://github.com/kubeflow/katib/pull/1123) by [@andreyvelich](https://github.com/andreyvelich))
- UI: Remove antd package ([#1117](https://github.com/kubeflow/katib/pull/1117) by [@andreyvelich](https://github.com/andreyvelich))
- Fix Katib e2e tests ([#1118](https://github.com/kubeflow/katib/pull/1118) by [@andreyvelich](https://github.com/andreyvelich))
- UI: Close menu on click ([#1114](https://github.com/kubeflow/katib/pull/1114) by [@andreyvelich](https://github.com/andreyvelich))
- UI: Change style for submit Experiment from YAML ([#1113](https://github.com/kubeflow/katib/pull/1113) by [@andreyvelich](https://github.com/andreyvelich))
- Refactor python files in cmd/suggestion ([#1105](https://github.com/kubeflow/katib/pull/1105) by [@c-bata](https://github.com/c-bata))
- Update prow config with the latest folders ([#1109](https://github.com/kubeflow/katib/pull/1109) by [@andreyvelich](https://github.com/andreyvelich))
- Fix logger namespace ([#1108](https://github.com/kubeflow/katib/pull/1108) by [@c-bata](https://github.com/c-bata))
- chore(deps): Bump tensorflow from 1.14.0 to 1.15.2 in /cmd/suggestion/nasrl/v1alpha3 ([#1035](https://github.com/kubeflow/katib/pull/1035) by [@dependabot[bot]](https://github.com/apps/dependabot))
- Refactor suggestion-internal-modules ([#1106](https://github.com/kubeflow/katib/pull/1106) by [@c-bata](https://github.com/c-bata))
- chore(deps): Bump psutil from 5.2.2 to 5.6.6 in /cmd/metricscollector/v1alpha3/tfevent-metricscollector ([#1085](https://github.com/kubeflow/katib/pull/1085) by [@dependabot[bot]](https://github.com/apps/dependabot))
- Fix custom Katib DB Manager env variables ([#1102](https://github.com/kubeflow/katib/pull/1102) by [@andreyvelich](https://github.com/andreyvelich))
- Refactor python files of suggestion services ([#1107](https://github.com/kubeflow/katib/pull/1107) by [@c-bata](https://github.com/c-bata))
- Add myself to approvers ([#1103](https://github.com/kubeflow/katib/pull/1103) by [@andreyvelich](https://github.com/andreyvelich))
- Enable to add Service Account Name in Katib config ([#1092](https://github.com/kubeflow/katib/pull/1092) by [@andreyvelich](https://github.com/andreyvelich))
- chore(deps): Bump tensorflow-gpu from 1.15.0 to 1.15.2 in /examples/v1alpha3/NAS-training-containers/RL-cifar10 ([#1034](https://github.com/kubeflow/katib/pull/1034) by [@dependabot[bot]](https://github.com/apps/dependabot))
- chore(deps): Bump tensorflow from 1.15.0 to 1.15.2 in /examples/v1alpha3/NAS-training-containers/RL-cifar10 ([#1036](https://github.com/kubeflow/katib/pull/1036) by [@dependabot[bot]](https://github.com/apps/dependabot))
- Add ghalton package to Chocolate Suggestion ([#1101](https://github.com/kubeflow/katib/pull/1101) by [@andreyvelich](https://github.com/andreyvelich))
- Enable to run Experiment without Goal ([#1065](https://github.com/kubeflow/katib/pull/1065) by [@andreyvelich](https://github.com/andreyvelich))
- UI: Support Trial Templates in all namespaces and all configMaps ([#1083](https://github.com/kubeflow/katib/pull/1083) by [@andreyvelich](https://github.com/andreyvelich))
- Fix Chocolate mocmaes algorithm name in Suggestion ([#1097](https://github.com/kubeflow/katib/pull/1097) by [@andreyvelich](https://github.com/andreyvelich))
- UI: Add Metrics Collector Spec to Submit Experiment ([#1096](https://github.com/kubeflow/katib/pull/1096) by [@andreyvelich](https://github.com/andreyvelich))
- UI: Add Experiment view for NAS Jobs ([#1077](https://github.com/kubeflow/katib/pull/1077) by [@andreyvelich](https://github.com/andreyvelich))
- Enable Prettier code formatting for the Katib UI ([#1078](https://github.com/kubeflow/katib/pull/1078) by [@andreyvelich](https://github.com/andreyvelich))
- Adding Karrot as adopter ([#1074](https://github.com/kubeflow/katib/pull/1074) by [@rky0930](https://github.com/rky0930))
- fix annotations ([#1072](https://github.com/kubeflow/katib/pull/1072) by [@sperlingxx](https://github.com/sperlingxx))
- Add more unit tests in Katib ([#1071](https://github.com/kubeflow/katib/pull/1071) by [@andreyvelich](https://github.com/andreyvelich))
- dynamic jobProvider and suggestionComposer registration ([#1069](https://github.com/kubeflow/katib/pull/1069) by [@sperlingxx](https://github.com/sperlingxx))
- UI: Update supported algorithms ([#1070](https://github.com/kubeflow/katib/pull/1070) by [@andreyvelich](https://github.com/andreyvelich))
- Fix TPE Suggestion ([#1063](https://github.com/kubeflow/katib/pull/1063) by [@andreyvelich](https://github.com/andreyvelich))
- Update Katib docs ([#1066](https://github.com/kubeflow/katib/pull/1066) by [@andreyvelich](https://github.com/andreyvelich))
- UI: Show best metrics in the Trial metrics information plot ([#1068](https://github.com/kubeflow/katib/pull/1068) by [@andreyvelich](https://github.com/andreyvelich))
- Update adopters ([#1064](https://github.com/kubeflow/katib/pull/1064) by [@janvdvegt](https://github.com/janvdvegt))
- Create Optimizer in BO Suggestion only for the first run ([#1057](https://github.com/kubeflow/katib/pull/1057) by [@andreyvelich](https://github.com/andreyvelich))
- Add missing GRPC health probe for arm64 to db-manager ([#1059](https://github.com/kubeflow/katib/pull/1059) by [@MrXinWang](https://github.com/MrXinWang))
- Change tell method for BO Suggestion ([#1055](https://github.com/kubeflow/katib/pull/1055) by [@andreyvelich](https://github.com/andreyvelich))
- MXNet -> Apache MXNet ([#1056](https://github.com/kubeflow/katib/pull/1056) by [@terrytangyuan](https://github.com/terrytangyuan))
- Adding error propagation for K8s client creation in KatibClient ([#1053](https://github.com/kubeflow/katib/pull/1053) by [@akirillov](https://github.com/akirillov))
- openAPI generation for katib resources ([#1054](https://github.com/kubeflow/katib/pull/1054) by [@sperlingxx](https://github.com/sperlingxx))
- Disable istio sidecar injection in Suggestion and Training Jobs ([#1050](https://github.com/kubeflow/katib/pull/1050) by [@andreyvelich](https://github.com/andreyvelich))
- UI: The best metrics in Trial table ([#1048](https://github.com/kubeflow/katib/pull/1048) by [@andreyvelich](https://github.com/andreyvelich))
- UI: Fix problem with equal time in different Trial metrics ([#1047](https://github.com/kubeflow/katib/pull/1047) by [@andreyvelich](https://github.com/andreyvelich))
- Adding Babylon Health as adopter ([#1046](https://github.com/kubeflow/katib/pull/1046) by [@jeremievallee](https://github.com/jeremievallee))
- Update adopter ([#1038](https://github.com/kubeflow/katib/pull/1038) by [@ywskycn](https://github.com/ywskycn))
- UI: Add Trial Status to HP Job Table ([#1032](https://github.com/kubeflow/katib/pull/1032) by [@andreyvelich](https://github.com/andreyvelich))
- UI: Experiment view in the Dialog ([#1031](https://github.com/kubeflow/katib/pull/1031) by [@andreyvelich](https://github.com/andreyvelich))
- record TrialNames group by stages in ExperimentStatus ([#1023](https://github.com/kubeflow/katib/pull/1023) by [@sperlingxx](https://github.com/sperlingxx))
- chore: Update docs ([#1024](https://github.com/kubeflow/katib/pull/1024) by [@gaocegege](https://github.com/gaocegege))
- fix: Ignore trials without metrics ([#1028](https://github.com/kubeflow/katib/pull/1028) by [@gaocegege](https://github.com/gaocegege))
- UI: Fix Trial Metric in the Plot ([#1027](https://github.com/kubeflow/katib/pull/1027) by [@andreyvelich](https://github.com/andreyvelich))
- feat: Add a flag to support removing GRPC probe at runtime ([#1020](https://github.com/kubeflow/katib/pull/1020) by [@gaocegege](https://github.com/gaocegege))
- Adding cisco in Katib adopters ([#1026](https://github.com/kubeflow/katib/pull/1026) by [@johnugeorge](https://github.com/johnugeorge))
- add antfin into list of adoptors ([#1025](https://github.com/kubeflow/katib/pull/1025) by [@sperlingxx](https://github.com/sperlingxx))
- Updated links and instructions for Katib API docs ([#1022](https://github.com/kubeflow/katib/pull/1022) by [@sarahmaddox](https://github.com/sarahmaddox))
- feat: Add adopters ([#1019](https://github.com/kubeflow/katib/pull/1019) by [@gaocegege](https://github.com/gaocegege))
- [FileMetricsCollector]skip line without metrics keywords ([#1018](https://github.com/kubeflow/katib/pull/1018) by [@sperlingxx](https://github.com/sperlingxx))
- Added version number and TODO descriptions to API proto ([#1017](https://github.com/kubeflow/katib/pull/1017) by [@sarahmaddox](https://github.com/sarahmaddox))
- fix: First check failed condition ([#1015](https://github.com/kubeflow/katib/pull/1015) by [@gaocegege](https://github.com/gaocegege))
- feat: Do not inject sh -c when it exists ([#1010](https://github.com/kubeflow/katib/pull/1010) by [@gaocegege](https://github.com/gaocegege))
- Nerual -> Neural ([#1000](https://github.com/kubeflow/katib/pull/1000) by [@tmielika](https://github.com/tmielika))
- [Feature] Enable imagePullPolicy in Katib Config ([#1013](https://github.com/kubeflow/katib/pull/1013) by [@andreyvelich](https://github.com/andreyvelich))
- fix: Avoid out-of-range exception ([#1012](https://github.com/kubeflow/katib/pull/1012) by [@gaocegege](https://github.com/gaocegege))
- E2E Test for NAS RL Suggestion ([#1011](https://github.com/kubeflow/katib/pull/1011) by [@andreyvelich](https://github.com/andreyvelich))
- Example with collecting timestamp of the metrics ([#970](https://github.com/kubeflow/katib/pull/970) by [@andreyvelich](https://github.com/andreyvelich))
- Add NAS RL training container to kubeflowkatib repository ([#1008](https://github.com/kubeflow/katib/pull/1008) by [@andreyvelich](https://github.com/andreyvelich))
- Fix number of Trials problem in NAS RL Suggestion ([#1009](https://github.com/kubeflow/katib/pull/1009) by [@andreyvelich](https://github.com/andreyvelich))
- Rename katib DB manager ([#1006](https://github.com/kubeflow/katib/pull/1006) by [@hougangliu](https://github.com/hougangliu))
- chore(deps): Bump tensorflow from 1.12.0 to 1.15.0 in /examples/v1alpha3/NAS-training-containers/RL-cifar10 ([#1005](https://github.com/kubeflow/katib/pull/1005) by [@dependabot[bot]](https://github.com/apps/dependabot))
- chore(deps): Bump tensorflow-gpu from 1.12.0 to 1.15.0 in /examples/v1alpha3/NAS-training-containers/RL-cifar10 ([#978](https://github.com/kubeflow/katib/pull/978) by [@dependabot[bot]](https://github.com/apps/dependabot))
- CPU example for NAS RL cifar10 training container ([#999](https://github.com/kubeflow/katib/pull/999) by [@andreyvelich](https://github.com/andreyvelich))
- Updated links to docs/github on Katib dashboard ([#1003](https://github.com/kubeflow/katib/pull/1003) by [@sarahmaddox](https://github.com/sarahmaddox))
- Fixed a few typos ([#1001](https://github.com/kubeflow/katib/pull/1001) by [@sarahmaddox](https://github.com/sarahmaddox))
- fix: Inherit labels and annotations from experiment ([#998](https://github.com/kubeflow/katib/pull/998) by [@gaocegege](https://github.com/gaocegege))
- Moved some content and added links to Kubeflow docs ([#990](https://github.com/kubeflow/katib/pull/990) by [@sarahmaddox](https://github.com/sarahmaddox))
- feat: Support resource in sidecar ([#991](https://github.com/kubeflow/katib/pull/991) by [@gaocegege](https://github.com/gaocegege))
- fix: Ignore the failure ([#996](https://github.com/kubeflow/katib/pull/996) by [@gaocegege](https://github.com/gaocegege))
- UI: Select namespace from Kubeflow dashboard ([#982](https://github.com/kubeflow/katib/pull/982) by [@andreyvelich](https://github.com/andreyvelich))
- feat: Add a flag to control the logic about sc ([#994](https://github.com/kubeflow/katib/pull/994) by [@gaocegege](https://github.com/gaocegege))
- Initialize securityContext in injected metrics container ([#964](https://github.com/kubeflow/katib/pull/964) by [@vpavlin](https://github.com/vpavlin))
- add disk setting into suggestionConfiguration ([#989](https://github.com/kubeflow/katib/pull/989) by [@sperlingxx](https://github.com/sperlingxx))
- Get dbUser from Env or default('root') ([#985](https://github.com/kubeflow/katib/pull/985) by [@UrmsOne](https://github.com/UrmsOne))
- feat(experiment_status): Add trial name ([#986](https://github.com/kubeflow/katib/pull/986) by [@gaocegege](https://github.com/gaocegege))
- feat(config): Add a new config for webhook ([#980](https://github.com/kubeflow/katib/pull/980) by [@gaocegege](https://github.com/gaocegege))
- add metrics for trial ([#974](https://github.com/kubeflow/katib/pull/974) by [@yeya24](https://github.com/yeya24))
- Use port higher than 1024 to be able to run as a non-root user ([#960](https://github.com/kubeflow/katib/pull/960) by [@vpavlin](https://github.com/vpavlin))
- Remove redundant serviceAccountName assignment ([#969](https://github.com/kubeflow/katib/pull/969) by [@hougangliu](https://github.com/hougangliu))
- Increase Suggestion memory limit ([#958](https://github.com/kubeflow/katib/pull/958) by [@andreyvelich](https://github.com/andreyvelich))
- User root user explicitely for DB readinessProbe ([#962](https://github.com/kubeflow/katib/pull/962) by [@vpavlin](https://github.com/vpavlin))
- Fix typo in getKabitJob function name ([#965](https://github.com/kubeflow/katib/pull/965) by [@vpavlin](https://github.com/vpavlin))
- Use port 8080 for Katib UI ([#967](https://github.com/kubeflow/katib/pull/967) by [@vpavlin](https://github.com/vpavlin))
- Validate experiment ([#957](https://github.com/kubeflow/katib/pull/957) by [@hougangliu](https://github.com/hougangliu))
- UI: Support namespace selection in experiment monitor ([#950](https://github.com/kubeflow/katib/pull/950) by [@andreyvelich](https://github.com/andreyvelich))
- Delete v1alpha2 api ([#953](https://github.com/kubeflow/katib/pull/953) by [@johnugeorge](https://github.com/johnugeorge))
- Resume experiment with extra trials from last checkpoint ([#952](https://github.com/kubeflow/katib/pull/952) by [@johnugeorge](https://github.com/johnugeorge))
- Add a gauge metric for current experiments ([#954](https://github.com/kubeflow/katib/pull/954) by [@yeya24](https://github.com/yeya24))
- feat: Support running ([#894](https://github.com/kubeflow/katib/pull/894) by [@gaocegege](https://github.com/gaocegege))
- Use kubeflowkatib repo as image repo of example ([#949](https://github.com/kubeflow/katib/pull/949) by [@hougangliu](https://github.com/hougangliu))
- Update API spec for early stopping ([#951](https://github.com/kubeflow/katib/pull/951) by [@richardsliu](https://github.com/richardsliu))
- rename counter metrics ([#942](https://github.com/kubeflow/katib/pull/942) by [@yeya24](https://github.com/yeya24))
- update deployment api version ([#937](https://github.com/kubeflow/katib/pull/937) by [@yeya24](https://github.com/yeya24))
- Fix: Empty Trial templates in Katib UI ([#938](https://github.com/kubeflow/katib/pull/938) by [@andreyvelich](https://github.com/andreyvelich))
- Implement metrics custom filters ([#947](https://github.com/kubeflow/katib/pull/947) by [@hougangliu](https://github.com/hougangliu))
- Remove katib webhook when undeploy ([#935](https://github.com/kubeflow/katib/pull/935) by [@hougangliu](https://github.com/hougangliu))
- Change web failPolicy to fail instead of default ingore ([#933](https://github.com/kubeflow/katib/pull/933) by [@hougangliu](https://github.com/hougangliu))
- feat: Add limit for suggestion pod ([#932](https://github.com/kubeflow/katib/pull/932) by [@gaocegege](https://github.com/gaocegege))
- Support multiple metric logs in one line ([#925](https://github.com/kubeflow/katib/pull/925) by [@hougangliu](https://github.com/hougangliu))
- Tfevent metriccollector fails when multiple files exist ([#920](https://github.com/kubeflow/katib/pull/920) by [@hougangliu](https://github.com/hougangliu))
- Handle metricscollector case worker container have no command ([#914](https://github.com/kubeflow/katib/pull/914) by [@hougangliu](https://github.com/hougangliu))
- tfevent-metricscollector support ppc64le ([#912](https://github.com/kubeflow/katib/pull/912) by [@hmtai](https://github.com/hmtai))
- Fix grid suggestion ValidateAlgorithmSettings return ([#913](https://github.com/kubeflow/katib/pull/913) by [@hougangliu](https://github.com/hougangliu))
- Fix wrong suggestion service endpoint ([#911](https://github.com/kubeflow/katib/pull/911) by [@hougangliu](https://github.com/hougangliu))
- Enable arm64 architecture support for katib images and fix grpc health probe multiarch error. ([#897](https://github.com/kubeflow/katib/pull/897) by [@MrXinWang](https://github.com/MrXinWang))
- feat: Support custom databases ([#910](https://github.com/kubeflow/katib/pull/910) by [@gaocegege](https://github.com/gaocegege))
- Enhance validation for metrics collector ([#909](https://github.com/kubeflow/katib/pull/909) by [@hougangliu](https://github.com/hougangliu))
- Support custom metrics collector kind ([#908](https://github.com/kubeflow/katib/pull/908) by [@hougangliu](https://github.com/hougangliu))
- support ppc64le ([#893](https://github.com/kubeflow/katib/pull/893) by [@hmtai](https://github.com/hmtai))
- fix: Add Suggestion into CI ([#907](https://github.com/kubeflow/katib/pull/907) by [@gaocegege](https://github.com/gaocegege))
- Validate algorithm ([#904](https://github.com/kubeflow/katib/pull/904) by [@hougangliu](https://github.com/hougangliu))
- Support restarting training job ([#901](https://github.com/kubeflow/katib/pull/901) by [@hougangliu](https://github.com/hougangliu))
- Fix katib-manager crash in kubeflow cluster ([#900](https://github.com/kubeflow/katib/pull/900) by [@hougangliu](https://github.com/hougangliu))
- Revert env for katib-db ([#899](https://github.com/kubeflow/katib/pull/899) by [@hougangliu](https://github.com/hougangliu))
- feat: Patch to fix running condition ([#895](https://github.com/kubeflow/katib/pull/895) by [@gaocegege](https://github.com/gaocegege))
- feat: Add quick start ([#878](https://github.com/kubeflow/katib/pull/878) by [@gaocegege](https://github.com/gaocegege))
- Pin operators to 0.7 branch ([#885](https://github.com/kubeflow/katib/pull/885) by [@johnugeorge](https://github.com/johnugeorge))
- fix: Use 64 instead of 32 since we are using float64 ([#883](https://github.com/kubeflow/katib/pull/883) by [@gaocegege](https://github.com/gaocegege))
- fix: Use as instead of , to support python 3 in tfevent metrics collector ([#881](https://github.com/kubeflow/katib/pull/881) by [@gaocegege](https://github.com/gaocegege))
- feat: Add event when the reconcile is failed ([#879](https://github.com/kubeflow/katib/pull/879) by [@gaocegege](https://github.com/gaocegege))
- feat: Add events in experiment ([#880](https://github.com/kubeflow/katib/pull/880) by [@gaocegege](https://github.com/gaocegege))
- Remove unused katib-manager-rest ([#876](https://github.com/kubeflow/katib/pull/876) by [@hougangliu](https://github.com/hougangliu))
- feat: Refactor to make it easy to extend new kinds ([#865](https://github.com/kubeflow/katib/pull/865) by [@gaocegege](https://github.com/gaocegege))
- feat: Support random state in random search ([#873](https://github.com/kubeflow/katib/pull/873) by [@gaocegege](https://github.com/gaocegege))
- Add prometheus metrics for experiment and trial ([#870](https://github.com/kubeflow/katib/pull/870) by [@hougangliu](https://github.com/hougangliu))
- fix: Use binary in test ([#875](https://github.com/kubeflow/katib/pull/875) by [@gaocegege](https://github.com/gaocegege))
- feat: Support env in mysql ([#868](https://github.com/kubeflow/katib/pull/868) by [@gaocegege](https://github.com/gaocegege))
- feat: Add liveness probe for DB ([#871](https://github.com/kubeflow/katib/pull/871) by [@gaocegege](https://github.com/gaocegege))
- Remove unused files ([#869](https://github.com/kubeflow/katib/pull/869) by [@hougangliu](https://github.com/hougangliu))
- feat: Add doc about algorithm ([#867](https://github.com/kubeflow/katib/pull/867) by [@gaocegege](https://github.com/gaocegege))
- feat: Add doc about how to add a new kind in trial ([#844](https://github.com/kubeflow/katib/pull/844) by [@gaocegege](https://github.com/gaocegege))
- Adding metric unavailability to events ([#864](https://github.com/kubeflow/katib/pull/864) by [@johnugeorge](https://github.com/johnugeorge))
- Fix worker error silent ([#863](https://github.com/kubeflow/katib/pull/863) by [@hougangliu](https://github.com/hougangliu))
- feat: Show experiment status in json ([#853](https://github.com/kubeflow/katib/pull/853) by [@gaocegege](https://github.com/gaocegege))
- Finish reconcile only after running trials are complete ([#861](https://github.com/kubeflow/katib/pull/861) by [@johnugeorge](https://github.com/johnugeorge))
- Update Readme ([#860](https://github.com/kubeflow/katib/pull/860) by [@johnugeorge](https://github.com/johnugeorge))
- fix: Fix docs about metrics collection and suggestion design ([#858](https://github.com/kubeflow/katib/pull/858) by [@gaocegege](https://github.com/gaocegege))
- Adding events to trials ([#852](https://github.com/kubeflow/katib/pull/852) by [@johnugeorge](https://github.com/johnugeorge))
- chore: Add dockerignore, enhance liveness for manager ([#851](https://github.com/kubeflow/katib/pull/851) by [@gaocegege](https://github.com/gaocegege))
- fix: Reorder to skip observation collection ([#847](https://github.com/kubeflow/katib/pull/847) by [@gaocegege](https://github.com/gaocegege))
- feat: Set default namespace and template for trial ([#850](https://github.com/kubeflow/katib/pull/850) by [@gaocegege](https://github.com/gaocegege))
- fix: Use namespace to get trial list ([#846](https://github.com/kubeflow/katib/pull/846) by [@gaocegege](https://github.com/gaocegege))
- [docs] Add suggestion proposal ([#726](https://github.com/kubeflow/katib/pull/726) by [@gaocegege](https://github.com/gaocegege))
- feat: Add doc for implementing new algorithms ([#769](https://github.com/kubeflow/katib/pull/769) by [@gaocegege](https://github.com/gaocegege))
- feat: Support namespace in NAS UI ([#839](https://github.com/kubeflow/katib/pull/839) by [@gaocegege](https://github.com/gaocegege))
- feat: Show all experiments in monitor ([#835](https://github.com/kubeflow/katib/pull/835) by [@gaocegege](https://github.com/gaocegege))
- Delete jobs when trials are completed ([#838](https://github.com/kubeflow/katib/pull/838) by [@johnugeorge](https://github.com/johnugeorge))
- Remove unused manager message definition ([#837](https://github.com/kubeflow/katib/pull/837) by [@hougangliu](https://github.com/hougangliu))
- Add tfjob and pytorch examples to e2e ([#820](https://github.com/kubeflow/katib/pull/820) by [@johnugeorge](https://github.com/johnugeorge))
- fix: Update liveness probe to avoid problems ([#833](https://github.com/kubeflow/katib/pull/833) by [@gaocegege](https://github.com/gaocegege))
- Remove used katib-manager code ([#836](https://github.com/kubeflow/katib/pull/836) by [@hougangliu](https://github.com/hougangliu))
- File metrics collector end to end test ([#832](https://github.com/kubeflow/katib/pull/832) by [@hougangliu](https://github.com/hougangliu))
- feat: support namespace for trial template ([#827](https://github.com/kubeflow/katib/pull/827) by [@gaocegege](https://github.com/gaocegege))
- Remove metrics in DB when delete trial ([#830](https://github.com/kubeflow/katib/pull/830) by [@hougangliu](https://github.com/hougangliu))
- Update status conditions during reconcile error ([#831](https://github.com/kubeflow/katib/pull/831) by [@johnugeorge](https://github.com/johnugeorge))
- feat: Use env var for namespace ([#829](https://github.com/kubeflow/katib/pull/829) by [@gaocegege](https://github.com/gaocegege))
- Make sure experiment namespace can inject metriccollector sidecar ([#828](https://github.com/kubeflow/katib/pull/828) by [@hougangliu](https://github.com/hougangliu))
- Doc about katib workflow design ([#824](https://github.com/kubeflow/katib/pull/824) by [@hougangliu](https://github.com/hougangliu))
- fix: Support multiple namespaces when using kubectl ([#826](https://github.com/kubeflow/katib/pull/826) by [@gaocegege](https://github.com/gaocegege))
- feat: Support step when using grid in UI ([#821](https://github.com/kubeflow/katib/pull/821) by [@gaocegege](https://github.com/gaocegege))
- fix: Build e2e-runner ([#822](https://github.com/kubeflow/katib/pull/822) by [@gaocegege](https://github.com/gaocegege))
- Fix stdout of worker container show nothing ([#819](https://github.com/kubeflow/katib/pull/819) by [@hougangliu](https://github.com/hougangliu))
- feat: Remove useless APIs ([#818](https://github.com/kubeflow/katib/pull/818) by [@gaocegege](https://github.com/gaocegege))
- feat: Add validation for grid ([#812](https://github.com/kubeflow/katib/pull/812) by [@gaocegege](https://github.com/gaocegege))
- Adding additional printer columns for better debugging ([#817](https://github.com/kubeflow/katib/pull/817) by [@johnugeorge](https://github.com/johnugeorge))
- metrics-collector role is not useful any more ([#816](https://github.com/kubeflow/katib/pull/816) by [@hougangliu](https://github.com/hougangliu))
- Rename algorithm deployment and service ([#814](https://github.com/kubeflow/katib/pull/814) by [@hougangliu](https://github.com/hougangliu))
- fix: Fix the type ([#813](https://github.com/kubeflow/katib/pull/813) by [@gaocegege](https://github.com/gaocegege))
- feat: Add tpe e2e test case ([#809](https://github.com/kubeflow/katib/pull/809) by [@gaocegege](https://github.com/gaocegege))
- Remove unused field from Experiment Spec ([#806](https://github.com/kubeflow/katib/pull/806) by [@johnugeorge](https://github.com/johnugeorge))
- feat: Add HyperBand ([#787](https://github.com/kubeflow/katib/pull/787) by [@gaocegege](https://github.com/gaocegege))
- Removing unnecessary config from examples ([#803](https://github.com/kubeflow/katib/pull/803) by [@johnugeorge](https://github.com/johnugeorge))
- feat: Add NAS RL based algorithm ([#793](https://github.com/kubeflow/katib/pull/793) by [@gaocegege](https://github.com/gaocegege))
- fix: Remove copy ([#802](https://github.com/kubeflow/katib/pull/802) by [@gaocegege](https://github.com/gaocegege))
- Using example as the default trial ([#801](https://github.com/kubeflow/katib/pull/801) by [@johnugeorge](https://github.com/johnugeorge))
- Removing metric collector templates from UI ([#800](https://github.com/kubeflow/katib/pull/800) by [@johnugeorge](https://github.com/johnugeorge))
- fix: Use commitid ([#799](https://github.com/kubeflow/katib/pull/799) by [@gaocegege](https://github.com/gaocegege))
- Use common metricsCollector struct ([#798](https://github.com/kubeflow/katib/pull/798) by [@hougangliu](https://github.com/hougangliu))
- build: Support arguments ([#795](https://github.com/kubeflow/katib/pull/795) by [@gaocegege](https://github.com/gaocegege))
- feat: Rename algorithms ([#794](https://github.com/kubeflow/katib/pull/794) by [@gaocegege](https://github.com/gaocegege))
- feat: Add events in suggestion ([#796](https://github.com/kubeflow/katib/pull/796) by [@gaocegege](https://github.com/gaocegege))
- UI: Fix problems ([#786](https://github.com/kubeflow/katib/pull/786) by [@gaocegege](https://github.com/gaocegege))
- Implement tfevent collector ([#792](https://github.com/kubeflow/katib/pull/792) by [@hougangliu](https://github.com/hougangliu))
- Run e2e tests parallel ([#790](https://github.com/kubeflow/katib/pull/790) by [@johnugeorge](https://github.com/johnugeorge))
- Mark trial as failed when job fails ([#791](https://github.com/kubeflow/katib/pull/791) by [@johnugeorge](https://github.com/johnugeorge))
- Adding javascripts locally ([#789](https://github.com/kubeflow/katib/pull/789) by [@johnugeorge](https://github.com/johnugeorge))
- feat: Add grid with the help of chocolate ([#780](https://github.com/kubeflow/katib/pull/780) by [@gaocegege](https://github.com/gaocegege))
- feat: Add bayesian ([#777](https://github.com/kubeflow/katib/pull/777) by [@gaocegege](https://github.com/gaocegege))
- Implement file metrics collector ([#783](https://github.com/kubeflow/katib/pull/783) by [@hougangliu](https://github.com/hougangliu))
- feat: Remove useless algorithms ([#782](https://github.com/kubeflow/katib/pull/782) by [@gaocegege](https://github.com/gaocegege))
- Adding algorithm deployment status to Suggestion status ([#784](https://github.com/kubeflow/katib/pull/784) by [@johnugeorge](https://github.com/johnugeorge))
- Wait for GRPC server to be up ([#785](https://github.com/kubeflow/katib/pull/785) by [@johnugeorge](https://github.com/johnugeorge))
- feat: Add GRPC health check in suggestions ([#779](https://github.com/kubeflow/katib/pull/779) by [@gaocegege](https://github.com/gaocegege))
- feat: Add more output in e2e test for debug purpose and fix test cases ([#775](https://github.com/kubeflow/katib/pull/775) by [@gaocegege](https://github.com/gaocegege))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.6.0-rc.0...v0.9.0)

# [v0.6.0-rc.0](https://github.com/kubeflow/katib/tree/v0.6.0-rc.0) (2019-06-28)

## Features and Bug Fixes

- Add npm build to the UI Dockerfile ([#665](https://github.com/kubeflow/katib/pull/665) by [@andreyvelich](https://github.com/andreyvelich))
- MetricController: Run only a single job per task ([#660](https://github.com/kubeflow/katib/pull/660) by [@epa095](https://github.com/epa095))
- Build images for nasrl training container ([#669](https://github.com/kubeflow/katib/pull/669) by [@andreyvelich](https://github.com/andreyvelich))
- UI: Add delete experiment functionality ([#654](https://github.com/kubeflow/katib/pull/654) by [@andreyvelich](https://github.com/andreyvelich))
- UI: Change adding a template ([#656](https://github.com/kubeflow/katib/pull/656) by [@andreyvelich](https://github.com/andreyvelich))
- UI: Select Objective Type from the list ([#653](https://github.com/kubeflow/katib/pull/653) by [@andreyvelich](https://github.com/andreyvelich))
- Add e2e test to presubmit ([#652](https://github.com/kubeflow/katib/pull/652) by [@johnugeorge](https://github.com/johnugeorge))
- fix: Do not use webhook in UT ([#657](https://github.com/kubeflow/katib/pull/657) by [@gaocegege](https://github.com/gaocegege))
- Enhancing katib client apis ([#650](https://github.com/kubeflow/katib/pull/650) by [@johnugeorge](https://github.com/johnugeorge))
- Wrong mock file name ([#651](https://github.com/kubeflow/katib/pull/651) by [@johnugeorge](https://github.com/johnugeorge))
- UI: Show only succeeded Trials ([#646](https://github.com/kubeflow/katib/pull/646) by [@andreyvelich](https://github.com/andreyvelich))
- v1alpha2 hyperband suggestion service validation ([#648](https://github.com/kubeflow/katib/pull/648) by [@hougangliu](https://github.com/hougangliu))
- refactor: Remove requests check for most test cases ([#626](https://github.com/kubeflow/katib/pull/626) by [@gaocegege](https://github.com/gaocegege))
- feat(experiment): Delete dup trials ([#647](https://github.com/kubeflow/katib/pull/647) by [@gaocegege](https://github.com/gaocegege))
- UI: Add bayesianoptimization algorithm in selectlist ([#645](https://github.com/kubeflow/katib/pull/645) by [@andreyvelich](https://github.com/andreyvelich))
- Fix v1alpha1 hyperband algorithm mismatch ([#634](https://github.com/kubeflow/katib/pull/634) by [@hougangliu](https://github.com/hougangliu))
- v1alpha2 hyperband suggestion service ([#631](https://github.com/kubeflow/katib/pull/631) by [@hougangliu](https://github.com/hougangliu))
- Upgrade Job operators to v1 ([#635](https://github.com/kubeflow/katib/pull/635) by [@johnugeorge](https://github.com/johnugeorge))
- Fix sql syntax for UpdateAlgorithmExtraSettings ([#633](https://github.com/kubeflow/katib/pull/633) by [@hougangliu](https://github.com/hougangliu))
- Update Algorithm extra settings during experiment creation ([#630](https://github.com/kubeflow/katib/pull/630) by [@johnugeorge](https://github.com/johnugeorge))
- Adding cascading delete of pods when jobs are deleted ([#632](https://github.com/kubeflow/katib/pull/632) by [@johnugeorge](https://github.com/johnugeorge))
- Add tests for grid suggestion algorithm ([#628](https://github.com/kubeflow/katib/pull/628) by [@johnugeorge](https://github.com/johnugeorge))
- Fixing tag for Suggestion BO ([#627](https://github.com/kubeflow/katib/pull/627) by [@johnugeorge](https://github.com/johnugeorge))
- Training Container for NAS RL Suggestion in v1alpha2 ([#614](https://github.com/kubeflow/katib/pull/614) by [@andreyvelich](https://github.com/andreyvelich))
- Implementing v1alpha2 grid search suggestion algorithm ([#622](https://github.com/kubeflow/katib/pull/622) by [@johnugeorge](https://github.com/johnugeorge))
- feat: Support bayesianoptimization in v1alpha2 ([#595](https://github.com/kubeflow/katib/pull/595) by [@gaocegege](https://github.com/gaocegege))
- NAS RL Suggestion for v1alpha2 ([#613](https://github.com/kubeflow/katib/pull/613) by [@andreyvelich](https://github.com/andreyvelich))
- Fix problems in the UI for v1alpha2 ([#623](https://github.com/kubeflow/katib/pull/623) by [@andreyvelich](https://github.com/andreyvelich))
- Updated help message for golint. ([#621](https://github.com/kubeflow/katib/pull/621) by [@gyliu513](https://github.com/gyliu513))
- Fix Scheme in Katib Client for v1alpha2 ([#620](https://github.com/kubeflow/katib/pull/620) by [@andreyvelich](https://github.com/andreyvelich))
- Set trial completion status only after metric collection ([#616](https://github.com/kubeflow/katib/pull/616) by [@johnugeorge](https://github.com/johnugeorge))
- go unit tests from presubmits ([#618](https://github.com/kubeflow/katib/pull/618) by [@johnugeorge](https://github.com/johnugeorge))
- Skip creating trials if add count is zero ([#617](https://github.com/kubeflow/katib/pull/617) by [@johnugeorge](https://github.com/johnugeorge))
- Fix nasrl example in v1alpha2 ([#609](https://github.com/kubeflow/katib/pull/609) by [@andreyvelich](https://github.com/andreyvelich))
- Enabled make check in travis. ([#608](https://github.com/kubeflow/katib/pull/608) by [@gyliu513](https://github.com/gyliu513))
- fix make check ([#606](https://github.com/kubeflow/katib/pull/606) by [@gyliu513](https://github.com/gyliu513))
- Fine-grained docker image build. ([#605](https://github.com/kubeflow/katib/pull/605) by [@gyliu513](https://github.com/gyliu513))
- Restructuring manifests ([#602](https://github.com/kubeflow/katib/pull/602) by [@johnugeorge](https://github.com/johnugeorge))
- Fixing latest tag ([#603](https://github.com/kubeflow/katib/pull/603) by [@johnugeorge](https://github.com/johnugeorge))
- Minor changes to metric collector manifest ([#601](https://github.com/kubeflow/katib/pull/601) by [@johnugeorge](https://github.com/johnugeorge))
- Mini fix for v1alpha1 metricsCollector ([#600](https://github.com/kubeflow/katib/pull/600) by [@hougangliu](https://github.com/hougangliu))
- Check error in OpenSQLConnection ([#588](https://github.com/kubeflow/katib/pull/588) by [@andreyvelich](https://github.com/andreyvelich))
- Fix issue of hyperband suggestion service cannot move on ([#596](https://github.com/kubeflow/katib/pull/596) by [@hougangliu](https://github.com/hougangliu))
- doc: Update readme ([#593](https://github.com/kubeflow/katib/pull/593) by [@gaocegege](https://github.com/gaocegege))
- Reverse logic of Less in hyperband v1alpha1 ([#592](https://github.com/kubeflow/katib/pull/592) by [@hougangliu](https://github.com/hougangliu))
- Mini fix for getExperimentConf ([#594](https://github.com/kubeflow/katib/pull/594) by [@hougangliu](https://github.com/hougangliu))
- feat: Add UI in manifests v1alpha2 ([#591](https://github.com/kubeflow/katib/pull/591) by [@gaocegege](https://github.com/gaocegege))
- feat: Support flags in UI ([#590](https://github.com/kubeflow/katib/pull/590) by [@gaocegege](https://github.com/gaocegege))
- Default make target to v1alpha2. ([#585](https://github.com/kubeflow/katib/pull/585) by [@gyliu513](https://github.com/gyliu513))
- Change undeploy script ([#587](https://github.com/kubeflow/katib/pull/587) by [@andreyvelich](https://github.com/andreyvelich))
- Added undeploy for katib. ([#579](https://github.com/kubeflow/katib/pull/579) by [@gyliu513](https://github.com/gyliu513))
- feat(trial): Add more failure test cases ([#570](https://github.com/kubeflow/katib/pull/570) by [@gaocegege](https://github.com/gaocegege))
- Add categories for katib CRDs ([#576](https://github.com/kubeflow/katib/pull/576) by [@hougangliu](https://github.com/hougangliu))
- Add Validate Algorithm Settings in v1alpha2 ([#574](https://github.com/kubeflow/katib/pull/574) by [@andreyvelich](https://github.com/andreyvelich))
- Updated makefile by adding more targets for developer. ([#575](https://github.com/kubeflow/katib/pull/575) by [@gyliu513](https://github.com/gyliu513))
- feat(experiment): Add more test cases ([#563](https://github.com/kubeflow/katib/pull/563) by [@gaocegege](https://github.com/gaocegege))
- refactor: Use manager client to get log for test ([#569](https://github.com/kubeflow/katib/pull/569) by [@gaocegege](https://github.com/gaocegege))
- Adding go tools scripts - part 1 ([#573](https://github.com/kubeflow/katib/pull/573) by [@gyliu513](https://github.com/gyliu513))
- Retain for job and metricsCollector ([#572](https://github.com/kubeflow/katib/pull/572) by [@hougangliu](https://github.com/hougangliu))
- Fix finalizer cannot work ([#571](https://github.com/kubeflow/katib/pull/571) by [@hougangliu](https://github.com/hougangliu))
- Implement GetExperimentInDB ([#558](https://github.com/kubeflow/katib/pull/558) by [@hougangliu](https://github.com/hougangliu))
- refactor: Unify the interface ([#568](https://github.com/kubeflow/katib/pull/568) by [@gaocegege](https://github.com/gaocegege))
- Implement trial observation metrics ([#564](https://github.com/kubeflow/katib/pull/564) by [@johnugeorge](https://github.com/johnugeorge))
- fix: Check if the deletion time is nil or zero ([#567](https://github.com/kubeflow/katib/pull/567) by [@gaocegege](https://github.com/gaocegege))
- feat(experiment-validator): Inject client ([#561](https://github.com/kubeflow/katib/pull/561) by [@gaocegege](https://github.com/gaocegege))
- Change path for yaml file and NAS training containers ([#566](https://github.com/kubeflow/katib/pull/566) by [@andreyvelich](https://github.com/andreyvelich))
- Added metric name to GetObservationLogRequest ([#559](https://github.com/kubeflow/katib/pull/559) by [@johnugeorge](https://github.com/johnugeorge))
- Reducing the Prow workflow name length ([#562](https://github.com/kubeflow/katib/pull/562) by [@johnugeorge](https://github.com/johnugeorge))
- chore: Add code coverage test ([#560](https://github.com/kubeflow/katib/pull/560) by [@gaocegege](https://github.com/gaocegege))
- feat(experiment): Add test cases ([#556](https://github.com/kubeflow/katib/pull/556) by [@gaocegege](https://github.com/gaocegege))
- chore: Refactor nasenvelopenet ([#492](https://github.com/kubeflow/katib/pull/492) by [@gaocegege](https://github.com/gaocegege))
- feat(trial): Refactor trial controller and add basic test cases ([#528](https://github.com/kubeflow/katib/pull/528) by [@gaocegege](https://github.com/gaocegege))
- Add status for experiment and trial in additionalPrinterColumns ([#555](https://github.com/kubeflow/katib/pull/555) by [@hougangliu](https://github.com/hougangliu))
- Fix default metricsController wrong args ([#550](https://github.com/kubeflow/katib/pull/550) by [@hougangliu](https://github.com/hougangliu))
- Add NAS RL yaml example for v1alpha2 ([#554](https://github.com/kubeflow/katib/pull/554) by [@andreyvelich](https://github.com/andreyvelich))
- Fix trial requestNumber error ([#553](https://github.com/kubeflow/katib/pull/553) by [@hougangliu](https://github.com/hougangliu))
- Adding test for random suggestion algorithm ([#552](https://github.com/kubeflow/katib/pull/552) by [@johnugeorge](https://github.com/johnugeorge))
- Adding minor styles changes ([#548](https://github.com/kubeflow/katib/pull/548) by [@johnugeorge](https://github.com/johnugeorge))
- Fix nil pointer error when create trial ([#547](https://github.com/kubeflow/katib/pull/547) by [@hougangliu](https://github.com/hougangliu))
- Used klog for katib - part 1. ([#526](https://github.com/kubeflow/katib/pull/526) by [@gyliu513](https://github.com/gyliu513))
- Implement GetSuggestions for general suggestion ([#546](https://github.com/kubeflow/katib/pull/546) by [@hougangliu](https://github.com/hougangliu))
- fix: Fix the conflicts in UI backend ([#545](https://github.com/kubeflow/katib/pull/545) by [@gaocegege](https://github.com/gaocegege))
- Earlystopping changes trigger CI based on version ([#544](https://github.com/kubeflow/katib/pull/544) by [@hougangliu](https://github.com/hougangliu))
- Adding manifests for manager rest ([#542](https://github.com/kubeflow/katib/pull/542) by [@johnugeorge](https://github.com/johnugeorge))
- Katib UI for v1alpha2 ([#486](https://github.com/kubeflow/katib/pull/486) by [@andreyvelich](https://github.com/andreyvelich))
- Enable suggestion-random image build and push in CI ([#543](https://github.com/kubeflow/katib/pull/543) by [@hougangliu](https://github.com/hougangliu))
- fix(status): Define status handler ([#518](https://github.com/kubeflow/katib/pull/518) by [@gaocegege](https://github.com/gaocegege))
- Include suggestion manager based on version in prow ([#541](https://github.com/kubeflow/katib/pull/541) by [@hougangliu](https://github.com/hougangliu))
- Adding random suggestion algorithm implementation and manifests ([#540](https://github.com/kubeflow/katib/pull/540) by [@johnugeorge](https://github.com/johnugeorge))
- fix: Add test cases for validator and manifest generator ([#508](https://github.com/kubeflow/katib/pull/508) by [@gaocegege](https://github.com/gaocegege))
- Update trial status DB operation ([#537](https://github.com/kubeflow/katib/pull/537) by [@hougangliu](https://github.com/hougangliu))
- [v1alpha2] Add labels for TFjob and PytorchJob in Metrics Collector ([#539](https://github.com/kubeflow/katib/pull/539) by [@richardsliu](https://github.com/richardsliu))
- v1alpha2 TFEvent metrics collector ([#538](https://github.com/kubeflow/katib/pull/538) by [@richardsliu](https://github.com/richardsliu))
- Register Trial in DB ([#530](https://github.com/kubeflow/katib/pull/530) by [@andreyvelich](https://github.com/andreyvelich))
- Restructuring docker files to build images per version ([#527](https://github.com/kubeflow/katib/pull/527) by [@johnugeorge](https://github.com/johnugeorge))
- Dep ensure to sync up vendor. ([#535](https://github.com/kubeflow/katib/pull/535) by [@gyliu513](https://github.com/gyliu513))
- fix: Avoid side effect ([#532](https://github.com/kubeflow/katib/pull/532) by [@gaocegege](https://github.com/gaocegege))
- Include vendor dir instead of Gopkg in prow config ([#536](https://github.com/kubeflow/katib/pull/536) by [@hougangliu](https://github.com/hougangliu))
- Update experiment status DB operation ([#534](https://github.com/kubeflow/katib/pull/534) by [@hougangliu](https://github.com/hougangliu))
- feat(api): Add total number of trials ([#501](https://github.com/kubeflow/katib/pull/501) by [@gaocegege](https://github.com/gaocegege))
- Fix wrong error-handling logic in db interface ([#529](https://github.com/kubeflow/katib/pull/529) by [@hougangliu](https://github.com/hougangliu))
- chore: Remove dep ensure in CI ([#525](https://github.com/kubeflow/katib/pull/525) by [@gaocegege](https://github.com/gaocegege))
- Delete experiment in DB if in need ([#519](https://github.com/kubeflow/katib/pull/519) by [@hougangliu](https://github.com/hougangliu))
- Support for Custom Job resources ([#512](https://github.com/kubeflow/katib/pull/512) by [@johnugeorge](https://github.com/johnugeorge))
- Fix ut test and enable ut-test of v1alpha2 ([#524](https://github.com/kubeflow/katib/pull/524) by [@hougangliu](https://github.com/hougangliu))
- godep: Remove useless dep ([#521](https://github.com/kubeflow/katib/pull/521) by [@gaocegege](https://github.com/gaocegege))
- Fix prow to trigger corresponding workflow ([#520](https://github.com/kubeflow/katib/pull/520) by [@hougangliu](https://github.com/hougangliu))
- create experiment in db ([#509](https://github.com/kubeflow/katib/pull/509) by [@hougangliu](https://github.com/hougangliu))
- refactor(suggestion): Use interface ([#502](https://github.com/kubeflow/katib/pull/502) by [@gaocegege](https://github.com/gaocegege))
- feat(CI): Run different flow according to version ([#516](https://github.com/kubeflow/katib/pull/516) by [@gaocegege](https://github.com/gaocegege))
- Added PR and Issue template. ([#505](https://github.com/kubeflow/katib/pull/505) by [@gyliu513](https://github.com/gyliu513))
- Enabled verbose logging for dev guide. ([#504](https://github.com/kubeflow/katib/pull/504) by [@gyliu513](https://github.com/gyliu513))
- v1alpha2 metrics collector - controller ([#496](https://github.com/kubeflow/katib/pull/496) by [@richardsliu](https://github.com/richardsliu))
- Update util for experiment in v1alpha2 ([#485](https://github.com/kubeflow/katib/pull/485) by [@andreyvelich](https://github.com/andreyvelich))
- add common package ([#491](https://github.com/kubeflow/katib/pull/491) by [@hougangliu](https://github.com/hougangliu))
- Add metrics collector spec and objective spec to Trial ([#489](https://github.com/kubeflow/katib/pull/489) by [@richardsliu](https://github.com/richardsliu))
- Prune katib OWNERS file ([#490](https://github.com/kubeflow/katib/pull/490) by [@richardsliu](https://github.com/richardsliu))
- Training container for NAS Envelopenet ([#429](https://github.com/kubeflow/katib/pull/429) by [@garganubhav](https://github.com/garganubhav))
- NAS Envelopenet Suggestion and Job Example ([#425](https://github.com/kubeflow/katib/pull/425) by [@garganubhav](https://github.com/garganubhav))
- V1alpha2 Metrics collector (part 1) ([#484](https://github.com/kubeflow/katib/pull/484) by [@richardsliu](https://github.com/richardsliu))
- enable test for katib-manager ([#478](https://github.com/kubeflow/katib/pull/478) by [@hougangliu](https://github.com/hougangliu))
- Remove outdated TODOs in README.md ([#468](https://github.com/kubeflow/katib/pull/468) by [@terrytangyuan](https://github.com/terrytangyuan))
- Get experiment config from the instance ([#474](https://github.com/kubeflow/katib/pull/474) by [@andreyvelich](https://github.com/andreyvelich))
- Fix KatibClient name in v1alpha2 ([#483](https://github.com/kubeflow/katib/pull/483) by [@andreyvelich](https://github.com/andreyvelich))
- Add Katib Client in v1alpha2 ([#480](https://github.com/kubeflow/katib/pull/480) by [@andreyvelich](https://github.com/andreyvelich))
- Add metrics collector spec to v1alpha2 API ([#481](https://github.com/kubeflow/katib/pull/481) by [@richardsliu](https://github.com/richardsliu))
- vizier-core does not need any role ([#482](https://github.com/kubeflow/katib/pull/482) by [@hougangliu](https://github.com/hougangliu))
- katib manager db error ([#476](https://github.com/kubeflow/katib/pull/476) by [@hougangliu](https://github.com/hougangliu))
- share one grpc-health-probe ([#477](https://github.com/kubeflow/katib/pull/477) by [@hougangliu](https://github.com/hougangliu))
- validation and mutating webhook for experiment ([#473](https://github.com/kubeflow/katib/pull/473) by [@hougangliu](https://github.com/hougangliu))
- enable test for v1alpha2 ([#465](https://github.com/kubeflow/katib/pull/465) by [@hougangliu](https://github.com/hougangliu))
- Add serviceAccountName in UI deployment ([#469](https://github.com/kubeflow/katib/pull/469) by [@andreyvelich](https://github.com/andreyvelich))
- chore: Skip test when code is not changed ([#467](https://github.com/kubeflow/katib/pull/467) by [@gaocegege](https://github.com/gaocegege))
- Adding initial v1alpha2 API controller ([#457](https://github.com/kubeflow/katib/pull/457) by [@johnugeorge](https://github.com/johnugeorge))
- v1alpha2 api server implementation ([#456](https://github.com/kubeflow/katib/pull/456) by [@YujiOshima](https://github.com/YujiOshima))
- fix(readme): Merge image directory ([#455](https://github.com/kubeflow/katib/pull/455) by [@gaocegege](https://github.com/gaocegege))
- Update REAME example links for v1alpha1 ([#452](https://github.com/kubeflow/katib/pull/452) by [@alexandraj777](https://github.com/alexandraj777))
- fix py client import error ([#453](https://github.com/kubeflow/katib/pull/453) by [@hougangliu](https://github.com/hougangliu))
- ClusterRoleBinding doesn't need namespace field ([#451](https://github.com/kubeflow/katib/pull/451) by [@hougangliu](https://github.com/hougangliu))
- Update API for NAS in v1alpha2 ([#450](https://github.com/kubeflow/katib/pull/450) by [@andreyvelich](https://github.com/andreyvelich))
- Restructuring test scripts for v1alpha1 and v1alpha2 ([#449](https://github.com/kubeflow/katib/pull/449) by [@johnugeorge](https://github.com/johnugeorge))
- Code restructuring to support V1alpha1 and V1alpha2 API ([#448](https://github.com/kubeflow/katib/pull/448) by [@johnugeorge](https://github.com/johnugeorge))
- Fix labels matching the job operator implementation ([#447](https://github.com/kubeflow/katib/pull/447) by [@johnugeorge](https://github.com/johnugeorge))
- Updating the pytorch example image ([#446](https://github.com/kubeflow/katib/pull/446) by [@johnugeorge](https://github.com/johnugeorge))
- Remove redundant lock ([#444](https://github.com/kubeflow/katib/pull/444) by [@mrkm4ntr](https://github.com/mrkm4ntr))
- add v1alpha2 grpc api ([#427](https://github.com/kubeflow/katib/pull/427) by [@YujiOshima](https://github.com/YujiOshima))
- Remove katibcli ([#436](https://github.com/kubeflow/katib/pull/436) by [@jdplatt](https://github.com/jdplatt))
- Change datadir for avoid failure due to lost+found ([#432](https://github.com/kubeflow/katib/pull/432) by [@mrkm4ntr](https://github.com/mrkm4ntr))
- fix demo link ([#434](https://github.com/kubeflow/katib/pull/434) by [@jq](https://github.com/jq))
- Add fault tolerance support for trial failure ([#424](https://github.com/kubeflow/katib/pull/424) by [@DeeperMind](https://github.com/DeeperMind))
- Test for Bayesian Optimization Algo ([#406](https://github.com/kubeflow/katib/pull/406) by [@jdplatt](https://github.com/jdplatt))
- Katib v1alpha2 API for CRDs ([#381](https://github.com/kubeflow/katib/pull/381) by [@richardsliu](https://github.com/richardsliu))
- Add NAS team as reviewers ([#419](https://github.com/kubeflow/katib/pull/419) by [@andreyvelich](https://github.com/andreyvelich))
- Multiple Trials for Reinforcement Learning Suggestion ([#416](https://github.com/kubeflow/katib/pull/416) by [@DeeperMind](https://github.com/DeeperMind))
- Fix the package version in training container ([#418](https://github.com/kubeflow/katib/pull/418) by [@DeeperMind](https://github.com/DeeperMind))
- Add validation for NAS job in Katib controller ([#398](https://github.com/kubeflow/katib/pull/398) by [@andreyvelich](https://github.com/andreyvelich))
- Fix path to API protobuf in developer guide ([#415](https://github.com/kubeflow/katib/pull/415) by [@andreyvelich](https://github.com/andreyvelich))
- Add support for parallel studyjobs ([#404](https://github.com/kubeflow/katib/pull/404) by [@DeeperMind](https://github.com/DeeperMind))
- Add separable/depthwise convolution, data augmentation and multiple GPU support ([#393](https://github.com/kubeflow/katib/pull/393) by [@DeeperMind](https://github.com/DeeperMind))
- Add create time to Trial API ([#410](https://github.com/kubeflow/katib/pull/410) by [@andreyvelich](https://github.com/andreyvelich))
- Metric collector must fail on error ([#405](https://github.com/kubeflow/katib/pull/405) by [@johnugeorge](https://github.com/johnugeorge))
- add latest tag for katib images ([#409](https://github.com/kubeflow/katib/pull/409) by [@hougangliu](https://github.com/hougangliu))
- add build and test for suggestion nasrl ([#401](https://github.com/kubeflow/katib/pull/401) by [@hougangliu](https://github.com/hougangliu))
- Database APIs for NAS updated ([#394](https://github.com/kubeflow/katib/pull/394) by [@Akado2009](https://github.com/Akado2009))
- Suggestion for Neural Architecture Search with Reinforcement Learning ([#339](https://github.com/kubeflow/katib/pull/339) by [@DeeperMind](https://github.com/DeeperMind))
- add validating webhook for studyJob ([#383](https://github.com/kubeflow/katib/pull/383) by [@hougangliu](https://github.com/hougangliu))
- Removing Operator specific handling during a StudyJob run ([#387](https://github.com/kubeflow/katib/pull/387) by [@johnugeorge](https://github.com/johnugeorge))
- Delete modeldb from unit tests ([#391](https://github.com/kubeflow/katib/pull/391) by [@andreyvelich](https://github.com/andreyvelich))
- show studyjob condition when run kubectl get ([#389](https://github.com/kubeflow/katib/pull/389) by [@hougangliu](https://github.com/hougangliu))
- Training Container with Model Constructor for cifar10 ([#345](https://github.com/kubeflow/katib/pull/345) by [@DeeperMind](https://github.com/DeeperMind))
- add studyjob python client ([#379](https://github.com/kubeflow/katib/pull/379) by [@hougangliu](https://github.com/hougangliu))
- fix wrong example ([#378](https://github.com/kubeflow/katib/pull/378) by [@hougangliu](https://github.com/hougangliu))
- Upgrading controller runtime and k8s to 1.11.2 ([#376](https://github.com/kubeflow/katib/pull/376) by [@johnugeorge](https://github.com/johnugeorge))
- Properly initialize CI cluster credential ([#360](https://github.com/kubeflow/katib/pull/360) by [@toshiiw](https://github.com/toshiiw))
- Include go dependencies in developer-guide.md ([#369](https://github.com/kubeflow/katib/pull/369) by [@alexandraj777](https://github.com/alexandraj777))
- fix invalid memory address ([#368](https://github.com/kubeflow/katib/pull/368) by [@hougangliu](https://github.com/hougangliu))
- Fix presubmits ([#363](https://github.com/kubeflow/katib/pull/363) by [@richardsliu](https://github.com/richardsliu))
- Katib 2019 Roadmap ([#348](https://github.com/kubeflow/katib/pull/348) by [@richardsliu](https://github.com/richardsliu))
- Update OWNERS ([#350](https://github.com/kubeflow/katib/pull/350) by [@richardsliu](https://github.com/richardsliu))
- Extend Katib API for NAS jobs ([#327](https://github.com/kubeflow/katib/pull/327) by [@andreyvelich](https://github.com/andreyvelich))
- ignore tfjob/pytorch job if corresponding CRD not created ([#335](https://github.com/kubeflow/katib/pull/335) by [@hougangliu](https://github.com/hougangliu))
- Clarify the example UI is generated by random-example. ([#333](https://github.com/kubeflow/katib/pull/333) by [@gyliu513](https://github.com/gyliu513))
- only try to delete study info in db when in need ([#342](https://github.com/kubeflow/katib/pull/342) by [@hougangliu](https://github.com/hougangliu))
- omit empty fields for studyjob status ([#336](https://github.com/kubeflow/katib/pull/336) by [@hougangliu](https://github.com/hougangliu))
- Update pytorch example with latest image ([#329](https://github.com/kubeflow/katib/pull/329) by [@TimZaman](https://github.com/TimZaman))
- Fix typo in json API ([#330](https://github.com/kubeflow/katib/pull/330) by [@richardsliu](https://github.com/richardsliu))
- Add information how to run TFjob and Pytorch examples in Katib ([#321](https://github.com/kubeflow/katib/pull/321) by [@andreyvelich](https://github.com/andreyvelich))
- Add xgboost example using Bayesian optimization ([#320](https://github.com/kubeflow/katib/pull/320) by [@richardsliu](https://github.com/richardsliu))
- katib should be able to be deployed in any namespace ([#324](https://github.com/kubeflow/katib/pull/324) by [@hougangliu](https://github.com/hougangliu))
- Adding distributed pytorch example for katib ([#309](https://github.com/kubeflow/katib/pull/309) by [@johnugeorge](https://github.com/johnugeorge))
- Minor fixes ([#307](https://github.com/kubeflow/katib/pull/307) by [@johnugeorge](https://github.com/johnugeorge))
- delete obsolete data in db ([#315](https://github.com/kubeflow/katib/pull/315) by [@hougangliu](https://github.com/hougangliu))
- add bestTrialId to statusJob status ([#312](https://github.com/kubeflow/katib/pull/312) by [@hougangliu](https://github.com/hougangliu))
- Add api doc ([#303](https://github.com/kubeflow/katib/pull/303) by [@YujiOshima](https://github.com/YujiOshima))
- validate studyJob when first reconcile it ([#308](https://github.com/kubeflow/katib/pull/308) by [@hougangliu](https://github.com/hougangliu))
- add hougangliu as a reviewer ([#310](https://github.com/kubeflow/katib/pull/310) by [@hougangliu](https://github.com/hougangliu))
- Adding to OWNERS file ([#304](https://github.com/kubeflow/katib/pull/304) by [@johnugeorge](https://github.com/johnugeorge))
- sync up worker status all the time ([#299](https://github.com/kubeflow/katib/pull/299) by [@hougangliu](https://github.com/hougangliu))
- studyJob with non-kubeflow namespace cannot work ([#302](https://github.com/kubeflow/katib/pull/302) by [@hougangliu](https://github.com/hougangliu))
- Adding master pod check for default metric collector ([#300](https://github.com/kubeflow/katib/pull/300) by [@johnugeorge](https://github.com/johnugeorge))
- reduce some redundant code ([#296](https://github.com/kubeflow/katib/pull/296) by [@hougangliu](https://github.com/hougangliu))
- Extend studyjob client API ([#288](https://github.com/kubeflow/katib/pull/288) by [@andreyvelich](https://github.com/andreyvelich))
- Use same deploy.sh when deploy katib components ([#284](https://github.com/kubeflow/katib/pull/284) by [@ytetra](https://github.com/ytetra))
- update Readme ([#295](https://github.com/kubeflow/katib/pull/295) by [@hougangliu](https://github.com/hougangliu))
- fix studyJob status suggestionCount mismatch error ([#290](https://github.com/kubeflow/katib/pull/290) by [@hougangliu](https://github.com/hougangliu))
- fix invalid worker kind issue ([#287](https://github.com/kubeflow/katib/pull/287) by [@hougangliu](https://github.com/hougangliu))
- get metricscollector by API ([#292](https://github.com/kubeflow/katib/pull/292) by [@YujiOshima](https://github.com/YujiOshima))
- Support Pytorch job in Katib ([#283](https://github.com/kubeflow/katib/pull/283) by [@johnugeorge](https://github.com/johnugeorge))
- Update k8s cluster version to 1.10 ([#286](https://github.com/kubeflow/katib/pull/286) by [@johnugeorge](https://github.com/johnugeorge))
- Enrich GUI ([#264](https://github.com/kubeflow/katib/pull/264) by [@YujiOshima](https://github.com/YujiOshima))
- update README ([#281](https://github.com/kubeflow/katib/pull/281) by [@hougangliu](https://github.com/hougangliu))
- fix typo error for MinikubeDemo ([#282](https://github.com/kubeflow/katib/pull/282) by [@hougangliu](https://github.com/hougangliu))
- fix typo error ([#280](https://github.com/kubeflow/katib/pull/280) by [@hougangliu](https://github.com/hougangliu))
- add e2eTest of each suggestion algorithm ([#265](https://github.com/kubeflow/katib/pull/265) by [@ytetra](https://github.com/ytetra))
- Allow studyjobcontroller to delete pods ([#278](https://github.com/kubeflow/katib/pull/278) by [@richardsliu](https://github.com/richardsliu))
- Fix katib ui resource paths ([#277](https://github.com/kubeflow/katib/pull/277) by [@richardsliu](https://github.com/richardsliu))
- Implement gRPC Health Checking Protocol + add readiness/liveness probes to vizier-core ([#270](https://github.com/kubeflow/katib/pull/270) by [@lkpdn](https://github.com/lkpdn))
- POC: Katib integration with tf-operator ([#267](https://github.com/kubeflow/katib/pull/267) by [@richardsliu](https://github.com/richardsliu))
- fix timing to determine slice size in grid search ([#271](https://github.com/kubeflow/katib/pull/271) by [@ytetra](https://github.com/ytetra))
- Add Update{Study,Trial} ([#269](https://github.com/kubeflow/katib/pull/269) by [@toshiiw](https://github.com/toshiiw))
- add Richard Liu to OWNERS ([#274](https://github.com/kubeflow/katib/pull/274) by [@YujiOshima](https://github.com/YujiOshima))
- fix uncompleted value in ui ([#238](https://github.com/kubeflow/katib/pull/238) by [@YujiOshima](https://github.com/YujiOshima))
- fix bayesian optimization suggestion ([#251](https://github.com/kubeflow/katib/pull/251) by [@YujiOshima](https://github.com/YujiOshima))
- Prevent pod restarts caused by slow db boot ([#261](https://github.com/kubeflow/katib/pull/261) by [@lkpdn](https://github.com/lkpdn))
- add UT of each suggestion algorithm ([#237](https://github.com/kubeflow/katib/pull/237) by [@ytetra](https://github.com/ytetra))
- Downgrade kubernetes dependency to 1.10.1 ([#256](https://github.com/kubeflow/katib/pull/256) by [@richardsliu](https://github.com/richardsliu))
- Fix incorrectly set namespace ([#260](https://github.com/kubeflow/katib/pull/260) by [@lkpdn](https://github.com/lkpdn))
- Set MYSQL_ROOT_PASSWORD via Secret ([#253](https://github.com/kubeflow/katib/pull/253) by [@lkpdn](https://github.com/lkpdn))
- update UI ([#255](https://github.com/kubeflow/katib/pull/255) by [@YujiOshima](https://github.com/YujiOshima))
- Refactor studyjobcontroller ([#254](https://github.com/kubeflow/katib/pull/254) by [@richardsliu](https://github.com/richardsliu))
- Change deploy.sh for Minikube example ([#252](https://github.com/kubeflow/katib/pull/252) by [@andreyvelich](https://github.com/andreyvelich))
- Add mysql based unit tests ([#243](https://github.com/kubeflow/katib/pull/243) by [@toshiiw](https://github.com/toshiiw))
- Update manifests ([#246](https://github.com/kubeflow/katib/pull/246) by [@YujiOshima](https://github.com/YujiOshima))
- Add texasmichelle as reviewer ([#247](https://github.com/kubeflow/katib/pull/247) by [@texasmichelle](https://github.com/texasmichelle))
- Tf event mc ([#235](https://github.com/kubeflow/katib/pull/235) by [@YujiOshima](https://github.com/YujiOshima))
- Fix typos for json and objective ([#242](https://github.com/kubeflow/katib/pull/242) by [@toshiiw](https://github.com/toshiiw))
- Add richardsliu to OWNERS/reviewer ([#239](https://github.com/kubeflow/katib/pull/239) by [@richardsliu](https://github.com/richardsliu))
- add starttime and completiontime to worker ([#236](https://github.com/kubeflow/katib/pull/236) by [@wukong1992](https://github.com/wukong1992))
- Fix typo ([#233](https://github.com/kubeflow/katib/pull/233) by [@ytetra](https://github.com/ytetra))
- More DB unit tests ([#234](https://github.com/kubeflow/katib/pull/234) by [@toshiiw](https://github.com/toshiiw))
- Fix the build script after #208 ([#231](https://github.com/kubeflow/katib/pull/231) by [@toshiiw](https://github.com/toshiiw))
- Only retry an INSERT operation on unique constraint violation ([#229](https://github.com/kubeflow/katib/pull/229) by [@toshiiw](https://github.com/toshiiw))
- New UI for Katib ([#208](https://github.com/kubeflow/katib/pull/208) by [@YujiOshima](https://github.com/YujiOshima))
- fix slice range ([#226](https://github.com/kubeflow/katib/pull/226) by [@ytetra](https://github.com/ytetra))
- More db tests ([#225](https://github.com/kubeflow/katib/pull/225) by [@toshiiw](https://github.com/toshiiw))
- Fix storelogs ([#222](https://github.com/kubeflow/katib/pull/222) by [@toshiiw](https://github.com/toshiiw))
- Check errors in order to avoid SEGV ([#219](https://github.com/kubeflow/katib/pull/219) by [@toshiiw](https://github.com/toshiiw))
- Fix reqest count ([#214](https://github.com/kubeflow/katib/pull/214) by [@YujiOshima](https://github.com/YujiOshima))

[Full Changelog](https://github.com/kubeflow/katib/compare/826657c14602a3f36263f3d6769451af0a75d18a...v0.6.0-rc.0)

# [0.2](https://github.com/kubeflow/katib/tree/0.2) (2018-08-20)

## Features

- pin mxnet/python image version ([#139](https://github.com/kubeflow/katib/pull/139) by [@mayankjuneja](https://github.com/mayankjuneja))
- Move the GKEDemo into kubeflow/examples ([#135](https://github.com/kubeflow/katib/pull/135) by [@jlewi](https://github.com/jlewi))
- update OWNERS ([#129](https://github.com/kubeflow/katib/pull/129) by [@mitake](https://github.com/mitake))
- Hyperband ([#124](https://github.com/kubeflow/katib/pull/124) by [@YujiOshima](https://github.com/YujiOshima))
- add releasing workflow ([#113](https://github.com/kubeflow/katib/pull/113) by [@YujiOshima](https://github.com/YujiOshima))
- API: Add WorkerStatus to GetMetrics and remove unused items ([#110](https://github.com/kubeflow/katib/pull/110) by [@YujiOshima](https://github.com/YujiOshima))
- Add e2e test ([#114](https://github.com/kubeflow/katib/pull/114) by [@YujiOshima](https://github.com/YujiOshima))
- use kubectl port-forward in demos ([#111](https://github.com/kubeflow/katib/pull/111) by [@YujiOshima](https://github.com/YujiOshima))
- docs: Generate CLI documentation ([#105](https://github.com/kubeflow/katib/pull/105) by [@gaocegege](https://github.com/gaocegege))
- changelog: Add ([#104](https://github.com/kubeflow/katib/pull/104) by [@gaocegege](https://github.com/gaocegege))

## Bug Fixes

- Corrected typos in hyperband example yml ([#146](https://github.com/kubeflow/katib/pull/146) by [@shibuiwilliam](https://github.com/shibuiwilliam))
- Update status of workers in GetWorkers ([#127](https://github.com/kubeflow/katib/pull/127) by [@YujiOshima](https://github.com/YujiOshima))
- fix doc link and kubectl port-forward command ([#120](https://github.com/kubeflow/katib/pull/120) by [@YujiOshima](https://github.com/YujiOshima))
- Fix typo ([#123](https://github.com/kubeflow/katib/pull/123) by [@mrkm4ntr](https://github.com/mrkm4ntr))
- Fix indentation to use spaces (instead of a mix of tabs and spaces) ([#121](https://github.com/kubeflow/katib/pull/121) by [@vinaykakade](https://github.com/vinaykakade))
- docs: Fix wrong command ([#108](https://github.com/kubeflow/katib/pull/108) by [@mrkm4ntr](https://github.com/mrkm4ntr))
- Remove dlk from manifests ([#107](https://github.com/kubeflow/katib/pull/107) by [@vinaykakade](https://github.com/vinaykakade))

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.1.2-alpha...826657c14602a3f36263f3d6769451af0a75d18a)

# [v0.1.2-alpha](https://github.com/kubeflow/katib/tree/v0.1.2-alpha) (2018-06-05)

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.1.1-alpha...v0.1.2-alpha)

**Closed issues:**

- \[request\] Invite libbyandhelen as reviewer for algorithm support [\#82](https://github.com/kubeflow/katib/issues/82)
- cli failed to connect [\#80](https://github.com/kubeflow/katib/issues/80)
- CreateStudy RPC error: Objective_Value_Name is required [\#73](https://github.com/kubeflow/katib/issues/73)
- \[cli\] Use cobra to refactor the cli [\#54](https://github.com/kubeflow/katib/issues/54)
- Reduce time it takes to build all images [\#50](https://github.com/kubeflow/katib/issues/50)
- \[release\] Ksonnet the katib [\#32](https://github.com/kubeflow/katib/issues/32)

**Merged pull requests:**

- update docs [\#103](https://github.com/kubeflow/katib/pull/103) ([YujiOshima](https://github.com/YujiOshima))
- cli: Rename to katib-cli [\#101](https://github.com/kubeflow/katib/pull/101) ([gaocegege](https://github.com/gaocegege))
- Fix dbbug [\#98](https://github.com/kubeflow/katib/pull/98) ([YujiOshima](https://github.com/YujiOshima))
- save logs after check status [\#95](https://github.com/kubeflow/katib/pull/95) ([YujiOshima](https://github.com/YujiOshima))
- Some fix to getting-started.md [\#94](https://github.com/kubeflow/katib/pull/94) ([lluunn](https://github.com/lluunn))
- Add katib-cli download command for Mac [\#93](https://github.com/kubeflow/katib/pull/93) ([vinaykakade](https://github.com/vinaykakade))
- fix get service-param-list bug [\#92](https://github.com/kubeflow/katib/pull/92) ([YujiOshima](https://github.com/YujiOshima))
- fix ui bug [\#91](https://github.com/kubeflow/katib/pull/91) ([YujiOshima](https://github.com/YujiOshima))
- Build parallel [\#85](https://github.com/kubeflow/katib/pull/85) ([YujiOshima](https://github.com/YujiOshima))
- reduce build time [\#84](https://github.com/kubeflow/katib/pull/84) ([YujiOshima](https://github.com/YujiOshima))
- OWNERS: Add @libbyandhelen [\#83](https://github.com/kubeflow/katib/pull/83) ([gaocegege](https://github.com/gaocegege))
- add random forest prior to Bayesian Optimization [\#81](https://github.com/kubeflow/katib/pull/81) ([libbyandhelen](https://github.com/libbyandhelen))
- workflows.libsonnet: Fix the image name [\#75](https://github.com/kubeflow/katib/pull/75) ([gaocegege](https://github.com/gaocegege))
- Refine API [\#74](https://github.com/kubeflow/katib/pull/74) ([YujiOshima](https://github.com/YujiOshima))
- worker: Rename worker_interface to worker [\#70](https://github.com/kubeflow/katib/pull/70) ([gaocegege](https://github.com/gaocegege))

# [v0.1.1-alpha](https://github.com/kubeflow/katib/tree/v0.1.1-alpha) (2018-04-24)

[Full Changelog](https://github.com/kubeflow/katib/compare/v0.1.0-alpha...v0.1.1-alpha)

**Closed issues:**

- \[upstream\] Update name in kubernetes/test-infra [\#63](https://github.com/kubeflow/katib/issues/63)
- \[go\] Update the package name, again [\#62](https://github.com/kubeflow/katib/issues/62)
- \[test\] Fix broken unit test cases [\#58](https://github.com/kubeflow/katib/issues/58)
- Provide a cli binary for macOS / darwin [\#57](https://github.com/kubeflow/katib/issues/57)
- Error running katib on latest master \(04/13\) [\#44](https://github.com/kubeflow/katib/issues/44)
- Upload existing models to modelDB interface [\#43](https://github.com/kubeflow/katib/issues/43)
- \[release\] Add cli to v0.1.0-alpha [\#31](https://github.com/kubeflow/katib/issues/31)
- \[discussion\] Find a new way to install CLI [\#26](https://github.com/kubeflow/katib/issues/26)
- \[maintainance\] Setup the repository [\#8](https://github.com/kubeflow/katib/issues/8)
- Existing approaches and design for hyperparameter-tuning [\#2](https://github.com/kubeflow/katib/issues/2)

**Merged pull requests:**

- Cobra cli [\#69](https://github.com/kubeflow/katib/pull/69) ([YujiOshima](https://github.com/YujiOshima))
- \*: Refactor the structure [\#65](https://github.com/kubeflow/katib/pull/65) ([gaocegege](https://github.com/gaocegege))
- \*: Update name [\#64](https://github.com/kubeflow/katib/pull/64) ([gaocegege](https://github.com/gaocegege))
- Replace kubeflow-images-staging with kubeflow-images-public [\#61](https://github.com/kubeflow/katib/pull/61) ([ankushagarwal](https://github.com/ankushagarwal))
- improve frontend [\#60](https://github.com/kubeflow/katib/pull/60) ([YujiOshima](https://github.com/YujiOshima))
- argo: Add unit test [\#56](https://github.com/kubeflow/katib/pull/56) ([gaocegege](https://github.com/gaocegege))
- main.go: Fix style [\#55](https://github.com/kubeflow/katib/pull/55) ([gaocegege](https://github.com/gaocegege))
- Fix modelsave [\#52](https://github.com/kubeflow/katib/pull/52) ([YujiOshima](https://github.com/YujiOshima))
- refactor Model API [\#51](https://github.com/kubeflow/katib/pull/51) ([YujiOshima](https://github.com/YujiOshima))
- improve test script [\#49](https://github.com/kubeflow/katib/pull/49) ([YujiOshima](https://github.com/YujiOshima))
- Add Model Management API [\#48](https://github.com/kubeflow/katib/pull/48) ([YujiOshima](https://github.com/YujiOshima))
- reviewers: Add @ddysher @jose5918 @mitake [\#45](https://github.com/kubeflow/katib/pull/45) ([gaocegege](https://github.com/gaocegege))
- add early stoppping service [\#41](https://github.com/kubeflow/katib/pull/41) ([YujiOshima](https://github.com/YujiOshima))
- bayesian optimization draft [\#38](https://github.com/kubeflow/katib/pull/38) ([libbyandhelen](https://github.com/libbyandhelen))
- Dockerfile: Use alpine as base image [\#37](https://github.com/kubeflow/katib/pull/37) ([gaocegege](https://github.com/gaocegege))
- docs: Update katib-cli [\#36](https://github.com/kubeflow/katib/pull/36) ([gaocegege](https://github.com/gaocegege))
- New db log schema [\#35](https://github.com/kubeflow/katib/pull/35) ([YujiOshima](https://github.com/YujiOshima))
- Fix CI failures [\#27](https://github.com/kubeflow/katib/pull/27) ([gaocegege](https://github.com/gaocegege))

# [v0.1.0-alpha](https://github.com/kubeflow/katib/tree/v0.1.0-alpha) (2018-04-10)

**Closed issues:**

- \[suggestion\] Move the logic about random service to `random` package [\#18](https://github.com/kubeflow/katib/issues/18)
- \[build-release\] Reuse the vendor during the image building process [\#14](https://github.com/kubeflow/katib/issues/14)
- \[go\] Rename the package from mlkube/katib to this repo [\#7](https://github.com/kubeflow/katib/issues/7)
- \[go\] Establish vendor dependencies for go [\#5](https://github.com/kubeflow/katib/issues/5)
- Rename to hyperparameter-tuning ? [\#1](https://github.com/kubeflow/katib/issues/1)

**Merged pull requests:**

- cleanup of README [\#30](https://github.com/kubeflow/katib/pull/30) ([ddutta](https://github.com/ddutta))
- delete unnecessary settings [\#29](https://github.com/kubeflow/katib/pull/29) ([YujiOshima](https://github.com/YujiOshima))
- Dockerfile: Support multiple stage build in dlk and frontend [\#28](https://github.com/kubeflow/katib/pull/28) ([YujiOshima](https://github.com/YujiOshima))
- Dockerfile: Support multiple stage build in manager and cli [\#25](https://github.com/kubeflow/katib/pull/25) ([gaocegege](https://github.com/gaocegege))
- Dockerfile: Use multiple stage builds [\#23](https://github.com/kubeflow/katib/pull/23) ([gaocegege](https://github.com/gaocegege))
- Ci setup [\#22](https://github.com/kubeflow/katib/pull/22) ([YujiOshima](https://github.com/YujiOshima))
- suggestion: Refactor [\#21](https://github.com/kubeflow/katib/pull/21) ([gaocegege](https://github.com/gaocegege))
- update packages [\#19](https://github.com/kubeflow/katib/pull/19) ([YujiOshima](https://github.com/YujiOshima))
- README: Add code quality badge [\#17](https://github.com/kubeflow/katib/pull/17) ([gaocegege](https://github.com/gaocegege))
- Fixing some basic typos in README [\#13](https://github.com/kubeflow/katib/pull/13) ([ddutta](https://github.com/ddutta))
- vendor: Add [\#12](https://github.com/kubeflow/katib/pull/12) ([gaocegege](https://github.com/gaocegege))
- ignore: Add macOS, Windows and Go ignore files [\#11](https://github.com/kubeflow/katib/pull/11) ([gaocegege](https://github.com/gaocegege))
- Rename packages and move dlk dir [\#10](https://github.com/kubeflow/katib/pull/10) ([YujiOshima](https://github.com/YujiOshima))
- doc: Refactor [\#9](https://github.com/kubeflow/katib/pull/9) ([gaocegege](https://github.com/gaocegege))
- add katib code [\#4](https://github.com/kubeflow/katib/pull/4) ([YujiOshima](https://github.com/YujiOshima))
- add OWNERS file [\#3](https://github.com/kubeflow/katib/pull/3) ([YujiOshima](https://github.com/YujiOshima))
