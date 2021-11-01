from distutils.core import setup
from setuptools import find_packages

setup(name='jointrpc',
      version='0.0.7',
      description='jointrpc python client',
      author='Zeng Ke',
      author_email='superisaac@gmail.com',
      packages=find_packages(),
      scripts=[],
      classifiers=[
          'Development Status :: 0 - Beta',
          'Environment :: Console',
          'Intended Audience :: Developers',
          'License :: MIT',
          'Programming Language :: Python :: 3.6',
          'Operating System :: POSIX',
          'Topic :: Micro-Services',
      ],
      install_requires=[
          'protobuf>=3.19.1',
          'grpclib>=0.4.2',
      ],
      python_requires='>=3.7',
)
