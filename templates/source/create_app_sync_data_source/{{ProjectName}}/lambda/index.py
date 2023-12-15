import json

from aws_lambda_powertools import Logger
from aws_lambda_powertools.utilities.data_classes import AppSyncResolverEvent, event_source
from aws_lambda_powertools.utilities.typing import LambdaContext

logger = Logger()


@event_source(data_class=AppSyncResolverEvent)
def lambda_handler(event: AppSyncResolverEvent, context: LambdaContext):
    logger.info(f"event {json.dumps(event.raw_event)}")
