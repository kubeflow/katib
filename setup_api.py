from setuptools import setup, find_packages

setup(name="katib_api",
      packages=find_packages("pkg/api"),
      package_dir={'':'pkg/api'})
