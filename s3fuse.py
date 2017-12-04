#!/usr/bin/env python
from __future__ import print_function, absolute_import, division
import sys
import logging
from functools  import lru_cache
from sys import argv, exit
from fuse import FUSE, Operations, LoggingMixIn, FuseOSError
import boto3
import os.path
from errno import ENOENT
from botocore.exceptions import ClientError

# Logging
logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)
ch = logging.StreamHandler()
ch.setLevel(logging.DEBUG)
ch.setFormatter(logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s"))
logger.addHandler(ch)

"""
With thanks to @perrygeo(https://gist.github.com/perrygeo)
For the s3 file-like wrapper code. Now torn apart but inspired from:
https://gist.github.com/perrygeo/9239b9ab64731cacbb35#file-s3reader-py
"""

s3 = boto3.resource('s3')

def open(path):
    return S3Reader(path)

@lru_cache(maxsize=32)
def get_s3_obj(path):
    logger.info("Creating S3 object for %s", path)
    bucket, key = parse_path(path)
    obj = s3.Object(bucket, key)
    obj.path = path
    return obj

def range_string(start, stop):
        return "bytes={}-{}".format(start, stop)

# TODO: worth caching? @lru_cache(maxsize=128)
def parse_path(path):
    path = path[1:] if path[0] == '/' else path
    parts = path.split("/")
    bucket = parts[0]
    key = "/".join(parts[1:])
    return bucket, key

def size_limited_caching_byte_request(path, start, stop):
    method = get_bytes if (stop - start < 17000) else  get_bytes.__wrapped__ # the limit should be just over a metadata/chunk request limit?
    return method(path, start, stop)

@lru_cache(maxsize=128)
def get_bytes(path, start, stop):
    rng=range_string(start, stop)
    logger.info("Request %s between %s", path, rng)
    return get_s3_obj(path).get(Range=rng)['Body'].read()


@lru_cache(maxsize=128)
def obj_type(path):
    """
    0 not found
    1 dir
    2 file
    """

    # Test if any object in bucket has prefix
    try:
        bucket, key = parse_path(path)
        if not len(key) > 0:
            return 1
        boto3.client('s3').list_objects_v2(Bucket=bucket,Prefix=key,MaxKeys=1)['Contents']
    except KeyError:
        raise FuseOSError(ENOENT)

    # Test if path represents a complete bucket, key pair.
    try:
        if get_s3_obj(path).content_length <= 0:
            raise ValueError("Content empty")
        return 2 # Object exists. It's a file.
    except ValueError as e:
        raise FuseOSError(ENOENT)
    except ClientError as e:
        if e.response['Error']['Code'] == "404":
            return 1 # The key doesn't exist so treat as a directory
        else:
            raise # Something else has gone wrong.


class S3Reader(object):
    def __init__(self, path):
        self.size = get_s3_obj(path).content_length
        self.pos = 0  # pointer to starting read position
        self.path = path

    def __enter__(self):
        return self

    def __exit__(self, type, value, traceback):
        self.pos = 0

    def read(self, nbytes=None):
        if not nbytes:
            nbytes = self.size - self.pos
        # TODO confirm that start and stop bytes are within 0 to size
        the_bytes =  size_limited_caching_byte_request(self.path, self.pos , self.pos + nbytes - 1)
        self.pos += nbytes
        return the_bytes

    def seek(self, offset, whence=0):
        self.pos = whence + offset


class S3FileSystemMount(Operations):

    def __init__(self):
        self.count = 0;
        self.openfh = {};

    def flush(self, path, fh):
        return None

    def getattr(self, path, fh=None):
        logger.info("you asked for %s", path)

        return  {'st_mode': 33188, 'st_size': open(path).size} if obj_type(path) == 2 else {'st_mode': 16877}


    def open(self, file, flags, mode=None):
        self.count += 1
        fh = self.count
        self.openfh[fh] = open(file)
        return fh

    def read(self, path, size, offset, fh):
        self.openfh[fh].seek(offset)
        return self.openfh[fh].read(size)

    def release(self, path, fh):
        del self.openfh[fh]
        return

    def readdir(self, path, fh):
        logger.info("Requested ls for %s", path)
        bucket, key = parse_path(path)
        logger.info("Requested ls for bucket %s , key %s", bucket, key)
        try:
            def parse(entry):
                prefix = key[:-1] if key[-1] == '/' else key
                s3_key = entry['Key']
                after_fix = s3_key[len(prefix):]
                if(after_fix[0] == '/'):
                    # show next level
                    return after_fix.split('/')[1]
                else :
                    # finish this level
                    # TODO: bug if key ends with '/' but who would do that!?
                    return prefix.split('/')[-1] + after_fix.split('/')[0]

            items = boto3.client('s3').list_objects_v2(Bucket=bucket,Prefix=key)['Contents']
            items = map(parse, items)
            items = list(set(items))
        except KeyError:
            items = []

        logger.info("Found %s for %s", items, path)

        return ['.', '..'] + items


if __name__ == '__main__':
    if len(argv) != 2:
        print('usage: %s <mountpoint>' % argv[0])
        exit(1)

    logger.info("Starting up s3 fuse at mount point %s", argv[1])
    fuse = FUSE(S3FileSystemMount(),argv[1], foreground=True)
