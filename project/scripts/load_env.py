import os 
import pathlib
import re
import boto3
import logging
from typing import List
from botocore.exceptions import ClientError

# Set up our logger
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger()

# Set session and client
session = boto3.session.Session(profile_name='default')
ecr = session.client('ecr', region_name='eu-central-1')

# Get account id
account_id = boto3.client('sts').get_caller_identity().get('Account')

# Paths
root_path = str(pathlib.Path(__file__).parents[2])
project_path = str(pathlib.Path(__file__).parents[1])

repos = []

def gen_repo_names() -> None:
    """Make docker repository name from project folder names"""

    for _, dirs, _ in os.walk(root_path):
        for dir in dirs:
                if '-service' in str(dir):
                    repos.append('experiment' + '/' + str(dir))   


def get_repositories() -> List[dict]:
    """Acquire info from aws ecr"""
    try:
        response = ecr.describe_repositories(
            registryId = account_id,
            repositoryNames= repos,
        )
        return response['repositories']
    except ClientError as error:
        if error.response['Error']['Code'] == 'RepositoryNotFoundException':
            logger.exception(error.response['Error']['Message'])
        if error.response['Error']['Code'] == 'InvalidParameterException':
            logger.exception(error.response['Error']['Message'])
        else:
            raise error
    

def create_env_var_name(repo_uri: str) -> str:
    """Generate environment variable name from the docker repository name"""

    return re.sub('-', '_', re.split(r"[\./]", repo_uri)[-1:][0].upper())
    

def populate_env_file(repositories: List[dict]) -> None:
    """Populate .env file with environment variables for the swarm config"""

    if repositories:
        repo_uris = [[v for k, v in repo.items() if k == 'repositoryUri'] for repo in repositories]
        with open(project_path + '/' + '.env', 'w') as f:
            for uri in repo_uris:
                f.write(create_env_var_name(uri[0]) + '=' + uri[0] + '\n')
        f.close()
    else:
        raise Exception("Problem with the repositories! None or wrong type.")



if __name__ == '__main__':
    gen_repo_names()
    repositories = get_repositories()
    populate_env_file(repositories)
