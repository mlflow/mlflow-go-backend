import json
import logging

from mlflow.entities.model_registry import ModelVersion, RegisteredModel
from mlflow.protos.model_registry_pb2 import (
    CreateRegisteredModel,
    DeleteModelVersion,
    DeleteModelVersionTag,
    DeleteRegisteredModel,
    DeleteRegisteredModelAlias,
    DeleteRegisteredModelTag,
    GetLatestVersions,
    GetModelVersion,
    GetModelVersionByAlias,
    GetModelVersionDownloadUri,
    GetRegisteredModel,
    RenameRegisteredModel,
    SetModelVersionTag,
    SetRegisteredModelAlias,
    SetRegisteredModelTag,
    TransitionModelVersionStage,
    UpdateModelVersion,
    UpdateRegisteredModel,
)

from mlflow_go_backend import is_go_enabled
from mlflow_go_backend.lib import get_lib
from mlflow_go_backend.store._service_proxy import _ServiceProxy

_logger = logging.getLogger(__name__)


class _ModelRegistryStore:
    def __init__(self, *args, **kwargs):
        store_uri = args[0] if len(args) > 0 else kwargs.get("db_uri", kwargs.get("root_directory"))
        config = json.dumps(
            {
                "model_registry_store_uri": store_uri,
                "log_level": logging.getLevelName(_logger.getEffectiveLevel()),
            }
        ).encode("utf-8")
        self.service = _ServiceProxy(get_lib().CreateModelRegistryService(config, len(config)))
        super().__init__(store_uri)

    def __del__(self):
        if hasattr(self, "service"):
            get_lib().DestroyModelRegistryService(self.service.id)

    def get_latest_versions(self, name, stages=None):
        request = GetLatestVersions(
            name=name,
            stages=stages,
        )
        response = self.service.call_endpoint(
            get_lib().ModelRegistryServiceGetLatestVersions, request
        )
        return [ModelVersion.from_proto(mv) for mv in response.model_versions]

    def update_registered_model(self, name, description):
        request = UpdateRegisteredModel(name=name, description=description)
        response = self.service.call_endpoint(
            get_lib().ModelRegistryServiceUpdateRegisteredModel, request
        )
        return RegisteredModel.from_proto(response.registered_model)

    def rename_registered_model(self, name, new_name):
        request = RenameRegisteredModel(name=name, new_name=new_name)
        response = self.service.call_endpoint(
            get_lib().ModelRegistryServiceRenameRegisteredModel, request
        )
        return RegisteredModel.from_proto(response.registered_model)

    def delete_registered_model(self, name):
        request = DeleteRegisteredModel(name=name)
        self.service.call_endpoint(get_lib().ModelRegistryServiceDeleteRegisteredModel, request)

    def get_registered_model(self, name):
        request = GetRegisteredModel(name=name)
        response = self.service.call_endpoint(
            get_lib().ModelRegistryServiceGetRegisteredModel, request
        )

        entity = RegisteredModel.from_proto(response.registered_model)
        if entity.description == "":
            entity.description = None

        # during conversion to proto, `version` value became a `string` value.
        # convert it back to `int` value again to satisfy all the Python tests and related logic.
        for key in entity.aliases:
            if entity.aliases[key].isnumeric():
                entity.aliases[key] = int(entity.aliases[key])

        return entity

    def delete_model_version(self, name, version):
        request = DeleteModelVersion(name=name, version=str(version))
        self.service.call_endpoint(get_lib().ModelRegistryServiceDeleteModelVersion, request)

    def get_model_version(self, name, version):
        request = GetModelVersion(name=name, version=str(version))
        response = self.service.call_endpoint(
            get_lib().ModelRegistryServiceGetModelVersion, request
        )
        entity = ModelVersion.from_proto(response.model_version)
        if entity.description == "":
            entity.description = None
        return entity

    def update_model_version(self, name, version, description=None):
        request = UpdateModelVersion(name=name, version=str(version), description=description)
        self.service.call_endpoint(get_lib().ModelRegistryServiceUpdateModelVersion, request)

    def transition_model_version_stage(self, name, version, stage, archive_existing_versions):
        request = TransitionModelVersionStage(
            name=name,
            version=str(version),
            stage=stage,
            archive_existing_versions=archive_existing_versions,
        )
        self.service.call_endpoint(
            get_lib().ModelRegistryServiceTransitionModelVersionStage, request
        )

    def set_registered_model_tag(self, name, tag):
        request = SetRegisteredModelTag(name=name, key=tag.key, value=tag.value)
        self.service.call_endpoint(get_lib().ModelRegistryServiceSetRegisteredModelTag, request)

    def create_registered_model(self, name, tags=None, description=None):
        request = CreateRegisteredModel(
            name=name,
            tags=[tag.to_proto() for tag in tags] if tags else [],
            description=description,
        )
        response = self.service.call_endpoint(
            get_lib().ModelRegistryServiceCreateRegisteredModel, request
        )
        entity = RegisteredModel.from_proto(response.registered_model)
        if not response.registered_model.HasField("description"):
            entity.description = None

        return entity

    def delete_registered_model_tag(self, name, key):
        request = DeleteRegisteredModelTag(name=name, key=key)
        self.service.call_endpoint(get_lib().ModelRegistryServiceDeleteRegisteredModelTag, request)

    def set_registered_model_alias(self, name, alias, version):
        request = SetRegisteredModelAlias(name=name, alias=alias, version=version)
        self.service.call_endpoint(get_lib().ModelRegistryServiceSetRegisteredModelAlias, request)

    def delete_registered_model_alias(self, name, alias):
        request = DeleteRegisteredModelAlias(name=name, alias=alias)
        self.service.call_endpoint(
            get_lib().ModelRegistryServiceDeleteRegisteredModelAlias, request
        )

    def delete_model_version_tag(self, name, version, key):
        request = DeleteModelVersionTag(name=name, version=str(version), key=key)
        self.service.call_endpoint(get_lib().ModelRegistryServiceDeleteModelVersionTag, request)

    def get_model_version_by_alias(self, name, alias):
        request = GetModelVersionByAlias(name=name, alias=alias)
        response = self.service.call_endpoint(
            get_lib().ModelRegistryServiceGetModelVersionByAlias, request
        )
        return ModelVersion.from_proto(response.model_version)

    def set_model_version_tag(self, name, version, tag):
        request = SetModelVersionTag(name=name, version=str(version), key=tag.key, value=tag.value)
        self.service.call_endpoint(get_lib().ModelRegistryServiceSetModelVersionTag, request)

    def get_model_version_download_uri(self, name, version):
        request = GetModelVersionDownloadUri(name=name, version=str(version))
        response = self.service.call_endpoint(
            get_lib().ModelRegistryServiceGetModelVersionDownloadUri, request
        )
        return response.artifact_uri


def ModelRegistryStore(cls):
    return type(cls.__name__, (_ModelRegistryStore, cls), {})


def _get_sqlalchemy_store(store_uri):
    from mlflow.store.model_registry.sqlalchemy_store import SqlAlchemyStore

    if is_go_enabled():
        SqlAlchemyStore = ModelRegistryStore(SqlAlchemyStore)

    return SqlAlchemyStore(store_uri)
