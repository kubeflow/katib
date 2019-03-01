pip install -r cmd/suggestion/bayesianoptimization/requirements.txt
pip install -r pkg/suggestion/test_requirements.txt
python setup.py develop
pylint pkg/suggestion/bayesianoptimization/src pkg/suggestion/bayesian_service.py --disable=fixme --exit-zero --reports=y
pytest pkg/suggestion/tests --verbose --cov=pkg/suggestion/bayesianoptimization/src --cov-report term-missing
