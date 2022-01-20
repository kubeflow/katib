from __future__ import absolute_import, division, print_function, unicode_literals
import os
import glob
from google.cloud import storage
import re
import logging

_GCS_PREFIX = "gs://"


class Storage(object):
    @staticmethod
    def upload(uri: str, out_dir: str = None) -> str:
        logging.info("Copying contents from %s to %s", uri, out_dir)

        if out_dir.startswith(_GCS_PREFIX):
            Storage._upload_gcs(uri, out_dir)
        else:
            raise Exception("Cannot recognize storage type for " + uri +
                            "\n'%s' are the current available storage type." %
                            (_GCS_PREFIX))

        logging.info("Successfully copied %s to %s", uri, out_dir)
        return out_dir
    
    @staticmethod
    def _upload_gcs(uri, out_dir: str):
        try:
            storage_client = storage.Client()
        except exceptions.DefaultCredentialsError:
            storage_client = storage.Client.create_anonymous_client()
        
        bucket_args = out_dir.replace(_GCS_PREFIX, "", 1).split("/", 1)
        bucket_name = bucket_args[0]
        gcs_path = bucket_args[1] if len(bucket_args) > 1 else ""
        bucket = storage_client.bucket(bucket_name)
        Storage.upload_local_directory_to_gcs(uri,bucket, gcs_path)
    
    @staticmethod
    def upload_local_directory_to_gcs(local_path, bucket, gcs_path):
        assert os.path.isdir(local_path)
        for local_file in glob.glob(local_path + '/**'):
            if not os.path.isfile(local_file):
                Storage.upload_local_directory_to_gcs(local_file, bucket, gcs_path + "/" + os.path.basename(local_file))
            else:
                remote_path = os.path.join(gcs_path, local_file[1 + len(local_path):])
                blob = bucket.blob(remote_path)
                blob.upload_from_filename(local_file)
  