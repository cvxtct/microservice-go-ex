from unittest.mock import patch
import json
from botocore.session import Session
from botocore.stub import Stubber
import load_env

session = Session()
STS_CLIENT = session.create_client('sts')
ECR_CLIENT = session.create_client('ecr')

f = open('./tests/responses/sts_get_caller_identity_response.json')
sts_get_caller_identity_response = json.load(f)

sts_get_caller_account_id = sts_get_caller_identity_response["Account"]

f = open('./tests/responses/ecr_repositories_response.json')
ecr_repositories = json.load(f)


@patch.object(load_env, 'get_sts_client', return_value=STS_CLIENT)
def test_get_identity_01(get_sts_client):
    STS_STUBBER = Stubber(STS_CLIENT)
    STS_STUBBER.add_response('get_caller_identity', sts_get_caller_identity_response)
    STS_STUBBER.activate()
    load_env_result = load_env.get_identity()
    assert(load_env_result == sts_get_caller_identity_response)
    STS_STUBBER.deactivate()

@patch.object(load_env, 'get_sts_client', return_value=STS_CLIENT)
def test_get_account_id_02(get_sts_client):
    STS_STUBBER = Stubber(STS_CLIENT)
    STS_STUBBER.add_response('get_caller_identity', sts_get_caller_identity_response)
    STS_STUBBER.activate()
    load_env_result = load_env.get_identity()["Account"]
    STS_STUBBER.assert_no_pending_responses()
    assert(load_env_result == sts_get_caller_account_id)
    STS_STUBBER.deactivate()

def test_gen_repo_names():
    expected = [
                'experiment/logger-service', 
                'experiment/broker-service', 
                'experiment/mailer-service', 
                'experiment/authentication-service', 
                'experiment/listener-service', 
                'experiment/front-end'
                ]
    repos = load_env.gen_repo_names()
    assert len(repos) == len(expected)


@patch.object(load_env, 'get_ecr_client', return_value=ECR_CLIENT)
def test_get_repositories_03(get_ecr_client):
    repos = load_env.gen_repo_names()
    with Stubber(ECR_CLIENT) as ECR_STUBBER:
        ECR_STUBBER.add_response('describe_repositories', ecr_repositories)
        load_env_result = load_env.get_repositories(repos)
        assert(load_env_result['repositories'] == ecr_repositories['repositories'])
        