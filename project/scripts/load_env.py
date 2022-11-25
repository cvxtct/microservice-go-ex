import os 
import pathlib
import re
import boto3
import logging
import time
import json
from typing import List
from botocore.exceptions import ClientError

# Set up our logger
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger()


def get_sts_client():
    return boto3.client('sts')

def get_identity():
    return get_sts_client().get_caller_identity()

def get_ecr_client():
    session = boto3.session.Session(profile_name='default')
    return session.client('ecr', region_name='eu-central-1') 

ecr = get_ecr_client()

def gen_repo_names() -> None:
    """Make docker repository name from project folder names"""
    repos = []
    root_path = str(pathlib.Path(__file__).parents[2])
    for _, dirs, _ in os.walk(root_path):
        for dir in dirs:
                if '-service' in str(dir) or 'front-end' in str(dir):
                    repos.append('experiment' + '/' + str(dir)) 
    
    return repos


def get_repositories(repos) -> List[dict]:
    
    """Acquire info from aws ecr"""
    try:
        response = ecr.describe_repositories(
            registryId = get_identity()["Account"],
            repositoryNames= repos,
        )
        # get rid of javascript date format
        json_string = json.dumps(response, indent=2, default=str)
        json_dict = json.loads(json_string)
       
        return json_dict
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
    project_path = str(pathlib.Path(__file__).parents[1])
    if repositories:
        repo_uris = [[v for k, v in repo.items() if k == 'repositoryUri'] for repo in repositories['repositories']]
        with open(project_path + '/' + '.env', 'w') as f:
            for uri in repo_uris:
                f.write(create_env_var_name(uri[0]) + '=' + uri[0] + '\n')
        f.close()
        logger.info(".env file successfully populated!")
    else:
        raise Exception("Problem with the repositories! None or wrong type!")



if __name__ == '__main__':
    repositories = get_repositories(gen_repo_names())
    
    populate_env_file(repositories)
