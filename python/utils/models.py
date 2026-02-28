import os
import sys
import random
from typing import Optional, Tuple, Dict, Type, Any

from pydantic import BaseModel, constr, Field

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..')))


class KafkaData(BaseModel):
    ego_id: Optional[str] = None

class CCData(KafkaData):
    mode: str
    ego_speed_new: float
    hostname: str

class SumoEgoData(KafkaData):
    ego_speed: float
    ego_position: Tuple[float, float]
    desired_speed: float
    leader_id: Optional[str]

class SumoLGapData(KafkaData):
    leader_id: Optional[str] = None
    leader_speed: Optional[float] = None
    leader_gap: Optional[float] = None

class SumoLPosData(KafkaData):
    leader_id: Optional[str] = None
    leader_speed: Optional[float] = None
    leader_position: Optional[Tuple[float, float]] = None


def create_model_instance(data: Optional[Dict], model_class: Type[BaseModel]) -> BaseModel:
    # Create a dictionary with None values for all fields in the model_class
    none_data = {field: None for field in model_class.__annotations__.keys()}

    # Update none_data with actual data if provided
    if data:
        none_data.update(data)

    # Return an instance of model_class with the combined data
    return model_class(**none_data)


class TASMessage(BaseModel):
    """
    A class to represent a TAS message using Pydantic.

    This class models a TAS message with both fixed and dynamic fields.
    The message can optionally include a subscriber topic.

    Fixed fields:
    - sender (str): The sender of the message, set to "application.ccam".
    - serviceType (str): The service type, set to "TAS".
    - responseTopic (str): The topic to which responses are sent, set to "application.ccam".

    Dynamic fields:
    - requestId (int): A unique request ID for the message, generated randomly.
    - messageType (str): The type of the message (e.g., "TAS_INIT_REQUEST").
    - message (Dict[str, Any]): The content of the message.

    Optional field:
    - subscriberTopic (str): The topic where the message should be sent if specified.

    Attributes:
    sender (str): The sender of the message.
    serviceType (str): The service type of the message.
    responseTopic (str): The response topic for the message.
    requestId (str): The request ID of the message.
    messageType (str): The type of the message.
    message (Dict[str, Any]): The content of the message.
    subscriberTopic (str): The topic for subscribers if provided.
    """

    sender: constr(strict=True, min_length=1) = Field(default="application.ccam")
    serviceType: constr(strict=True, min_length=1) = Field(default="TAS")
    responseTopic: constr(strict=True, min_length=1) = Field(default="application.ccam")
    subscriberTopic: constr(strict=True, min_length=1) = Field(default="application.ccam")
    requestId: str = Field(default_factory=lambda: str(random.randint(1000, 9999)))
    messageType: str
    message: Dict[str, Any]

    def __repr__(self):
        return (f"TASMessage(sender={self.sender}, serviceType={self.serviceType}, "
                f"responseTopic={self.responseTopic}, requestId={self.requestId}, "
                f"messageType={self.messageType}, message={self.message}, "
                f"subscriberTopic={self.subscriberTopic})")

class CustomMessage(BaseModel):
    """
    A class to represent a custom message using Pydantic.

    This class models a custom message with fixed and dynamic fields.

    Fixed fields:
    - sender (str): The sender of the message, set to "application.ccam".
    - serviceType (str): The service type, set to "CUSTOM".
    - responseTopic (str): The topic to which responses are sent, set to "application.ccam".

    Dynamic fields:
    - requestId (int): A unique request ID for the message, generated randomly.
    - messageType (str): The type of the message (e.g., "CUSTOM_MESSAGE").
    - message (Dict[str, Any]): The content of the message.

    Attributes:
    sender (str): The sender of the message.
    serviceType (str): The service type of the message.
    responseTopic (str): The response topic for the message.
    requestId (str): The request ID of the message.
    messageType (str): The type of the message.
    message (Dict[str, Any]): The content of the message.
    """

    sender: constr(strict=True, min_length=1) = Field(default="application.ccam")
    serviceType: constr(strict=True, min_length=1) = Field(default="CUSTOM")
    responseTopic: constr(strict=True, min_length=1) = Field(default="application.ccam")
    requestId: str = Field(default_factory=lambda: str(random.randint(1000, 9999)))
    messageType: str
    message: Dict[str, Any]

    def __repr__(self):
        return (f"CustomMessage(sender={self.sender}, serviceType={self.serviceType}, "
                f"responseTopic={self.responseTopic}, requestId={self.requestId}, "
                f"messageType={self.messageType}, message={self.message})")