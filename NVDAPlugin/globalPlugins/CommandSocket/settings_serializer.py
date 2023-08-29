import json
from typing import Any
from config import featureFlag


def _serialize_feature_flag(flag: featureFlag.FeatureFlag) -> str:
    return flag.value.name


class NVDASettingsSerializer(json.JSONEncoder):
    def default(self, o) -> Any:
        if isinstance(o, featureFlag.FeatureFlag):
            return _serialize_feature_flag(o)
        return json.JSONEncoder.default(self, o)
