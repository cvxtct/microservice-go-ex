import os 
import pathlib
import re
import boto3


session = boto3.session.Session(profile_name='default')
ecr = session.client('ecr', region_name='eu-central-1')
account_id = boto3.client('sts').get_caller_identity().get('Account')

root_path = str(pathlib.Path(__file__).parents[2])
project_path = str(pathlib.Path(__file__).parents[1])

repos = []

def gen_repo_names() -> None:
    """Make docker repository name from project folder names"""

    for _, dirs, _ in os.walk(root_path):
        for dir in dirs:
                if '-service' in str(dir):
                    repos.append('experiment' + '/' + str(dir))   


def get_repositories() -> dict:
    """Acquire info from aws ecr"""

    response = ecr.describe_repositories(
        registryId = account_id,
        repositoryNames= repos,
    )
    
    return response['repositories']

def create_env_var_name(repo_uri: str) -> str:
    """Generate environment variable name from the docker repository name"""

    return re.sub('-', '_', re.split(r"[\./]", repo_uri)[-1:][0].upper())
    

def populate_env_file(repositories: dict) -> None:
    """Populate .env file with environment variables for the swarm config"""

    repo_uris = [[v for k, v in repo.items() if k == 'repositoryUri'] for repo in repositories]
    with open(project_path + '/' + '.env', 'w') as f:
         for uri in repo_uris:
            f.write(create_env_var_name(uri[0]) + '=' + uri[0] + '\n')
    f.close()



if __name__ == '__main__':
    gen_repo_names()
    repositories = get_repositories()
    populate_env_file(repositories)
