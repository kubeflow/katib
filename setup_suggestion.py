from setuptools import setup, find_packages

setup(name="katib_suggestion",
      packages=find_packages("pkg/suggestion/bayesianoptimization"),
      package_dir={'':'pkg/suggestion/bayesianoptimization'})
