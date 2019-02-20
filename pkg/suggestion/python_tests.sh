pylint bayesianoptimization/src --disable=fixme --exit-zero --reports=y
pytest tests --verbose --cov=bayesianoptimization/src --cov-report term-missing
